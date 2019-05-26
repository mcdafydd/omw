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
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Flag vars
var From string
var To string

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Create a simple report of your most recent task entries",
	Long: `Report provides options for creating a simple, formatted view
	of a portion of the tasks in your timesheet.  The default command will 
	show today's tasks, but you may also specify 
	
	--from YYYY-MM-DD --to YYYY-MM-DD 

	to provide start and end dates for the report.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("report called")
		client.Report()
	},
}

func init() {
	now := time.Now()
	defaultTs := strings.Fields(now.String())[0] // Should be YYYY-MM-DD

	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringVarP(&From, "from", "f", defaultTs, "Beginning date for report period - beginning today if not specified")
	reportCmd.Flags().StringVarP(&To, "to", "t", defaultTs, "End date for report period - end of today if not specified")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
