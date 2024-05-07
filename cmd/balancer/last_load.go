package balancer

import (
	fibHeap "github.com/starwander/GoFibonacciHeap"
	"sync"
)

func init() {
	factories[LeastLoadBalancer] = NewLastLoad
}

func (h *host) Tag() interface{} { return h.name }

func (h *host) Key() float64 { return float64(h.load) }

type LeastLoad struct {
	sync.RWMutex
	heap *fibHeap.FibHeap
}

// NewLastLoad 创建一个最小负载均衡器
func NewLastLoad(hosts []string) Balancer {
	ll := &LeastLoad{heap: fibHeap.NewFibHeap()}
	for _, h := range hosts {
		ll.Add(h)
	}
	return ll
}

// Add 添加节点
func (l *LeastLoad) Add(node string) {
	l.Lock()
	defer l.Unlock()
	if ok := l.heap.GetValue(node); ok != nil {
		return
	}
	_ = l.heap.InsertValue(&host{node, 0})
}

// Remove 删除节点
func (l *LeastLoad) Remove(node string) {
	l.Lock()
	defer l.Unlock()
	if ok := l.heap.GetValue(node); ok == nil {
		return
	}
	_ = l.heap.Delete(node)
}

// Balance 是最小负载均衡算法的实现
func (l *LeastLoad) Balance(_ string) (string, error) {
	l.RLock()
	defer l.RUnlock()
	if l.heap.Num() == 0 {
		return "", NodeHostErr
	}
	return l.heap.MinimumValue().Tag().(string), nil
}

// Inc 增加节点负载
func (l *LeastLoad) Inc(hostName string) {
	l.Lock()
	defer l.Unlock()
	if ok := l.heap.GetValue(hostName); ok == nil {
		return
	}
	h := l.heap.GetValue(hostName)
	h.(*host).load++
	_ = l.heap.IncreaseKeyValue(h)
}

// Done 减少节点负载
func (l *LeastLoad) Done(hostName string) {
	l.Lock()
	defer l.Unlock()
	if ok := l.heap.GetValue(hostName); ok == nil {
		return
	}
	h := l.heap.GetValue(hostName)
	h.(*host).load--
	_ = l.heap.DecreaseKeyValue(h)
}
