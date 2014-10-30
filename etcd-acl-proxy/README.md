This is a simple proxy that will deny non-GET requests to /v2/keys/read-only
and deny GET requests to /v2/keys/write-only. It uses the new members API
exposed in etcd 0.5.0 so a build of etcd from git is required.

Examples:

```
philips:874c7de0d742fb9fcfcc/ (master*) $ curl --verbose -X PUT  http://localhost:7777/v2/keys/read-only -d value=asdf
* Hostname was NOT found in DNS cache
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 7777 (#0)
> PUT /v2/keys/read-only HTTP/1.1
> User-Agent: curl/7.37.1
> Host: localhost:7777
> Accept: */*
> Content-Length: 10
> Content-Type: application/x-www-form-urlencoded
>
* upload completely sent off: 10 out of 10 bytes
< HTTP/1.1 501 Not Implemented
< Date: Thu, 30 Oct 2014 00:48:08 GMT
< Content-Length: 0
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host localhost left intact
philips:874c7de0d742fb9fcfcc/ (master) $ curl --verbose  http://localhost:7777/v2/keys/write-only/asdf
* Hostname was NOT found in DNS cache
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 7777 (#0)
> GET /v2/keys/write-only/asdf HTTP/1.1
> User-Agent: curl/7.37.1
> Host: localhost:7777
> Accept: */*
>
< HTTP/1.1 501 Not Implemented
< Date: Thu, 30 Oct 2014 00:48:21 GMT
< Content-Length: 0
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host localhost left intact
```
