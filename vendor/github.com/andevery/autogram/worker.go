package autogram

type Reporter interface {
	Report(id int64, report map[string]int)
	Error(id int64, err error)
	Fatal(id int64, err error)
}

type BackgroundTask interface {
	Start()
	Stop()
}
