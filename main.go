package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const VERSION = "1.0"

type (
	Header struct {
		Unknown1stByte int8
		Unknown2ndByte int8
		Filesize       uint32
		FilenameSize   int8
		MaybePadding   int8
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

func ParseEntry(file io.Reader) (*Entry, error) {
	entry := Entry{}
	err := binary.Read(file, binary.LittleEndian, &entry.Header.Unknown1stByte)
	if err == io.EOF {
		return nil, err
	}

	binary.Read(file, binary.LittleEndian, &entry.Header.Unknown2ndByte)

	binary.Read(file, binary.LittleEndian, &entry.Header.Filesize)

	binary.Read(file, binary.LittleEndian, &entry.Header.FilenameSize)
	binary.Read(file, binary.LittleEndian, &entry.Header.MaybePadding)

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

func main() {
	fmt.Printf("Stargazer v%s\n", VERSION)
	flag.Parse()
	if flag.NArg() < 1 || flag.NArg() > 2 {
		fmt.Println("Usage: stargazer <file> [output dir (optional)]")
		os.Exit(1)
	}
	inputFile := flag.Arg(0)
	outputDir := flag.Arg(1)
	if outputDir == "" {
		outputDir = fmt.Sprintf("%s_extracted", filepath.Base(strings.TrimSuffix(inputFile, filepath.Ext(inputFile))))
	}

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	var star Star

	log.Printf("Parsing '%s'...\n", inputFile)
	for {
		entry, err := ParseEntry(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalln(err)
		}

		star.Entries = append(star.Entries, *entry)
	}

	log.Printf("Extracting to '%s'...\n", outputDir)
	for _, entry := range star.Entries {
		log.Printf("Extracting %s...\n", entry.GetFileName())
		err := entry.Extract(outputDir)
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.Println("Extraction complete!")
}
