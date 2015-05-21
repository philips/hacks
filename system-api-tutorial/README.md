## syscalls

```
toolbox
ps aux | grep ssh
strace -p 
```

Inside toolbox:
```
md5sum /usr/lib64/libm.so.6
```

Outside toolbox: 
```
md5sum /usr/lib64/libm.so.6
```

## omaha

```
journalctl -u update-engine
```

application: CoreUpdate, devserver

## namespaces

```
ps aux | grep nginx
nsenter -n -t 847
```

application: docker, nspawn, rkt

## cgroups

```
cd /sys/fs/cgroups/
cat /sys/fs/cgroup/memory/system.slice/etcd2.service/memory.limit_in_bytes
cat /sys/fs/cgroup/memory/system.slice/etcd2.service/memory.memsw.usage_in_bytes
```

applications: cAdvisor, rkt, nspawn, docker

## docker 

```
systemd-run ncat -vlk 1111 -c 'ncat -U /var/run/docker.sock'
```

```
docker ps
curl localhost:1111/containers/json
```

```
docker images
curl localhost:1111/containers/json
```

applications: kubernetes, docker client

## dbus + systemd

```
dbus-monitor --system
sudo systemd-run rkt run coreos.com/etcd,version=v2.0.10
```

applications: systemctl, fleet, kubernetes

## host configuration

vagrant:
```
cat /var/lib/coreos-vagrant/vagrantfile-user-data
```

aws:
```
curl http://169.254.169.254/latest/meta-data/
```

## etcd

```
etcdctl --debug -o json set foobar baz
etcdctl --debug -o json set --swap-with-index 11282 foobar baz2
```

applications: locksmith, kubernetes, vulcan, etc

## k8s

```
kubectl run-container my-nginx --image=nginx --replicas=1 --port=80
```

```
curl http://192.168.168.154:59101/api/v1beta3/pods
curl http://192.168.168.154:59101/api/v1beta3/nodes
```

```
kubectl resize --replicas=2 rc my-nginx
```


```
kubectl expose rc my-nginx --port=80
```

applications: kubectl, dashboards, automation software

## fleet

```
systemd-run ncat -vlk 1337 -c 'ncat -U /var/run/fleet.sock'
```

```
fleetctl list-units
curl http://localhost:1337/fleet/v1/state?alt=json
```

applications: CoreGI, fleetctl

