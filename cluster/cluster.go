package cluster

import (
	"github.com/hashicorp/memberlist"
	"io/ioutil"
	"log"
	"stathat.com/c/consistent"
	"time"
)

type Node interface {
	ShouldProcess(key string) (string, bool) //节点是否应该自己处理该请求
	Members() []string                       //提供整个集群的节点列表
	Addr() string                            //获取本节点的地址
}

type node struct {
	*consistent.Consistent
	addr string
}

func (n *node) Addr() string {
	return n.addr
}

func (n *node) ShouldProcess(key string) (string, bool) {
	addr, _ := n.Get(key)
	return addr, addr == n.addr
}

func New(addr, cluster string) (Node, error) {
	conf := memberlist.DefaultLocalConfig()
	conf.Name = addr
	conf.BindAddr = addr
	conf.LogOutput = ioutil.Discard //就先把日志丢掉吧
	list, err := memberlist.Create(conf)
	if err != nil {
		return nil, err
	}
	if cluster == "" {
		cluster = addr
	}
	clu := []string{cluster}
	log.Println(cluster, clu)
	numJoined, err := list.Join(clu)
	log.Printf("the number of hosts successfully contacted: %d.\n", numJoined)
	if err != nil {
		return nil, err
	}

	circle := consistent.New()
	circle.NumberOfReplicas = 256 //设置虚拟节点个数
	//每隔1秒将memberlist中集群节点列表更新到circle中
	go func() {
		for {
			m := list.Members()
			nodes := make([]string, len(m))
			for i, n := range m {
				nodes[i] = n.Name
			}
			circle.Set(nodes)
			time.Sleep(time.Second)
		}
	}()
	return &node{circle, addr}, nil
}
