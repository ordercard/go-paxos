package obc

import (
	"fmt"
	"go-paxos/dor/paxos"
	"strconv"
	"strings"
	"sync"
)

type Proposer struct {
	// 序列号
	MyID int
	// 提议者的名字
	Name string
	// 提议者的提案
	Proposal *paxos.Proposal

	// 是否已经有提案获得大多数决策者确认
	Voted bool
	// 大多数决策者的最小个数
	HalfCount int

	ProposerCount int

	NumCycle  int
	Acceptors []*Acceptor
	*sync.RWMutex
}

func NewProposer(propId int, name string, proposerCount int, acceptors []*Acceptor) Proposer {
	p := Proposer{propId, name,
		&paxos.Proposal{strconv.Itoa(propId), paxos.GenerateId(propId, 0, proposerCount), strconv.Itoa(propId) + ":Leader"},
		false, len(acceptors)/2 + 1, proposerCount, 0, acceptors, &sync.RWMutex{}}
	p.NumCycle++
	return p
}

func (p *Proposer) Prepare() bool {
	//System.out.println(proposal.getId() + "提案准备阶段开始");
	var prepareResult *paxos.PrepareResult

	isContinue := true
	// 已获得承诺个数
	promisedCount := 0

	for isContinue {
		// 决策者已承诺的提案集合
		//var  promisedProposals[]paxos.Proposal = []paxos.Proposal{}
		// 决策者已接受(确定)的提案集合
		var acceptedProposals []*paxos.Proposal = []*paxos.Proposal{}
		// 承诺数量
		promisedCount = 0

		for _, acceptor := range p.Acceptors {
			//System.out.println(proposal.getId() + "的prepare阶段,当前acceptor为:" + acceptor.getAcceptor_id());
			// 发送准备提案给决策者
			prepareResult = acceptor.Prepare(p.Proposal)
			// 随机休眠一段时间，模拟网络延迟。
			paxos.SleepRandom()

			// 模拟网络异常
			if nil == prepareResult {
				//System.out.println(proposal.getId() + "提案者 prepare 阶段网络故障,acceptor:"+acceptor.getAcceptor_id());
				continue
			}

			// 获得承诺
			if prepareResult.IsPromised {
				//System.out.println(proposal.getId() + "在acceptor:"+acceptor.getAcceptor_id()+"acceptor获取一个承诺, 结果" + prepareResult);
				promisedCount++
			} else {

				// 决策者已经给了更高id题案的承诺
				//					if (prepareResult.getAcceptorStatus() == AcceptorStatus.PROMISED) {
				//						//System.out.println(proposal.getId() + "的提案承诺失败,已有更高级提案被承诺了");
				//						promisedProposals.add(prepareResult.getProposal());
				//					}
				// 决策者已经通过了一个题案
				if prepareResult.AcceptorStatus == paxos.ACCEPTED {
					//System.out.println(proposal.getId() + "的提案承诺失败,"+acceptor+"已通过了一个题案,无法更改" + prepareResult.getProposal());
					acceptedProposals = append(acceptedProposals, prepareResult.Proposal)
				}
			}
			//System.out.println(proposal.getId() + "的prepare阶段结束,当前acceptor为:" + acceptor.getAcceptor_id());
		} // end of for

		// 获得多数决策者的承诺
		// 可以进行第二阶段：题案提交
		if promisedCount >= p.HalfCount {
			//System.out.println(proposal.getId() + "的'承诺'阶段通过,数量:" + promisedCount + "-----最小提案通过数:" + halfCount);
			break
		}
		//	Proposal votedProposal = votedEnd(acceptedProposals);
		// 决策者已经半数通过题案 可删除
		//			if (votedProposal != null) {
		//				System.out.println(myID + " : 决策已经投票结束:" + votedProposal);
		//				return true;
		//			}

		maxIdAcceptedProposal := getMaxIdProposal(acceptedProposals)
		// 在已经被决策者通过题案中选择序列号最大的决策,作为自己的决策。
		if maxIdAcceptedProposal != nil {
			p.Proposal.Id = paxos.GenerateId(p.MyID, p.NumCycle, p.ProposerCount)
			p.Proposal.Value = maxIdAcceptedProposal.Value
			//System.out.println(proposal.getId() + "获取更大ID,并接受maxIdAcceptedProposal提案中的value" + maxIdAcceptedProposal);
		} else {
			//System.out.println(proposal.getId() + "当前maxIdAcceptedProposal中没有被'确认'的提案");
			p.Proposal.Id = paxos.GenerateId(p.MyID, p.NumCycle, p.ProposerCount)
		}

		p.NumCycle++
	}
	return false
}

