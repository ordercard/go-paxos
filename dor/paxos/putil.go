package paxos

import (
	"math/rand"
	"time"
)

/*
 * @brief 提案序列号生成：保证唯一且递增。参考chubby中提议生成算法
 * @param myID 提议者的ID
 * @param numCycle 生成提议的轮次
 * @param 提议者个数
 * @return 生成的提案id
 *
 */
func GenerateId( myID int,  numCycle int,  proposerCount int ) int{
 id := numCycle * proposerCount + myID;
return id;
}


//	随机休眠，模拟网络延迟
func SleepRandom() int32{
	r := rand.New(rand.NewSource(time.Now().Unix()))
	timeInMs := r.Int31n(10) + 10
	time.Sleep(time.Duration(timeInMs) * time.Millisecond)
	return timeInMs;

}

func IsCrashed() bool{
r := rand.New(rand.NewSource(time.Now().Unix()))
x := r.Int31n(4)
if	(0==x){
	return true;
}
return false;
}
