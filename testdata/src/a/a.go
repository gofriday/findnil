package a

func f() error {
	var err error

	if err != nil {
		return nil // want "NG"
	}

	if nil != err {
		return nil // want "NG"
	}

	if "hoge" == "fuga" {
		return nil // OK
	}

	return nil
}

func g() int {
	return 0
}
