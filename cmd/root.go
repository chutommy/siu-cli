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
	"os/exec"
	"runtime"
	"strings"

	"github.com/chutified/siu/db"
	"github.com/chutified/siu/models"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "siu",
	Short: "Siu is a very fast and minimalist urls opener",
	RunE: func(cmd *cobra.Command, args []string) error {
		motions, err := getMotionsToRun()
		if err != nil {
			return err
		}

		if err := runMotions(motions); err != nil {
			return err
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.siu.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
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
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".siu" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".siu")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getMotionsToRun() ([]models.Motion, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n")

	// get input and trim it
	fmt.Printf("RUN: ")
	items, err := reader.ReadString('\n')
	if err != nil {
		return []models.Motion{}, fmt.Errorf("Could not read what to run: %v", err)
	}
	items = strings.TrimSpace(strings.TrimSuffix(items, "\n"))
	fmt.Printf("\n")

	// search for each motion
	searches := strings.Split(items, " ")
	var motions []models.Motion
	for _, search := range searches {

		// skip if extra spaces
		if search == "" {
			continue
		}
		m, err := db.ReadOne(search)
		if err != nil {
			// if not found log and skip
			fmt.Printf("Motion %v not found...\n", search)
			continue
		}

		motions = append(motions, m)
	}

	return motions, nil
}

func runMotions(motions []models.Motion) error {
	for _, m := range motions {
		fmt.Printf("Openning %v ...\n", m.URL)
		if err := openBrowser(m.URL); err != nil {
			fmt.Printf("Could not open: %v ...\n", m.URL)
		}
	}
	fmt.Printf("\n")
	return nil
}

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
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return err
	}

	return nil
}
