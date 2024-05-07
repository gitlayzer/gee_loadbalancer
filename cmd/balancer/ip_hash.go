package balancer

import "hash/crc32"

func init() {
	factories[IPHashBalancer] = NewIPHash
}

// IPHash will choose a host based on the client's IP address
type IPHash struct {
	BaseBalancer
}

// NewIPHash create new IPHash balancer
func NewIPHash(nodes []string) Balancer {
	return &IPHash{
		BaseBalancer: BaseBalancer{
			nodes: nodes,
		},
	}
}

// Balance 是 IPHash 负载均衡算法的实现
func (r *IPHash) Balance(key string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.nodes) == 0 {
		return "", NodeHostErr
	}
	value := crc32.ChecksumIEEE([]byte(key)) % uint32(len(r.nodes))
	return r.nodes[value], nil
}

// Add 添加节点
func (r *IPHash) Add(node string) {}

// Remove 删除节点
func (r *IPHash) Remove(node string) {}
