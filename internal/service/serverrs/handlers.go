package serverrs

import (
	"app/internal/repo/repoerrs"
	"errors"
)

func HandleRepoNotFound(err error, notFoundErr, genericErr error) error {
	if errors.Is(err, repoerrs.ErrNotFound) {
		return notFoundErr
	}
	if err != nil {
		return genericErr
	}
	return nil
}
