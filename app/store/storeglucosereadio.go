package store

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type DataStoreGlucoseReadBatchWriter struct {
	c context.Context
	k *datastore.Key
}

// NewDataStoreGlucoseReadBatchWriter creates a new GlucoseReadBatchWriter that persists to the datastore
func NewDataStoreGlucoseReadBatchWriter(context context.Context, userProfileKey *datastore.Key) *DataStoreGlucoseReadBatchWriter {
	w := new(DataStoreGlucoseReadBatchWriter)
	w.c = context
	w.k = userProfileKey
	return w
}

func (w *DataStoreGlucoseReadBatchWriter) WriteGlucoseReadBatches(p []apimodel.DayOfGlucoseReads) (glukitio.GlucoseReadBatchWriter, error) {
	if _, err := StoreDaysOfReads(w.c, w.k, p); err != nil {
		return w, err
	} else {
		return w, nil
	}
}

func (w *DataStoreGlucoseReadBatchWriter) WriteGlucoseReadBatch(p []apimodel.GlucoseRead) (glukitio.GlucoseReadBatchWriter, error) {
	dayOfGlucoseReads := make([]apimodel.DayOfGlucoseReads, 1)
	dayOfGlucoseReads[0] = apimodel.NewDayOfGlucoseReads(p)
	return w.WriteGlucoseReadBatches(dayOfGlucoseReads)
}

func (w *DataStoreGlucoseReadBatchWriter) Flush() (glukitio.GlucoseReadBatchWriter, error) {
	return w, nil
}
