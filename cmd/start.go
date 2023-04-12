/*
Copyright © 2023 zer062 <zero62@zero62.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start dockenv services",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		startMainService()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func startMainService() {
	startCommand := exec.Command("docker", "run", "-d", "--rm", "--name=dockenv-proxy", "--net", "dockenv-network", "-p", "80:80", "-v", "/var/run/docker.sock:/tmp/docker.sock:ro", "jwilder/nginx-proxy")

	if startCommandError := startCommand.Run(); startCommandError != nil {
		fmt.Fprintln(os.Stderr, "Unable to start dockenv services: ", startCommandError)
		return
	}

	fmt.Println("dockenv started successfully")
}
