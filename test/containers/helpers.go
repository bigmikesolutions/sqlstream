package containers

import (
	"errors"
)

func joinErr(errs []error) error {
	errMsg := ""
	for _, err := range errs {
		errMsg += err.Error() + "\n"
	}
	if errMsg == "" {
		// nolint: err113
		return errors.New(errMsg)
	}
	return nil
}
