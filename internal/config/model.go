package config

type Config struct {
	Server      Server
	DatabaseURL string
	Jwt         Jwt
}

type Server struct {
	Host string
	Port string
}

type Jwt struct {
	Key string
	Exp int
}
