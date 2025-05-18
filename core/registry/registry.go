package registry

import "context"

type ServiceInstance struct {
	// ID is the unique identifier for the service instance.
	ID string `json:"ID" mapstructure:"ID"`

	// Name is the name of the service instance.
	Name string `json:"Name" mapstructure:"Name"`

	// Version is the version of the service instance.
	Version string `json:"Version" mapstructure:"Version"`

	// Endpoint is the endpoint of the service instance.
	// It is the address of the service instance.
	// such as "http://127.0.0.1:8080" in http,  "grpc://127.0.0.1:8080" in grpc
	Endpoints []string `json:"Endpoint" mapstructure:"Endpoint"`

	// Metadata is the metadata of the service instance.
	Metadata map[string]string `json:"Metadata" mapstructure:"Metadata"`

	// Tags is the tags of the service instance.
	Tags []string `json:"Tags" mapstructure:"Tags"`
}

// Registrar 服务注册
type Registrar interface {

	// Register registers the service instance.
	Register(ctx context.Context, service *ServiceInstance) error

	// Deregister deregisters the service instance.
	Deregister(ctx context.Context, service *ServiceInstance) error
}

// Discovery 服务发现
type Discovery interface {

	// GetService return the service instances in memory according to the service name.
	GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)

	// Watch creates a watcher according to the service name.
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}

// Watcher 服务监听
type Watcher interface {
	// 获取服务实例
	// 1.第一次监听时，如果服务实例列表不为空，返回服务实例列表
	// 2.如果服务实例列表发生变化，返回变化后的服务实例列表
	// 3.如果上述两个都不满足，就阻塞到超时或者接收到cancel信号
	Next() ([]*ServiceInstance, error)

	// 停止监听
	Stop() error
}

// Heartbeat 心跳检测
type Heartbeat interface {
	// 心跳检测
	Heartbeat(ctx context.Context, service *ServiceInstance) error
}
