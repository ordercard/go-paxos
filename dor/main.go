package main

import (
	"fmt"
	"go-paxos/dor/obc"
	"strconv"
	"sync"
)
import "go-paxos/dor/paxos"

func foo(p paxos.AcceptorStatus) {
	fmt.Printf("enum value: %v\n", p)
}

func main() {
	foo(paxos.PROMISED)
	var wg =sync.WaitGroup{}
	var acceptedProposals []*obc.Acceptor = []*obc.Acceptor{}
	for i:=0;i<15;i++{
		acceptedProposals=append(acceptedProposals,obc.NewAcceptor(i))
	}
	for i:=0;i<10;i++{
		wg.Add(1)
		go func(i int) {

			pp :=obc.NewProposer(i,strconv.Itoa(i)+"~Proposer",10,acceptedProposals)
			pp.Prepare()
			pp.Commit()
			wg.Done()
		}(i)
	}

	wg.Wait() //阻塞直到所有任务完成
	fmt.Printf("leader节点已经选出！")
}

