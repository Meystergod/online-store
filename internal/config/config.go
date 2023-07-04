package config

type Config struct {
	Log struct {
		LogLevel string `envconfig:"LOG_LEVEL" default:"debug"`
		LogFile  string `envconfig:"LOG_FILE"`
		LogSize  int    `envconfig:"LOG_SIZE" default:"10"`
		LogAge   int    `envconfig:"LOG_AGE" default:"28"`
	}

	HTTPServer struct {
		Address string `envconfig:"HTTP_ADDR" default:"0.0.0.0:8000"`
	}

	Database struct {
		Host     string `envconfig:"DB_HOST" default:"localhost"`
		Port     string `envconfig:"DB_PORT" default:"27017"`
		Username string `envconfig:"DB_USERNAME"`
		Password string `envconfig:"DB_PASSWORD"`
		Auth     string `envconfig:"DB_AUTH"`
		Name     string `envconfig:"DB_NAME" default:"onlinestoredb"`
	}

	Application struct {
		Name    string `envconfig:"APP_VERSION" default:"online store"`
		Version string `envconfig:"APP_VERSION" default:"v0.0.1"`
	}
}
