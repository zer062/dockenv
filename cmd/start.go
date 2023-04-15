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
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
)

var showLogs bool = false

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
	rootCmd.PersistentFlags().BoolVarP(&showLogs, "logs", "l", false, "Define if show dockenv logs")
}

func startMainService() {
	var containerName string = "dockenv-proxy"
	appContext := context.Background()
	dockerClient, dockerClientError := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if dockerClientError != nil {
		panic(dockerClientError)
	}

	containersList, containersListError := dockerClient.ContainerList(appContext, types.ContainerListOptions{})

	if containersListError != nil {
		panic(containersListError)
	}

	var dockenvServiceIsRunning bool = false

	for _, container := range containersList {
		for _, name := range container.Names {
			if name == "/"+containerName {
				if container.State == "running" {
					dockenvServiceIsRunning = true
				}
			}
		}
	}

	if dockenvServiceIsRunning {
		fmt.Println("Service already running")
		return
	}

	network := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"dockenv-network": {
				NetworkID: "dockenv-network",
			},
		},
	}

	ports := nat.PortMap{
		"80/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "80",
			},
		},
	}

	dockenvServiceContainer, dockenvServiceContainerError := dockerClient.ContainerCreate(appContext, &container.Config{
		Image:        "jwilder/nginx-proxy",
		ExposedPorts: nat.PortSet{"80/tcp": struct{}{}},
	}, &container.HostConfig{
		NetworkMode:  container.NetworkMode("dockenv-network"),
		Binds:        []string{"/var/run/docker.sock:/tmp/docker.sock:ro"},
		PortBindings: ports,
		AutoRemove:   true,
	}, network, &v1.Platform{}, "dockenv-proxy")

	if dockenvServiceContainerError != nil {
		panic(dockenvServiceContainerError)
	}

	if dockenvStartContainerError := dockerClient.ContainerStart(context.Background(), dockenvServiceContainer.ID, types.ContainerStartOptions{}); dockenvStartContainerError != nil {
		panic(dockenvStartContainerError)
	}

	fmt.Println("dockenv started successfully")
	fmt.Println("")

	if showLogs {
		out, err := dockerClient.ContainerLogs(context.Background(), dockenvServiceContainer.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
		if err != nil {
			panic(err)
		}
		defer out.Close()

		for {
			_, err := io.Copy(os.Stdout, out)
			if err != nil {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}
