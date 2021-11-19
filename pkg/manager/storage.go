package manager

import (
	"time"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/model"
)

type AgentStorage struct {
	agentCache *Cache
}

func NewAgentStorage(address string, port int, password string) *AgentStorage {
	s := &AgentStorage{}
	lock := ctx.Get(CtxCacheLock)
	if lock != nil {
		s.agentCache = NewCache(address, port, password)
	}

	return s
}

func (a *AgentStorage) AddAgent(ctx *common.Context, tx *Tx, agent *model.Agents) error {
	agent.CreatedAt = time.Now().UTC()
	err := tx.addAgent(agent)

	if err == nil && a.agentCache != nil {
		a.agentCache.UpdateAgent(ctx, agent.GroupId, agent)
	}

	if err != nil {
		return err
	}

	return nil
}

func (a *AgentStorage) UpdateAgent(ctx *common.Context, tx *Tx, agent *model.Agents) {
	agent.UpdatedAt = time.Now().UTC()
	tx.updateAgent(agent)

	if a.agentCache != nil {
		//logger.Debugf("##### AgentStorage::UpdateAgent(grouID:%d", agent.GroupId)
		a.agentCache.UpdateAgent(ctx, agent.GroupId, agent)
	}
}

func (a *AgentStorage) UpdateZoneStatus(ctx *common.Context, tx *Tx, zoneID uint64, arrAgent []model.Agents) {
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

func (a *AgentStorage) UpdateAgentStatus(ctx *common.Context, tx *Tx, inactiveAgents *[]model.Agents, ids []uint64) {
	if a.agentCache != nil {
		a.agentCache.UpdateAgentDisabledStatus(ctx, inactiveAgents)
	}

	tx.updateAgentStatus(ids)
}

func (a *AgentStorage) GetAgentsForInactive(ctx *common.Context, tx *Tx, before time.Time) (int64, *[]model.Agents) {
	if a.agentCache != nil {
		cnt, agents := a.agentCache.GetAgentsForInactive(ctx, before)
		if cnt > 0 {
			return cnt, agents
		}
	}

	cnt, agents := tx.getAgentsForInactive(before)

	return cnt, agents
}

func (a *AgentStorage) GetAgentsByZoneID(ctx *common.Context, tx *Tx, zoneID uint64) (int64, *[]model.Agents) {
	if a.agentCache != nil {
		cnt, agents := a.agentCache.GetAgentsByZoneID(ctx, zoneID)
		if cnt > 0 {
			return cnt, agents
		}
	}

	cnt, agents := tx.getAgentsByGroupId(zoneID)

	return cnt, agents
}

func (a *AgentStorage) GetAgentByID(ctx *common.Context, tx *Tx, zoneID uint64, agentID uint64) *model.Agents {
	if a.agentCache != nil {
		agent := a.agentCache.GetAgentByID(ctx, zoneID, agentID)
		if agent != nil {
			return agent
		}
	}

	agent := tx.getAgentByID(agentID)

	return agent
}

func (a *AgentStorage) GetAgentByAgentKey(ctx *common.Context, tx *Tx, agentKey string, zoneID uint64) *model.Agents {
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
