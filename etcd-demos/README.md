slides: https://speakerdeck.com/philips/etcd-at-containercon-2015

# locksmith - coordinate cluster reboots

- Homepage: https://github.com/coreos/locksmith

We will use a copy of coreos-vagrant to test this out. Follow instructions here: https://github.com/coreos/coreos-vagrant

```
cd coreos-vagrant
for i in `seq 2 3` ; do vagrant ssh core-0${i} -c 'locksmithctl send-need-reboot'; done
```

```
etcdctl get /coreos.com/updateengine/rebootlock/semaphore
```

# skydns - server up DNS from configuration files

- Homepage: https://github.com/skynetservices/skydns

First, download the latest release and follow the instructions at the bottom of the page for a single member cluster running on localhost. https://github.com/coreos/etcd/releases


```
export ETCD_MACHINES='http://127.0.0.1:2379'
skydns -addr 127.0.0.1:5355
```

```
etcdctl set /skydns/local/skydns/east/production/rails     '{"host":"service6.example.com","priority":20}'
dig -p 5355 @127.0.0.1 SRV rails.production.east.skydns.local
```

```
etcdctl set /skydns/local/skydns/east/production/mail     '{"host":"mail.example.com","priority":20}'
dig -p 5355 @127.0.0.1 SRV mail.production.east.skydns.local
```

```
dig -p 5355 @127.0.0.1 SRV production.east.skydns.local
```

# confd - generate configuration files from etcd

- Homepage: https://github.com/kelseyhightower/confd

```
$ cat confd/conf.d/myconfig.toml
[template]
src = "myconfig.conf.tmpl"
dest = "/tmp/myconfig.conf"
keys = [
    "/myapp/database/url",
     "/myapp/database/user",
]
```

```
$ cat confd/templates/myconfig.conf.tmpl
[myconfig]
database_url = {{getv "/myapp/database/url"}}
database_user = {{getv "/myapp/database/user"}}
```

```
confd -onetime -backend etcd -node 127.0.0.1:2379 -confdir confd
```

# vulcand - an http load balancer

- Homepage: https://github.com/mailgun/vulcand

```
vulcand  -etcd=http://localhost:2379 -logSeverity=INFO
```

```
etcdctl set /vulcand/backends/b1/servers/srv1 '{"URL": "http://localhost:5000"}'
etcdctl set /vulcand/frontends/f1/frontend '{"Type": "http", "BackendId": "b1", "Route": "Path(`/`)"}'
```

```
vctl top
```

```
python -m SimpleHTTPServer 5000
```

```
boom http://localhost:8181
```

```
etcdctl set /vulcand/backends/b1/servers/srv2 '{"URL": "http://localhost:3000"}'
```

```
go run server.go
```

```
boom http://localhost:8181
```

# kubernetes - service discovery and load balancing

- Homepage: kubernetes.io
- get up and running: https://coreos.com/blog/introducing-the-kubelet-in-coreos/
- example app: https://github.com/kubernetes/kubernetes/tree/master/examples/https-nginx

```
./kubectl create -f secret.json
./kubectl create -f nginx.yaml
```

```
./kubectl get service
```

```
curl -k https://<IP>
```
