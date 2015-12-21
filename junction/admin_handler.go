package junction

import "github.com/travis-ci/junction/database"

// An AdminHandler handles requests from administrators. This mostly contains
// diagnostic and informational commands.
type AdminHandler struct {
	database database.Database
	auth     *AuthService
}

func (h *AdminHandler) ListWorkers(token string) ([]database.Worker, error) {
	if !h.auth.AuthenticateAdmin(token) {
		return nil, ErrAuthenticationError
	}

	return h.database.ListWorkers()
}
