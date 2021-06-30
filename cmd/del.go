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
	"log"
	"os"
	"strings"

	"github.com/chutommy/siu/db"
	"github.com/chutommy/siu/models"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

// delCmd represents the del command.
var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Deletes one or multiple motions",
	RunE:  del,
}

func init() {
	setCmd.AddCommand(delCmd)
}

func del(*cobra.Command, []string) error {
	// get motion to get id
	m, err := getMotionToDel()
	if err != nil {
		return fmt.Errorf("failed to load motion: %w", err)
	}

	if err := db.Delete(m.ID); err != nil {
		return fmt.Errorf("unable to delete a motion: %w", err)
	}

	printDeleted(m)

	return nil
}

func getMotionToDel() (models.Motion, error) {
	reader := bufio.NewReader(os.Stdin)

	log.Printf("\nDeleting [ID/Name/URL/Shortcut]: ")

	search, err := reader.ReadString('\n')
	if err != nil {
		return models.Motion{}, fmt.Errorf("could not read identificator: %w", err)
	}

	search = strings.TrimSuffix(search, "\n")

	m, err := db.ReadOne(search)
	if err != nil {
		return models.Motion{}, fmt.Errorf("can't read motion: %w", err)
	}

	return m, nil
}

func printDeleted(m models.Motion) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	log.Printf("\nMotion to delete:\n")

	t.AppendHeader(table.Row{"NAME", "URL", "SHORTCUT", "USAGE", "ID"})
	t.AppendRow(table.Row{m.Name, m.URL, m.Shortcut, m.Usage, m.ID})

	t.Render()
}
