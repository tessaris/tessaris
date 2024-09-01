package config

type Config struct {
	Port   int
	Routes []interface{}
}

func CreateConfig() *Config {
	return &Config{
		Port:   3000,
		Routes: []interface{}{},
	}
}
