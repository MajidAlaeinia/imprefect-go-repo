package exceptions

func PanicException(err error) {
	if err != nil {
		panic(err)
	}
}
