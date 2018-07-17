# doc
doc
导读：众所周知，K8S安装难点在于镜像下载。查询网上许多关于K8S安装部署的文章，作者要么是将镜像下载使用的工作交给了读者完成，要么是放在共享云盘中，由读者下载后导入使用。都知道，如果你不是该云盘提供商的会员，镜像下载速度堪比蜗牛。总之，三个字：”不方便“。基于最大化便利刚接触K8S的同行快速上手实践的目的，在参阅了众多同行关于K8S部署的文章，并经过反复实验验证后，将本人的实验成果在这里分享，希望可以帮助有需要的朋友。由于时间仓促，一些对本文撰写有价值的文章没有仔细考证作者和出处，在本文末尾的参阅文章中可能没有注明。如你发现本文中一些内容的原创属于你本人，请邮件联系本人处理。并在此感谢对本文创作产生帮助的作者。谢谢！

本文K8S集群高可用方案采用Keepalived。

实验环境：
1、3台centos 1611版本虚拟机，mini安装。Linux localhost 3.10.0-514.el7.x86_64 #1 SMP Tue Nov 22 16:42:41 UTC 2016 x86_64 x86_64 x86_64 GNU/Linux

2、docker version
Client:
Version: 1.13.1
API version: 1.26
Package version: <unknown>
Go version: go1.8.3
Git commit: 774336d/1.13.1
Built: Wed Mar 7 17:06:16 2018
OS/Arch: linux/amd64

Server:
Version: 1.13.1
API version: 1.26 (minimum version 1.12)
Package version: <unknown>
Go version: go1.8.3
Git commit: 774336d/1.13.1
Built: Wed Mar 7 17:06:16 2018
OS/Arch: linux/amd64
Experimental: false

3、etcd Version: 3.1.13

4、kubeadm，kubelet，kubectl，kube-cni版本如下：
-rw-r–r– 1 root root 18176678 Mar 30 00:08 5844c6be68e95a741f1ad403e5d4f6962044be060bc6df9614c2547fdbf91ae5-kubelet-1.10.0-0.x86_64.rpm
-rw-r–r– 1 root root 17767206 Mar 30 00:07 8b0cb6654a2f6d014555c6c85148a5adc5918de937f608a30b0c0ae955d8abce-kubeadm-1.10.0-0.x86_64.rpm
-rw-r–r– 1 root root 7945458 Mar 30 00:07 9ff2e28300e012e2b692e1d4445786f0bed0fd5c13ef650d61369097bfdd0519-kubectl-1.10.0-0.x86_64.rpm
-rw-r–r– 1 root root 9008838 Mar 5 21:56 fe33057ffe95bfae65e2f269e1b05e99308853176e24a4d027bc082b471a07c0-kubernetes-cni-0.6.0-0.x86_64.rpm

5、k8s网络组件：flannel:v0.10.0-amd64

6、实验网络规划如下：
host1 172.18.0.154/22
host2 172.18.0.155/22
host3 172.18.0.156/22
VIP 172.18.0.192/22

文章视频链接：https://pan.baidu.com/s/1XVagd765eGacuoR_cgesiQ

安装步骤：

0、请先从该链接下载后面步骤所需脚本https://pan.baidu.com/s/1oK7PRLeeYHrouNCRgIQlcQ

1、在3台主机中执行基础环境配置脚本 base-env-config.sh。

2、在主机1执行脚本 host1-base-env.sh

3、在主机2执行脚本 host2-base-env.sh

4、在主机3执行脚本 host3-base-env.sh


免密码登陆
ssh-keygen -t rsa
scp -p ~/.ssh/id_rsa.pub root@193.112.34.152:/root/.ssh/authorized_keys
scp -p ~/.ssh/id_rsa.pub root@120.77.41.248:/root/.ssh/authorized_keys

5、在host1主机执行如下命令
scp -r /etc/etcd/ssl root@193.112.34.152:/etc/etcd/
scp -r /etc/etcd/ssl root@120.77.41.248:/etc/etcd/

6、在3台主机中分别执行脚本 etcd.sh


查看错误 journalctl -xe

7、查看keepalived状态
systemctl status keepalived

8、查看etcd运行状态
在host1,host2,host3分别执行如下命令：[参数名称前面是两个’-‘，注意]
etcdctl --endpoints=https://${NODE_IP}:2379 --ca-file=/etc/etcd/ssl/ca.pem --cert-file=/etc/etcd/ssl/etcd.pem --key-file=/etc/etcd/ssl/etcd-key.pem cluster-health

etcdctl --endpoints=https://${NODE_IP}:2379 --ca-file=/etc/etcd/ssl/ca.pem --cert-file=/etc/etcd/ssl/etcd.pem --key-file=/etc/etcd/ssl/etcd-key.pem  member list

#sed -i -E 's/--cgroup-driver=systemd/--cgroup-driver=cgroupfs/' /etc/systemd/system/kubelet.service.d/10-kubeadm.conf

