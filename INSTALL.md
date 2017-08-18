# OpenSDS Installation with Kubernetes Service Catalog
In this tutorial, we want to show how OpenSDS provide Ceph storage as a service for Kubernetes users through Service Catalog. Now enjoy your trip!

## Pre-configuration

### Check it out about your os (very important)
Please NOTICE that the installation tutorial is tested on Ubuntu17.04, and we SUGGEST you follow our styles and use Ubuntu16.04+.

### Download and install Golang
```
wget https://storage.googleapis.com/golang/go1.7.6.linux-amd64.tar.gz
tar xvf go1.7.6.linux-amd64.tar.gz -C /usr/local/
mkdir -p $HOME/gopath/src
mkdir -p $HOME/gopath/bin
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/gopath/bin' >> /etc/profile
echo 'export GOPATH=$HOME/gopath' >> /etc/profile
source /etc/profile
go version (check if go has been installed)
```

### Install some package dependencies
```
apt-get install gcc make libc-dev docker.io
```
If the docker command doesn't work, try to restart it:
```
sudo service docker stop
sudo nohup docker daemon -H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock &
```

### Install containerized Ceph cluster (If your machine has already a ceph cluster installed, just skip this step)
Since this is a test, we choose to install a simple containerized ceph single-node cluster:
```
sudo docker run -d --net=host -v /etc/ceph:/etc/ceph -e MON_IP=your_host_ip -e CEPH_PUBLIC_NETWORK=your_host_ip/24 ceph/demo
```
NOTICE that ```your_host_ip``` means the real ip address of your machine. After this container has been running, Add ```rbd_default_features = 1``` one line in ```/etc/ceph/ceph.conf``` file:
```
echo 'rbd_default_features = 1' >> /etc/ceph/ceph.conf
```

## Local Kubernetes Cluster Setup

### Download and unpackage k8s and etcd source code
```
wget https://github.com/kubernetes/kubernetes/archive/v1.6.0.tar.gz
wget https://github.com/coreos/etcd/releases/download/v3.2.0/etcd-v3.2.0-linux-amd64.tar.gz
tar -zxvf etcd-v3.2.0-linux-amd64.tar.gz
cp etcd-v3.2.0-linux-amd64/etcd /usr/local/bin
tar -zxvf v1.6.0.tar.gz
```

### Setup your cluster (suggest run as ```root```)
```
source /etc/profile
go get -u github.com/cloudflare/cfssl/cmd/...
KUBE_ENABLE_CLUSTER_DNS=true kubernetes-1.6.0/hack/local-up-cluster.sh -O
```
The setup process may run for a while, and after it done, some tips will occur to help you configure your cluster correctly.

### Check the status of cluster (open another terminal)
Follow the tips and configure your cluster, and then add another path ```your_path_to/kubernetes-1.6.0/hack/local-up-cluster.sh``` in
the ```PATH``` variable configured in ```/etc/profile/ file.

After that, run ```kubectl.sh get po``` to check the status of kubernetes cluster.

## Service Catalog Setup

### Install helm (from scipt)
```
curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get > get_helm.sh
chmod 700 get_helm.sh
./get_helm.sh

kubectl.sh get po -n kube-system (check if the till-deploy pod is running)
```

### Service Catalog download and install
```
mkdir -p gopath/src/github.com/kubernetes-incubator
wget https://github.com/kubernetes-incubator/service-catalog/archive/v0.0.13.tar.gz
tar xvf service-catalog-0.0.13.tar.gz -C gopath/src/github.com/kubernetes-incubator
helm install gopath/src/github.com/kubernetes-incubator/service-catalog-0.0.13/charts/catalog --name catalog --namespace catalog

kubectl.sh get po -n catalog (check if the catalog api-server and controller-manager pod are running)
```

## OpenSDS Service Broker Setup

### OpenSDS cluster install
To ensure service broker connecting to OpenSDS api-service, you probably need to configure your service ip:
```
docker run -it --net=host -v /var/log/opensds:/var/log/opensds leonwanghui/opensds-controller:v1alpha1 /usr/bin/osdslet --api-endpoint=your_host_ip:50040
docker run -d --net=host -v /var/log/opensds:/var/log/opensds -v /etc/opensds:/etc/opensds -v /etc/ceph:/etc/ceph leonwanghui/opensds-dock:v1alpha1
```

### OpenSDS service broker install
```
mkdir -p gopath/src/github.com/leonwanghui
git clone https://github.com/leonwanghui/opensds-broker.git gopath/src/github.com/leonwanghui
cd gopath/src/github.com/leonwanghui/opensds-broker
helm install charts/opensds-broker --name opensds-broker --namespace opensds-broker

kubectl get po -n opensds-broker (check if opensds broker pod is running)
```
