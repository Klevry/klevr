package manager

import "time"

type AgentStorage struct {
	agentCache ICache
}

func NewAgentStorage() *AgentStorage {
	manager := ctx.Get(CtxServer).(*KlevrManager)

	s := &AgentStorage{}
	if manager.Config.DB.Cache == true {
		s.agentCache = NewCache(manager.Config.Cache.Address,
			manager.Config.Cache.Port,
			manager.Config.Cache.Password)
	}

	return s
}

func (a *AgentStorage) AddAgent(tx *Tx, agent *Agents) error {
	err := tx.addAgent(agent)

	if err == nil && a.agentCache != nil {
		_, agents := tx.getAgentsByGroupId(agent.GroupId)
		a.agentCache.SyncAgent(agent.GroupId, agents)
	}

	if err != nil {
		return err
	}

	return nil
}

func (a *AgentStorage) UpdateAgent(tx *Tx, agent *Agents) {
	tx.updateAgent(agent)

	if a.agentCache != nil {
		_, agents := tx.getAgentsByGroupId(agent.GroupId)
		a.agentCache.SyncAgent(agent.GroupId, agents)
	}
}

func (a *AgentStorage) UpdateZoneStatus(tx *Tx, zoneID uint64, arrAgent []Agents) {
	tx.updateZoneStatus(arrAgent)

	if a.agentCache != nil {
		a.agentCache.UpdateZoneStatus(zoneID, arrAgent)
	}
}

func (a *AgentStorage) DeleteAgent(tx *Tx, zoneID uint64) {
	tx.deleteAgent(zoneID)

	if a.agentCache != nil {
		a.agentCache.DeleteAllAgent(zoneID)
	}
}

func (a *AgentStorage) UpdateAccessAgent(tx *Tx, zoneID uint64, agentKey string, accessTime time.Time) int64 {
	if a.agentCache != nil {
		a.agentCache.UpdateAccessAgent(zoneID, agentKey, accessTime)
	}

	cnt := tx.updateAccessAgent(agentKey, accessTime)

	return cnt
}

func (a *AgentStorage) UpdateAgentStatus(tx *Tx, inactiveAgents *[]Agents, ids []uint64) {
	if a.agentCache != nil {
		a.agentCache.UpdateAgentDisabledStatus(inactiveAgents)
	}

	tx.updateAgentStatus(ids)
}

func (a *AgentStorage) GetAgentsForInactive(tx *Tx, before time.Time) (int64, *[]Agents) {
	if a.agentCache != nil {
		cnt, agents := a.agentCache.GetAgentsForInactive(before)
		if cnt > 0 {
			return cnt, agents
		}
	}

	cnt, agents := tx.getAgentsForInactive(before)

	return cnt, agents
}

func (a *AgentStorage) GetAgentsByZoneID(tx *Tx, zoneID uint64) (int64, *[]Agents) {
	if a.agentCache != nil {
		cnt, agents := a.agentCache.GetAgentsByZoneID(zoneID)
		if cnt > 0 {
			return cnt, agents
		}
	}

	cnt, agents := tx.getAgentsByGroupId(zoneID)

	return cnt, agents
}

func (a *AgentStorage) GetAgentByID(tx *Tx, zoneID uint64, agentID uint64) *Agents {
	if a.agentCache != nil {
		agent := a.agentCache.GetAgentByID(zoneID, agentID)
		if agent != nil {
			return agent
		}
	}

	agent := tx.getAgentByID(agentID)

	return agent
}

func (a *AgentStorage) GetAgentByAgentKey(tx *Tx, agentKey string, zoneID uint64) *Agents {
	if a.agentCache != nil {
		agent := a.agentCache.GetAgentByAgentKey(zoneID, agentKey)
		if agent != nil {
			return agent
		}
	}

	agent := tx.getAgentByAgentKey(agentKey, zoneID)

	return agent
}
