# single master cluster

cluster bootstrap
```
# kubeadm init --pod-network-cidr=10.224.0.0/16 --apiserver-advertise-address=<ipaddress>
# mkdir $HOME/.kube
# cp /etc/kubernetes/admin.conf $HOME/.kube/config
# chown $(id -u):$(id -g) $HOME/.kube/config
# exprot KUBECONFIG=/etc/kubernetes/admin.conf
```

network addon install
```
# kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/a70459be0084506e4ec919aa1c114638878db11b/Documentation/kube-flannel.yml
```

Isolating control plane nodes
```
kubectl taint nodes --all node-role.kubernetes.io/master-
```


















参考URL: https://kubernetes.io/ja/docs/setup/independent/create-cluster-kubeadm/
