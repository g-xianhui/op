package main

import (
	"sync"
)

type AgentCenter struct {
	lock     sync.RWMutex
	accounts map[string]*Agent
	agents   map[uint32]*Agent

	namelock sync.Mutex
	// names booked but not insert into db yet will got 0 as roleid
	name2roleid map[string]uint32
}

func (c *AgentCenter) init() {
	c.accounts = make(map[string]*Agent)
	c.agents = make(map[uint32]*Agent)
	c.name2roleid = make(map[string]uint32)
}

func (c *AgentCenter) addByAccount(accountName string, agent *Agent) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.accounts[accountName] = agent
}

func (c *AgentCenter) find(roleid uint32) *Agent {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.agents[roleid]
}

func (c *AgentCenter) findByAccount(accountName string) *Agent {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.accounts[accountName]
}

func (c *AgentCenter) bookName(name string) bool {
	c.namelock.Lock()
	defer c.namelock.Unlock()
	_, found := c.name2roleid[name]
	if !found {
		c.name2roleid[name] = 0
		return true
	}
	return false
}

func (c *AgentCenter) unbookName(name string) {
	c.namelock.Lock()
	defer c.namelock.Unlock()
	delete(c.name2roleid, name)
}

func (c *AgentCenter) confirmName(name string, roleid uint32) {
	c.namelock.Lock()
	defer c.namelock.Unlock()
	c.name2roleid[name] = roleid
}
