package shared

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
