package e

import "github.com/pkg/errors"

func S(err error) error {
	return errors.WithStack(err)
}
