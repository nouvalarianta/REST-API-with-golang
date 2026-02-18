package config

type Config struct {
	Server      Server
	DatabaseURL string
	Jwt         Jwt
	Storage     Storage
}

type Server struct {
	Host  string
	Port  string
	Asset string
}

type Jwt struct {
	Key string
	Exp int
}

type Storage struct {
	BasePath string
}
