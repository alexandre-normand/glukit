package fetcher

import (
	"net/http"
	"fmt"
	"log"
	"drive"
	"strings"
	"time"
	"timeutils"
	"io"
)

func SearchDataFiles(client *http.Client, lastUpdate time.Time) (file []*drive.File, err error) {
	var files []*drive.File

	if service, err := drive.New(client); err != nil {
		return nil, err
	} else {
		query := fmt.Sprintf("fullText contains \"<Glucose\" and fullText contains \"<Patient Id=\" and modifiedDate > '%s'", lastUpdate.Format(timeutils.DRIVE_TIMEFORMAT))
		call := service.Files.List().MaxResults(10).Q(query)
		if filelist, err := call.Do(); err != nil {
			return nil, err
		} else {
			for i := range filelist.Items {
				file := filelist.Items[i]
				if strings.HasSuffix(file.OriginalFilename, ".xml") {
					log.Printf("Found match: %s\n", file)
					files = append(files, file)
				} else {
					log.Printf("Skipping search result item: %s\n", file)
				}
			}
		}
	}

	return files, nil
}

// Caller should call Close() when done
func GetFileReader(client http.RoundTripper, file *drive.File) (reader io.ReadCloser, err error) {
	// t parameter should use an oauth.Transport
	downloadUrl := file.DownloadUrl
	if downloadUrl == "" {
		// If there is no downloadUrl, there is no body
		log.Printf("An error occurred: File is not downloadable")
		return nil, nil
	}
	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return nil, err
	}
	resp, err := client.RoundTrip(req)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return nil, err
	}

	return resp.Body, nil
}
