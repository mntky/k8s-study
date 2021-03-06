# kube-controller-manager

# 1.Main() 

```go:cmd/kube-controller-manager/controller-manager.go
rand.Seed(time.Now().UnixNano())

//デフォルトのパラメータで*cobra.Commandオブジェクト作成。
command := app.NewControllerManagerCommand()

//Run(*config.CompletedConfig, stopCh<-chan struct{})が呼ばれる。
if err := command.Execute(); err != nil {
		os.Exit(1)
	}
```

---

# 2.NewControllerManagerCommand
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
kubernetesがよく使用するCLIツール(cobra)を使ってインタフェースを作成している。
https://github.com/spf13/cobra

---
  
# 2-1.NewKubeControllerManagerOptions
```go:/cmd/kube-controller-manager/app/options/options.go
type KubeControllerManagerOptions struct {
	Generic           *cmoptions.GenericControllerManagerConfigurationOptions
	KubeCloudShared   *cmoptions.KubeCloudSharedOptions
	ServiceController *cmoptions.ServiceControllerOptions

	AttachDetachController           *AttachDetachControllerOptions
	CSRSigningController             *CSRSigningControllerOptions
	DaemonSetController              *DaemonSetControllerOptions
	DeploymentController             *DeploymentControllerOptions
	StatefulSetController            *StatefulSetControllerOptions
	DeprecatedFlags                  *DeprecatedControllerOptions
	EndpointController               *EndpointControllerOptions
	EndpointSliceController          *EndpointSliceControllerOptions
	GarbageCollectorController       *GarbageCollectorControllerOptions
	HPAController                    *HPAControllerOptions
	JobController                    *JobControllerOptions
	NamespaceController              *NamespaceControllerOptions
	NodeIPAMController               *NodeIPAMControllerOptions
	NodeLifecycleController          *NodeLifecycleControllerOptions
	PersistentVolumeBinderController *PersistentVolumeBinderControllerOptions
	PodGCController                  *PodGCControllerOptions
	ReplicaSetController             *ReplicaSetControllerOptions
	ReplicationController            *ReplicationControllerOptions
	ResourceQuotaController          *ResourceQuotaControllerOptions
	SAController                     *SAControllerOptions
	TTLAfterFinishedController       *TTLAfterFinishedControllerOptions

	SecureServing *apiserveroptions.SecureServingOptionsWithLoopback
	// TODO: remove insecure serving mode
	InsecureServing *apiserveroptions.DeprecatedInsecureServingOptionsWithLoopback
	Authentication  *apiserveroptions.DelegatingAuthenticationOptions
	Authorization   *apiserveroptions.DelegatingAuthorizationOptions

	Master     string
	Kubeconfig string
}
```
kube-controller-managerのメインコンテキストオブジェクト

