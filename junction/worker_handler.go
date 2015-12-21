package junction

import (
	"github.com/hashicorp/go-uuid"
	"github.com/travis-ci/junction/database"
)

// A WorkerHandler handles requests from workers.
type WorkerHandler struct {
	database database.Database
	auth     *AuthService
}

// Create creates an internal record of a worker with the given attributes, and
// returns the workers UUID or an error.
func (h *WorkerHandler) Create(token string, queue string, maxJobCount int, attributes map[string]string) (string, error) {
	if !h.auth.AuthenticateWorker(token) {
		return "", ErrAuthenticationError
	}

	worker := database.Worker{
		ID:          uuid.GenerateUUID(),
		Queue:       queue,
		MaxJobCount: maxJobCount,
		Attributes:  attributes,
	}

	err := h.database.Create(worker)
	if err != nil {
		return "", err
	}

	return worker.ID, nil
}
