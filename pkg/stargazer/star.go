package stargazer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type (
	Star struct {
		Entries []Entry
	}

	Entry struct {
		Header
		Content []byte
		SHA1    [20]byte
	}

	Header struct {
		Headersize   uint8
		Padding1     uint8
		Filesize     uint32
		FilenameSize uint8
		Padding2     uint8
		Filename     string
	}
)

func (e *Entry) SHA1String() string {
	return fmt.Sprintf("%x", e.SHA1)
}

func (e *Entry) Unpack(outputDir string) error {
	fp := filepath.Join(outputDir, e.Filename)
	err := os.MkdirAll(filepath.Dir(fp), os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(e.Content)
	if err != nil {
		return err
	}
	return nil
}

func (e *Entry) Info() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s\n", e.Filename))
	sb.WriteString(fmt.Sprintf("  Filesize: %d bytes\n", e.Header.Filesize))
	sb.WriteString(fmt.Sprintf("  SHA1: %s\n", e.SHA1String()))

	return sb.String()
}

func (s *Star) Unpack(outputDir string) error {
	for _, e := range s.Entries {
		err := e.Unpack(outputDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Star) Info() string {
	var sb strings.Builder

	var contentSize uint64
	for _, e := range s.Entries {
		contentSize += uint64(e.Header.Filesize)
		sb.WriteString(e.Info())
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("Total contents: %d\n", len(s.Entries)))
	sb.WriteString(fmt.Sprintf("Total content size: %d bytes\n", contentSize))

	return sb.String()
}
