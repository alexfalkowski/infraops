package runtime

// Must panics if we have an error.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// ConvertRecover to an error.
func ConvertRecover(rec any) error {
	return rec.(error)
}
