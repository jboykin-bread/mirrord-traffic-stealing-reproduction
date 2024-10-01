## Steps to Reproduce

Build the local service image
```shell
$ docker build -t local-dev/color-server:latest .
```

Create a local kind cluster and install istio 1.23
```shell
$ kind create cluster --name=mirrord-stealing-test

$ istioctl install --context kind-mirrord-stealing-test
```

Load the local docker image into the kind cluster
```shell
$ kind load docker-image local-dev/color-server:latest --name=mirrord-stealing-test
```

Deploy the test server
```shell
$ kubectl --context kind-mirrord-stealing-test apply -f kube-manifests.yaml
```

Ensure it's healthy once deployed
```shell
$ kubectl --context kind-mirrord-stealing-test -n color-server logs color-server-77ff58464f-xd2tg

2024/10/01 19:12:06 INFO Starting server on :8000
```

Open up a test shell with networking tools, and hit the service
```shell
$ kubectl --context kind-mirrord-stealing-test -n color-server run test-shell --rm -i --tty --image nicolaka/netshoot

 test-shell  ~  curl -m 10 -vv http://color-server:8000/color
* Host color-server:8000 was resolved.
* IPv6: (none)
* IPv4: 10.96.153.133
*   Trying 10.96.153.133:8000...
* Connected to color-server (10.96.153.133) port 8000
> GET /color HTTP/1.1
> Host: color-server:8000
> User-Agent: curl/8.7.1
> Accept: */*
> 
* Request completely sent off
< HTTP/1.1 200 OK
< content-type: text/plain
< date: Tue, 01 Oct 2024 19:16:49 GMT
< content-length: 4
< x-envoy-upstream-service-time: 4
< server: envoy
< 
* Connection #0 to host color-server left intact
blue
```

In another terminal window, build the go service, and run mirrord
```shell
$ go build -o color-server

$ mirrord-3.118.0 exec --context kind-mirrord-stealing-test -n color-server -e \
  --steal --target pod/color-server-77ff58464f-xd2tg -- ./color-server -color=green
* service mesh detected: istio
* Running binary "./color-server" with arguments: ["-color=green"].
* mirrord will target: pod/color-server-77ff58464f-xd2tg, no configuration file was loaded
* operator: the operator will be used if possible
* env: all environment variables will be fetched
* fs: file operations will default to read only from the remote
* incoming: incoming traffic will be stolen
* outgoing: forwarding is enabled on TCP and UDP
* dns: DNS will be resolved remotely
⠲ mirrord exec
    ✓ running on latest (3.118.0)!
    ✓ ready to launch process
      ✓ layer extracted
      ✓ operator not found
      ✓ container created
      ✓ container is ready
    ✓ config summary
      2024/10/01 15:20:54 INFO Starting server on :8000
```
Requests to `color-server:8000` will either time out or will get a response from the Pod, rather 
than the local service. You can validate this in the temp shell:
```shell
 test-shell  ~  curl -m 10 -vv http://color-server:8000/color
* Host color-server:8000 was resolved.
* IPv6: (none)
* IPv4: 10.96.153.133
*   Trying 10.96.153.133:8000...
* Connected to color-server (10.96.153.133) port 8000
> GET /color HTTP/1.1
> Host: color-server:8000
> User-Agent: curl/8.7.1
> Accept: */*
> 
* Request completely sent off
* Operation timed out after 10010 milliseconds with 0 bytes received
* Closing connection
curl: (28) Operation timed out after 10010 milliseconds with 0 bytes received

 test-shell  ~  curl -m 10 -vv http://color-server:8000/color
* Host color-server:8000 was resolved.
* IPv6: (none)
* IPv4: 10.96.153.133
*   Trying 10.96.153.133:8000...
* Connected to color-server (10.96.153.133) port 8000
> GET /color HTTP/1.1
> Host: color-server:8000
> User-Agent: curl/8.7.1
> Accept: */*
> 
* Request completely sent off
< HTTP/1.1 200 OK
< content-type: text/plain
< date: Tue, 01 Oct 2024 19:25:05 GMT
< content-length: 4
< x-envoy-upstream-service-time: 1
< server: envoy
< 
* Connection #0 to host color-server left intact
blue# 
```

Note that localhost requests work fine while mirrord is running.
```shell
❯ curl localhost:8000/color
green
```

Stop the mirrord session locally and exit the 'test-shell'. 

Then, disable istio-injection and cycle the pods. Create a new 'test-shell' for testing.
```shell
$ kubectl --context kind-mirrord-stealing-test label ns color-server istio-injection=disabled --overwrite
$ kubectl --context kind-mirrord-stealing-test delete po -n color-server color-server-77ff58464f-xd2tg test-shell

$ kubectl --context kind-mirrord-stealing-test -n color-server run test-shell-again --rm -i --tty --image nicolaka/netshoot
```

Start up mirrord, and traffic stealing will work now that istio is gone. 
```shell
❯ ~/.slice/bin/mirrord-3.118.0 exec --context kind-mirrord-stealing-test -n color-server -e --steal --target pod/color-server-77ff58464f-rgnpv -- ./color-server -color=green
* Running binary "./color-server" with arguments: ["-color=green"].
* mirrord will target: pod/color-server-77ff58464f-rgnpv, no configuration file was loaded
* operator: the operator will be used if possible
* env: all environment variables will be fetched
* fs: file operations will default to read only from the remote
* incoming: incoming traffic will be stolen
* outgoing: forwarding is enabled on TCP and UDP
* dns: DNS will be resolved remotely
⠐ mirrord exec
    ✓ running on latest (3.118.0)!
    ✓ ready to launch process
      ✓ layer extracted
      ✓ operator not found
      ✓ container created
      ✓ container is ready
    ✓ config summary                                                                2024/10/01 15:46:43 INFO Starting server on :8000
2024/10/01 15:46:47 INFO HTTP Request method=GET path=/color protocol=HTTP/1.1 remote_addr=10.244.0.6:37620


 test-shell-again  ~  curl -m 10 -vv http://color-server:8000/color
* Host color-server:8000 was resolved.
* IPv6: (none)
* IPv4: 10.96.153.133
*   Trying 10.96.153.133:8000...
* Connected to color-server (10.96.153.133) port 8000
> GET /color HTTP/1.1
> Host: color-server:8000
> User-Agent: curl/8.7.1
> Accept: */*
> 
* Request completely sent off
< HTTP/1.1 200 OK
< Content-Type: text/plain
< Date: Tue, 01 Oct 2024 19:46:52 GMT
< Content-Length: 5
< 
* Connection #0 to host color-server left intact
green#
```