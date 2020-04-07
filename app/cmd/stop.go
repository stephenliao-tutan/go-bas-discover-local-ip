/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"github.com/BASChain/go-bas-dns-server/app/cmdclient"
	"github.com/BASChain/go-bas-dns-server/app/cmdcommon"
	"github.com/spf13/cobra"

	"log"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop basd",
	Long:  `stop basd`,
	Run: func(cmd *cobra.Command, args []string) {

		if _, err := cmdcommon.IsProcessStarted(); err != nil {
			log.Println(err)
			return
		}

		cmdclient.DefaultCmdSend("", cmdcommon.CMD_STOP)

		//config.GetBasDCfg().Save()
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

}
