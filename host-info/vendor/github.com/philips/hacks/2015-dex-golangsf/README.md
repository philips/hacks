## Get Dex Running

```
git clone https://github.com/coreos/dex
cd dex
./build
./bin/dex-worker  --no-db
```


```
curl localhost:5556/.well-known/openid-configuration | python -m json.tool
```


## Get App Running

```
./bin/example-app --cent-id=XXX \
  --client-secret=secrete \
  --discovery=http://127.0.0.1:5556
```

```
open http://127.0.0.1:5555
```


