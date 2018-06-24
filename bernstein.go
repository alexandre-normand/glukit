package main

import (
	"bytes"
	"fmt"
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/engine"
	"github.com/alexandre-normand/glukit/app/importer"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	GLUKIT_BERNSTEIN_EMAIL = "dr.bernstein@glukit.com"
	PERFECT_SCORE          = 83
)

var BERNSTEIN_EARLIEST_READ, _ = time.Parse(util.TIMEFORMAT_NO_TZ, "2014-06-01 12:00:00")
var BERNSTEIN_MOST_RECENT_READ_TIME, _ = time.Parse(util.TIMEFORMAT_NO_TZ, "2015-01-01 12:00:00")
var BERNSTEIN_MOST_RECENT_READ = apimodel.GlucoseRead{apimodel.Time{BERNSTEIN_EARLIEST_READ.Unix(), "America/New_York"}, apimodel.MG_PER_DL, PERFECT_SCORE}
var BERNSTEIN_BIRTH_DATE, _ = time.Parse(util.TIMEFORMAT_NO_TZ, "1934-06-17 00:00:00")

// initializeGlukitBernstein does lazy initialization of the "perfect" glukit user.
// It's called Glukit Bernstein because much of this comes from Dr. Berstein himself.
func initializeGlukitBernstein(writer http.ResponseWriter, reader *http.Request) {
	context := appengine.NewContext(reader)

	_, _, _, err := store.GetUserData(context, GLUKIT_BERNSTEIN_EMAIL)
	if err == datastore.ErrNoSuchEntity {
		log.Infof(context, "No data found for glukit bernstein user [%s], creating it", GLUKIT_BERNSTEIN_EMAIL)
		userProfileKey, err := store.StoreUserProfile(context, time.Now(),
			model.GlukitUser{GLUKIT_BERNSTEIN_EMAIL, "Glukit", "Bernstein", BERNSTEIN_BIRTH_DATE, model.DIABETES_TYPE_1, "America/New_York", time.Now(),
				BERNSTEIN_MOST_RECENT_READ, model.UNDEFINED_SCORE, model.UNDEFINED_SCORE, true, "", time.Now(), model.UNDEFINED_A1C_ESTIMATE})
		if err != nil {
			util.Propagate(err)
		}

		fileReader := generateBernsteinData(context)
		lastReadTime, err := importer.ParseContent(context, fileReader, userProfileKey, util.GLUKIT_EPOCH_TIME,
			store.StoreDaysOfReads, store.StoreDaysOfMeals, store.StoreDaysOfInjections, store.StoreDaysOfExercises)

		if err != nil {
			util.Propagate(err)
		}

		store.LogFileImport(context, userProfileKey, model.FileImportLog{Id: "bernstein", Md5Checksum: "dummychecksum",
			LastDataProcessed: lastReadTime, ImportResult: "Success"})

		if glukitUser, err := store.GetUserProfile(context, userProfileKey); err != nil {
			log.Warningf(context, "Error getting retrieving GlukitUser [%s], this needs attention: [%v]", GLUKIT_BERNSTEIN_EMAIL, err)
		} else {
			// Start batch calculation of the glukit scores
			err := engine.StartGlukitScoreBatch(context, glukitUser)

			if err != nil {
				log.Warningf(context, "Error starting batch calculation of GlukitScores for [%s], this needs attention: [%v]", GLUKIT_BERNSTEIN_EMAIL, err)
			}

			err = engine.StartA1CCalculationBatch(context, glukitUser)
			if err != nil {
				log.Warningf(context, "Error starting batch calculation of a1cs for [%s], this needs attention: [%v]", GLUKIT_BERNSTEIN_EMAIL, err)
			}
		}
	} else if err != nil {
		util.Propagate(err)
	} else {
		log.Infof(context, "Data already stored for user [%s], continuing...", GLUKIT_BERNSTEIN_EMAIL)
	}
}

// generateBernsteinData generates an in-memory dexcom file for the user Glukit Bernstein.
func generateBernsteinData(context context.Context) (reader io.Reader) {
	buffer := new(bytes.Buffer)
	buffer.WriteString("<Patient Id=\"{E1B2FE4C-35F0-40B8-A15A-D3CBCA27B666}\" SerialNumber=\"sm11111111\" IsDataBlinded=\"0\" IsKeepPrivate=\"1\">\n")
	buffer.WriteString("<MeterReadings></MeterReadings>\n")
	buffer.WriteString("<GlucoseReadings>\n")

	startTime := BERNSTEIN_EARLIEST_READ
	endTime := BERNSTEIN_MOST_RECENT_READ_TIME

	log.Debugf(context, "Data for bernstein from %s to %s:", startTime.In(time.UTC).Format(util.TIMEFORMAT_NO_TZ),
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
