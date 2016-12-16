package store

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type DataStoreInjectionBatchWriter struct {
	c context.Context
	k *datastore.Key
}

// NewDataStoreInjectionBatchWriter creates a new InjectionBatchWriter that persists to the datastore
func NewDataStoreInjectionBatchWriter(context context.Context, userProfileKey *datastore.Key) *DataStoreInjectionBatchWriter {
	w := new(DataStoreInjectionBatchWriter)
	w.c = context
	w.k = userProfileKey
	return w
}

func (w *DataStoreInjectionBatchWriter) WriteInjectionBatches(p []apimodel.DayOfInjections) (glukitio.InjectionBatchWriter, error) {
	if _, err := StoreDaysOfInjections(w.c, w.k, p); err != nil {
		return w, err
	} else {
		return w, nil
	}
}

func (w *DataStoreInjectionBatchWriter) WriteInjectionBatch(p []apimodel.Injection) (glukitio.InjectionBatchWriter, error) {
	dayOfInjections := make([]apimodel.DayOfInjections, 1)
	dayOfInjections[0] = apimodel.NewDayOfInjections(p)
	return w.WriteInjectionBatches(dayOfInjections)
}

func (w *DataStoreInjectionBatchWriter) Flush() (glukitio.InjectionBatchWriter, error) {
	return w, nil
}