```go:/cmd/kube-controller-manager/app/options/options.go
func NewKubeControllerManagerOptions() (*KubeControllerManagerOptions, error) {
	componentConfig, err := NewDefaultComponentConfig(ports.InsecureKubeControllerManagerPort)
	if err != nil {
		return nil, err
	}

	s := KubeControllerManagerOptions{
		Generic:         cmoptions.NewGenericControllerManagerConfigurationOptions(&componentConfig.Generic),
		KubeCloudShared: cmoptions.NewKubeCloudSharedOptions(&componentConfig.KubeCloudShared),
		ServiceController: &cmoptions.ServiceControllerOptions{
			ServiceControllerConfiguration: &componentConfig.ServiceController,
		},
		AttachDetachController: &AttachDetachControllerOptions{
			&componentConfig.AttachDetachController,
		},
		CSRSigningController: &CSRSigningControllerOptions{
			&componentConfig.CSRSigningController,
		},
		DaemonSetController: &DaemonSetControllerOptions{
			&componentConfig.DaemonSetController,
		},
		DeploymentController: &DeploymentControllerOptions{
			&componentConfig.DeploymentController,
		},
		StatefulSetController: &StatefulSetControllerOptions{
			&componentConfig.StatefulSetController,
		},
		DeprecatedFlags: &DeprecatedControllerOptions{
			&componentConfig.DeprecatedController,
		},
		EndpointController: &EndpointControllerOptions{
			&componentConfig.EndpointController,
		},
		EndpointSliceController: &EndpointSliceControllerOptions{
			&componentConfig.EndpointSliceController,
		},
		GarbageCollectorController: &GarbageCollectorControllerOptions{
			&componentConfig.GarbageCollectorController,
		},
		HPAController: &HPAControllerOptions{
			&componentConfig.HPAController,
		},
		JobController: &JobControllerOptions{
			&componentConfig.JobController,
		},
		NamespaceController: &NamespaceControllerOptions{
			&componentConfig.NamespaceController,
		},
		NodeIPAMController: &NodeIPAMControllerOptions{
			&componentConfig.NodeIPAMController,
		},
		NodeLifecycleController: &NodeLifecycleControllerOptions{
			&componentConfig.NodeLifecycleController,
		},
		PersistentVolumeBinderController: &PersistentVolumeBinderControllerOptions{
			&componentConfig.PersistentVolumeBinderController,
		},
		PodGCController: &PodGCControllerOptions{
			&componentConfig.PodGCController,
		},
		ReplicaSetController: &ReplicaSetControllerOptions{
			&componentConfig.ReplicaSetController,
		},
		ReplicationController: &ReplicationControllerOptions{
			&componentConfig.ReplicationController,
		},
		ResourceQuotaController: &ResourceQuotaControllerOptions{
			&componentConfig.ResourceQuotaController,
		},
		SAController: &SAControllerOptions{
			&componentConfig.SAController,
		},
		TTLAfterFinishedController: &TTLAfterFinishedControllerOptions{
			&componentConfig.TTLAfterFinishedController,
		},
	}
	~~~
	return &s, nil
}
```
デフォルトのコンフィグでKubeControllerManagerOptionsを作成

---

## 2-2.Config
```go:/cmd/kube-controller-manager/app/options/options.go
func (s KubeControllerManagerOptions) Config(allControllers []string, disabledByDefaultControllers []string) (*kubecontrollerconfig.Config, error) {

	//コントローラマネージャを起動する前にOptionとConfigを検証する。
	///cmd/kube-controller-manager/app/options/options.go:354
	if err := s.Validate(allControllers, disabledByDefaultControllers); err != nil {
		return nil, err
	}

	//自己証明書を生成する。
	///staging/src/k8s.io/apiserver/pkg/server/options/serving.go
	if err := s.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	//BuildConfigFromFlagsは、masterURLまたはkubeconfigファイル(path)からコンフィグを構築するヘルパ関数。
	//クラスタコンポーネントのコマンドラインフラグとして渡される。
	//masterURL及びkubeconfigファイルパスも渡されない場合inClusterConfig(staging/src/k8s.io/client-go/rest/config.go:410)にフォールバックする。
	//inVlusterConfigが失敗した場合はDefaultconfigにフォールバックする。
	kubeconfig, err := clientcmd.BuildConfigFromFlags(s.Master, s.Kubeconfig)
	if err != nil {
		return nil, err
	}
	kubeconfig.DisableCompression = true
	kubeconfig.ContentConfig.ContentType = s.Generic.ClientConnection.ContentType
	kubeconfig.QPS = s.Generic.ClientConnection.QPS
	kubeconfig.Burst = int(s.Generic.ClientConnection.Burst)

	client, err := clientset.NewForConfig(restclient.AddUserAgent(kubeconfig, KubeControllerManagerUserAgent))
	if err != nil {
		return nil, err
	}

	// shallow copy, do not modify the kubeconfig.Timeout.
	config := *kubeconfig
	config.Timeout = s.Generic.LeaderElection.RenewDeadline.Duration
	leaderElectionClient := clientset.NewForConfigOrDie(restclient.AddUserAgent(&config, "leader-election"))

	eventRecorder := createRecorder(client, KubeControllerManagerUserAgent)

	//controller managerのメインコンテキスト(cmd/kube-controller-manager/app/config/config.go)
	c := &kubecontrollerconfig.Config{
		Client:               client,
		Kubeconfig:           kubeconfig,
		EventRecorder:        eventRecorder,
		LeaderElectionClient: leaderElectionClient,
	}
	if err := s.ApplyTo(c); err != nil {
		return nil, err
	}

	return c, nil
}

```

