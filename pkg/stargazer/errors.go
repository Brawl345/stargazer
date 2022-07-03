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

	ErrFilenameTooLong struct {
		Filename string
	}

	ErrNotAFile struct {
		Filename string
	}

	ErrNotADirectory struct {
		Path string
	}
)

func (e ErrSHA1Mismatch) Error() string {
	return fmt.Sprintf("SHA1 mismatch on file '%s': expected %s, actual %s", e.Filename, e.Expected, e.Actual)
}

func (e ErrFilenameTooLong) Error() string {
	return fmt.Sprintf("filename '%s' is too long, needs to be < %d characters.", e.Filename, MaxFilenameSize)
}

func (e ErrNotAFile) Error() string {
	return fmt.Sprintf("'%s' is not a file.", e.Filename)
}

func (e ErrNotADirectory) Error() string {
	return fmt.Sprintf("'%s' is not a directory.", e.Path)
}
