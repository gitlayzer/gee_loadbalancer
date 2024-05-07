package balancer

import "github.com/lafikl/consistent"

func init() {
	factories[BoundedBalancer] = NewBounded
}

type Bounded struct {
	ch *consistent.Consistent
}

// NewBounded 用于创建一个新的 Bounded 负载均衡器
func NewBounded(nodes []string) Balancer {
	c := &Bounded{
		ch: consistent.New(),
	}

	// 循环添加节点
	for _, node := range nodes {
		c.ch.Add(node)
	}

	return c
}

// Add 用于向 Bounded 负载均衡器中添加一个节点
func (b *Bounded) Add(node string) {
	b.ch.Add(node)
}

// Remove 用于从 Bounded 负载均衡器中移除一个节点
func (b *Bounded) Remove(node string) {
	b.ch.Remove(node)
}

// Balance 用于根据 key 选择一个节点
func (b *Bounded) Balance(key string) (string, error) {
	if len(b.ch.Hosts()) == 0 {
		return "", NodeHostErr
	}

	return b.ch.GetLeast(key)
}

// Inc 用于增加一个节点的权重
func (b *Bounded) Inc(node string) {
	b.ch.Inc(node)
}

// Done 用于完成一次请求
func (b *Bounded) Done(node string) {
	b.ch.Done(node)
}
