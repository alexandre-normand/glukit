package util

// Propagate aborts the current execution if err is non-nil.
func Propagate(err error) {
	if err != nil {
		panic(err)
	}
}
