package stargazer

import (
	"errors"
	"fmt"
)

var ErrNoEntries = errors.New("invalid STAR file - no file entries found")

type (
	ErrSHA1Mismatch struct {
		Expected string
		Actual   string
		Filename string
	}
)

func (e ErrSHA1Mismatch) Error() string {
	return fmt.Sprintf("SHA1 mismatch on file '%s': expected %s, actual %s", e.Filename, e.Expected, e.Actual)
}
