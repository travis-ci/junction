package junction

import "github.com/travis-ci/junction/database"

// Core is used as a central manager for Junction activity. It is the primary
// point of interface for API handlers.
type Core struct {
	database      database.Database
	WorkerHandler *WorkerHandler
	auth          *AuthService
}

type CoreConfig struct {
	Database     database.Database
	WorkerTokens []string
}

func NewCore(conf *CoreConfig) (*Core, error) {
	c := &Core{
		database: conf.Database,
		auth: &AuthService{
			workerTokens: conf.WorkerTokens,
		},
	}

	c.WorkerHandler = &WorkerHandler{
		database: c.database,
		auth:     c.auth,
	}

	return c, nil
}
