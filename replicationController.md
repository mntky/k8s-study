# ReplicationController

定義されたReplicaSetオブジェクトを実行中のPodと同期させる。
ReplicationControllerはrolling-updateの機能があるが、ReplicaSetControllerには無い。

[startRelicationController]
cmd/kube-controller-manager/app/core.go:335
↓
[NewReplicationManager]
pkg/controller/replication/replication_controller:52
↓
[NewBaseController]
pkg/controller/replicaset/replica_set.go:126
ReplicaSetControllerの構造体使う
↓
[Run]
pkg/controller/replicaset/replica_set.go:177

※controllerContext(cmd/kube-controller-manager/app/controllermanager.go:289)


エントリーポイント
```go:cmd/kube-controller-manager/app/core.go
func startReplicationController(ctx ControllerContext) (http.Handler, bool, error) {
	go replicationcontroller.NewReplicationManager(
		ctx.InformerFactory.Core().V1().Pods(),
		ctx.InformerFactory.Core().V1().ReplicationControllers(),
		ctx.ClientBuilder.ClientOrDie("replication-controller"),
		replicationcontroller.BurstReplicas,
	).Run(int(ctx.ComponentConfig.ReplicationController.ConcurrentRCSyncs), ctx.Stop)
	return nil, true, nil
}
```

NewReplicationManager
```go:pkg/controller/replication/replication_controller.go
func NewReplicationManager(podInformer coreinformers.PodInformer, rcInformer coreinformers.ReplicationControllerInformer, kubeClient clientset.Interface, burstReplicas int) *ReplicationManager {
	eventBroadcaster := record.NewBroadcaster()

	//受信したイベントを指定した関数へ送信する（klog.Infof)
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})

	//RelicationManagerはreplicaset.ReplicaSetController。ReplicasetControllerのラッパー (次のコード)
	return &ReplicationManager{
		*replicaset.NewBaseController(informerAdapter{rcInformer}, podInformer, clientsetAdapter{kubeClient}, burstReplicas,
			v1.SchemeGroupVersion.WithKind("ReplicationController"),
			"replication_controller",
			"replicationmanager",
			podControlAdapter{controller.RealPodControl{
				KubeClient: kubeClient,
				Recorder:   eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "replication-controller"}),
			}},
		),
	}
}
```


```go:pkg/controller/replicaset/replica_set.go
type ReplicaSetController struct {
	//GroupVersionKindはコントローラータイプを示す。
	//この構造体の異なるインスタンスは異なるGVKを処理できる。
	//この構造体はReplicationControllerを処理するために使用する。
	schema.GroupVersionKind

	kubeClient clientset.Interface
	podControl controller.PodControlInterface

	//ReplicaSetはこれらのレアプリカを作成/削除した後、一時的に中断される。
	//それらの監視イベントを見た後、通常の運転を再開する。
	burstReplicas int
	//テスト用にsyncReplicaSetの挿入を許可する。
	syncHandler func(rcKey string) error

	//podのTTLCacheは各ReplicationControllerが期待するものを作成/削除する。
	expectations *controller.UIDTrackingControllerExpectations

	//NewReplicaSetControllerに渡された、共有インフォーマによって設定されたReplicaSetのストア
	rsLister appslisters.ReplicaSetLister

	//もし少なくとも一回pod storeが同期されていればrsListerSyncedはtrueを返す。
	//テスト用の挿入許可を構造体のメンバーとして追加された。
	rsListerSynced cache.InformerSynced

	//同期の必要があるController
	queue workqueue.RateLimitingInterface
}
```

pod数の監視と同期
```go:pkg/controller/replicaset/replica_set.go
//workersは同時に同期できるreplication-controllerの数。
//大きい数の場合は応答性の高いレプリカ管理が可能。(CPU負荷は高くなる)
func (rsc *ReplicaSetController) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer rsc.queue.ShutDown()

	//ハードコードで"ReplicaSet"を使用するのではなくrsc.Kindを使用している。
	controllerName := strings.ToLower(rsc.Kind)
	klog.Infof("Starting %v controller", controllerName)
	defer klog.Infof("Shutting down %v controller", controllerName)

	//rsc.Kindはコントローラの名前を保持しており,
	//その名前で識別される呼び出し元が同期を待機しており、
	//同期の成功か失敗を示すログメッセージを生成する。
	if !cache.WaitForNamedCacheSync(rsc.Kind, stopCh, rsc.podListerSynced, rsc.rsListerSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.Until(rsc.worker, time.Second, stopCh)
	}

	<-stopCh
}
```

