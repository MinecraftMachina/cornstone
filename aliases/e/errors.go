package e

import (
	"fmt"
	"github.com/pkg/errors"
)

func S(err error) error {
	return errors.WithStack(err)
}

func P(err error) string {
	return fmt.Sprintf("%+v\n", err)
}
