# kube-controller-manager

1. Main() 

```go:cmd/kube-controller-manager/controller-manager.go
rand.Seed(time.Now().UnixNano())

//デフォルトのパラメータで*cobra.Commandオブジェクト作成。
command := app.NewControllerManagerCommand()

//Run(*config.CompletedConfig, stopCh<-chan struct{})が呼ばれる。
if err := command.Execute(); err != nil {
		os.Exit(1)
	}
```

1. NewControllerManagerCommand
```go:/cmd/kube-controller-manager/app/controllermanager.go
func NewControllerManagerCommand() *cobra.Command {
	s, err := options.NewKubeControllerManagerOptions()
	if err != nil {
		klog.Fatalf("unable to initialize command options: %v", err)
	}

	//cobraのオブジェクト作成
	cmd := &cobra.Command{
		Use: "kube-controller-manager",
		Long: `The Kubernetes controller manager is a daemon that embeds
the core control loops shipped with Kubernetes. In applications of robotics and
automation, a control loop is a non-terminating loop that regulates the state of
the system. In Kubernetes, a controller is a control loop that watches the shared
state of the cluster through the apiserver and makes changes attempting to move the
current state towards the desired state. Examples of controllers that ship with
Kubernetes today are the replication controller, endpoints controller, namespace
controller, and serviceaccounts controller.`,
		Run: func(cmd *cobra.Command, args []string) {
			verflag.PrintAndExitIfRequested()
			utilflag.PrintFlags(cmd.Flags())

			c, err := s.Config(KnownControllers(), ControllersDisabledByDefault.List())
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			//Run()の呼び出し。
			if err := Run(c.Complete(), wait.NeverStop); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}
	~~~
	return cmd
}

```
kubernetesがよく使っているCLIツールを使って作成。
Run()を呼んでいる。