---
## 2-3.Run
```go:cmd/kube-controller-manager/app/controllermanager.go
func Run(c *config.CompletedConfig, stopCh <-chan struct{}) error {
	// To help debugging, immediately log version
	klog.Infof("Version: %+v", version.Get())

	//ConfigzName(kubecontrollermanager.config.k8s.io)のConfigオブジェクト作成。
	//各Configはこのconfigzパッケージの/configzハンドラに登録される。
	if cfgz, err := configz.New(ConfigzName); err == nil {
		cfgz.Set(c.ComponentConfig)
	} else {
		klog.Errorf("unable to register configz: %v", err)
	}

	//使用するヘルスチェックの設定。
	var checks []healthz.HealthChecker
	var electionChecker *leaderelection.HealthzAdaptor

	//リーダ選出が有効(Raft Algorithm)
	if c.ComponentConfig.Generic.LeaderElection.LeaderElect {
		electionChecker = leaderelection.NewLeaderHealthzAdaptor(time.Second * 20)
		/leaderのhealthcheckをcheckに追加
		checks = append(checks, electionChecker)
	}

	// Start the controller manager HTTP server
	// unsecuredMux is the handler for these controller *after* authn/authz filters have been applied
	//unsecuredMuxは以下のフィルタが適用されたコントローラのハンドラ
	// authn = AutheNtication = 認証
	// authz = AuthoriZation = 承認
	var unsecuredMux *mux.PathRecorderMux

	//HTTPS
	if c.SecureServing != nil {
		unsecuredMux = genericcontrollermanager.NewBaseHandler(&c.ComponentConfig.Generic.Debugging, checks...)
		handler := genericcontrollermanager.BuildHandlerChain(unsecuredMux, &c.Authorization, &c.Authentication)
		// TODO: handle stoppedCh returned by c.SecureServing.Serve
		if _, err := c.SecureServing.Serve(handler, 0, stopCh); err != nil {
			return err
		}
	}

	//HTTP
	if c.InsecureServing != nil {
		unsecuredMux = genericcontrollermanager.NewBaseHandler(&c.ComponentConfig.Generic.Debugging, checks...)
		insecureSuperuserAuthn := server.AuthenticationInfo{Authenticator: &server.InsecureSuperuser{}}
		handler := genericcontrollermanager.BuildHandlerChain(unsecuredMux, nil, &insecureSuperuserAuthn)
		if err := c.InsecureServing.Serve(handler, 0, stopCh); err != nil {
			return err
		}
	}

	run := func(ctx context.Context) {
		rootClientBuilder := controller.SimpleControllerClientBuilder{
			ClientConfig: c.Kubeconfig,
		}
		//クライアントとコントローラのconfigを取得
		var clientBuilder controller.ControllerClientBuilder

		//コントローラを個々のサービスアカウント資格情報で実行するか。
		if c.ComponentConfig.KubeCloudShared.UseServiceAccountCredentials {
			//privateRSAキーがあるかどうか。
			if len(c.ComponentConfig.SAController.ServiceAccountKeyFile) == 0 {
				//別のコントローラプロセスがトークンを作成している可能性がある。
				//クライアントビルダーがトークンを作成できない場合タイムアウトして終了する。
				klog.Warningf("--use-service-account-credentials was specified without providing a --service-account-private-key-file")
			}

			if shouldTurnOnDynamicClient(c.Client) {
				klog.V(1).Infof("using dynamic client builder")
				//Dynamic builder will use TokenRequest feature and refresh service account token periodically
				//DynamicbuilderはTokenRequest機能を利用し、サービスアカウントトークンを定期的に更新する。
				clientBuilder = controller.NewDynamicClientBuilder(
					restclient.AnonymousClientConfig(c.Kubeconfig),
					c.Client.CoreV1(),		//
					"kube-system"					//namespace)

			} else {
				klog.V(1).Infof("using legacy client builder")
				//サービスアカウントとして識別するクライアントを返す。
				clientBuilder = controller.SAControllerClientBuilder{
					ClientConfig:         restclient.AnonymousClientConfig(c.Kubeconfig),
					CoreClient:           c.Client.CoreV1(),
					AuthenticationClient: c.Client.AuthenticationV1(),
					Namespace:            "kube-system",
				}
			}

		//UseServiceAccountCredentials使わないとき
		} else {
			clientBuilder = rootClientBuilder //:198
		}
		controllerContext, err := CreateControllerContext(c, rootClientBuilder, clientBuilder, ctx.Done())
		if err != nil {
			klog.Fatalf("error building controller context: %v", err)
		}
		saTokenControllerInitFunc := serviceAccountTokenControllerStarter{rootClientBuilder: rootClientBuilder}.startServiceAccountTokenController

		//各コントローラの起動。NewControllerInitializersで作ったmapをStartControllersに渡して動かしている。
		if err := StartControllers(controllerContext, saTokenControllerInitFunc, NewControllerInitializers(controllerContext.LoopMode), unsecuredMux); err != nil {
			klog.Fatalf("error starting controllers: %v", err)
		}

		controllerContext.InformerFactory.Start(controllerContext.Stop)
		controllerContext.ObjectOrMetadataInformerFactory.Start(controllerContext.Stop)
		close(controllerContext.InformersStarted)

		select {}
	}

	if !c.ComponentConfig.Generic.LeaderElection.LeaderElect {
		run(context.TODO())
		panic("unreachable")
	}

	id, err := os.Hostname()
	if err != nil {
		return err
	}

	// add a uniquifier so that two processes on the same host don't accidentally both become active
	id = id + "_" + string(uuid.NewUUID())

	rl, err := resourcelock.New(c.ComponentConfig.Generic.LeaderElection.ResourceLock,
		c.ComponentConfig.Generic.LeaderElection.ResourceNamespace,
		c.ComponentConfig.Generic.LeaderElection.ResourceName,
		c.LeaderElectionClient.CoreV1(),
		c.LeaderElectionClient.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity:      id,
			EventRecorder: c.EventRecorder,
		})
	if err != nil {
		klog.Fatalf("error creating lock: %v", err)
	}

	leaderelection.RunOrDie(context.TODO(), leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: c.ComponentConfig.Generic.LeaderElection.LeaseDuration.Duration,
		RenewDeadline: c.ComponentConfig.Generic.LeaderElection.RenewDeadline.Duration,
		RetryPeriod:   c.ComponentConfig.Generic.LeaderElection.RetryPeriod.Duration,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: run,
			OnStoppedLeading: func() {
				klog.Fatalf("leaderelection lost")
			},
		},
		WatchDog: electionChecker,
		Name:     "kube-controller-manager",
	})
	panic("unreachable")
}
```
RunはKubeControllerManagerOptionsを走らせる。この処理は終わらない(はず)