9、在3台主机上安装kubeadm,kubelet,kubctl,docker
yum install kubelet kubeadm kubectl kubernetes-cni docker -y

10、在3台主机禁用docker启动项参数关于SELinux的设置[参数名称前面是两个’-‘，注意]
sed -i 's/--selinux-enabled/--selinux-enabled=false/g' /etc/sysconfig/docker

11、在3台主机的kubelet配置文件中添加如下参数[参数名称前面是两个’-‘，注意]
sed -i '9a\Environment="KUBELET_EXTRA_ARGS=--pod-infra-container-image=registry.cn-hangzhou.aliyuncs.com/osoulmate/pause-amd64:3.0"' /etc/systemd/system/kubelet.service.d/10-kubeadm.conf


sed -i '9a\Environment="KUBELET_CGROUP_ARGS=--cgroup-driver=systemd""' /etc/systemd/system/kubelet.service.d/10-kubeadm.conf


遇到问题首先查看kubelet程序

systemctl status kubelet
journalctl -xeu kubelet

12、在3台主机添加docker加速器配置（可选步骤）
 #请自行申请阿里云账号获取镜像加速链接
cat <<EOF > /etc/docker/daemon.json
{
"registry-mirrors": ["https://lmu8xt7e.mirror.aliyuncs.com"]
}
EOF

--registry-mirror=https://lmu8xt7e.mirror.aliyuncs.com --insecure-registry gcr.io

13、在3台主机分别执行以下命令
systemctl daemon-reload
systemctl enable docker && systemctl restart docker
systemctl enable kubelet && systemctl restart kubelet

14、在3台主机中分别执行kubeadmconfig.sh生成配置文件config.yaml

15、在host1主机中首先执行kubeadm初始化操作
命令如下：
kubeadm init --config config.yaml --ignore-preflight-errors all

16、在host1主机中执行初始化后操作
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

17、将主机host1中kubeadm初始化后生成的证书和密钥文件拷贝至host2,host3相应目录下
scp -r /etc/kubernetes/ root@120.77.41.248:/etc/kubernetes/
scp -r /etc/kubernetes/ root@193.112.34.152:/etc/kubernetes/
		
        #
18、为主机host1安装网络组件 podnetwork【这里选用flannel】
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
systemctl stop kubelet    #由于kubelet会调用docker到默认url【谷歌】下载镜像，所以先禁用
systemctl restart docker
docker pull registry.cn-hangzhou.aliyuncs.com/osoulmate/flannel:v0.10.0-amd64
systemctl start kubelet

kubectl get nodes

kubectl get pods --all-namespaces -o wide


kubectl describe pods 

19、在host2,host3上执行如下命令 


kubeadm init --config config.yaml --ignore-preflight-errors all

注意 切换到普通用户
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config


注意 切换到root
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
systemctl stop kubelet #如果提示需要systemctl daemon-reload，则将守护进程重启后再停止kubelet服务。
systemctl restart docker
docker pull registry.cn-hangzhou.aliyuncs.com/osoulmate/flannel:v0.10.0-amd64
systemctl start kubelet


获取正在运行的容器XCXXX的 IP。
docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' kube-flannel-ds


systemctl stop kubelet #如果提示需要systemctl daemon-reload，则将守护进程重启后再停止kubelet服务。
systemctl restart docker
systemctl start kubelet

20、查看集群各节点状态【这里仅在host1主机查看结果】
kubectl -s https://47.75.152.154:6443 get nodes

kubectl describe pods/kubernetes-dashboard-7b44ff9b77-g4mhz --namespace="kube-system" 

kubectl logs kubernetes-dashboard-7b44ff9b77-g4mhz -n kube-system


kubectl get pods --all-namespaces -o wide

查看分配的 NodePort
$ kubectl get services kubernetes-dashboard -n kube-system
NAME                   CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
kubernetes-dashboard   10.254.224.130   <nodes>       80:30312/TCP   25s
检查 controller
$ kubectl get deployment kubernetes-dashboard  -n kube-system
NAME                   DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
kubernetes-dashboard   1         1         1            1           3m
$ kubectl get pods  -n kube-system | grep dashboard
kubernetes-dashboard-1339745653-pmn6z   1/1       Running   0          4m

kubectl proxy --address='47.75.152.154' --port=8086 --accept-hosts='^*$'


kubectl delete -f kubernetes-dashboard.yaml

kubectl delete pod kubernetes-dashboard-7b44ff9b77-m7khw --grace-period=0 --force


https://47.75.152.154:6443/api/v1/proxy/namespaces/kube-system/services/kubernetes-dashboard

http://47.75.152.154:8080/api/v1/proxy/namespaces/kube-system/services/kubernetes-dashboard

https://47.75.152.154:6443/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/

https://47.75.152.154:6443/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/login

curl -H "Content-type: application/json" -X GET -d '{"srcRef":"1002"}' --cacert /etc/kubernetes/pki/front-proxy-ca.crt --cert /etc/kubernetes/pki/front-proxy-client.crt  --key /etc/kubernetes/pki/front-proxy-client.key "https://47.75.152.154:6443/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/"