// 获得大多数决策者承诺后，开始进行提案确认
func (p *Proposer) Commit() bool {
	//System.out.println(proposal.getId() + "提案确认阶段开始");
	isContinue := true

	// 已获得接受该提案的决策者个数
	acceptedCount := 0
	for isContinue {
		// 决策者已接受(确定)的提案集合
		var acceptedProposals []*paxos.Proposal = []*paxos.Proposal{}
		acceptedCount = 0
		// 分别给决策者发送提案
		for _, acceptor := range p.Acceptors {
			//System.out.println(proposal.getId() + "的commit阶段,当前acceptor:" + acceptor.getAcceptor_id());
			// 决策者返回的提案结果
			commitResult := acceptor.OnCommit(p.Proposal)
			// 模拟网络延迟
			paxos.SleepRandom()

			// 模拟网络异常
			if nil == commitResult {
				//System.out.println(proposal.getId() + "提案者 commit 阶段网络故障,acceptor:"+acceptor.getAcceptor_id());
				continue
			}

			// 题案被决策者接受。
			if commitResult.Accepted {
				//System.out.println(proposal.getId() + "提案被接受一次," + "接受者:" + acceptor.getAcceptor_id());
				acceptedCount++
			} else {
				//System.out.println(proposal.getId() + "提案未被接受," + "不接受者:" + acceptor.getAcceptor_id());
				// 将未接受的提案 放入集合中
				acceptedProposals = append(acceptedProposals, commitResult.Proposal)
			}
			//System.out.println(proposal.getId() + "的commit阶段结束,当前acceptor为:" + acceptor.getAcceptor_id());
		}
		// 题案被半数以上决策者接受，说明题案已经被选出来。
		if acceptedCount >= p.HalfCount {
			fmt.Printf(strconv.Itoa(p.MyID) + " : 题案已经投票选出:" + p.Proposal.Value + "\n")
			vv, _ := strconv.ParseInt(strings.Split(p.Proposal.Value, ":")[0], 10, 32)
			if int64(p.MyID) == vv {
				fmt.Printf(strconv.Itoa(p.MyID) + "保存状态为:leader" + "这是第多少次的投票" + strconv.Itoa(p.Proposal.Id) + "\n")
			} else {
				fmt.Printf(strconv.Itoa(p.MyID) + "保存状态为:follower" + "\n")
			}
			return true
		} else {
			maxIdAcceptedProposal := getMaxIdProposal(acceptedProposals)
			// 在已经被决策者通过题案中选择序列号最大的决策,重新生成递增id，改变自己的value为序列号最大的value。
			// 这是一种预测，预测此maxIdAccecptedProposal最有可能被超过半数的决策者接受。并没什么用
			if maxIdAcceptedProposal != nil {
				p.Proposal.Id = paxos.GenerateId(p.MyID, p.NumCycle, p.ProposerCount)
				//	proposal.setValue(maxIdAcceptedProposal.getValue());
				//System.out.println(proposal.getId() + "为更大ID,并接受maxIdAcceptedProposal提案中的value" + maxIdAcceptedProposal);
			} else {
				p.Proposal.Id = paxos.GenerateId(p.MyID, p.NumCycle, p.ProposerCount)
				//System.out.println(proposal.getId() + "提案号加大了,做重审准备" + proposal);
			}

			p.NumCycle++
			// 回退到决策准备阶段  delete if
			if p.Prepare() {
				return true
			}
		}

	}
	return true
}

func getMaxIdProposal(acceptedProposals []*paxos.Proposal) *paxos.Proposal {
	if len(acceptedProposals) > 0 {
		retProposal := acceptedProposals[0]
		for _, proposal := range acceptedProposals {
			//System.out.println("寻找最大提案号:" + proposal.getId());
			if proposal.Id > retProposal.Id {
				retProposal = proposal
			}
			//System.out.println(proposal + "接受名单中最大编号的提案" + retProposal);
		}

		return retProposal
	}
	return nil
}
