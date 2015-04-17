# App Container Demo

### Packaging a Go Application

cat main.go

#### Build a statically linked Go binary

```
CGO_ENABLED=0 GOOS=linux go build -o etcd-backup -a -tags netgo -ldflags '-w' .
```

```
file etcd-backup
```

#### Create the image manifest

Edit manifest.json

##### Validate the image manifest

```
actool -debug validate manifest.json
```

#### Create the layout and the rootfs

```
mkdir -p etcd-backup-layout/rootfs/bin
```

Copy the image manifest

```
cp manifest.json etcd-backup-layout/manifest
```

Copy the etcd-backup binary

```
cp etcd-backup etcd-backup-layout/rootfs/bin/
```

#### Build the application image

```
actool build etcd-backup-layout/ etcd-backup-0.0.1-linux-amd64.aci
```

##### Validate the application image

```
actool -debug validate etcd-backup-0.0.1-linux-amd64.aci
```

### Signing the ACI

#### Generate a gpg signing key

```
cat gpg-batch
```

```
gpg --batch --gen-key gpg-batch
```

List the keys.

```
gpg --no-default-keyring --secret-keyring ./rkt.sec --keyring ./rkt.pub --list-keys
```

Export the public key.

```
gpg --no-default-keyring --armor \
--secret-keyring ./rkt.sec --keyring ./rkt.pub \
--output pubkeys.gpg \
--export "<brandon.philips@coreos.com>"
```

#### Create the detached signature

```
gpg --no-default-keyring --armor \
--secret-keyring ./rkt.sec --keyring ./rkt.pub \
--output etcd-backup-0.0.1-linux-amd64.aci.asc \
--detach-sig etcd-backup-0.0.1-linux-amd64.aci
```

```
etcd-backup-0.0.1-linux-amd64.aci.asc
etcd-backup-0.0.1-linux-amd64.aci
pubkeys.gpg
```

```
cp etcd-backup-0.0.1-linux-amd64.aci etcd-backup-0.0.1-linux-amd64.aci.asc /opt/images/example.com/
```

```
cp pubkeys.gpg /opt/images
```

#### Test the discover end-point

```
actool discover --insecure example.com/etcd-backup:0.0.1,os=linux,arch=amd64
```

```
http://example.com/images/example.com/etcd-backup-0.0.1-linux-amd64.aci.asc
http://example.com/images/example.com/etcd-backup-0.0.1-linux-amd64.aci
http://example.com/pubkeys.gpg
```

# Rocket Demo

### Fetch the example.com/etcd-backup:0.0.1 ACI

```
rkt fetch example.com/etcd-backup:0.0.1
```

### Trust the example.com signing key

```
rkt trust --prefix example.com/etcd-backup
```

### Run the example.com/etcd-backup:0.0.1 ACI

```
rkt run example.com/etcd-backup:0.0.1
```

### Listing Containers

```
rkt list
```

After the container has stopped:

```
rkt list
```

### rkt gc

Cleaning up old containers with rkt gc

```
rkt gc
```

```
rkt gc
```

Noting happens? Why?

```
rkt gc --grace-period=10s
```

#### Create a systemd units

```
systemctl cat rkt-gc.timer
```
