package config

type Config struct {
	App struct {
		Name string
		Env  string
	}

	HTTP struct {
		Addr string
	}

	Postgres struct {
		Username string
		Password string
		Host     string
		Port     string
		Database string
	}
}
