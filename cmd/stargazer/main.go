package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Brawl345/stargazer/pkg/stargazer"
	"github.com/urfave/cli/v2"
)

var input string
var output string
var quiet bool

func main() {
	app := &cli.App{
		Name:                 "stargazer",
		Usage:                "A tool to handle PSX STAR files",
		Version:              "2.0.0",
		Suggest:              true,
		EnableBashCompletion: true,
		HideHelpCommand:      true,
		Authors: []*cli.Author{
			{
				Name: "Brawl345",
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "quiet",
				Aliases:     []string{"q"},
				Usage:       "Do not print any messages",
				Destination: &quiet,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "unpack",
				Aliases: []string{"u"},
				Usage:   "Unpacks files from a STAR file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "input",
						Aliases:     []string{"i"},
						Required:    true,
						Usage:       "Path to STAR file",
						Destination: &input,
					},
					&cli.StringFlag{
						Name:        "output",
						Aliases:     []string{"o"},
						Required:    false,
						Usage:       "Path to output directory. Defaults to '<input file without .star>_extracted'",
						Destination: &output,
					},
				},
				Action: unpack,
			},
			{
				Name:    "pack",
				Aliases: []string{"p"},
				Usage:   "Pack a folder into a STAR file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "input",
						Aliases:     []string{"i"},
						Required:    true,
						Usage:       "Path to a folder",
						Destination: &input,
					},
					&cli.StringFlag{
						Name:        "output",
						Aliases:     []string{"o"},
						Required:    false,
						Usage:       "Output path of the STAR file. Defaults to '<input folder>_packed.star'",
						Destination: &output,
					},
				},
				Action: pack,
			},
			{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "Shows information about a STAR file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "input",
						Aliases:     []string{"i"},
						Required:    true,
						Usage:       "Path to STAR file",
						Destination: &input,
					},
				},
				Action: info,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(1)
	}
}

func unpack(_ *cli.Context) error {
	star, err := stargazer.LoadSTARFromFile(input)
	if err != nil {
		return err
	}

	if output == "" {
		output = fmt.Sprintf("%s_extracted", filepath.Base(strings.TrimSuffix(input, filepath.Ext(input))))
	}

	if !quiet {
		log.Printf("Will unpack to '%s'", output)
	}

	for _, entry := range star.Entries {
		if !quiet {
			log.Printf("Unpacking '%s'...\n", entry.Filename)
		}
		err := entry.Unpack(output)
		if err != nil {
			return err
		}
	}

	return nil
}

func pack(_ *cli.Context) error {
	if output == "" {
		output = fmt.Sprintf("%s_packed.star", filepath.Base(input))
	}

	if !quiet {
		log.Printf("Will pack to '%s'", output)
	}

	if !quiet {
		log.Printf("Reading '%s'...", input)
	}

	star, err := stargazer.NewSTARFileFromDirectory(input)

	if err != nil {
		return err
	}

	if !quiet {
		log.Printf("Writing to '%s'...\n", output)
	}

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	_, err = star.WriteTo(writer)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func info(_ *cli.Context) error {
	star, err := stargazer.LoadSTARFromFile(input)
	if err != nil {
		return err
	}

	fmt.Println(star.Info())
	return nil
}
