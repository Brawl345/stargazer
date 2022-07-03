package stargazer

import (
	"errors"
	"fmt"
	"math"
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

	ErrFileTooLarge struct {
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

func (e ErrFileTooLarge) Error() string {
	return fmt.Sprintf("file '%s' is too large, needs to be < %d bytes.", e.Filename, uint32(math.MaxUint32))
}

func (e ErrNotAFile) Error() string {
	return fmt.Sprintf("'%s' is not a file.", e.Filename)
}

func (e ErrNotADirectory) Error() string {
	return fmt.Sprintf("'%s' is not a directory.", e.Path)
}
