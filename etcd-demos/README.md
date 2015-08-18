# skydns

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

# confd


confd -onetime -backend etcd -node 127.0.0.1:4001
