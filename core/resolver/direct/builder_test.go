package direct

import (
	"Malt/pkg/errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

// 创建 ClientConn 的 mock
type mockClientConn struct {
	mock.Mock
}

// ParseServiceConfig implements resolver.ClientConn.
func (m *mockClientConn) ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult {
	panic("unimplemented")
}

func (m *mockClientConn) UpdateState(state resolver.State) error {
	args := m.Called(state)
	return args.Error(0)
}

func (m *mockClientConn) ReportError(err error) {
	m.Called(err)
}

func (m *mockClientConn) NewAddress(addresses []resolver.Address) {
	m.Called(addresses)
}

func (m *mockClientConn) NewServiceConfig(serviceConfig string) {
	m.Called(serviceConfig)
}

// 测试 NewDirectBuilder 函数
func TestNewDirectBuilder(t *testing.T) {
	builder := NewDirectBuilder()
	assert.NotNil(t, builder, "Builder 不应该为 nil")
	assert.IsType(t, &directBuilder{}, builder, "Builder 应该是 *directBuilder 类型")
}

// 测试 Scheme 方法
func TestDirectBuilder_Scheme(t *testing.T) {
	builder := NewDirectBuilder()
	scheme := builder.Scheme()
	assert.Equal(t, "direct", scheme, "Scheme 应该是 'direct'")
}

// 测试 Build 方法 - 成功情况
func TestDirectBuilder_Build_Success(t *testing.T) {
	// 准备测试数据
	builder := NewDirectBuilder()
	mockCC := new(mockClientConn)

	// 创建 target
	u, _ := url.Parse("direct:///127.0.0.1:8080,127.0.0.1:8081")
	target := resolver.Target{URL: *u}

	// 设置 mock 期望
	expectedAddrs := []resolver.Address{
		{Addr: "127.0.0.1:8080"},
		{Addr: "127.0.0.1:8081"},
	}
	expectedState := resolver.State{Addresses: expectedAddrs}
	mockCC.On("UpdateState", expectedState).Return(nil)

	// 调用被测试的方法
	res, err := builder.Build(target, mockCC, resolver.BuildOptions{})

	// 验证结果
	assert.NoError(t, err, "Build 不应该返回错误")
	assert.NotNil(t, res, "Resolver 不应该为 nil")

	// 验证 mock 是否按预期被调用
	mockCC.AssertExpectations(t)
}

// 测试 Build 方法 - 无端点情况
func TestDirectBuilder_Build_NoEndpoints(t *testing.T) {
	// 准备测试数据
	builder := NewDirectBuilder()
	mockCC := new(mockClientConn)

	// 创建空路径的 target
	u, _ := url.Parse("direct:///")
	target := resolver.Target{URL: *u}

	// 调用被测试的方法
	res, err := builder.Build(target, mockCC, resolver.BuildOptions{})

	// 验证结果
	assert.Error(t, err, "当没有提供端点时，Build 应该返回错误")
	assert.Nil(t, res, "当发生错误时，Resolver 应该为 nil")
	assert.Contains(t, err.Error(), "no endpoints provided", "错误消息应该提及缺少端点")

	// 不需要验证 mock 调用，因为在错误情况下不应该调用 UpdateState
}

// 测试 Build 方法 - UpdateState 失败情况
func TestDirectBuilder_Build_UpdateStateFailed(t *testing.T) {
	// 准备测试数据
	builder := NewDirectBuilder()
	mockCC := new(mockClientConn)

	// 创建 target
	u, _ := url.Parse("direct:///127.0.0.1:8080")
	target := resolver.Target{URL: *u}

	// 设置 mock 期望 - UpdateState 返回错误
	expectedAddrs := []resolver.Address{
		{Addr: "127.0.0.1:8080"},
	}
	expectedState := resolver.State{Addresses: expectedAddrs}
	mockCC.On("UpdateState", expectedState).Return(errors.New("update state failed"))

	// 调用被测试的方法
	res, err := builder.Build(target, mockCC, resolver.BuildOptions{})

	// 验证结果
	assert.Error(t, err, "当 UpdateState 失败时，Build 应该返回错误")
	assert.Nil(t, res, "当发生错误时，Resolver 应该为 nil")
	assert.Contains(t, err.Error(), "update state failed", "错误消息应该包含原始错误")

	// 验证 mock 是否按预期被调用
	mockCC.AssertExpectations(t)
}

// 测试 Build 方法 - 空端点被过滤
func TestDirectBuilder_Build_EmptyEndpointsFiltered(t *testing.T) {
	// 准备测试数据
	builder := NewDirectBuilder()
	mockCC := new(mockClientConn)

	// 创建包含空端点的 target
	u, _ := url.Parse("direct:///127.0.0.1:8080,,127.0.0.1:8081")
	target := resolver.Target{URL: *u}

	// 设置 mock 期望 - 空端点应被过滤
	expectedAddrs := []resolver.Address{
		{Addr: "127.0.0.1:8080"},
		{Addr: "127.0.0.1:8081"},
	}
	expectedState := resolver.State{Addresses: expectedAddrs}
	mockCC.On("UpdateState", expectedState).Return(nil)

	// 调用被测试的方法
	res, err := builder.Build(target, mockCC, resolver.BuildOptions{})

	// 验证结果
	assert.NoError(t, err, "Build 不应该返回错误")
	assert.NotNil(t, res, "Resolver 不应该为 nil")

	// 验证 mock 是否按预期被调用
	mockCC.AssertExpectations(t)
}
