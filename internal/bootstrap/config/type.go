package config

type Config struct {
	App struct {
		Name string
		Env  string
	}

	HTTP struct {
		Addr string
	}

	Mongo struct {
		URI  string
		Name string
	}
}
