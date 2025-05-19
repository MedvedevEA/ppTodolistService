package config

type Config struct {
	Environment string
	Server      Server
	Store       Store
}
type Store struct {
}
type Server struct {
}

func MustNew() *Config {
	return nil
}
