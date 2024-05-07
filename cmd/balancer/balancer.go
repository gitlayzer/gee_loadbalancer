package balancer

import (
	"errors"
)

var (
	NodeHostErr                = errors.New("no host")
	AlgorithmNotSupportedError = errors.New("algorithm not supported")
)

// Balancer 是一个负载均衡器的接口
type Balancer interface {
	// Add 添加一个节点
	Add(string)
	// Remove 移除一个节点
	Remove(string)
	// Balance 负载均衡
	Balance(string) (string, error)
	// Inc 增加一个节点的权重
	Inc(string)
	// Done 一个节点完成了一次请求
	Done(string)
}

// Factory 创建一个负载均衡器的工厂函数
type Factory func([]string) Balancer

// factories 负载均衡器工厂集合
var factories = make(map[string]Factory)

// Build 根据算法和节点列表创建负载均衡器
func Build(algorithm string, nodes []string) (Balancer, error) {
	factory, ok := factories[algorithm]
	if !ok {
		return nil, AlgorithmNotSupportedError
	}
	return factory(nodes), nil
}