//CreateControllerContext一旦ストップ
## 2-3-1.CreateControllerContext
```go:cmd/kube-controller-manager/app/controllermanager.go
	//TODO:versionedClient辺をよく調べる。
	versionedClient := rootClientBuilder.ClientOrDie("shared-informers")

	//すべての名前空間に対してsharedinformaerfactoryの新しいインスタンスを構築する。
	//https://github.com/kubernetes/sample-controller/blob/master/docs/controller-client-go.md#client-go-components
	sharedInformers := informers.NewSharedInformerFactory(versionedClient, ResyncPeriod(s)())

	metadataClient := metadata.NewForConfigOrDie(rootClientBuilder.ConfigOrDie("metadata-informers"))
	metadataInformers := metadatainformer.NewSharedInformerFactory(metadataClient, ResyncPeriod(s)())

	//apiserverが動いていない場合、待ってから(10秒)失敗する。
	//apiserverとcontrollermanagerを同時に起動する場合に特に重要。
	//WaitForAPIServerは、APIサーバの/healthzエンドポイントがタイムアウトでokを報告するのを待つ。(今回の場合10秒)
	if err := genericcontrollermanager.WaitForAPIServer(versionedClient, 10*time.Second); err != nil {
		return ControllerContext{}, fmt.Errorf("failed to wait for apiserver being healthy: %v", err)
	}

	// Use a discovery client capable of being refreshed.
	discoveryClient := rootClientBuilder.ClientOrDie("controller-discovery")
	cachedClient := cacheddiscovery.NewMemCacheClient(discoveryClient.Discovery())
	restMapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedClient)
	go wait.Until(func() {
		restMapper.Reset()
	}, 30*time.Second, stop)

	availableResources, err := GetAvailableResources(rootClientBuilder)
	if err != nil {
		return ControllerContext{}, err
	}

	cloud, loopMode, err := createCloudProvider(s.ComponentConfig.KubeCloudShared.CloudProvider.Name, s.ComponentConfig.KubeCloudShared.ExternalCloudVolumePlugin,
		s.ComponentConfig.KubeCloudShared.CloudProvider.CloudConfigFile, s.ComponentConfig.KubeCloudShared.AllowUntaggedCloud, sharedInformers)
	if err != nil {
		return ControllerContext{}, err
	}

	ctx := ControllerContext{
		ClientBuilder:                   clientBuilder,
		InformerFactory:                 sharedInformers,
		ObjectOrMetadataInformerFactory: controller.NewInformerFactory(sharedInformers, metadataInformers),
		ComponentConfig:                 s.ComponentConfig,
		RESTMapper:                      restMapper,
		AvailableResources:              availableResources,
		Cloud:                           cloud,
		LoopMode:                        loopMode,
		Stop:                            stop,
		InformersStarted:                make(chan struct{}),
		ResyncPeriod:                    ResyncPeriod(s),
	}
	return ctx, nil
}
```

