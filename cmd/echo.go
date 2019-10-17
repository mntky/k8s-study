package commands

import (
	"github.com/spf13/cobra"
	"fmt"
)

var times int
var cmdEcho = &cobra.Command{
	//使用方法の定義,最初のフィールドはコマンドの名前
	Use:	"echo [string to echo]",
	Short:	"echo anything to the screen",
	Long:	`echo is echoing to the screen.
				yeee
				`,
	Run:	echoRun,
}

func init() {
	RootCmd.AddCommand(cmdEcho)
	cmdEcho.Flags().IntVarP(&times, "times", "n", 1, "times to echo")
}

func echoRun(cmd *cobra.Command, args []string) {
	for i:=0; i<times; i++{
		fmt.Println(args)
	}
}
