package glukit

import (
	"app/engine"
	"app/importer"
	"app/model"
	"app/store"
	"app/util"
	"appengine"
	"appengine/channel"
	"appengine/datastore"
	"appengine/delay"
	"appengine/taskqueue"
	"appengine/urlfetch"
	"bufio"
	"lib/drive"
	"lib/goauth2/oauth"
	"os"
	"time"
)

var processFile = delay.Func("processSingleFile", processSingleFile)
var processDemoFile = delay.Func("processDemoFile", processStaticDemoFile)
var refreshUserData = delay.Func(REFRESH_USER_DATA_FUNCTION_NAME, func(context appengine.Context, userEmail string,
	autoScheduleNextRun bool) {
	context.Criticalf("This function purely exists as a workaround to the \"initialization loop\" error that " +
		"shows up because the function calls itself. This implementation defines the same signature as the " +
		"real one which we define in init() to override this implementation!")
})

const (
	REFRESH_USER_DATA_FUNCTION_NAME = "refreshUserData"
)

// updateUserData is an async task that searches on Google Drive for dexcom files. It handles some high
// watermark of the last import to avoid downloading already imported files (unless they've been updated).
// It also schedules itself to run again the next day unless the token is invalid.
func updateUserData(context appengine.Context, userEmail string, autoScheduleNextRun bool) {
	glukitUser, userProfileKey, _, _, err := store.GetUserData(context, userEmail, model.DEFAULT_LOOKBACK_PERIOD)
	if _, ok := err.(store.StoreError); err != nil && !ok {
		context.Errorf("We're trying to run an update data task for user [%s] that doesn't exist. "+
			"Got error: %v", userEmail, err)
		return
	}

	transport := &oauth.Transport{
		Config: config(),
		Transport: &urlfetch.Transport{
			Context: context,
		},
		Token: &glukitUser.Token,
	}

	// If the token is expired, try to get a fresh one by doing a refresh (which should use the refresh_token
	if glukitUser.Token.Expired() {
		transport.Token.RefreshToken = glukitUser.RefreshToken
		err := transport.Refresh(context)
		if err != nil {
			context.Errorf("Error updating token for user [%s], let's hope he comes back soon so we can "+
				"get a fresh token: %v", userEmail, err)
			return
		}

		// Update the user with the new token
		context.Infof("Token refreshed, updating user [%s] with token [%v]", userEmail, glukitUser.Token)
		store.StoreUserProfile(context, time.Now(), *glukitUser)
	}

	// Next update in one day
	nextUpdate := time.Now().AddDate(0, 0, 1)
	files, err := importer.SearchDataFiles(transport.Client(), glukitUser.MostRecentRead.GetTime())
	if err != nil {
		context.Warningf("Error while searching for files on google drive for user [%s]: %v", userEmail, err)
	} else {
		switch {
		case len(files) == 0:
			context.Infof("No new or updated data found for existing user [%s]", userEmail)
		case len(files) > 0:
			context.Infof("Found new data files for user [%s], downloading and storing...", userEmail)
			processFileSearchResults(&glukitUser.Token, files, context, userEmail, userProfileKey)
		}
	}

	engine.CalculateGlukitScoreBatch(context, glukitUser)

	if autoScheduleNextRun {
		task, err := refreshUserData.Task(userEmail, autoScheduleNextRun)
		if err != nil {
			context.Criticalf("Couldn't schedule the next execution of the data refresh for user [%s]. "+
				"This breaks background updating of user data!: %v", userEmail, err)
		}
		task.ETA = nextUpdate
		taskqueue.Add(context, task, "refresh")

		context.Infof("Scheduled next data update for user [%s] at [%s]", userEmail, nextUpdate.Format(util.TIMEFORMAT))
	} else {
		context.Infof("Not scheduling a the next refresh as requested by autoScheduleNextRun [%t]", autoScheduleNextRun)
	}
}