curl -H "Content-type: application/json" -X POST -d '{"srcRef":"1002"}' --no-check-certificate "https://47.75.152.154:6443/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/"

--no-check-certificate


看整体的全部信息
kubectl cluster-info

[root@host1 ~]# kubectl get nodes
NAME STATUS ROLES AGE VERSION
host1 Ready master 5m v1.10.0
host2 Ready master 1m v1.10.0
host3 Ready master 1m v1.10.0
[root@host1 ~]# kubectl get po –all-namespaces
NAMESPACE NAME READY STATUS RESTARTS AGE
kube-system coredns-7997f8864c-k9dcx 1/1 Running 0 5m
kube-system coredns-7997f8864c-sv9rv 1/1 Running 0 5m
kube-system kube-apiserver-host1 1/1 Running 1 4m
kube-system kube-apiserver-host2 1/1 Running 0 1m
kube-system kube-apiserver-host3 1/1 Running 0 1m
kube-system kube-controller-manager-host1 1/1 Running 1 4m
kube-system kube-controller-manager-host2 1/1 Running 0 1m
kube-system kube-controller-manager-host3 1/1 Running 0 1m
kube-system kube-flannel-ds-88tz5 1/1 Running 0 1m
kube-system kube-flannel-ds-g9dpj 1/1 Running 0 2m
kube-system kube-flannel-ds-h58tp 1/1 Running 0 1m
kube-system kube-proxy-6fsvq 1/1 Running 1 5m
kube-system kube-proxy-g8xnb 1/1 Running 1 1m
kube-system kube-proxy-gmqv9 1/1 Running 1 1m
kube-system kube-scheduler-host1 1/1 Running 1 5m
kube-system kube-scheduler-host2 1/1 Running 1 1m
kube-system kube-scheduler-host3 1/1 Running 0 1m

21、高可用验证
将host1关机，在host3上执行
while true; do sleep 1; kubectl get node;date; done
在host2上观察keepalived是否已切换为主状态。

Q&A

1、为什么在kubeadm init时出现kubelet 版本不支持系统中安装的etcd的报错？

因为本文k8s管理组件kubeadm,kubectl,kubelet的安装源为阿里云源，阿里云源会和最新k8s版本保持同步。如出现版本不兼容的问题，请按照报错提示安装相应版本的etcd或kubelet，kubeadm,kubectl组件。

2、为什么安装时间同步软件chrony？

由于集群采用keepalived检测集群各节点的活动状态，如不能保证各节点时间同步，会导致心跳异常，进而影响集群的故障倒换。当然，你也可以采用其它时间同步措施，只要能保证各节点之间的时间同步即可。

3、步骤18中为什么用kubectl应用了网络组件后还要用docker pull从阿里云拉取镜像呢？

kubectl应用了flannel组件后默认会从谷歌镜像库中拉取镜像，所以要先停止kubelet服务，使其停止对docker的调用，在我们手动从阿里云拉取镜像后，再重启kubelet服务，k8s相关服务会自动识别镜像。在host2,host3主机kubeadm init完成后可能还会出现其它镜像包未拉取完成的情况，这时也可以采用这种方法：即停止kubelet服务，重启docker服务，手动拉取镜像【确定需要拉取那些镜像可先在主机上使用 kubectl get po –all-namespaces命令获取各主机镜像的当前状态。如READY列显示0/1表示镜像仍在拉取/创建中，可使用你下载的k8s压缩包中名称为阿里云镜像包的txt文档中的相应命令】，之后再启动kubelet服务。除flannel镜像外，理论上所有镜像在kubeadm init中都会从阿里云镜像库中拉取，所以，如果host2,host3在kubeadm init时有镜像没有拉取完成，可等待1-2分钟，如还未成功，直接重启kubelet服务即可。

go-mesher

func pingInstance() {
	it := registry.MicroserviceInstanceCache.Items()
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)

	for k, v := range it {
		mics, ok := v.Object.([]*registry.MicroServiceInstance)
		if !ok {
			continue
		}

		okMicroServiceInstance := make([]*registry.MicroServiceInstance, 0)
		select {
		case <-ctx.Done(): //time out
			lager.Logger.Debugf("time out !!!")
			return
		default:
			for _, v := range mics {
				if err := ping(v.DefaultEndpoint); err == nil { //ok
					lager.Logger.Debugf("ping %+v %+v ", v.DefaultEndpoint, "ok")
					okMicroServiceInstance = append(okMicroServiceInstance, v)
				}

			}
		}

		registry.MicroserviceInstanceCache.Set(k, okMicroServiceInstance, 0)

	}

}

func ping(add string) error {
	conn, err := net.DialTimeout("tcp", add, time.Second*5)
	if err != nil {
		lager.Logger.Debugf("ping %+v %+v ", add, err)
		return err
	}
	defer conn.Close()
	return nil
}
