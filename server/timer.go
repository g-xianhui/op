package main

import (
	"sort"
	"sync"
	"time"
)

type DelayTask struct {
	guid     uint
	deadline time.Time
	interval time.Duration
	task     func()
}

type ByTime []DelayTask

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].deadline.Before(a[j].deadline) }

type TimeMgr struct {
	Ticker              *time.Ticker
	guid                uint
	looplock, delaylock sync.Mutex
	looplist, oncelist  []DelayTask
}

func (mgr *TimeMgr) Run() {
	mgr.Ticker = time.NewTicker(time.Second)
	for now := range mgr.Ticker.C {
		mgr.delaylock.Lock()
		j := 0
		for i := range mgr.oncelist {
			if now.Before(mgr.oncelist[i].deadline) {
				break
			}
			mgr.oncelist[i].task()
			j++
		}
		left := copy(mgr.oncelist, mgr.oncelist[j:])
		mgr.oncelist = mgr.oncelist[:left]
		mgr.delaylock.Unlock()

		mgr.looplock.Lock()
		for i, t := range mgr.looplist {
			if now.Before(t.deadline) {
				break
			}
			t.task()
			mgr.looplist[i].deadline = t.deadline.Add(t.interval)
		}
		sort.Sort(ByTime(mgr.looplist))
		mgr.looplock.Unlock()
	}
}

func (mgr *TimeMgr) AddLoop(beginTime time.Time, interval time.Duration, task func()) uint {
	now := time.Now()
	for now.After(beginTime) {
		beginTime = beginTime.Add(interval)
	}
	delaytask := DelayTask{0, beginTime, interval, task}

	mgr.looplock.Lock()
	mgr.guid++
	delaytask.guid = mgr.guid
	mgr.looplist = append(mgr.looplist, delaytask)

	index := len(mgr.looplist) - 1
	for ; index > 0; index-- {
		if !delaytask.deadline.Before(mgr.looplist[index-1].deadline) {
			break
		}
		mgr.looplist[index] = mgr.looplist[index-1]
	}
	mgr.looplist[index] = delaytask

	mgr.looplock.Unlock()
	return delaytask.guid
}

func (mgr *TimeMgr) DelLoop(guid uint) {
	for i, t := range mgr.looplist {
		if t.guid == guid {
			mgr.looplist = append(mgr.looplist[:i], mgr.looplist[i+1:]...)
			break
		}
	}
}

func (mgr *TimeMgr) AddCall(delay time.Duration, task func()) uint {
	delaytask := DelayTask{deadline: time.Now().Add(delay), task: task}

	mgr.delaylock.Lock()
	mgr.guid++
	delaytask.guid = mgr.guid
	mgr.oncelist = append(mgr.oncelist, delaytask)
	index := len(mgr.oncelist) - 1
	for ; index > 0; index-- {
		if !delaytask.deadline.Before(mgr.oncelist[index-1].deadline) {
			break
		}
		mgr.oncelist[index] = mgr.oncelist[index-1]
	}
	mgr.oncelist[index] = delaytask
	mgr.delaylock.Unlock()
	return delaytask.guid
}

func (mgr *TimeMgr) DelCall(guid uint) {
	for i, t := range mgr.oncelist {
		if t.guid == guid {
			mgr.oncelist = append(mgr.oncelist[:i], mgr.oncelist[i+1:]...)
			break
		}
	}
}
