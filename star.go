package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type (
	Header struct {
		Headersize   uint8
		Padding1     uint8
		Filesize     uint32
		FilenameSize uint8
		Padding2     uint8
	}

	Entry struct {
		Header
		FileName []byte
		Content  []byte
		Sha1     [20]byte
	}

	Star struct {
		Entries []Entry
	}
)

func (e *Entry) GetFileName() string {
	return string(e.FileName[:])
}

func (e *Entry) CalculateSha1() []byte {
	h := sha1.New()
	h.Write(e.Content)
	return h.Sum(nil)
}

func (e *Entry) GetSha1() string {
	return fmt.Sprintf("%x", e.Sha1)
}

func (e *Entry) Extract(outputDir string) error {
	fp := filepath.Join(outputDir, e.GetFileName())
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

func (e *Entry) Pack(file *os.File) error {
	err := binary.Write(file, binary.LittleEndian, e.Header)
	if err != nil {
		return err
	}
	_, err = file.Write(e.FileName)
	if err != nil {
		return err
	}
	_, err = file.Write(e.Content)
	if err != nil {
		return err
	}
	_, err = file.Write(e.Sha1[:])
	if err != nil {
		return err
	}
	return nil
}

func ParseEntry(file io.Reader) (*Entry, error) {
	entry := Entry{}
	err := binary.Read(file, binary.LittleEndian, &entry.Header.Headersize)
	if err == io.EOF {
		return nil, err
	}

	binary.Read(file, binary.LittleEndian, &entry.Header.Padding1)

	binary.Read(file, binary.LittleEndian, &entry.Header.Filesize)

	binary.Read(file, binary.LittleEndian, &entry.Header.FilenameSize)
	binary.Read(file, binary.LittleEndian, &entry.Header.Padding2)

	filename := make([]byte, entry.Header.FilenameSize)
	binary.Read(file, binary.LittleEndian, &filename)
	entry.FileName = filename

	entry.Content = make([]byte, entry.Header.Filesize)
	binary.Read(file, binary.LittleEndian, &entry.Content)

	binary.Read(file, binary.LittleEndian, &entry.Sha1)

	calculatedHash := entry.CalculateSha1()

	if !bytes.Equal(calculatedHash, entry.Sha1[:]) {
		log.Fatalln("Hash mismatch")
	}

	return &entry, nil
}
