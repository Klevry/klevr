package manager

import (
	"time"

	"github.com/Klevry/klevr/pkg/common"
)

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

func (a *AgentStorage) AddAgent(ctx *common.Context, tx *Tx, agent *Agents) error {
	err := tx.addAgent(agent)

	if err == nil && a.agentCache != nil {
		_, agents := tx.getAgentsByGroupId(agent.GroupId)
		a.agentCache.SyncAgent(ctx, agent.GroupId, agents)
	}

	if err != nil {
		return err
	}

	return nil
}

func (a *AgentStorage) UpdateAgent(ctx *common.Context, tx *Tx, agent *Agents) {
	tx.updateAgent(agent)

	if a.agentCache != nil {
		_, agents := tx.getAgentsByGroupId(agent.GroupId)
		a.agentCache.SyncAgent(ctx, agent.GroupId, agents)
	}
}

func (a *AgentStorage) UpdateZoneStatus(ctx *common.Context, tx *Tx, zoneID uint64, arrAgent []Agents) {
	if len(arrAgent) == 0 {
		return
	}

	tx.updateZoneStatus(arrAgent)

	if a.agentCache != nil {
		a.agentCache.UpdateZoneStatus(ctx, zoneID, arrAgent)
	}
}

func (a *AgentStorage) DeleteAgent(ctx *common.Context, tx *Tx, zoneID uint64) {
	tx.deleteAgent(zoneID)

	if a.agentCache != nil {
		a.agentCache.DeleteAllAgent(ctx, zoneID)
	}
}

func (a *AgentStorage) UpdateAccessAgent(ctx *common.Context, tx *Tx, zoneID uint64, agentKey string, accessTime time.Time) int64 {
	if a.agentCache != nil {
		a.agentCache.UpdateAccessAgent(ctx, zoneID, agentKey, accessTime)
	}

	cnt := tx.updateAccessAgent(agentKey, accessTime)

	return cnt
}

func (a *AgentStorage) UpdateAgentStatus(ctx *common.Context, tx *Tx, inactiveAgents *[]Agents, ids []uint64) {
	if a.agentCache != nil {
		a.agentCache.UpdateAgentDisabledStatus(ctx, inactiveAgents)
	}

	tx.updateAgentStatus(ids)
}

func (a *AgentStorage) GetAgentsForInactive(ctx *common.Context, tx *Tx, before time.Time) (int64, *[]Agents) {
	if a.agentCache != nil {
		cnt, agents := a.agentCache.GetAgentsForInactive(ctx, before)
		if cnt > 0 {
			return cnt, agents
		}
	}

	cnt, agents := tx.getAgentsForInactive(before)

	return cnt, agents
}

func (a *AgentStorage) GetAgentsByZoneID(ctx *common.Context, tx *Tx, zoneID uint64) (int64, *[]Agents) {
	if a.agentCache != nil {
		cnt, agents := a.agentCache.GetAgentsByZoneID(ctx, zoneID)
		if cnt > 0 {
			return cnt, agents
		}
	}

	cnt, agents := tx.getAgentsByGroupId(zoneID)

	return cnt, agents
}

func (a *AgentStorage) GetAgentByID(ctx *common.Context, tx *Tx, zoneID uint64, agentID uint64) *Agents {
	if a.agentCache != nil {
		agent := a.agentCache.GetAgentByID(ctx, zoneID, agentID)
		if agent != nil {
			return agent
		}
	}

	agent := tx.getAgentByID(agentID)

	return agent
}

func (a *AgentStorage) GetAgentByAgentKey(ctx *common.Context, tx *Tx, agentKey string, zoneID uint64) *Agents {
	if a.agentCache != nil {
		agent := a.agentCache.GetAgentByAgentKey(ctx, zoneID, agentKey)
		if agent != nil {
			return agent
		}
	}

	agent := tx.getAgentByAgentKey(agentKey, zoneID)

	return agent
}

func (a *AgentStorage) Close() error {
	if a.agentCache != nil {
		return a.agentCache.Close()
	}

	return nil
}
