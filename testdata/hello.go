package testdata

import (
	"github.com/pkg/errors"
)

// Hello says hello world!
func Hello() error {
	return errors.New("hello world")
}
