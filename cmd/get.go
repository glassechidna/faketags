// Copyright © 2017 Aidan Steele <aidan.steele@glassechidna.com.au>
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

	"github.com/spf13/cobra"
	"github.com/aws/aws-sdk-go/aws/session"
	"encoding/json"
	"github.com/glassechidna/faketags/faketags"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.PersistentFlags().GetString("id")
		sess, _ := session.NewSessionWithOptions(session.Options{})

		f := faketags.NewWithNamespace(sess, cliNamespace)
		tagMap, _ := f.TagsForId(id)
		data, _ := json.MarshalIndent(tagMap, "", "  ")
		fmt.Println(string(data))
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().String("id", "", "")
}
