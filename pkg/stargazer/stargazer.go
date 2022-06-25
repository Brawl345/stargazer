package stargazer

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func parseEntry(file io.Reader) (*Entry, error) {
	entry := Entry{}
	err := binary.Read(file, binary.LittleEndian, &entry.Header.Headersize)
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("error reading headersize: %v", err)
	}

	err = binary.Read(file, binary.LittleEndian, &entry.Header.Padding1)
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("error reading padding1: %v", err)
	}

	err = binary.Read(file, binary.LittleEndian, &entry.Header.Filesize)
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("error reading filesize: %v", err)
	}

	err = binary.Read(file, binary.LittleEndian, &entry.Header.FilenameSize)
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("error reading filename size: %v", err)
	}

	err = binary.Read(file, binary.LittleEndian, &entry.Header.Padding2)
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("error reading padding2: %v", err)
	}

	filename := make([]byte, entry.Header.FilenameSize)
	err = binary.Read(file, binary.LittleEndian, &filename)
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("error reading filename: %v", err)
	}
	entry.Header.Filename = string(filename)

	entry.Content = make([]byte, entry.Header.Filesize)
	err = binary.Read(file, binary.LittleEndian, &entry.Content)
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("error reading content: %v", err)
	}

	err = binary.Read(file, binary.LittleEndian, &entry.SHA1)
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("error reading sha1: %v", err)
	}

	h := sha1.New()
	h.Write(entry.Content)
	calculatedHash := h.Sum(nil)

	if !bytes.Equal(calculatedHash, entry.SHA1[:]) {
		return nil, ErrSHA1Mismatch{
			Filename: entry.Header.Filename,
			Expected: hex.EncodeToString(entry.SHA1[:]),
			Actual:   hex.EncodeToString(calculatedHash),
		}
	}

	return &entry, nil
}

func LoadSTARFromFile(fp string) (*Star, error) {
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return LoadSTAR(file)
}

func LoadSTAR(file io.Reader) (*Star, error) {
	star := &Star{}
	for {
		entry, err := parseEntry(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		star.Entries = append(star.Entries, *entry)
	}

	if len(star.Entries) == 0 {
		return nil, ErrNoEntries
	}

	return star, nil
}
