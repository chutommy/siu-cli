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
	"os/exec"
	"runtime"
	"strings"

	"github.com/chutified/siu/db"
	"github.com/chutified/siu/models"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var errUnsupportedPlatform = errors.New("platform is not supported")

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "siu",
	Short: "Siu is a very fast and minimalist urls opener",
	RunE: func(cmd *cobra.Command, args []string) error {
		motions, err := getMotionsToRun()
		if err != nil {
			return err
		}

		return runMotions(motions)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.siu.yaml)")

	// flags
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with a name ".siu" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".siu")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// getMotionsToRun reads the user's input.
func getMotionsToRun() ([]models.Motion, error) {
	reader := bufio.NewReader(os.Stdin)

	// get input and trim it
	log.Printf("RUN: ")

	items, err := reader.ReadString('\n')
	if err != nil {
		return []models.Motion{}, fmt.Errorf("could not read what to run: %w", err)
	}

	items = strings.TrimSpace(strings.TrimSuffix(items, "\n"))

	log.Printf("\n")

	// search for each motion
	searches := strings.Split(items, " ")
	motions := make([]models.Motion, len(searches))

	for i, search := range searches {
		// skip if extra spaces
		if search == "" {
			continue
		}

		m, err := db.ReadOne(search)
		if err != nil {
			// if not found a log and skip
			log.Printf("Motion %v not found...\n", search)

			continue
		}

		motions[i] = m
	}

	return motions, nil
}

// runMotions takes the motions and handles them.
func runMotions(motions []models.Motion) error {
	for _, m := range motions {
		log.Printf("Openning %v ...\n", m.URL)

		if err := openBrowser(m.URL); err != nil {
			log.Printf("Could not open: %v ...\n", m.URL)
		}

		if err := db.IncMotionUsage(m); err != nil {
			return fmt.Errorf("internal application error, failed to increment usage: %w", err)
		}
	}

	log.Printf("\n")

	return nil
}

// openBrowser opens the url in a new browser
// If a browser window exists it opens the url in a new tab.
func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = errUnsupportedPlatform
	}

	if err != nil {
		return err
	}

	return nil
}
