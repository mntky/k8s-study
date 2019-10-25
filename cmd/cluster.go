package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type Options struct{
	optint int
	optstr string
}

var (
	o = &Options{}
)

var clusterCmd = &cobra.Command{
	Use:		"cluster",
	Short:	"create kubernetes cluster",
	Long:		`This command is Create a kubernetes cluster with kubeadm
						it will be use sh
					`,
	Run:		func(cmd *cobra.Command, args []string){
		out, err := exec.Command("sh", "/root/k8s-study/doc/startup.sh").Output()
		fmt.Println(string(out))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.Flags().StringVarP(&o.optstr, "str", "s", "degault", "string option")
}
