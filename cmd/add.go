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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add argument <task> and current time to end of timesheet",
	Long: `Add <task> should be run at the end of a task before switching focus.
	Add '**' after your task to categorize it as break time (ie: lunch)
	Add '***' after your task to categorize it as time to ignore (ie: commuting)
	`,
	Example: `
	omw add finish meeting with team
	omw add break **
	omw add commuting ***
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("add called")
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Missing task after add command!\n")
			os.Exit(1)
		}
		client.Add(args)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
