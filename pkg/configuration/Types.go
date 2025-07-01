package configuration

type Configuration struct {
	AllowOrigin string
	MasterPort  string
	Port        string
	Certificate string
	Key         string
	Environment string
}

const (
	DEVELOPMENT_ENV = "develop"
	PRODUCTION_ENV  = "production"
)
