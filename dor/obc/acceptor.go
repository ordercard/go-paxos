package obc

import (
	"go-paxos/dor/paxos"
	"sync"
)
/*
接受者的接口定义已经 处理第一第一阶段和第二阶段的逻辑
 */
type Acceptor struct {
	  Acceptor_id int
	// 记录已处理提案的状态
	  Status paxos.AcceptorStatus
	// 记录最新承诺的提案
	  PromisedProposal *paxos.Proposal
	// 记录最新批准的提案
	  AcceptedProposal *paxos.Proposal
	  *sync.RWMutex
}

func NewAcceptor(accid int) *Acceptor {
	return &Acceptor{accid,paxos.NONE,&paxos.Proposal{},&paxos.Proposal{},&sync.RWMutex{}}
}

func (acceptor *Acceptor)Prepare( szProposal *paxos.Proposal) *paxos.PrepareResult {
	acceptor.Lock()

	 var prepareResult  = paxos.NewPrepareResult()
	if (paxos.IsCrashed()) {
		acceptor.Unlock()
		return nil;
	}


	switch (acceptor.Status) {
	// NONE表示之前没有承诺过任何提议者
	// 此时，接受提案
	case paxos.NONE:
		prepareResult.AcceptorStatus=paxos.NONE
		prepareResult.IsPromised=true
		prepareResult.Proposal = nil
		// 转换自身的状态，已经承诺了提议者，并记录承诺的提案。
		acceptor.Status = paxos.PROMISED;
		acceptor.PromisedProposal.CopyFromInstance(szProposal);
		acceptor.Unlock()
		//System.out.println(this.getAcceptor_id()+"acceptor说:"+szProposal+"当前为NONE，承诺了提议者，并记录承诺的提案");
		return prepareResult;
	// 已经承诺过任意提议者
	case paxos.PROMISED:
		// 判断提案的先后顺序，只承诺相对较新的提案
		if (acceptor.PromisedProposal.Id > szProposal.Id) {
			prepareResult.AcceptorStatus=acceptor.Status
			prepareResult.IsPromised=false
			prepareResult.Proposal=acceptor.PromisedProposal
			acceptor.Unlock()
			//System.out.println(this.getAcceptor_id()+"acceptor说:"+szProposal+"太低的编号,当前为" +promisedProposal.getId()+ ",所以不接受承诺");
			return prepareResult
		} else {
			acceptor.PromisedProposal.CopyFromInstance(szProposal);
			prepareResult.AcceptorStatus=acceptor.Status
			prepareResult.IsPromised=true
			prepareResult.Proposal=acceptor.PromisedProposal
			acceptor.Unlock()
			//System.out.println(this.getAcceptor_id()+"acceptor说:"+"承诺了"+szProposal + "的提案");
			return prepareResult
		}
		// 已经批准过提案
	case paxos.ACCEPTED:
		// 如果是同一个提案，只是序列号增大
		// 批准提案，更新序列号。
		if (acceptor.PromisedProposal.Id  < szProposal.Id && acceptor.PromisedProposal.Value == szProposal.Value) {
		acceptor.PromisedProposal.Id=szProposal.Id
		prepareResult.AcceptorStatus=acceptor.Status
		prepareResult.IsPromised=true
		prepareResult.Proposal=acceptor.PromisedProposal
			acceptor.Unlock()
		//System.out.println(this.getAcceptor_id()+"acceptor说:"+"已经批准过"+szProposal+"更新序列号,返回:" + prepareResult);
		return prepareResult
	} else { // 否则，不予批准
		prepareResult.AcceptorStatus=acceptor.Status
		prepareResult.IsPromised=false
		prepareResult.Proposal=acceptor.AcceptedProposal
			acceptor.Unlock()
		//System.out.println(this.getAcceptor_id()+"acceptor说:"+"不予批准"+szProposal+"返回:" + prepareResult);
		return prepareResult
	}
	default:
		acceptor.Unlock()
		// return nil
	}

	return nil
}

// 加锁此提交函数，不允许同时访问，模拟单个决策者串行决策
 func (acceptor *Acceptor)OnCommit( szProposal *paxos.Proposal) *paxos.CommitResult {
	 acceptor.Lock()
var commitResult = paxos.NewCommitResult()
if (paxos.IsCrashed()) {
	acceptor.Unlock()
//System.out.println(this.getAcceptor_id()+"acceptor说:"+szProposal+"Acceptor Commit阶段 网络出现故障"+this);
return nil
}
switch (acceptor.Status) {
// 不可能存在此状态
case paxos.NONE:
	acceptor.Unlock()
return nil
// 已经承诺过提案
case paxos.PROMISED:
// 判断commit提案和承诺提案的序列号大小
// 大于，接受提案。
//System.out.println(this+"acceptor说:"+szProposal.getId() + "与PROMISED阶段的编号比较" + promisedProposal.getId());
if (szProposal.Id >= acceptor.PromisedProposal.Id) {
//promisedProposal.copyFromInstance(szProposal)
	acceptor.AcceptedProposal.CopyFromInstance(szProposal)
	acceptor.Status = paxos.ACCEPTED
commitResult.Accepted=true
commitResult.AcceptorStatus=acceptor.Status
commitResult.Proposal=acceptor.PromisedProposal
	acceptor.Unlock()
//System.out.println(this.getAcceptor_id()+"acceptor说:"+szProposal+"号大 接受"+this);
return commitResult;

} else { // 小于，回绝提案
commitResult.Accepted=false
commitResult.AcceptorStatus=acceptor.Status
commitResult.Proposal=acceptor.PromisedProposal
	acceptor.Unlock()
//System.out.println(this.getAcceptor_id()+"acceptor说:"+szProposal+"号小 拒绝");
return commitResult
}
// 已接受过提案
case paxos.ACCEPTED:
// 同一提案，序列号较大，接受
//System.out.println(szProposal.getId() + "与ACCEPTED阶段的编号比较" + promisedProposal.getId());
if (szProposal.Id > acceptor.AcceptedProposal.Id && szProposal.Value == acceptor.AcceptedProposal.Value) {
	acceptor.AcceptedProposal.Id=szProposal.Id
	commitResult.Accepted=true
	commitResult.AcceptorStatus=acceptor.Status
	commitResult.Proposal=acceptor.AcceptedProposal
	acceptor.Unlock()
//System.out.println(this.getAcceptor_id()+"acceptor说:"+szProposal+"号大 接受 状态是ACCEPTED,只更新编号,不更新值 ");
return commitResult;
} else { // 否则，回绝提案
	commitResult.Accepted=true
	commitResult.AcceptorStatus=acceptor.Status
	commitResult.Proposal=acceptor.AcceptedProposal
	acceptor.Unlock()
//System.out.println(this.getAcceptor_id()+"acceptor说:"+szProposal+"号太小 拒绝");
return commitResult
}
}
return nil
}