## 2-3-2.NewControllerInitializers
```go:cmd/kube-controller-manager/app/controllermanager.go
func NewControllerInitializers(loopMode ControllerLoopMode) map[string]InitFunc {
	controllers := map[string]InitFunc{}
	controllers["endpoint"] = startEndpointController
	controllers["endpointslice"] = startEndpointSliceController
	controllers["replicationcontroller"] = startReplicationController
	controllers["podgc"] = startPodGCController
	controllers["resourcequota"] = startResourceQuotaController
	controllers["namespace"] = startNamespaceController
	controllers["serviceaccount"] = startServiceAccountController
	controllers["garbagecollector"] = startGarbageCollectorController
	controllers["daemonset"] = startDaemonSetController
	controllers["job"] = startJobController
	controllers["deployment"] = startDeploymentController
	controllers["replicaset"] = startReplicaSetController
	controllers["horizontalpodautoscaling"] = startHPAController
	controllers["disruption"] = startDisruptionController
	controllers["statefulset"] = startStatefulSetController
	controllers["cronjob"] = startCronJobController
	controllers["csrsigning"] = startCSRSigningController
	controllers["csrapproving"] = startCSRApprovingController
	controllers["csrcleaner"] = startCSRCleanerController
	controllers["ttl"] = startTTLController
	controllers["bootstrapsigner"] = startBootstrapSignerController
	controllers["tokencleaner"] = startTokenCleanerController
	controllers["nodeipam"] = startNodeIpamController
	controllers["nodelifecycle"] = startNodeLifecycleController
	if loopMode == IncludeCloudLoops {
		controllers["service"] = startServiceController
		controllers["route"] = startRouteController
		controllers["cloud-node-lifecycle"] = startCloudNodeLifecycleController
		// TODO: volume controller into the IncludeCloudLoops only set.
	}
	controllers["persistentvolume-binder"] = startPersistentVolumeBinderController
	controllers["attachdetach"] = startAttachDetachController
	controllers["persistentvolume-expander"] = startVolumeExpandController
	controllers["clusterrole-aggregation"] = startClusterRoleAggregrationController
	controllers["pvc-protection"] = startPVCProtectionController
	controllers["pv-protection"] = startPVProtectionController
	controllers["ttl-after-finished"] = startTTLAfterFinishedController
	controllers["root-ca-cert-publisher"] = startRootCACertPublisher

	return controllers
}

```

