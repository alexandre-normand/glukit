package glukit

import (
	"app/engine"
	"app/importer"
	"app/model"
	"app/store"
	"app/util"
	"appengine"
	"appengine/datastore"
	"bytes"
	"fmt"
	"io"
	"lib/goauth2/oauth"
	"net/http"
	"strings"
	"time"
)

const (
	GLUKIT_BERNSTEIN_EMAIL = "dr.bernstein@glukit.com"
	PERFECT_SCORE          = 83
)

var BERNSTEIN_EARLIEST_READ, _ = time.Parse(util.TIMEFORMAT_NO_TZ, "2013-06-01 12:00:00")
var BERNSTEIN_MOST_RECENT_READ_TIME, _ = time.Parse(util.TIMEFORMAT_NO_TZ, "2014-03-11 12:00:00")
var BERNSTEIN_MOST_RECENT_READ = model.GlucoseRead{BERNSTEIN_EARLIEST_READ.Format(util.TIMEFORMAT), model.Timestamp(BERNSTEIN_EARLIEST_READ.Unix()), PERFECT_SCORE}
var BERNSTEIN_BIRTH_DATE, _ = time.Parse(util.TIMEFORMAT_NO_TZ, "1934-06-17 00:00:00")

// initializeGlukitBernstein does lazy initialization of the "perfect" glukit user.
// It's called Glukit Bernstein because much of this comes from Dr. Berstein himself.
func initializeGlukitBernstein(writer http.ResponseWriter, reader *http.Request) {
	context := appengine.NewContext(reader)

	_, _, _, err := store.GetUserData(context, GLUKIT_BERNSTEIN_EMAIL)
	if err == datastore.ErrNoSuchEntity {
		context.Infof("No data found for glukit bernstein user [%s], creating it", GLUKIT_BERNSTEIN_EMAIL)
		dummyToken := oauth.Token{"", "", util.GLUKIT_EPOCH_TIME}
		userProfileKey, err := store.StoreUserProfile(context, time.Now(),
			model.GlukitUser{GLUKIT_BERNSTEIN_EMAIL, "Glukit", "Bernstein", BERNSTEIN_BIRTH_DATE, model.DIABETES_TYPE_1, "America/New_York", time.Now(),
				BERNSTEIN_MOST_RECENT_READ, dummyToken, "", model.UNDEFINED_SCORE, model.UNDEFINED_SCORE, true, ""})
		if err != nil {
			util.Propagate(err)
		}

		reader := generateBernsteinData(context)
		lastReadTime := importer.ParseContent(context, reader, importer.IMPORT_BATCH_SIZE, userProfileKey, util.GLUKIT_EPOCH_TIME,
			store.StoreDaysOfReads, store.StoreDaysOfCarbs, store.StoreDaysOfInjections, store.StoreDaysOfExercises)
		store.LogFileImport(context, userProfileKey, model.FileImportLog{Id: "bernstein", Md5Checksum: "dummychecksum",
			LastDataProcessed: lastReadTime})

		if glukitUser, err := store.GetUserProfile(context, userProfileKey); err != nil {
			context.Warningf("Error getting retrieving GlukitUser [%s], this needs attention: [%v]", GLUKIT_BERNSTEIN_EMAIL, err)
		} else {
			// Start batch calculation of the glukit scores
			err := engine.CalculateGlukitScoreBatch(context, glukitUser)

			if err != nil {
				context.Warningf("Error starting batch calculation of GlukitScores for [%s], this needs attention: [%v]", GLUKIT_BERNSTEIN_EMAIL, err)
			}
		}
	} else if err != nil {
		util.Propagate(err)
	} else {
		context.Infof("Data already stored for user [%s], continuing...", GLUKIT_BERNSTEIN_EMAIL)
	}

	writer.Header().Set("Content-Type", "text/plain")
	writer.Write([]byte("dr.bernstein@glukit.com has been initialized.\n"))
}

// generateBernsteinData generates an in-memory dexcom file for the user Glukit Bernstein.
func generateBernsteinData(context appengine.Context) (reader io.Reader) {
	buffer := new(bytes.Buffer)
	buffer.WriteString("<Patient Id=\"{E1B2FE4C-35F0-40B8-A15A-D3CBCA27B666}\" SerialNumber=\"sm11111111\" IsDataBlinded=\"0\" IsKeepPrivate=\"1\">\n")
	buffer.WriteString("<MeterReadings></MeterReadings>\n")
	buffer.WriteString("<GlucoseReadings>\n")

	startTime := BERNSTEIN_EARLIEST_READ
	endTime := BERNSTEIN_MOST_RECENT_READ_TIME

	context.Debugf("Data for bernstein from %s to %s:", startTime.In(time.UTC).Format(util.TIMEFORMAT_NO_TZ),
		endTime.In(time.UTC).Format(util.TIMEFORMAT_NO_TZ))
	for currentTime := startTime; !currentTime.After(endTime); currentTime = currentTime.Add(time.Duration(5 * time.Minute)) {
		line := fmt.Sprintf("<Glucose InternalTime=\"%s\" DisplayTime=\"%s\" Value=\"%d\"/>\n",
			currentTime.In(time.UTC).Format(util.TIMEFORMAT_NO_TZ), currentTime.In(time.UTC).Format(util.TIMEFORMAT_NO_TZ),
			PERFECT_SCORE)
		buffer.WriteString(line)
	}

	buffer.WriteString("</GlucoseReadings>\n")
	buffer.WriteString("<EventMarkers></EventMarkers>\n")
	buffer.WriteString("</Patient>\n")

	return strings.NewReader(buffer.String())
}
