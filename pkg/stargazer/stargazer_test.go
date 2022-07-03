package stargazer

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSTARFromFile(t *testing.T) {
	got, err := LoadSTARFromFile(filepath.Join("..", "..", "testdata", "testfile.star"))
	if err != nil {
		t.Error(err)
		return
	}

	if len(got.Entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(got.Entries))
		return
	}

	// Check first entry
	if got.Entries[0].Header.Headersize != 22 {
		t.Errorf("Expected header size of first entry to be 22 bytes long, got %d", got.Entries[0].Header.Headersize)
		return
	}

	if got.Entries[0].Header.Padding1 != 0 {
		t.Errorf("Expected first padding of first entry to be 0, got %d", got.Entries[0].Header.Padding1)
		return
	}

	if got.Entries[0].Header.Filesize != 10928 {
		t.Errorf("Expected filesize of first entry to be 10928, got %d", got.Entries[0].Header.Filesize)
		return
	}

	if got.Entries[0].Header.FilenameSize != 14 {
		t.Errorf("Expected filename size of first entry to be 14, got %d", got.Entries[0].Header.FilenameSize)
		return
	}

	if got.Entries[0].Header.Padding2 != 0 {
		t.Errorf("Expected second padding of first entry to be 0, got %d", got.Entries[0].Header.Padding2)
		return
	}

	if got.Entries[0].Filename != "NulledFile.rel" {
		t.Errorf("Expected filename of first entry to be 'NulledFile.rel', got '%s'", got.Entries[0].Filename)
		return
	}

	if got.Entries[0].Content == nil {
		t.Errorf("Expected content of first entry to be non-nil")
		return
	}

	if uint32(len(got.Entries[0].Content)) != got.Entries[0].Header.Filesize {
		t.Errorf("Expected content of first entry to be %d bytes long, got %d", got.Entries[0].Header.Filesize, len(got.Entries[0].Content))
		return
	}

	if got.Entries[0].SHA1String() != "3d433fcbe9585b05ea877814bad60774ff8a9e7c" {
		t.Errorf("Expected SHA1 of first entry to be '3d433fcbe9585b05ea877814bad60774ff8a9e7c', got '%s'", got.Entries[0].SHA1String())
		return
	}

	// Check second entry
	if got.Entries[1].Header.Headersize != 20 {
		t.Errorf("Expected header size of second entry to be 22 bytes long, got %d", got.Entries[1].Header.Headersize)
		return
	}

	if got.Entries[1].Header.Padding1 != 0 {
		t.Errorf("Expected first padding of second entry to be 0, got %d", got.Entries[1].Header.Padding1)
		return
	}

	if got.Entries[1].Header.Filesize != 313 {
		t.Errorf("Expected filesize of second entry to be 313, got %d", got.Entries[1].Header.Filesize)
		return
	}

	if got.Entries[1].Header.FilenameSize != 12 {
		t.Errorf("Expected filename size of second entry to be 12, got %d", got.Entries[1].Header.FilenameSize)
		return
	}

	if got.Entries[1].Header.Padding2 != 0 {
		t.Errorf("Expected second padding of second entry to be 0, got %d", got.Entries[1].Header.Padding2)
		return
	}

	if got.Entries[1].Filename != "metadata.xml" {
		t.Errorf("Expected filename of second entry to be 'metadata.xml', got '%s'", got.Entries[1].Filename)
		return
	}

	if got.Entries[1].Content == nil {
		t.Errorf("Expected content of second entry to be non-nil")
		return
	}

	if uint32(len(got.Entries[1].Content)) != got.Entries[1].Header.Filesize {
		t.Errorf("Expected content of second entry to be %d bytes long, got %d", got.Entries[1].Header.Filesize, len(got.Entries[1].Content))
		return
	}

	if got.Entries[1].SHA1String() != "2e59ec1846a50fb75042c6786299d13d8f5e39b6" {
		t.Errorf("Expected SHA1 of second entry to be '2e59ec1846a50fb75042c6786299d13d8f5e39b6', got '%s'", got.Entries[1].SHA1String())
		return
	}

	// Check third entry
	if got.Entries[2].Header.Headersize != 19 {
		t.Errorf("Expected header size of third entry to be 19 bytes long, got %d", got.Entries[2].Header.Headersize)
		return
	}

	if got.Entries[2].Header.Padding1 != 0 {
		t.Errorf("Expected first padding of third entry to be 0, got %d", got.Entries[2].Header.Padding1)
		return
	}

	if got.Entries[2].Header.Filesize != 411 {
		t.Errorf("Expected filesize of third entry to be 411, got %d", got.Entries[2].Header.Filesize)
		return
	}

	if got.Entries[2].Header.FilenameSize != 11 {
		t.Errorf("Expected filename size of third entry to be 11, got %d", got.Entries[2].Header.FilenameSize)
		return
	}

	if got.Entries[2].Header.Padding2 != 0 {
		t.Errorf("Expected second padding of third entry to be 0, got %d", got.Entries[2].Header.Padding2)
		return
	}

	if got.Entries[2].Filename != "install.txt" {
		t.Errorf("Expected filename of third entry to be 'install.txt', got '%s'", got.Entries[2].Filename)
		return
	}

	if got.Entries[2].Content == nil {
		t.Errorf("Expected content of third entry to be non-nil")
		return
	}

	if uint32(len(got.Entries[2].Content)) != got.Entries[2].Header.Filesize {
		t.Errorf("Expected content of third entry to be %d bytes long, got %d", got.Entries[2].Header.Filesize, len(got.Entries[2].Content))
		return
	}

	if got.Entries[2].SHA1String() != "6c5768c3c82a174f0ea264c1c0e80450648da4c5" {
		t.Errorf("Expected SHA1 of third entry to be '6c5768c3c82a174f0ea264c1c0e80450648da4c5', got '%s'", got.Entries[2].SHA1String())
		return
	}
}

