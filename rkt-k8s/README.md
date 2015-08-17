## Build a Container

Make the binary

```
mkdir -p example.com/host-info/rootfs
GOOS=linux go build -o example.com/host-info/rootfs/host-info github.com/philips/hacks/host-info
```

Build the ACI

```
actool build --overwrite example.com/host-info www-data/host-info-latest-linux-amd64.aci
```

Sign the ACI

```
gpg --sign --detach www-data/host-info-latest-linux-amd64.aci
```

## Build up the index.html

```
cat www-data/index.html
<head>
<meta name="ac-discovery" content="example.com/host-info http://example.com/host-info-{version}-{os}-{arch}.{ext}">
<meta name="ac-discovery-pubkeys" content="example.com/host-info http://example.com/pubkey.gpg">
</head>

<h1>welcome to example.com</h1>
```

Export pubkeys

```
gpg --export --armor brandon.philips@coreos.com > www-data/pubkey.gpg
mv www-data/host-info-latest-linux-amd64.aci.sig www-data/host-info-latest-linux-amd64.aci.asc
```

## Start an HTTP server

```
cd www-data
sudo python -m SimpleHTTPServer 80
```

## Try it all out

Trust the public key

```
sudo ./rkt trust --insecure-allow-http --prefix example.com/host-info
```

Fetch the image

```
sudo ./rkt fetch example.com/host-info
```


```
sudo ./rkt run example.com/host-info
```

Visit host-info

http://192.168.168.134:5483/
