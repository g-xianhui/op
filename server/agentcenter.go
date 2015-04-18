package main

import (
	"sync"
)

type AgentCenter struct {
	lock     sync.RWMutex
	accounts map[string]uint32
	agents   map[uint32]*Agent
}

func (c *AgentCenter) init() {
	c.accounts = make(map[string]uint32)
	c.agents = make(map[uint32]*Agent)
}

func (c *AgentCenter) add(accountName string, roleid uint32, agent *Agent) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.accounts[accountName] = roleid
	c.agents[roleid] = agent
}

func (c *AgentCenter) find(roleid uint32) *Agent {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.agents[roleid]
}

func (c *AgentCenter) findByAccount(accountName string) *Agent {
	c.lock.RLock()
	defer c.lock.RUnlock()
	roleid, found := c.accounts[accountName]
	if !found {
		return nil
	}
	return c.agents[roleid]
}
