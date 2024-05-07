package balancer

import "sync/atomic"

func init() {
	factories[RoundRobinBalancer] = NewRoundRobin
}

type RoundRobin struct {
	BaseBalancer
	i atomic.Uint64
}

// NewRoundRobin 创建一个轮询负载均衡器
func NewRoundRobin(nodes []string) Balancer {
	return &RoundRobin{
		i: atomic.Uint64{},
		BaseBalancer: BaseBalancer{
			nodes: nodes,
		},
	}
}

// Balance 轮询负载均衡器
func (r *RoundRobin) Balance(_ string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.nodes) == 0 {
		return "", NodeHostErr
	}
	node := r.nodes[r.i.Add(1)%uint64(len(r.nodes))]
	return node, nil
}

// Add 添加节点
func (r *RoundRobin) Add(node string) {}

// Remove 移除节点
func (r *RoundRobin) Remove(node string) {}
