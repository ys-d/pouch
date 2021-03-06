package main

import (
	"fmt"
	"strings"

	"github.com/alibaba/pouch/pkg/reference"
	"github.com/spf13/cobra"
)

// createDescription is used to describe create command in detail and auto generate command doc.
var createDescription = "Create a static container object in Pouchd. " +
	"When creating, all configuration user input will be stored in memory store of Pouchd. " +
	"This is useful when you wish to create a container configuration ahead of time so that Pouchd will preserve the resource in advance. " +
	"The container you created is ready to start when you need it."

// CreateCommand use to implement 'create' command, it create a container.
type CreateCommand struct {
	container
	baseCommand
}

// Init initialize create command.
func (cc *CreateCommand) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:   "create [image]",
		Short: "Create a new container with specified image",
		Long:  createDescription,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cc.runCreate(args)
		},
		Example: createExample(),
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *CreateCommand) addFlags() {
	flagSet := cc.cmd.Flags()
	flagSet.SetInterspersed(false)
	flagSet.StringVar(&cc.name, "name", "", "Specify name of container")
	flagSet.BoolVarP(&cc.tty, "tty", "t", false, "Allocate a tty device")
	flagSet.StringSliceVarP(&cc.volume, "volume", "v", nil, "Bind mount volumes to container")
	flagSet.StringVar(&cc.runtime, "runtime", "", "Specify oci runtime")
	flagSet.StringSliceVarP(&cc.env, "env", "e", nil, "Set environment variables for container")
	flagSet.StringSliceVarP(&cc.labels, "label", "l", nil, "Set label for a container")
	flagSet.StringVar(&cc.entrypoint, "entrypoint", "", "Overwrite the default entrypoint")
	flagSet.StringVarP(&cc.workdir, "workdir", "w", "", "Set the working directory in a container")
	flagSet.StringVar(&cc.hostname, "hostname", "", "Set container's hostname")
	flagSet.Int64Var(&cc.cpushare, "cpu-share", 0, "CPU shares")
	flagSet.StringVar(&cc.cpusetcpus, "cpuset-cpus", "", "CPUs in cpuset")
	flagSet.StringVar(&cc.cpusetmems, "cpuset-mems", "", "MEMs in cpuset")
}

// runCreate is the entry of create command.
func (cc *CreateCommand) runCreate(args []string) error {
	config, err := cc.config()
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}

	ref, err := reference.Parse(args[0])
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}
	config.Image = ref.String()

	if len(args) > 1 {
		config.Cmd = args[1:]
	}
	containerName := cc.name

	apiClient := cc.cli.Client()
	result, err := apiClient.ContainerCreate(config.ContainerConfig, config.HostConfig, containerName)
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}

	if len(result.Warnings) != 0 {
		fmt.Printf("WARNING: %s \n", strings.Join(result.Warnings, "\n"))
	}
	fmt.Printf("container ID: %s, name: %s \n", result.ID, result.Name)
	return nil
}

// createExample shows examples in create command, and is used in auto-generated cli docs.
func createExample() string {
	return `$ pouch create --name foo busybox:latest
container ID: e1d541722d68dc5d133cca9e7bd8fd9338603e1763096c8e853522b60d11f7b9, name: foo`
}
