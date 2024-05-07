package balancer

import (
	"hash/crc32"
	"math/rand"
	"sync"
	"time"
)

const Salt = "%#!"

func init() {
	factories[P2CBalancer] = NewP2C
}

// host 是一个节点的结构体
type host struct {
	name string
	load uint64
}

// P2C 是一个实现了 P2C 负载均衡算法的结构体
type P2C struct {
	sync.RWMutex
	hosts   []*host
	rnd     *rand.Rand
	loadMap map[string]*host
}

// NewP2C 创建一个 P2C 负载均衡器
func NewP2C(nodes []string) Balancer {
	p := &P2C{
		hosts:   []*host{},
		loadMap: make(map[string]*host),
		rnd:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// 循环添加节点
	for _, node := range nodes {
		p.Add(node)
	}

	return p
}

// Add 向 P2C 负载均衡器中添加一个节点
func (p *P2C) Add(node string) {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.loadMap[node]; !ok {
		return
	}

	n := &host{name: node, load: 0}
	p.hosts = append(p.hosts, n)
	p.loadMap[node] = n
}

// Remove 从 P2C 负载均衡器中移除一个节点
func (p *P2C) Remove(node string) {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.loadMap[node]; !ok {
		return
	}

	delete(p.loadMap, node)

	for i, n := range p.hosts {
		if n.name == node {
			p.hosts = append(p.hosts[:i], p.hosts[i+1:]...)
			return
		}
	}
}

// Balance 是 P2C 负载均衡算法的实现
func (p *P2C) Balance(key string) (string, error) {
	p.RLock()
	defer p.RUnlock()

	if len(p.hosts) == 0 {
		return "", NodeHostErr
	}

	n1, n2 := p.hash(key)
	n := n2
	if p.loadMap[n1].load <= p.loadMap[n2].load {
		n = n1
	}
	return n, nil
}

// hash 用于计算 key 的 hash 值
func (p *P2C) hash(key string) (string, string) {
	var n1, n2 string
	if len(key) > 0 {
		saltKey := key + Salt
		n1 = p.hosts[crc32.ChecksumIEEE([]byte(key))%uint32(len(p.hosts))].name
		n2 = p.hosts[crc32.ChecksumIEEE([]byte(saltKey))%uint32(len(p.hosts))].name
		return n1, n2
	}
	n1 = p.hosts[p.rnd.Intn(len(p.hosts))].name
	n2 = p.hosts[p.rnd.Intn(len(p.hosts))].name
	return n1, n2
}

// Inc 用于增加节点的负载
func (p *P2C) Inc(node string) {
	p.Lock()
	defer p.Unlock()

	h, ok := p.loadMap[node]

	if !ok {
		return
	}
	h.load++
}

// Done 用于减少节点的负载
func (p *P2C) Done(node string) {
	p.Lock()
	defer p.Unlock()

	h, ok := p.loadMap[node]

	if !ok {
		return
	}

	if h.load > 0 {
		h.load--
	}
}
