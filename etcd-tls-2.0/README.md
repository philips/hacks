Put a etcd 2.0 build into the bin directory.

Initial setup:

```
etcd-ca  init
etcd-ca new-cert --help
etcd-ca new-cert infra2
etcd-ca new-cert infra3
etcd-ca sign infra1
etcd-ca sign infra2
etcd-ca sign infra3
etcd-ca export --insecure infra1 | tar xvf -
etcd-ca export --insecure infra2 | tar xvf -
etcd-ca export --insecure infra3 | tar xvf -
etcd-ca export | tar xzvf -
```
