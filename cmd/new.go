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
	"fmt"
	"os"
	"strings"

	"github.com/chutified/siu/db"
	"github.com/chutified/siu/models"
	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

// newCmd represents the new command.
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new motion",
	RunE:  new,
}

func init() {
	setCmd.AddCommand(newCmd)
}

func new(*cobra.Command, []string) error {
	m, err := getNewMotionToCreate()
	if err != nil {
		return err
	}

	if collision, bad := db.CheckCollision(m, models.Motion{}); bad {
		return fmt.Errorf("invalid motion, reusing values: %v", collision)
	}

	if err := db.Create(m); err != nil {
		return fmt.Errorf("failed to create a new motion: %w", err)
	}

	printCreated(m)

	return nil
}

func getNewMotionToCreate() (models.Motion, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\n")

	// get name
	fmt.Print("Name: ")

	name, err := reader.ReadString('\n')
	if err != nil {
		return models.Motion{}, fmt.Errorf("could not get motion's name: %w", err)
	}

	if len(strings.Split(name, " ")) != 1 || name == "\n" {
		return models.Motion{}, fmt.Errorf("invalid name: %s", name)
	}

	// get url
	fmt.Print("URL: ")

	url, err := reader.ReadString('\n')
	if err != nil {
		return models.Motion{}, fmt.Errorf("could not get motion's name: %w", err)
	}

	if len(strings.Split(url, " ")) != 1 || url == "\n" {
		return models.Motion{}, fmt.Errorf("invalid url: %s", url)
	}

	// get shortcut
	fmt.Print("Shortcut: ")

	shortcut, err := reader.ReadString('\n')
	if err != nil {
		return models.Motion{}, fmt.Errorf("could not get motion's name: %w", err)
	}

	if len(strings.Split(shortcut, " ")) != 1 || shortcut == "\n" {
		return models.Motion{}, fmt.Errorf("invalid shortcut: %s", shortcut)
	}

	return models.Motion{
		ID:       uuid.New().String(),
		Name:     strings.TrimSuffix(name, "\n"),
		URL:      strings.TrimSuffix(url, "\n"),
		Shortcut: strings.TrimSuffix(shortcut, "\n"),
		Usage:    0,
	}, nil
}

func printCreated(m models.Motion) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	fmt.Printf("\nNew motion created:\n")
	t.AppendHeader(table.Row{"NAME", "URL", "SHORTCUT", "ID"})
	t.AppendRow(table.Row{m.Name, m.URL, m.Shortcut, m.ID})

	t.Render()
	fmt.Printf("\n")
}
