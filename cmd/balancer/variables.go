package balancer

var (
	P2CBalancer            = "p2c"             // p2c 的全称是 power of two choices 负载均衡算法,原理是从两个节点中选择负载较小的节点作为服务端
	IPHashBalancer         = "ip_hash"         // ip_hash 的全称是 ip hash 负载均衡算法,原理是根据客户端的 IP 地址进行哈希计算,将请求分配给固定的一台服务器
	ConsistentHashBalancer = "consistent_hash" // consistent_hash 的全称是 consistent hash 负载均衡算法,原理是根据请求的 key 进行哈希计算,将请求分配给哈希环上最近的一台服务器
	RandomBalancer         = "random"          // random 的全称是 random 负载均衡算法,原理是随机选择一台服务器作为服务端
	RoundRobinBalancer     = "round_robin"     // round_robin 的全称是 round robin 负载均衡算法,原理是按照节点列表的顺序依次选择一台服务器作为服务端
	LeastLoadBalancer      = "least_load"      // least_load 的全称是 least load 负载均衡算法,原理是选择负载最小的节点作为服务端
	BoundedBalancer        = "bounded"         // bounded 的全称是 bounded 负载均衡算法,原理是限制服务器的最大并发连接数,防止服务器过载
)
