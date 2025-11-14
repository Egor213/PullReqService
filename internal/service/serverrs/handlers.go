package serverrs

import (
	"app/internal/repo/repoerrs"
	errutils "app/pkg/errors"
	"errors"

	log "github.com/sirupsen/logrus"
)

func HandleRepoNotFound(err error, notFoundErr, genericErr error) error {
	if errors.Is(err, repoerrs.ErrNotFound) {
		return notFoundErr
	}
	if err != nil {
		log.Error(errutils.WrapPathErr(err))
		return genericErr
	}
	return nil
}