func TestLoadSTARFromFileNotExisting(t *testing.T) {
	_, err := LoadSTARFromFile("invalid.star")
	if err == nil {
		t.Errorf("Expected error, got nil")
		return
	}
}

func TestLoadSTARFromFileNotAFile(t *testing.T) {
	_, err := LoadSTARFromFile(".")
	if err == nil {
		t.Errorf("Expected error, got nil")
		return
	}

	var expectedErr ErrNotAFile
	if !errors.As(err, &expectedErr) {
		t.Errorf("Expected error to be of type ErrNotAFile, got %T", err)
		return
	}
}

func TestStar_Unpack(t *testing.T) {
	got, err := LoadSTARFromFile(filepath.Join("..", "..", "testdata", "testfile.star"))
	if err != nil {
		t.Error(err)
		return
	}

	outputDir := t.TempDir()
	err = got.Unpack(outputDir)
	if err != nil {
		t.Error(err)
		return
	}

	if !fileExists(filepath.Join(outputDir, "NulledFile.rel")) {
		t.Errorf("Expected file 'NulledFile.rel' to exist in output directory")
		return
	}

	if !fileExists(filepath.Join(outputDir, "metadata.xml")) {
		t.Errorf("Expected file 'metadata.xml' to exist in output directory")
		return
	}

	if !fileExists(filepath.Join(outputDir, "install.txt")) {
		t.Errorf("Expected file 'install.txt' to exist in output directory")
		return
	}
}

