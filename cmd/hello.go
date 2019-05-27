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

// helloCmd represents the hello command
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Start your day with the current time and word 'hello'",
	Long: `Hello adds a blank line from tne end of yesterday's timesheet 
	
	and then adds a line with the current timestamp and a task of 'hello'. 
	It should be run at the beginning of a new work day to signify the 
	start of your first task.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello called")
		if len(args) > 0 {
			fmt.Fprintf(os.Stderr, "Unused arguments provided after hello command\n")
			os.Exit(1)
		}
		client.Hello()
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)
}
