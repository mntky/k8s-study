package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewControllerManagerCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:		"version",
		Short:	"Print the version nubmer of k8s-study",
		Long:		`long description
							aaa
						`,
		Run:		func(cmd *cobra.Command, args []string){
			fmt.Println("k8s-study version 0.1")
		},
	}

	return cmd
}
