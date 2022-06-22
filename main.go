package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const VERSION = "2.0"

func usage() {
	fmt.Println("Usage: stargazer <operation> <arguments>")
	fmt.Println("")
	fmt.Println("  To extract files:")
	fmt.Println("    stargazer x <star file> [output dir (optional)]")
	fmt.Println("")
	fmt.Println("  To pack a folder:")
	fmt.Println("    stargazer p <input dir> <star file>")
	fmt.Println("")
	os.Exit(1)
}

func extract() {
	inputFile := flag.Arg(1)
	outputDir := flag.Arg(2)
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

func pack() {
	log.Printf("WARNING!!! Packing is experimental and may not work properly!\n")
	inputDir := flag.Arg(1)
	outputFile := flag.Arg(2)
	if outputFile == "" {
		outputFile = fmt.Sprintf("%s_packed.star", filepath.Base(inputDir))
	}

	log.Println("Reading files...")
	var star Star
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		entry := Entry{}
		fp := strings.TrimPrefix(path, inputDir)
		fp = strings.TrimPrefix(fp, string(os.PathSeparator))
		entry.FileName = []byte(strings.ReplaceAll(fp, "\\", "/"))
		entry.Content, err = ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		entry.Header.Headersize = uint8(8 + len(entry.GetFileName()))
		entry.Header.Filesize = uint32(len(entry.Content))
		entry.Header.FilenameSize = uint8(len(entry.GetFileName()))

		copy(entry.Sha1[:], entry.CalculateSha1())
		star.Entries = append(star.Entries, entry)
		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Packing...")
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	for _, entry := range star.Entries {
		log.Printf("Packing %s...\n", entry.GetFileName())
		err := entry.Pack(file)
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Println("Packing complete!")
}

func main() {
	fmt.Printf("Stargazer v%s\n", VERSION)
	flag.Parse()
	if flag.NArg() < 1 || flag.NArg() > 3 {
		usage()
		os.Exit(1)
	}
	operation := flag.Arg(0)
	switch operation {
	case "x":
		extract()
	case "p":
		pack()
	default:
		usage()
	}
}
