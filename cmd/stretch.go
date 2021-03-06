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

	"github.com/spf13/cobra"
)

// stretchCmd represents the stretch command
var stretchCmd = &cobra.Command{
	Use:   "stretch",
	Short: "Stretch adds a copy of the most recent task to the timesheet",
	Long: `Stretch creates a copy of the last entry on your timesheet
	with the current time, effectively 'stretching' it's total time.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			fmt.Fprintf(os.Stderr, "Unused arguments provided after stretch command\n")
			os.Exit(1)
		}
		return server.Stretch()
	},
}

func init() {
	rootCmd.AddCommand(stretchCmd)
}
