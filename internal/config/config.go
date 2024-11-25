package config

// import config file
// database (SQLITE)

type Config struct {
	DBPath string
}

func (c *Config) New() *Config {
	return &Config{
		DBPath: c.DBPath,
	}
}
