package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "varconf",
	PostRun: func(cmd *cobra.Command, args []string) {
	},
	Short: "The Varconf is a lite powerful config center platform.",
	Long: `
 __      __     _____   _____ ____  _   _ ______ 
 \ \    / /\   |  __ \ / ____/ __ \| \ | |  ____|
  \ \  / /  \  | |__) | |   | |  | |  \| | |__   
   \ \/ / /\ \ |  _  /| |   | |  | | .   |  __|  
    \  / ____ \| | \ \| |___| |__| | |\  | |     
     \/_/    \_\_|  \_\\_____\____/|_| \_|_|     

The Varconf is a lite powerful config center platform.                       
	`,
	Version: `0.0.1`,
	Example: `varconf start -c ./config.json`,
}

var pidFile = "./pid.lock"
var (
	configFile string
	daemon     bool
	initFile   string
	startCmd   = &cobra.Command{
		Use:   "start",
		Short: "start varconf server",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting!")
			if daemon {
				cmd := exec.Command(os.Args[0], "--config", configFile, "--init", initFile, "start")
				err := cmd.Start()
				if err == nil {
					fmt.Printf("PID %d is running...\n", cmd.Process.Pid)
				} else {
					fmt.Println("Start failed!", err.Error())
				}
			} else {
				pid := fmt.Sprintf("%d", os.Getpid())
				err := os.WriteFile(pidFile, []byte(pid), 0666)
				if err != nil {
					fmt.Println("Start failed!", err.Error())
					return
				}
				err = Start(configFile, initFile)
				if err != nil {
					fmt.Println("Start failed!", err.Error())
					os.Remove(pidFile)
				}
			}
		},
	}
	stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "stop varconf server",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Stopping!")
			pidBytes, err := os.ReadFile(pidFile)
			if err != nil {
				fmt.Println("Read PID error!", err)
				return
			}

			pid := string(pidBytes)
			killCmd := new(exec.Cmd)
			if runtime.GOOS == "windows" {
				killCmd = exec.Command("taskkill", "/f", "/pid", pid)
			} else {
				killCmd = exec.Command("kill", pid)
			}

			err = killCmd.Start()
			if err == nil {
				fmt.Printf("PID %s has been stopped!\n", pid)
				os.Remove(pidFile)
			} else {
				fmt.Println("PID "+pid+" stop failed! %s\n", err)
			}
		},
	}
)

func Execute() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := Cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	startCmd.Flags().StringVarP(&configFile, "config", "c", "./config.json", "config file path")
	startCmd.Flags().BoolVarP(&daemon, "daemon", "d", false, "is daemon?")
	startCmd.Flags().StringVarP(&initFile, "init", "i", "./init.json", "init file path")

	Cmd.AddCommand(startCmd)
	Cmd.AddCommand(stopCmd)

}
