package stargazer

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
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
		Filename string
		Content  []byte
		SHA1     [20]byte
	}

	Header struct {
		Headersize   uint8
		Padding1     uint8
		Filesize     uint32
		FilenameSize uint8
		Padding2     uint8
	}
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
	entry.Filename = string(filename)

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
			Filename: entry.Filename,
			Expected: hex.EncodeToString(entry.SHA1[:]),
			Actual:   hex.EncodeToString(calculatedHash),
		}
	}

	return &entry, nil
}

//NewSTARFileFromDirectory creates a new STAR from a given directory.
func NewSTARFileFromDirectory(dir string) (*Star, error) {
	// TODO: Order files with install.txt and metadata.txt at end

	if !isDir(dir) {
		return nil, ErrNotADirectory{
			Path: dir,
		}
	}

	star := &Star{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relativePath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		entry, err := NewEntryFromFile(dir, relativePath)
		if err != nil {
			return err
		}
		star.Entries = append(star.Entries, *entry)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return star, nil
}

//NewEntryFromFile creates a new STAR entry from a file.
//Second parameter is the relative filename of the file in the archive.
func NewEntryFromFile(dir string, f string) (*Entry, error) {
	fp := filepath.Join(dir, f)
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entry := &Entry{}

	filename := strings.TrimPrefix(f, string(os.PathSeparator))
	filename = strings.ReplaceAll(filename, "\\", "/")

	if len(filename) > MaxFilenameSize {
		return nil, ErrFilenameTooLong{
			Filename: filename,
		}
	}

	entry.Filename = filename
	entry.Header.FilenameSize = uint8(len(entry.Filename))

	entry.Content, err = io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(entry.Content) > math.MaxUint32 {
		return nil, ErrFileTooLarge{
			Filename: filename,
		}
	}

	entry.Header.Filesize = uint32(len(entry.Content))

	entry.Header.Headersize = uint8(1 + 1 + 4 + 1 + 1 + len(entry.Filename))

	h := sha1.New()
	h.Write(entry.Content)
	copy(entry.SHA1[:], h.Sum(nil))

	return entry, nil
}

//LoadSTARFromFile loads a STAR file from a filepath.
func LoadSTARFromFile(fp string) (*Star, error) {
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	if !isFile(fp) {
		return nil, ErrNotAFile{
			Filename: fp,
		}
	}

	return LoadSTAR(file)
}

//LoadSTAR loads a STAR file from an io.Reader.
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

//SHA1String returns the SHA1 of the file as a hex string.
func (e *Entry) SHA1String() string {
	return fmt.Sprintf("%x", e.SHA1)
}

//Unpack unpacks the entry to the given directory.
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

//WriteTo writes the entry to the given io.Writer.
func (e *Entry) WriteTo(w io.Writer) (int64, error) {
	var total int64
	err := binary.Write(w, binary.LittleEndian, &e.Header)
	if err != nil {
		return 0, err
	}
	total += int64(e.Header.Headersize)

	err = binary.Write(w, binary.LittleEndian, []byte(e.Filename))
	if err != nil {
		return 0, err
	}
	total += int64(e.Header.FilenameSize)

	err = binary.Write(w, binary.LittleEndian, e.Content)
	if err != nil {
		return 0, err
	}
	total += int64(e.Header.Filesize)

	err = binary.Write(w, binary.LittleEndian, e.SHA1)
	if err != nil {
		return 0, err
	}
	total += int64(20)

	return total, nil
}

//Info returns a string with various information of the entry.
func (e *Entry) Info() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s\n", e.Filename))
	sb.WriteString(fmt.Sprintf("  Filesize: %d bytes\n", e.Header.Filesize))
	sb.WriteString(fmt.Sprintf("  SHA1: %s\n", e.SHA1String()))

	return sb.String()
}

//Unpack unpacks all entries from the STAR to the given directory.
func (s *Star) Unpack(outputDir string) error {
	for _, e := range s.Entries {
		err := e.Unpack(outputDir)
		if err != nil {
			return err
		}
	}
	return nil
}

//WriteTo writes the STAR with all entries to the given io.Writer.
func (s *Star) WriteTo(w io.Writer) (int64, error) {
	var total int64
	for _, e := range s.Entries {
		n, err := e.WriteTo(w)
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

//Info returns a string with various information of all entries.
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
