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
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit your current timesheet",
	Long:  `Opens a new window to view/edit your current timesheet using your default editor.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		reopen, err := server.Edit()
		for reopen {
			reopen, err = server.Edit()
			if err != nil {
				break
			}
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
