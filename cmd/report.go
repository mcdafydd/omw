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

// From specifies the start date of the report output
var From string

// To specified the end date of the report output
var To string

// Format defines the string output format for the report (text or json)
var Format = "text"

var defaultTs string

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Create a simple report of your most recent task entries",
	Long: `Report provides options for creating a simple, formatted view
	of a portion of the tasks in your timesheet.  The default command will 
	show today's tasks, but you may also specify 
	
	--from YYYY-MM-DD --to YYYY-MM-DD 

	to provide start and end dates for the report.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		output, err := client.Report(From, To, Format)
		if err != nil {
			return err
		}
		fmt.Printf("\n%+v\n", output)
		return nil
	},
}

func init() {
	now := time.Now()
	defaultTs = strings.Fields(now.String())[0] // Should be YYYY-MM-DD
	reportCmd.Flags().StringVarP(&From, "from", "f", defaultTs, "Beginning date for report output - beginning today if not specified")
	reportCmd.Flags().StringVarP(&To, "to", "t", defaultTs, "End date for report output - end of today if not specified")
	reportCmd.Flags().StringVarP(&Format, "format", "a", "text", "Format for report output - valid values are \"text\" or \"json\"")
	rootCmd.AddCommand(reportCmd)
}
