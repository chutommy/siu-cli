// Package cmd provides control for the SIU
/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chutified/siu/db"
	"github.com/chutified/siu/models"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

var (
	errInvalidMotion     = errors.New("motion is not available")
	errInvalidIdentifier = errors.New("identifier can not be recognized")
	errMotionNotFound    = errors.New("motion can not be found")
)

// updCmd represents the upd command.
var updCmd = &cobra.Command{
	Use:   "upd",
	Short: "Updates the motion",
	RunE: func(cmd *cobra.Command, args []string) error {
		old, err := getOldMotionToUpd()
		if err != nil {
			return err
		}

		m, err := getNewMotionToUpd(old)
		if err != nil {
			return err
		}

		if _, bad := db.CheckCollision(m, old); bad {
			return errInvalidMotion
		}

		printUpdated(m)

		if err := db.Update(old.ID, m); err != nil {
			return fmt.Errorf("failed to update a motion: %w", err)
		}

		return nil
	},
}

func init() {
	setCmd.AddCommand(updCmd)
}

// getOldMotionToUpd updates the motion.
func getOldMotionToUpd() (models.Motion, error) {
	reader := bufio.NewReader(os.Stdin)

	log.Printf("\n")

	// get search
	log.Print("Updating [ID/Name/URL/Shortcut]: ")

	search, err := reader.ReadString('\n')
	if err != nil {
		return models.Motion{}, fmt.Errorf("could not get motion's identificator: %w", err)
	}

	if len(strings.Split(search, " ")) != 1 || search == "\n" {
		return models.Motion{}, errInvalidIdentifier
	}

	search = strings.TrimSuffix(search, "\n")

	m, err := db.ReadOne(search)
	if err != nil {
		return models.Motion{}, errMotionNotFound
	}

	return m, nil
}

func getNewMotionToUpd(old models.Motion) (models.Motion, error) {
	reader := bufio.NewReader(os.Stdin)

	log.Printf("\n")

	// get name
	log.Printf("Name [%v]: ", old.Name)

	name, err := reader.ReadString('\n')
	if err != nil {
		return models.Motion{}, fmt.Errorf("could not get motion's name: %w", err)
	}
	// if empty use the previous
	if name == "\n" {
		name = old.Name
	}

	if len(strings.Split(name, " ")) != 1 {
		return models.Motion{}, errInvalidName
	}

	// get url
	log.Printf("URL [%v]: ", old.URL)

	url, err := reader.ReadString('\n')
	if err != nil {
		return models.Motion{}, fmt.Errorf("could not get motion's name: %w", err)
	}
	// if empty use the previous
	if url == "\n" {
		url = old.URL
	}

	if len(strings.Split(url, " ")) != 1 {
		return models.Motion{}, errInvalidURL
	}

	// get shortcut
	log.Printf("Shortcut [%v]: ", old.Shortcut)

	shortcut, err := reader.ReadString('\n')
	if err != nil {
		return models.Motion{}, fmt.Errorf("could not get motion's name: %w", err)
	}
	// if empty use the previous
	if shortcut == "\n" {
		shortcut = old.Shortcut
	}

	if len(strings.Split(shortcut, " ")) != 1 {
		return models.Motion{}, errInvalidShortcut
	}

	return models.Motion{
		ID:       old.ID,
		Name:     strings.TrimSuffix(name, "\n"),
		URL:      strings.TrimSuffix(url, "\n"),
		Shortcut: strings.TrimSuffix(shortcut, "\n"),
		Usage:    old.Usage,
	}, nil
}

func printUpdated(m models.Motion) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	log.Printf("\nMotion updated:\n")
	t.AppendHeader(table.Row{"NAME", "URL", "SHORTCUT", "USAGE", "ID"})
	t.AppendRow(table.Row{m.Name, m.URL, m.Shortcut, m.Usage, m.ID})

	t.Render()
	log.Printf("\n")
}
