package importer

import (
	"github.com/alexandre-normand/glukit/app/util"
	"appengine"
	"fmt"
	"io"
	"github.com/alexandre-normand/glukit/lib/drive"
	"log"
	"net/http"
	"strings"
	"time"
)

// SearchDataFiles does a search on GoogleDrive for any file that look like it's a Dexcom xml file.
// The search is restricted to files that have a modified date after the given last update time.
func SearchDataFiles(client *http.Client, lastUpdate time.Time) (file []*drive.File, err error) {
	var files []*drive.File

	if service, err := drive.New(client); err != nil {
		return nil, err
	} else {
		query := fmt.Sprintf("fullText contains \"<Glucose\" and fullText contains \"<Patient Id=\" and trashed=false and modifiedDate > '%s'", lastUpdate.Format(util.DRIVE_TIMEFORMAT))
		call := service.Files.List().MaxResults(100).Q(query)
		if filelist, err := call.Do(); err != nil {
			return nil, err
		} else {
			for i := range filelist.Items {
				file := filelist.Items[i]
				if strings.HasSuffix(file.OriginalFilename, ".xml") {
					log.Printf("Found match with id [%s] and original filename [%s]\n", file.Id, file.OriginalFilename)
					files = append(files, file)
				} else {
					log.Printf("Skipping search result item: %s\n", file)
				}
			}
		}
	}

	return files, nil
}

// GetFileReader returns the file reader for the GoogleDrive file. The caller is responsible for calling Close() when done.
func GetFileReader(context appengine.Context, client http.RoundTripper, file *drive.File) (reader io.ReadCloser, err error) {
	// t parameter should use an oauth.Transport
	downloadUrl := file.DownloadUrl
	if downloadUrl == "" {
		// If there is no downloadUrl, there is no body
		context.Errorf("An error occurred: File is not downloadable")
		return nil, nil
	}
	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		context.Errorf("An error occurred: %v\n", err)
		return nil, err
	}
	// Request for compressed files to make download faster
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("User-agent", "glukit (gzip)")

	resp, err := client.RoundTrip(req)
	if err != nil {
		context.Errorf("An error occurred: %v\n", err)
		return nil, err
	}

	return resp.Body, nil
}
