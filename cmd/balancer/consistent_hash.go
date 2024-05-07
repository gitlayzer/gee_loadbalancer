package balancer

import "github.com/lafikl/consistent"

func init() {
	factories[ConsistentHashBalancer] = NewConsistentHash
}

type Consistent struct {
	BaseBalancer
	ch *consistent.Consistent
}

// NewConsistentHash 创建一个一致性哈希的负载均衡器
func NewConsistentHash(nodes []string) Balancer {
	c := &Consistent{
		ch: consistent.New(),
	}

	// 循环添加节点
	for _, n := range nodes {
		c.ch.Add(n)
	}

	return c
}

// Add 用于向一致性哈希负载均衡器中添加一个节点
func (c *Consistent) Add(node string) {
	c.ch.Add(node)
}

// Remove 用于从一致性哈希负载均衡器中移除一个节点
func (c *Consistent) Remove(node string) {
	c.ch.Remove(node)
}

// Balance 用于从一致性哈希负载均衡器中选择一个节点
func (c *Consistent) Balance(key string) (string, error) {
	if len(c.ch.Hosts()) == 0 {
		return "", NodeHostErr
	}
	return c.ch.Get(key)
}
