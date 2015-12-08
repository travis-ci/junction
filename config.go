package junction

type Config struct {
	// Addr is the address the API should listen on
	Addr string

	// DatabaseURL is the PostgresSQL database URL
	DatabaseURL string
}
