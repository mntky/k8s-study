# ReplicasetController

定義されたReplicaSetオブジェクトを実行中のPodと同期させる。


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


