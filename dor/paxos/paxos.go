package paxos



type AcceptorStatus string

const (
	PROMISED      AcceptorStatus = "PROMISED"
	ACCEPTED      AcceptorStatus = "ACCEPTED"
	NONE      AcceptorStatus = "NONE"
)

func (p AcceptorStatus) String() string {
	switch (p) {
	case PROMISED: return "PROMISED"
	case ACCEPTED: return "ACCEPTED"
	case NONE: return "NONE"
	default:         return "UNKNOWN"
	}
}

type Proposal struct {
	//	提案名称
	 Name  string
	//	提案的序列号
	 Id int
	//	提案的值
	 Value string
}
type CommitResult struct {
	//	提案名称
	Accepted  bool
	//	提案的序列号
	 AcceptorStatus AcceptorStatus
	//	提案的值
	Proposal *Proposal
}

type PrepareResult struct {
	//	提案名称
	IsPromised  bool
	//	提案的序列号
	AcceptorStatus AcceptorStatus
	//	提案的值
	Proposal *Proposal
}
func NewCommitResult () *CommitResult {
	return &CommitResult{false,NONE,nil,}
}

func NewPrepareResult () *PrepareResult {

	return &PrepareResult{false,NONE,nil}
}

func (this *Proposal)CopyFromInstance( proposal *Proposal){
this.Id = proposal.Id;
this.Name = proposal.Name;
this.Value = proposal.Value;
}