package config

type Configuration struct {
	AppName string
}

var AppConfig Configuration = Configuration{
	AppName: "shelf",
}
