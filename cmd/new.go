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
	table "github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new motion",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := getNewMotionToCreate()
		if err != nil {
			return err
		}

		if colision, bad := db.CheckCollision(m, models.Motion{}); bad {
			return fmt.Errorf("Invalid motion. Reusing values: %v", colision)
		}

		if err := db.Create(m); err != nil {
			return err
		}

		printCreated(m)

		return nil
	},
}

func init() {
	setCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getNewMotionToCreate() (models.Motion, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n")

	// get name
	fmt.Print("Name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return models.Motion{}, fmt.Errorf("Could not get motion's name: %v", err)
	}
	if len(strings.Split(name, " ")) != 1 || name == "\n" {
		return models.Motion{}, fmt.Errorf("Invalid name: %v", name)
	}

	// get url
	fmt.Print("URL: ")
	url, err := reader.ReadString('\n')
	if err != nil {
		return models.Motion{}, fmt.Errorf("Could not get motion's name: %v", err)
	}
	if len(strings.Split(url, " ")) != 1 || url == "\n" {
		return models.Motion{}, fmt.Errorf("Invalid url: %v", url)
	}

	// get shortcut
	fmt.Print("Shortcut: ")
	shortcut, err := reader.ReadString('\n')
	if err != nil {
		return models.Motion{}, fmt.Errorf("Could not get motion's name: %v", err)
	}
	if len(strings.Split(shortcut, " ")) != 1 || shortcut == "\n" {
		return models.Motion{}, fmt.Errorf("Invalid shortcut: %v", shortcut)
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
