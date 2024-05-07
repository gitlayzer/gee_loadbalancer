package balancer

import "sync"

// BaseBalancer 是一个负载均衡器的基础实现
type BaseBalancer struct {
	sync.RWMutex
	nodes []string
}

// AddNode 向负载均衡器中添加一个节点
func (b *BaseBalancer) AddNode(node string) {
	b.Lock()
	defer b.Unlock()
	for _, n := range b.nodes {
		if n == node {
			return
		}
	}
	b.nodes = append(b.nodes, node)
}

// RemoveNode 从负载均衡器中移除一个节点
func (b *BaseBalancer) removeNode(node string) {
	b.Lock()
	defer b.Unlock()
	for i, n := range b.nodes {
		if n == node {
			b.nodes = append(b.nodes[:i], b.nodes[i+1:]...)
			return
		}
	}
}

// Balance 负载均衡算法的实现
func (b *BaseBalancer) Balance(key string) (string, error) {
	return "", nil
}

// Inc 增加一个节点的权重
func (b *BaseBalancer) Inc(_ string) {}

// Done 一个节点完成了一次请求
func (b *BaseBalancer) Done(_ string) {}
