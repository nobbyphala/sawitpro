package constant

import "os"

var (
	EnvJWTSecretKey     = os.Getenv("JWT_KEY")
	EnvPostgresHost     = os.Getenv("PGHOST")
	EnvPostgresPort     = os.Getenv("PGPORT")
	EnvPostgresUser     = os.Getenv("PGUSER")
	EnvPostgresDatabase = os.Getenv("PGDATABASE")
	EnvPostgresPassword = os.Getenv("PGPASSWORD")
)
