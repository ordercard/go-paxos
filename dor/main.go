package main

import (
	"fmt"
	"go-paxos/dor/obc"
	"net/http"
	"strconv"
	"sync"
)
import "go-paxos/dor/paxos"

func foo(p paxos.AcceptorStatus) {
	fmt.Printf("enum value: %v\n", p)
}


/*
测试 go的函数集合-和函数转成接口，为函数添加方法和 函数类型转化
*/
type Greeting func(name string) http.HandlerFunc
func (g Greeting) ServeHTTP( x http.ResponseWriter, y *http.Request) {
	g("nimade")
}
func say(g Greeting, n string) {
	fmt.Println(g(n))
}

func english(name string) http.HandlerFunc {
	fmt.Printf("xxx")
	return func( x http.ResponseWriter, y *http.Request) {
		fmt.Printf("Hello, " + name)
	}
}

func main() {
	//var gg Greeting = func(name string) http.HandlerFunc {
	//	return func(writer http.ResponseWriter, request *http.Request) {
	//
	//	}
	//}


	//var g Greeting = english;

	//var ggg http.Handler = g
	//ggg.ServeHTTP(nil,nil)


	//g("xx")
	//var g http.Handler = Greeting(english)
	//g.ServeHTTP(nil,nil)





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

