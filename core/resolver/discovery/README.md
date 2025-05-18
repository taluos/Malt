# HowTo

## 1, grpc服务发现和负载均衡原理

 1）grpc客户端启动后，首先通过 Resolver（解析器） 获取服务的地址信息（比如通过 DNS、etcd、Consul 等）。

 2）Resolver 将服务地址（如 IP:Port）传递给 Balancer（负载均衡器）。

 3）Balancer 接收到地址信息后，建立多个 conn（连接），这些连接指向多个 grpc 服务端 实例。同时 Balancer 会维护一个 Address Cache（地址缓存），存储当前可用的服务端地址。

 4）根据不同的策略（轮询、最少连接、权重等），Balancer 更新一个 picker（选择器），这个 picker 决定每个请求应该发往哪一个 conn。

 5）当客户端发起一个请求时，Balancer 会根据 picker 的策略，选择一个可用的 conn 来发送请求。请求会被写入到 client stream（客户端流），最终发送给目标 grpc 服务端（帧传输）。

 6）当 grpc 服务端返回响应时，Balancer 会根据响应的状态来更新地址缓存。如果某个服务端实例不可用，Balancer 会将其从地址缓存中移除，从而避免将请求发送到不可用的实例。

## 2, 在 Malt 中的实现

discovery中resolver.go 和 builder.go 为抽象化的接口，为了后续rpc服务实现时对Consul、etcd等服务的解耦。rpc服务在服务注册和发现时，只需要实现相关接口即可。

在builder.go中NewBuilder()方法中，会根据用户传入的discovery(registry.Discovery)来创建对应的resolver。 这里的discovery是在Malt/registry/registry.go中创建的。

在client.go中NewClient()方法中，使用 grpc.WithResolvers 将自定义的 resolver 注册到 gRPC 客户端中。(Malt/client/rpcClient/client.go 68-72)

如果用户没有传入discovery，那么会默认使用directDiscovery。(Malt/client/rpcClient/client.go 73-76)
