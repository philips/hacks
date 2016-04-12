## Tectonic Console Rolling Update

Create the service

```
kubectl create -f tectonic-console-service.json
```

Deploy the container

```
kubectl create -f tectonic-console-v0.1.4.yaml
```

Checkout the version number

```
while true ; do curl --silent http://172.17.4.99:31190/version | jq .version; sleep 1; done
```

Roll to the new version

```
kubectl rolling-update tectonic-console-v0.1.4 -f tectonic-console-v0.1.5.yaml --update-period=5s
```


