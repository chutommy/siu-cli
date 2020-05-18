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
	table "github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

// updCmd represents the upd command
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

		if colision, bad := db.CheckCollision(m, old); bad {
			return fmt.Errorf("Invalid motion. Reusing values: %v", colision)
		}

		printUpdated(m)

		db.Update(old.ID, m)
		return nil
	},
}

func init() {
	setCmd.AddCommand(updCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getOldMotionToUpd() (models.Motion, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n")

	// get search
	fmt.Print("Updating [ID/Name/URL/Shortcut]: ")
	search, err := reader.ReadString('\n')
	if err != nil {
		return models.Motion{}, fmt.Errorf("Could not get motion's identificator: %v", err)
	}
	if len(strings.Split(search, " ")) != 1 || search == "\n" {
		return models.Motion{}, fmt.Errorf("Invalid identificator: %v", search)
	}
	search = strings.TrimSuffix(search, "\n")

	m, err := db.ReadOne(search)
	if err != nil {
		return models.Motion{}, fmt.Errorf("Could not find motion: %v", search)
	}
	return m, nil
}

func getNewMotionToUpd(old models.Motion) (models.Motion, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n")

	// get name
	fmt.Printf("Name [%v]: ", old.Name)
	name, err := reader.ReadString('\n')
	if err != nil {
		return models.Motion{}, fmt.Errorf("Could not get motion's name: %v", err)
	}
	// if empty use the previous
	if name == "\n" {
		name = old.Name
	}
	if len(strings.Split(name, " ")) != 1 {
		return models.Motion{}, fmt.Errorf("Invalid name: %v", name)
	}

	// get url
	fmt.Printf("URL [%v]: ", old.URL)
	url, err := reader.ReadString('\n')
	if err != nil {
		return models.Motion{}, fmt.Errorf("Could not get motion's name: %v", err)
	}
	// if empty use the previous
	if url == "\n" {
		url = old.URL
	}
	if len(strings.Split(url, " ")) != 1 {
		return models.Motion{}, fmt.Errorf("Invalid url: %v", url)
	}

	// get shortcut
	fmt.Printf("Shortcut [%v]: ", old.Shortcut)
	shortcut, err := reader.ReadString('\n')
	if err != nil {
		return models.Motion{}, fmt.Errorf("Could not get motion's name: %v", err)
	}
	// if empty use the previous
	if shortcut == "\n" {
		shortcut = old.Shortcut
	}
	if len(strings.Split(shortcut, " ")) != 1 {
		return models.Motion{}, fmt.Errorf("Invalid shortcut: %v", shortcut)
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

	fmt.Printf("\nMotion updated:\n")
	t.AppendHeader(table.Row{"NAME", "URL", "SHORTCUT", "USAGE", "ID"})
	t.AppendRow(table.Row{m.Name, m.URL, m.Shortcut, m.Usage, m.ID})

	t.Render()
	fmt.Printf("\n")
}
