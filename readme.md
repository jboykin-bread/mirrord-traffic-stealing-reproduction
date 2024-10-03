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

 test-shell  ~  grpcurl -vv -max-time 5 -plaintext color-server:9000 ColorService.GetColor

Resolved method descriptor:
rpc GetColor ( .ColorRequest ) returns ( .ColorResponse );

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc
date: Wed, 02 Oct 2024 16:42:54 GMT
server: envoy
x-envoy-upstream-service-time: 0

Estimated response size: 6 bytes

Response contents:
{
  "color": "blue"
}

Response trailers received:
(empty)
Sent 0 requests and received 1 response
Timing Data: 7.658416ms
  Dial: 2.052083ms
    BlockingDial: 2.043583ms
  InvokeRPC: 4.645458ms
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

🐛 All grpc requests to color-server:9000 begin to fail

```shell
 test-shell  ~  grpcurl -vv -max-time 5 -plaintext color-server:9000 ColorService.GetColor
Error invoking method "ColorService.GetColor": rpc error: code = DeadlineExceeded desc = failed to query for service descriptor "ColorService": context deadline exceeded
```

Note that localhost requests work fine while mirrord is running.
```shell
❯ grpcurl -plaintext localhost:9000 ColorService.GetColor
{
  "color": "green"
}
```

