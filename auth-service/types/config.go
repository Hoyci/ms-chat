package types

type Config struct {
	Port                          int    `env:"PORT" envDefault:"8080"`
	Environment                   string `env:"ENVIRONMENT" envDefault:"development"`
	DatabaseURL                   string `env:"DATABASE_URL" envDefault:"postgresql://user:password@localhost:5432/postgres?sslmode=disable"`
	AccessJWTSecret               string `env:"ACCESS_JWT_SECRET" envDefault:"UM_ACCESS_TOKEN_MTO_DIFICIL"`
	AccessJWTExpirationInSeconds  int    `env:"ACCESS_JWT_EXPIRATION" envDefault:"604800"`
	RefreshJWTSecret              string `env:"REFRESH_JWT_SECRET" envDefault:"UM_REFRESH_TOKEN_MTO_DIFICIL"`
	RefreshJWTExpirationInSeconds int    `env:"REFRESH_JWT_EXPIRATION" envDefault:"9072000"`
}
