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
	"runtime"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dockenv services",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Installing dockenv services")

		os := runtime.GOOS

		switch os {
		case "darwin":
			macosInstallation()
		default:
			fmt.Println("OS not supported")
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func macosInstallation() {
	brewIsInstalledCommand := exec.Command("which", "brew")

	if brewIsInstalledError := brewIsInstalledCommand.Run(); brewIsInstalledError != nil {
		fmt.Println("Homebrew is not installed on this system. Please install it. https://brew.sh")
		return
	}

	dockerIsInstalledCommand := exec.Command("which", "docker")

	if dockerIsInstalledCommandError := dockerIsInstalledCommand.Run(); dockerIsInstalledCommandError != nil {
		dockerInstallationCommand := exec.Command("brew", "install", "--cask", "docker")

		if dockerInstallationCommandError := dockerInstallationCommand.Run(); dockerInstallationCommandError != nil {
			fmt.Fprintln(os.Stderr, "One error ha happened when try to install docker:", dockerInstallationCommandError)
			return
		}
	}

	createDockenvNetwork()
	pullNginxProxyImage()

	fmt.Println("dockenv was installed successfully")
}

func pullNginxProxyImage() {
	pullNginxImageCommand := exec.Command("docker", "image", "pull", "jwilder/nginx-proxy")

	if pullNginxImageCommandError := pullNginxImageCommand.Run(); pullNginxImageCommandError != nil {
		fmt.Fprintln(os.Stderr, "One error ha happened when try to install proxy image: ", pullNginxImageCommandError)
		return
	}
}

func createDockenvNetwork() {
	fmt.Println("Creating networks")
	pullNginxImageCommand := exec.Command("docker", "network", "create", "dockenv-network")

	if pullNginxImageCommandError := pullNginxImageCommand.Run(); pullNginxImageCommandError != nil {
		fmt.Fprintln(os.Stderr, "One error ha happened when try to create network: ", pullNginxImageCommandError)
		return
	}
}
