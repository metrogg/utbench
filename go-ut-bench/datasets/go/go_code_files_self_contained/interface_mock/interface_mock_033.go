package main

func Monitor() {
	if err := recover(); err != nil {
		DefaultClient.Notify(newError(err, 2))
		panic(err)
	}
}
