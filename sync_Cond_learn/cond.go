package main

import (
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type User struct {
	RegTime time.Time //注册时间
	Score   int       //积分
}

var (
	users        []*User    //用户集合
	mu           sync.Mutex //读写users之前先加锁，保证它的并发安全性
	listenNumber int32      //监听全局变量users的次数
	prized       bool       //是否已经执行过积分奖励
)

const (
	H = 3
	M = 8
	N = 5
)

func InitGlobalVar() {
	users = make([]*User, 0, 500)
	mu = sync.Mutex{}
	listenNumber = 0
	prized = false
}

// 给前H个用户积分奖励
func prize() {
	//按注册时间排序
	sort.Slice(users, func(i, j int) bool { return users[i].RegTime.Before(users[j].RegTime) })
	//把前H个用户的Score加1
	for _, user := range users[:H] {
		user.Score += 1
	}
}

// 业务模型介绍。
// 上游：M个协程接收用户注册，把用户添加到集合users中。
// 下游：N个协程监听users，当发现它里面攒够100个用户时，给前3(常量H)个注册的用户发放积分，然后协程终止。（这部分逻辑上游不关心，相关代码不能放到上游协程里实现）
func BusinessModel() {
	InitGlobalVar()
	downstreamOver := false
	//上游
	for i := 0; i < M; i++ {
		go func() {
			for { //不停地注册新用户
				if downstreamOver {
					break
				}
				mu.Lock()
				users = append(users, &User{RegTime: time.Now()}) //注册用户
				mu.Unlock()
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100))) //随机休息一段时间，再注册下一个用户
			}
		}()
	}
	//下游
	wg := sync.WaitGroup{}
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for {
				mu.Lock()
				if !prized {
					atomic.AddInt32(&listenNumber, 1)
					if len(users) >= 100 {
						prize()
						prized = true
					}
				}
				mu.Unlock()
				if prized {
					break
				}
			}
		}()
	}
	wg.Wait()
	downstreamOver = true
}

// 减少对全局变量users的监听次数。上游每次改变users时向一个channel里发送一条数据
func SignalWithChannel() {
	InitGlobalVar()
	ch := make(chan struct{}, 10*N)
	downstreamOver := false
	//上游
	for i := 0; i < M; i++ {
		go func() {
			for { //不停地注册新用户
				if downstreamOver {
					break
				}
				mu.Lock()
				users = append(users, &User{RegTime: time.Now()}) //注册用户
				mu.Unlock()
				ch <- struct{}{}
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100))) //随机休息一段时间，再注册下一个用户
			}
		}()
	}
	//下游
	wg := sync.WaitGroup{}
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for {
				<-ch //阻塞，直到users有改变
				mu.Lock()
				if !prized {
					atomic.AddInt32(&listenNumber, 1)
					if len(users) >= 100 {
						prize()
						prized = true
					}
				}
				mu.Unlock()
				if prized {
					break
				}
			}
		}()
	}
	wg.Wait()
	downstreamOver = true
}

func SignalWithCond() {
	InitGlobalVar()
	cond := sync.NewCond(&mu) //cond.L等价于mu
	downstreamOver := false
	//上游
	for i := 0; i < M; i++ {
		go func() {
			for { //不停地注册新用户
				if downstreamOver {
					break
				}
				mu.Lock()
				users = append(users, &User{RegTime: time.Now()}) //注册用户
				mu.Unlock()
				cond.Signal()                                                //通知别人users有变化。Signal只能通知到一个协程
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100))) //随机休息一段时间，再注册下一个用户
			}
		}()
	}
	//下游
	wg := sync.WaitGroup{}
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for {
				mu.Lock()   //等价于cond.L.Lock()
				cond.Wait() //阻塞，直到接收到通知。Wait内部会先执行mu.Unlock()，等接收到信号后再执行mu.Lock()，所以在调Wait()之前需要先上锁
				if !prized {
					atomic.AddInt32(&listenNumber, 1)
					if len(users) >= 100 {
						prize()
						prized = true
					}
				}
				mu.Unlock() //等价于cond.L.Unlock()
				if prized {
					break
				}
			}
		}()
	}
	wg.Wait()
	downstreamOver = true
}

func BroadcastWithChannel() {
	InitGlobalVar()
	ch := make(chan struct{}, 10*N)
	downstreamOver := false
	//上游
	for i := 0; i < M; i++ {
		go func() {
			for { //不停地注册新用户
				if downstreamOver {
					break
				}
				mu.Lock()
				users = append(users, &User{RegTime: time.Now()}) //注册用户
				mu.Unlock()
				//把n个下游协程全部通知一遍。close channel也能实现通知的功能，但是一个channl只能close一次，本业务中我们需要多次通知。实际中上游一般不知道下游协程的数目，这种情况下只能用cond.Broadcast()
				for j := 0; j < N; j++ {
					ch <- struct{}{}
				}
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100))) //随机休息一段时间，再注册下一个用户
			}
		}()
	}
	//下游
	wg := sync.WaitGroup{}
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for {
				<-ch //阻塞，直到users有改变
				atomic.AddInt32(&listenNumber, 1)
				mu.Lock()
				done := false
				if len(users) >= 100 {
					prize()
					done = true
				}
				mu.Unlock()
				if done {
					break
				}
			}
		}()
	}
	wg.Wait()
	downstreamOver = true
}

func BroadcastWithCond() {
	InitGlobalVar()
	cond := sync.NewCond(&mu) //cond.L等价于mu
	downstreamOver := false
	//上游
	for i := 0; i < M; i++ {
		go func() {
			for { //不停地注册新用户
				if downstreamOver {
					break
				}
				mu.Lock()
				users = append(users, &User{RegTime: time.Now()}) //注册用户
				mu.Unlock()
				cond.Broadcast()                                             //通知所有下游协程
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100))) //随机休息一段时间，再注册下一个用户
			}
		}()
	}
	//下游
	wg := sync.WaitGroup{}
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for {
				mu.Lock()
				cond.Wait()
				atomic.AddInt32(&listenNumber, 1)
				done := false
				if len(users) >= 100 {
					prize()
					done = true
				}
				mu.Unlock()
				if done {
					break
				}
			}
		}()
	}
	wg.Wait()
	downstreamOver = true
}

//作者：高性能golang https://www.bilibili.com/read/cv25426306/?spm_id_from=333.999.0.0 出处：bilibili
