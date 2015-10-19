# Slides

https://speakerdeck.com/philips/coreos-plus-kubernetes-at-linuxcon-europe

# Setup CoreOS + Kubernetes

Find the docs and set it up:

- https://coreos.com
- https://coreos.com/docs/
- https://coreos.com/kubernetes/docs/latest/
- https://coreos.com/kubernetes/docs/latest/kubernetes-on-vagrant.html


# Run our first service

Use the right cluster:

```
kubectl config use-context vagrant-single
```

Start it up:
```
kubectl run host-info --image=quay.io/philips/host-info
kubectl expose rc host-info --port=80 --target-port=5483 --type=NodePort
kubectl describe service host-info
kubectl get pods
```


Scale it:
```
kubectl scale rc host-info --replicas=2
watch -n 0.2 curl http://172.17.4.99:32688/
```

Shut it down

```
kubectl stop rc host-info
kubectl stop service host-info
```

# Create a new namespace

kubectl config use-context vagrant

kubectl create -f namespace.yaml

# Run a more complex service

This is going to setup a three tier app with a load balancer, web app, and database.

```
for i in *.json; do kubectl create --namespace=guestbook -f ${i}; done
```

Visit http://172.17.4.202:30002/