## 2-3-3.StartControllers
```go:cmd/kube-controller-manager/app/controllermanager.go
	//コントローラのセットを開始する。
	// Always start the SA token controller first using a full-power client, since it needs to mint tokens for the rest
	// If this fails, just return here and fail since other controllers won't be able to get credentials.
	//full-powerクライアントを使って、サービスアカウントトークンコントローラの起動。
	if _, _, err := startSATokenController(ctx); err != nil {
		return err
	}

	// Initialize the cloud provider with a reference to the clientBuilder only after token controller
	// has started in case the cloud provider uses the client builder.
	//クラウドプロバイダ(kubeadm,aws,azure...)がclientBuilderを使用する場合,
	//TokenControllerが起動した後のみclientBuilderへの参照でクラウドプロバイダを初期化する。
	if ctx.Cloud != nil {
		ctx.Cloud.Initialize(ctx.ClientBuilder, ctx.Stop)
	}

	//NewControllerInitializersで作られたmapに入ってる各コントローラ
	for controllerName, initFn := range controllers {
		if !ctx.IsControllerEnabled(controllerName) {
			//僕の場合"endpointslice"だけdisabledだった。
			klog.Warningf("%q is disabled", controllerName)
			continue
		}

		time.Sleep(wait.Jitter(ctx.ComponentConfig.Generic.ControllerStartInterval.Duration, ControllerStartJitter))

		//↓このログ確認できない:thinkingface
		klog.V(1).Infof("Starting %q", controllerName)
		//initFnに各controllerのスタートする関数名が入ってるのでここでスタート
		debugHandler, started, err := initFn(ctx)
		if err != nil {
			klog.Errorf("Error starting %q", controllerName)
			return err
		}
		if !started {
			klog.Warningf("Skipping %q", controllerName)
			continue
		}
		if debugHandler != nil && unsecuredMux != nil {
			basePath := "/debug/controllers/" + controllerName
			unsecuredMux.UnlistedHandle(basePath, http.StripPrefix(basePath, debugHandler))
			unsecuredMux.UnlistedHandlePrefix(basePath+"/", http.StripPrefix(basePath, debugHandler))
		}
		klog.Infof("Started %q", controllerName)
	}

	return nil
}
```

### 上記initFn(ctx)がreplicasetcontrollerの場合
```go:/vmd/kube-controller-manager/app/apps.go
func startReplicaSetController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "replicasets"}] {
		return nil, false, nil
	}
	go replicaset.NewReplicaSetController(
		ctx.InformerFactory.Apps().V1().ReplicaSets(),
		ctx.InformerFactory.Core().V1().Pods(),
		ctx.ClientBuilder.ClientOrDie("replicaset-controller"),
		replicaset.BurstReplicas,
	).Run(int(ctx.ComponentConfig.ReplicaSetController.ConcurrentRSSyncs), ctx.Stop)
	return nil, true, nil
}
```
