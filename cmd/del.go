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

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Deletes one or multiple motions",
	RunE: func(cmd *cobra.Command, args []string) error {
		// get motion to get id
		m, err := getMotionToDel()
		if err != nil {
			return err
		}

		if err := db.Delete(m.ID); err != nil {
			return err
		}

		printDeleted(m)

		return nil
	},
}

func init() {
	setCmd.AddCommand(delCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// delCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// delCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getMotionToDel() (models.Motion, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n")

	// get, check, trim
	fmt.Print("Deleting [ID/Name/URL/Shortcut]: ")
	search, err := reader.ReadString('\n')
	if err != nil {
		return models.Motion{}, fmt.Errorf("Could not read identificator: %v", err)
	}
	search = strings.TrimSuffix(search, "\n")

	return db.ReadOne(search)
}

func printDeleted(m models.Motion) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	fmt.Printf("\nMotion deleted:\n")
	t.AppendHeader(table.Row{"NAME", "URL", "SHORTCUT", "USAGE", "ID"})
	t.AppendRow(table.Row{m.Name, m.URL, m.Shortcut, m.Usage, m.ID})

	t.Render()
	fmt.Printf("\n")
}
