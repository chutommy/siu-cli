/*
Package cmd provides control for the SIU

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
	"fmt"
	"os"

	"github.com/chutified/siu/db"
	"github.com/chutified/siu/models"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

// listCmd represents the list command.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows all available motions",
	RunE: func(cmd *cobra.Command, args []string) error {
		// get all records
		list, err := db.ReadAll()
		if err != nil {
			return fmt.Errorf("failed to read all motions: %w", err)
		}

		// get the table
		t := getTableFromList(list)
		t.Render()
		fmt.Printf("\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func getTableFromList(list []models.Motion) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	fmt.Printf("\n")

	// header
	t.AppendHeader(table.Row{"#", "NAME", "URL", "SHORTCUT", "USAGE", "ID"})

	// records and count total
	var total int32

	for index, m := range list {
		t.AppendRow(table.Row{index, m.Name, m.URL, m.Shortcut, m.Usage, m.ID})
		total += m.Usage
	}

	// total
	t.AppendFooter(table.Row{"", "", "", "Total", total, ""})

	return t
}
