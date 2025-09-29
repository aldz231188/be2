package db

type Config struct {
	DSN string
}

// TODO: Не хардкодить DSN в коде: подними из env/файла в другом месте (например, infra/config).

func NewPGConfig() Config {
	return Config{
		DSN: "postgres://postgres:Qwaszx_1@localhost:5432/shopdb",
	}
}
