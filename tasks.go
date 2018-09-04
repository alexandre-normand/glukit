package main

import (
	"bufio"
	"github.com/alexandre-normand/glukit/app/engine"
	"github.com/alexandre-normand/glukit/app/importer"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"context"
	"google.golang.org/appengine/channel"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
	"os"
)

var processDemoFile = delay.Func("processDemoFile", processStaticDemoFile)

const (
	DATASTORE_WRITES_QUEUE_NAME = "datastore-writes"
)

func disabledUpdateUserData(context context.Context, userEmail string, autoScheduleNextRun bool) {
	// noop
}

// processStaticDemoFile imports the static resource included with the app for the demo user
func processStaticDemoFile(context context.Context, userProfileKey *datastore.Key) {

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

	lastReadTime, err := importer.ParseContent(context, reader, userProfileKey, util.GLUKIT_EPOCH_TIME,
		store.StoreDaysOfReads, store.StoreDaysOfMeals, store.StoreDaysOfInjections, store.StoreDaysOfExercises)

	if err != nil {
		util.Propagate(err)
	}

	store.LogFileImport(context, userProfileKey, model.FileImportLog{Id: "demo", Md5Checksum: "dummychecksum",
		LastDataProcessed: lastReadTime, ImportResult: "Success"})

	if userProfile, err := store.GetUserProfile(context, userProfileKey); err != nil {
		log.Warningf(context, "Error while persisting score for %s: %v", DEMO_EMAIL, err)
	} else {
		if err := engine.StartGlukitScoreBatch(context, userProfile); err != nil {
			log.Warningf(context, "Error while starting batch calculation of glukit scores for %s: %v", DEMO_EMAIL, err)
		}

		err = engine.StartA1CCalculationBatch(context, userProfile)
		if err != nil {
			log.Warningf(context, "Error starting a1c calculation batch for user [%s]: %v", DEMO_EMAIL, err)
		}
	}

	channel.Send(context, DEMO_EMAIL, "Refresh")
}
