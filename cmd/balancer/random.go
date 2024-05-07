package balancer

import (
	"math/rand"
	"time"
)

func init() {
	factories[RandomBalancer] = NewRandom
}

type Random struct {
	BaseBalancer
	rnd *rand.Rand
}

// NewRandom 创建一个随机负载均衡器
func NewRandom(nodes []string) Balancer {
	return &Random{
		BaseBalancer: BaseBalancer{
			nodes: nodes,
		},
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Balance 实现负载均衡接口
func (r *Random) Balance(_ string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.nodes) == 0 {
		return "", NodeHostErr
	}
	return r.nodes[r.rnd.Intn(len(r.nodes))], nil
}

// Add 添加节点
func (r *Random) Add(node string) {}

// Remove 移除节点
func (r *Random) Remove(node string) {}