func TestNewEntryFromFile(t *testing.T) {
	got, err := NewEntryFromFile(filepath.Join("..", "..", "testdata", "unpacked"), filepath.Join("dummyfolder", "NulledFile.rel"))
	if err != nil {
		t.Error(err)
		return
	}

	if got.Filename != "dummyfolder/NulledFile.rel" {
		t.Errorf("Expected filename to be 'dummyfolder/NulledFile', got '%s'", got.Filename)
		return
	}

	if got.Header.Headersize != 34 {
		t.Errorf("Expected header size to be 34, got %d", got.Header.Headersize)
		return
	}

	if got.Header.Padding1 != 0 {
		t.Errorf("Expected first padding to be 0, got %d", got.Header.Padding1)
		return
	}

	if got.Header.Filesize != 10928 {
		t.Errorf("Expected filesize to be 10928, got %d", got.Header.Filesize)
		return
	}

	if got.Header.FilenameSize != 26 {
		t.Errorf("Expected filename size to be 26, got %d", got.Header.FilenameSize)
		return
	}

	if got.Header.Padding2 != 0 {
		t.Errorf("Expected second padding to be 0, got %d", got.Header.Padding2)
		return
	}

	if got.SHA1String() != "3d433fcbe9585b05ea877814bad60774ff8a9e7c" {
		t.Errorf("Expected SHA1 to be '3d433fcbe9585b05ea877814bad60774ff8a9e7c', got '%s'", got.SHA1String())
		return
	}

	if got.Content == nil {
		t.Errorf("Expected content to be non-nil")
		return
	}

	if uint32(len(got.Content)) != got.Header.Filesize {
		t.Errorf("Content size in header is set to %d bytes, but the file is actually %d bytes long", got.Header.Filesize, len(got.Content))
		return
	}
}

func TestNewSTARFileFromDirectory(t *testing.T) {
	got, err := NewSTARFileFromDirectory(filepath.Join("..", "..", "testdata", "unpacked"))
	if err != nil {
		t.Error(err)
		return
	}

	if len(got.Entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(got.Entries))
		return
	}

	// TODO: Test entries when the correct order is implemented

	// Check first entry
	if got.Entries[0].Header.Headersize != 34 {
		t.Errorf("Expected header size of first entry to be 34 bytes long, got %d", got.Entries[0].Header.Headersize)
		return
	}

	if got.Entries[0].Header.Padding1 != 0 {
		t.Errorf("Expected first padding of first entry to be 0, got %d", got.Entries[0].Header.Padding1)
		return
	}

	if got.Entries[0].Header.Filesize != 10928 {
		t.Errorf("Expected filesize of first entry to be 10928, got %d", got.Entries[0].Header.Filesize)
		return
	}

	if got.Entries[0].Header.FilenameSize != 26 {
		t.Errorf("Expected filename size of first entry to be 26, got %d", got.Entries[0].Header.FilenameSize)
		return
	}

	if got.Entries[0].Header.Padding2 != 0 {
		t.Errorf("Expected second padding of first entry to be 0, got %d", got.Entries[0].Header.Padding2)
		return
	}

	if got.Entries[0].Filename != "dummyfolder/NulledFile.rel" {
		t.Errorf("Expected filename of first entry to be 'dummyfolder/NulledFile.rel', got '%s'", got.Entries[0].Filename)
		return
	}

	if got.Entries[0].Content == nil {
		t.Errorf("Expected content of first entry to be non-nil")
		return
	}

	if uint32(len(got.Entries[0].Content)) != got.Entries[0].Header.Filesize {
		t.Errorf("Expected content of first entry to be %d bytes long, got %d", got.Entries[0].Header.Filesize, len(got.Entries[0].Content))
		return
	}

	if got.Entries[0].SHA1String() != "3d433fcbe9585b05ea877814bad60774ff8a9e7c" {
		t.Errorf("Expected SHA1 of first entry to be '3d433fcbe9585b05ea877814bad60774ff8a9e7c', got '%s'", got.Entries[0].SHA1String())
		return
	}
}

func TestNewSTARFileFromDirectoryNotADirectory(t *testing.T) {
	_, err := NewSTARFileFromDirectory(filepath.Join("..", "..", "testdata", "install.txt"))
	if err == nil {
		t.Errorf("Expected error, got nil")
		return
	}

	var expectedErr ErrNotADirectory
	if !errors.As(err, &expectedErr) {
		t.Errorf("Expected error to be of type ErrNotADirectory, got %T", err)
		return
	}
}

// TODO: More tests for failures (file size, file name size, etc.)

func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	}
	return false
}
