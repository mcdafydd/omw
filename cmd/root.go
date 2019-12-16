// Copyright © 2019 David McPike
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/inconshreveable/mousetrap"
	"github.com/mcdafydd/omw/backend"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

const (
	// DefaultDir is the default directory inside the user's home directory
	// that will store omw data files
	DefaultDir = ".local/share/omw"
	// DefaultFile is the default filename for the primary time tracking data log
	DefaultFile = "omw.toml"
)

var server *backend.Backend

// MousetrapHelpText Set MousetrapHelpText to an empty string to disable Cobra's
// automatic display of a warning to Windows users who double-click the binary
// from Windows Explorer.  We want to have our own mousetrap and alias it to
// 'omw server'.
var MousetrapHelpText string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "omw",
	Short: "Omw - Out of My Way Time Tracker",
	Long: `A minimalist time tracker inspired by the Ultimate Time Tracker (UTT).

	The primary purposes of this tool are:

	1. Help a user track time and tasks without getting in the way of flow
	2. Provide a simple, extendable reporting interface to help transfer
	tasks to an external system`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if mousetrap.StartedByExplorer() {
			err = serverCmd.RunE(cmd, args)
			fmt.Println("running backend from explorer")
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		return err
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return server.Close()
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

	home, err := homedir.Dir()
	if err != nil {
		errors.Wrap(err, "homedir.Dir() returned error")
	}

	fm := os.FileMode(0700)
	omwDir := fmt.Sprintf("%s/%s", home, DefaultDir)
	err = os.MkdirAll(omwDir, fm)
	if err != nil {
		errors.Wrapf(err, "MkdirAll %s", omwDir)
	}

	omwFile := fmt.Sprintf("%s/%s", omwDir, DefaultFile)
	if _, err := os.Stat(omwFile); os.IsNotExist(err) {
		fmt.Println("file does not exist - creating file", omwFile)
		fp, err := os.OpenFile(omwFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			errors.Wrapf(err, "Can't open or create %s", omwFile)
		}
		fp.Close()
	}

	server = backend.Create(nil, omwDir, omwFile)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.omw.yaml)")
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

		// Search config in home directory with name ".omw" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".omw")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
