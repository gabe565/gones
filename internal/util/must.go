package util

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Must2[T any](val T, err error) T { //nolint:ireturn,nolintlint
	if err != nil {
		panic(err)
	}
	return val
}
