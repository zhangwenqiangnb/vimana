package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"vimana/cmd/utils"

	"github.com/spf13/cobra"
)

type AvailLightCommander struct {
	BaseCommander
}

func NewAvailLightCommander(node_type string) *AvailLightCommander {
	return &AvailLightCommander{
		BaseCommander: BaseCommander{NodeType: "light"},
	}
}

func (a *AvailLightCommander) AddFlags(cmd *cobra.Command) {
}

func (a *AvailLightCommander) Init(cmd *cobra.Command, args []string, mode Mode, node_info string) error {
	utils.ExecBashCmd(exec.Command("bash", mode.Download), node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	a.initComponentManager("avail", mode.Binary)
	return a.componentMgr.InitializeConfig()
}

func (a *AvailLightCommander) Run(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	cmdexecute := a.componentMgr.GetStartCmd()
	utils.ExecBinaryCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *AvailLightCommander) Start(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	// check if daemon already running.
	PIDFile := utils.GetPIDFileName(node_info)
	if _, err := os.Stat(PIDFile); err == nil {
		fmt.Println("Already running or " + PIDFile + " file exist.")
		return
	}

	node_info_arr := strings.Split(node_info, "-")
	a.Init(cmd, args, mode, node_info_arr[0])
	fmt.Println(a.componentMgr)
	fmt.Println("executing start command")
	cmdexecute := a.componentMgr.GetStartCmd()
	fmt.Println(cmdexecute)
	utils.ExecBashCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *AvailLightCommander) Stop(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	fmt.Println("Stopping Celestia bridge node")

	PIDFile := utils.GetPIDFileName(node_info)
	if _, err := os.Stat(PIDFile); err == nil {
		data, err := ioutil.ReadFile(PIDFile)
		if err != nil {
			fmt.Println("Not running")
			return
		}
		ProcessID, err := strconv.Atoi(string(data))

		if err != nil {
			fmt.Println("Unable to read and parse process id found in ", PIDFile)
			return
		}

		process, err := os.FindProcess(ProcessID)

		if err != nil {
			fmt.Printf("Unable to find process ID [%v] with error %v \n", ProcessID, err)
			return
		}
		// remove PID file
		os.Remove(PIDFile)

		node_info_arr := strings.Split(node_info, "-")
		fmt.Println("Stopping " + node_info_arr[0] + " " + node_info_arr[1] + " node")
		// kill process and exit immediately
		err = process.Kill()

		if err != nil {
			fmt.Printf("Unable to kill process ID [%v] with error %v \n", ProcessID, err)
		} else {
			fmt.Printf("Killed process ID [%v]\n", ProcessID)
		}
	} else {
		fmt.Println("Not running.")
	}
}

func (a *AvailLightCommander) Status(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	node_info_arr := strings.Split(node_info, "-")
	fmt.Println("Getting status of " + node_info_arr[0] + " " + node_info_arr[1] + " node")

	PIDFile := utils.GetPIDFileName(node_info)
	if _, err := os.Stat(PIDFile); err == nil {
		_, err := ioutil.ReadFile(PIDFile)
		if err != nil {
			fmt.Println("Not running")
		} else {
			fmt.Println(node_info_arr[0] + " " + node_info_arr[1] + " node is running")
		}
	} else {
		fmt.Println("Not running.")
	}
}

func (a *AvailLightCommander) Install(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	return
}
