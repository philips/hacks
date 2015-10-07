# Slides

https://speakerdeck.com/philips/coreos-plus-kubernetes-at-linuxcon-europe

# Setup CoreOS + Kubernetes

Find the docs and set it up:

- https://coreos.com
- https://coreos.com/docs/
- https://coreos.com/kubernetes/docs/latest/
- https://coreos.com/kubernetes/docs/latest/kubernetes-on-vagrant.html


# Run our first service

```
kubectl run host-info --image=quay.io/philips/host-info
kubectl expose rc host-info --port=80 --target-port=5483 --type=NodePort
kubectl describe service host-info
kubectl scale rc host-info --replicas=5
kubectl get pods
kubectl stop host-info
watch -n 0.2 curl http://172.17.4.202:32430/
```

# Create a new namespace

kubectl create -f namespace/namespace-guestbook.yaml

# Run a more complex service

This is going to setup a three tier app with a load balancer, web app, and database.

```
for i in *.json; do kubectl create --namespace=guestbook -f ${i}; done
```

Visit http://172.17.4.202:30002/
