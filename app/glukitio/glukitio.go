package glukitio

import (
	"errors"
	"github.com/alexandre-normand/glukit/app/model"
)

// ErrShortWrite means that a write accepted fewer bytes than requested
// but failed to return an explicit error.
var ErrShortWrite = errors.New("short write")

// CalibrationBatchWriter is the interface that wraps the basic
// WriteCalibrationBatch and WriteCalibrationBatches methods.
//
// WriteCalibrationBatch writes len(p) model.CalibrationRead from p to the
// underlying data stream. It returns the number of elements written
// from p (0 <= n <= len(p)) and any error encountered that caused the
// write to stop early. Write must return a non-nil error if it returns n < len(p).
//
// WriteCalibrationBatches writes len(p) model.DayOfCalibrationReads from p to the
// underlying data stream. It returns the number of batch elements written
// from p (0 <= n <= len(p)) and any error encountered that caused the
// write to stop early. Write must return a non-nil error if it returns n < len(p).
type CalibrationBatchWriter interface {
	WriteCalibrationBatch(p []model.CalibrationRead) (w CalibrationBatchWriter, err error)
	WriteCalibrationBatches(p []model.DayOfCalibrationReads) (w CalibrationBatchWriter, err error)
	Flush() (w CalibrationBatchWriter, err error)
}

// GlucoseReadBatchWriter is the interface that wraps the basic
// WriteGlucoseReadBatch and WriteGlucoseReadBatches methods.
//
// WriteGlucoseReadBatch writes len(p) model.CalibrationRead from p to the
// underlying data stream. It returns the number of elements written
// from p (0 <= n <= len(p)) and any error encountered that caused the
// write to stop early. Write must return a non-nil error if it returns n < len(p).
//
// WriteGlucoseReadBatches writes len(p) model.DayOfCalibrationReads from p to the
// underlying data stream. It returns the number of batch elements written
// from p (0 <= n <= len(p)) and any error encountered that caused the
// write to stop early. Write must return a non-nil error if it returns n < len(p).
type GlucoseReadBatchWriter interface {
	WriteGlucoseReadBatch(p []model.GlucoseRead) (w GlucoseReadBatchWriter, err error)
	WriteGlucoseReadBatches(p []model.DayOfGlucoseReads) (w GlucoseReadBatchWriter, err error)
	Flush() (w GlucoseReadBatchWriter, err error)
}

// InjectionBatchWriter is the interface that wraps the basic
// WriteInjectionBatch and WriteInjectionBatches methods.
//
// WriteInjectionBatch writes len(p) model.Injection from p to the
// underlying data stream. It returns the number of elements written
// from p (0 <= n <= len(p)) and any error encountered that caused the
// write to stop early. Write must return a non-nil error if it returns n < len(p).
//
// WriteInjectionBatches writes len(p) model.DayOfInjectionss from p to the
// underlying data stream. It returns the number of batch elements written
// from p (0 <= n <= len(p)) and any error encountered that caused the
// write to stop early. Write must return a non-nil error if it returns n < len(p).
type InjectionBatchWriter interface {
	WriteInjectionBatch(p []model.Injection) (w InjectionBatchWriter, err error)
	WriteInjectionBatches(p []model.DayOfInjections) (w InjectionBatchWriter, err error)
	Flush() (w InjectionBatchWriter, err error)
}

// CarbBatchWriter is the interface that wraps the basic
// WriteCarbBatch and WriteCarbBatches methods.
//
// WriteCarbBatch writes len(p) model.Carb from p to the
// underlying data stream. It returns the number of elements written
// from p (0 <= n <= len(p)) and any error encountered that caused the
// write to stop early. Write must return a non-nil error if it returns n < len(p).
//
// WriteCarbBatches writes len(p) model.DayOfCarbss from p to the
// underlying data stream. It returns the number of batch elements written
// from p (0 <= n <= len(p)) and any error encountered that caused the
// write to stop early. Write must return a non-nil error if it returns n < len(p).
type CarbBatchWriter interface {
	WriteCarbBatch(p []model.Carb) (n int, err error)
	WriteCarbBatches(p []model.DayOfCarbs) (n int, err error)
	Flush() error
}

// ExerciseBatchWriter is the interface that wraps the basic
// WriteExerciseBatch and WriteExerciseBatches methods.
//
// WriteExerciseBatch writes len(p) model.Exercise from p to the
// underlying data stream. It returns the number of elements written
// from p (0 <= n <= len(p)) and any error encountered that caused the
// write to stop early. Write must return a non-nil error if it returns n < len(p).
//
// WriteExerciseBatches writes len(p) model.DayOfExercisess from p to the
// underlying data stream. It returns the number of batch elements written
// from p (0 <= n <= len(p)) and any error encountered that caused the
// write to stop early. Write must return a non-nil error if it returns n < len(p).
type ExerciseBatchWriter interface {
	WriteExerciseBatch(p []model.Exercise) (n int, err error)
	WriteExerciseBatches(p []model.DayOfExercises) (n int, err error)
	Flush() error
}