// processFileSearchResults reads the list of files detected on google drive and kicks off a new queued task
// to process each one
func processFileSearchResults(token *oauth.Token, files []*drive.File, context appengine.Context, userEmail string,
	userProfileKey *datastore.Key) {
	// TODO : Look at recent file import log for that file and skip to the new data. It would be nice to be able to
	// use the Http Range header but that's unlikely to be possible since new event/read data is spreadout in the
	// file
	for i := range files {
		task, err := processFile.Task(token, files[i], userEmail, userProfileKey)
		if err != nil {
			util.Propagate(err)
		}
		taskqueue.Add(context, task, "store")
	}
}

// processSingleFile handles the import of a single file. It deals with:
//    1. Logging the file import operation
//    2. Calculating and updating the new GlukitScore
//    3. Sending a "refresh" message to any connected client
func processSingleFile(context appengine.Context, token *oauth.Token, file *drive.File, userEmail string,
	userProfileKey *datastore.Key) {
	t := &oauth.Transport{
		Config: config(),
		Transport: &urlfetch.Transport{
			Context: context,
		},
		Token: token,
	}

	reader, err := importer.GetFileReader(context, t, file)
	if err != nil {
		context.Infof("Error reading file %s, skipping: [%v]", file.OriginalFilename, err)
	} else {
		// Default to beginning of time
		startTime := util.GLUKIT_EPOCH_TIME
		if lastFileImportLog, err := store.GetFileImportLog(context, userProfileKey, file.Id); err == nil {
			startTime = lastFileImportLog.LastDataProcessed
			context.Infof("Reloading data from file [%s]-[%s] starting at date [%s]...", file.Id,
				file.OriginalFilename, startTime.Format(util.TIMEFORMAT))
		} else if err == datastore.ErrNoSuchEntity {
			context.Debugf("First import of file [%s]-[%s]...", file.Id, file.OriginalFilename)
		} else if err != nil {
			util.Propagate(err)
		}

		lastReadTime := importer.ParseContent(context, reader, importer.IMPORT_BATCH_SIZE, userProfileKey, startTime,
			store.StoreDaysOfReads, store.StoreDaysOfCarbs, store.StoreDaysOfInjections, store.StoreDaysOfExercises)
		store.LogFileImport(context, userProfileKey, model.FileImportLog{Id: file.Id, Md5Checksum: file.Md5Checksum,
			LastDataProcessed: lastReadTime})
		reader.Close()

		if glukitUser, err := store.GetUserProfile(context, userProfileKey); err != nil {
			context.Warningf("Error getting retrieving GlukitUser [%s], this needs attention: [%v]", userEmail, err)
		} else {
			// Calculate Glukit Score batch here for the newly imported data
			err := engine.CalculateGlukitScoreBatch(context, glukitUser)
			if err != nil {
				context.Warningf("Error starting batch calculation of GlukitScores for [%s], this needs attention: [%v]", userEmail, err)
			}
		}
	}
	channel.Send(context, userEmail, "Refresh")
}

// processStaticDemoFile imports the static resource included with the app for the demo user
func processStaticDemoFile(context appengine.Context, userProfileKey *datastore.Key) {

	// open input file
	fi, err := os.Open("data.xml")
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if fi.Close() != nil {
			panic(err)
		}
	}()
	// make a read buffer
	reader := bufio.NewReader(fi)

	lastReadTime := importer.ParseContent(context, reader, importer.IMPORT_BATCH_SIZE, userProfileKey, util.GLUKIT_EPOCH_TIME,
		store.StoreDaysOfReads, store.StoreDaysOfCarbs, store.StoreDaysOfInjections, store.StoreDaysOfExercises)
	store.LogFileImport(context, userProfileKey, model.FileImportLog{Id: "demo", Md5Checksum: "dummychecksum",
		LastDataProcessed: lastReadTime})

	if userProfile, err := store.GetUserProfile(context, userProfileKey); err != nil {
		context.Warningf("Error while persisting score for %s: %v", DEMO_EMAIL, err)
	} else {
		if err := engine.CalculateGlukitScoreBatch(context, userProfile); err != nil {
			context.Warningf("Error while starting batch calculation of glukit scores for %s: %v", DEMO_EMAIL, err)
		}
	}

	channel.Send(context, DEMO_EMAIL, "Refresh")
}
