package config

type Config struct {
	App struct {
		Name string
		Env  string
	}

	HTTP struct {
		Host string
		Port string
	}

	Mongo struct {
		URI  string
		Name string
	}
}
