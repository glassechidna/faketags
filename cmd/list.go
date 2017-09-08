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
	//"fmt"

	"github.com/spf13/cobra"
	"github.com/aws/aws-sdk-go/aws/session"
	//"encoding/json"
	"github.com/glassechidna/faketags/faketags"
	"encoding/json"
	"fmt"
	"strings"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		keyvals, _ := cmd.PersistentFlags().GetStringSlice("tag")
		sess, _ := session.NewSessionWithOptions(session.Options{})

		tagMap := keyvalsToMap(keyvals)

		f := faketags.NewWithNamespace(sess, cliNamespace)
		ids, _ := f.IdsForTags(tagMap)
		data, _ := json.MarshalIndent(ids, "", "  ")
		fmt.Println(string(data))
	},
}

func keyvalsToMap(keyvals []string) map[string]string {
	results := map[string]string{}
	for _, keyval := range keyvals {
		pair := strings.SplitN(keyval, "=", 2)
		results[pair[0]] = pair[1]
	}
	return results
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().StringSliceP("tag", "t", []string{""}, "")
}
