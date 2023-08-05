package repositoryimpl

type Config struct {
	Table Table `json:"table" required:"true"`
}

type Table struct {
	AsyncTask string `json:"async_task" required:"true"`
	Access    string `json:"access"     required:"true"`
}
