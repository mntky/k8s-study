package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var CfgFile string
var Verbose bool

var RootCmd = &cobra.Command{
	Use:	"k8study",
	Short:	"test k8study command",
	Long:	`test k8study command.
				yeeeeeeeee.
			`,
	Run:	func(cmd *cobra.Command, args []string){
				fmt.Println("k8study test command dayo")
			},
	}

//func init() {
//	//持続的なフラグ
//	RootCmd.PersistentFlags().StringVar(CfgFile, "config", "", "$HOME/k8s-study/cmd/config.yaml")
//	//RootCmd.PersistentFlags().String("
//}
