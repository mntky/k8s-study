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



### Run
```go:pkg/controller/replicaset/replica_set.go:177
func (rsc *ReplicaSetController) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer rsc.queue.ShutDown()

	controllerName := strings.ToLower(rsc.Kind)
	klog.Infof("Starting %v controller", controllerName)
	defer klog.Infof("Shutting down %v controller", controllerName)

	if !cache.WaitForNamedCacheSync(rsc.Kind, stopCh, rsc.podListerSynced, rsc.rsListerSynced) {
		return
	}

	//wait.UntilはstopChが閉じるまでループして、毎周rsc.workerを実行する。
	for i := 0; i < workers; i++ {
		go wait.Until(rsc.worker, time.Second, stopCh)
	}

	<-stopCh
}
```

↑のwait.Untilで呼ばれている関数
```
func JitterUntil(f func(), period time.Duration, jitterFactor float64, sliding bool, stopCh <-chan struct{}) {
	var t *time.Timer
	var sawTimeout bool

	for {
		//stopChに値が入っていない場合はreturn、入っている場合はdefault以下の処理が実行される。
		//参考:https://qiita.com/najeira/items/71a0bcd079c9066347b4
		select {
		case <-stopCh:
			return
		default:
		}

		jitteredPeriod := period
		if jitterFactor > 0.0 {
			jitteredPeriod = Jitter(period, jitterFactor)
		}

		if !sliding {
			t = resetOrReuseTimer(t, jitteredPeriod, sawTimeout)
		}

		func() {
			defer runtime.HandleCrash()
			f()
		}()

		if sliding {
			t = resetOrReuseTimer(t, jitteredPeriod, sawTimeout)
		}

		// NOTE: b/c there is no priority selection in golang
		// it is possible for this to race, meaning we could
		// trigger t.C and stopCh, and t.C select falls through.
		// In order to mitigate we re-check stopCh at the beginning
		// of every loop to prevent extra executions of f().
		select {
		case <-stopCh:
			return
		case <-t.C:
			sawTimeout = true
		}
	}
}
```

