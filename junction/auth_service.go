package junction

import "errors"

var (
	// ErrAuthenticationError is returned if the given authentication details for
	// an operation were incorrect. This could mean the token itself isn't valid,
	// or that the token was not authorized for that operation.
	ErrAuthenticationError = errors.New("could not authenticate")
)

type AuthService struct {
	workerTokens []string
	adminTokens  []string
}

// AuthenticateWorker returns true if the given token is a valid token for a
// worker.
func (as *AuthService) AuthenticateWorker(token string) bool {
	for _, workerToken := range as.workerTokens {
		if workerToken == token {
			return true
		}
	}

	return false
}

// AuthenticateAdmin returns true if the given token is a valid token for an
// administrator.
func (as *AuthService) AuthenticateAdmin(token string) bool {
	for _, adminToken := range as.adminTokens {
		if adminToken == token {
			return true
		}
	}

	return false
}
