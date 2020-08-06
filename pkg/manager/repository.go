package manager

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
)

type Tx struct {
	*common.Session
}

func (tx *Tx) getPrimaryAgent(zoneID uint64) *PrimaryAgents {
	var pa PrimaryAgents

	common.CheckGetQuery(tx.Where("group_id = ?", zoneID).Get(&pa))
	logger.Debugf("Selected PrimaryAgent : %v", pa)

	return &pa
}

func (tx *Tx) insertPrimaryAgent(pa *PrimaryAgents) (int64, error) {
	return tx.Insert(pa)
}

func (tx *Tx) deletePrimaryAgentIfOld(zoneID uint64, agentID uint64, time time.Duration) {
	sql := "delete p from PRIMARY_AGENTS p join AGENT a on a.ID = p.AGENT_ID where p.GROUP_ID = ? and p.AGENT_ID = ? and a.LAST_ACCESS_TIME < ?"
	res, err := tx.Exec(sql, zoneID, agentID, time)
	if err != nil {
		logger.Warningf("%+v", errors.Wrap(err, "sql error"))
	}

	logger.Debug(res)

}

func (tx *Tx) getAgentByAgentKey(agentKey string) *Agents {
	var a Agents

	common.CheckGetQuery(tx.Where("agent_key = ?", agentKey).Get(&a))
	logger.Debugf("Selected Agent : %v", a)

	return &a
}

func (tx *Tx) getAgentByID(id uint64) *Agents {
	var a Agents

	common.CheckGetQuery(tx.ID(id).Get(&a))
	logger.Debugf("Selected Agent : %v", a)

	return &a
}

func (tx *Tx) getAgentsForInactive(before time.Time) (int64, *[]Agents) {
	var agents []Agents

	cnt, err := tx.Table(&Agents{}).
		Join("INNER", "PRIMARY_AGENTS", "AGENTS.ID = PRIMARY_AGENTS.AGENT_ID").
		Where("AGENTS.IS_ACTIVE = ?", true).And("AGENTS.LAST_ACCESS_TIME < ?", before).
		Cols("AGENTS.ID").FindAndCount(&agents)
	if err != nil {
		panic(err)
	}

	return cnt, &agents
}

func (tx *Tx) updateAgentStatus(ids []uint64) {
	cnt, err := tx.Table(new(Agents)).In("id", ids).Update(map[string]interface{}{"IS_ACTIVE": false})

	if err != nil {
		panic(err)
	} else if cnt != int64(len(ids)) {
		common.PanicForUpdate("updated", cnt, int64(len(ids)))
	}
}

func (tx *Tx) updateZoneStatus(arrAgent *[]Agents) {
	for _, a := range *arrAgent {
		cnt, err := tx.Where("AGENT_KEY = ?", a.AgentKey).
			Cols("LAST_ALIVE_CHECK_TIME", "IS_ACTIVE", "CPU", "MEMORY", "DISK").
			Update(a)

		if err != nil {
			panic(err)
		} else if cnt != 1 {
			common.PanicForUpdate(fmt.Sprintf("updated - agentKey : %s", a.AgentKey), cnt, 1)
		}
	}
}

func (tx *Tx) getAgentGroups() *[]AgentGroups {
	var ags []AgentGroups

	err := tx.Find(&ags)
	if err != nil {
		panic(err)
	}

	return &ags
}

func (tx *Tx) getAgentGroup(zoneID uint64) (*AgentGroups, bool) {
	var ag AgentGroups

	exist := common.CheckGetQuery(tx.ID(zoneID).Get(&ag))
	logger.Debugf("Selected AgentGroup : %v", ag)

	return &ag, exist
}

func (tx *Tx) addAgentGroup(ag *AgentGroups) {
	cnt, err := tx.Insert(ag)

	logger.Debugf("Inserted AgentGroup(%d) : %v", cnt, ag)

	if err != nil {
		panic(err)
	} else if cnt != 1 {
		common.PanicForUpdate("inserted", cnt, 1)
	}
}

func (tx *Tx) addAgent(a *Agents) {
	cnt, err := tx.Insert(a)
	logger.Debugf("Inserted Agent(%d) : %v", cnt, a)

	if err != nil {
		panic(err)
	} else if cnt != 1 {
		common.PanicForUpdate("inserted", cnt, 1)
	}
}

func (tx *Tx) updateAgent(a *Agents) {
	cnt, err := tx.Where("id = ?", a.Id).Update(a)
	logger.Debugf("Updated Agent(%d) : %v", cnt, a)

	if err != nil {
		panic(err)
	} else if cnt != 1 {
		common.PanicForUpdate("updated", cnt, 1)
	}
}

func (tx *Tx) addAPIKey(auth *ApiAuthentications) {
	cnt, err := tx.Insert(auth)
	if err != nil {
		panic(err)
	} else if cnt != 1 {
		common.PanicForUpdate("inserted", cnt, 1)
	}
}

func (tx *Tx) updateAPIKey(auth *ApiAuthentications) {
	cnt, err := tx.Where("group_id = ?", auth.GroupId).Update(auth)
	logger.Debugf("Updated APIKey(%d) : %v", cnt, auth)

	if err != nil {
		panic(err)
	} else if cnt != 1 {
		common.PanicForUpdate("updated", cnt, 1)
	}
}

func (tx *Tx) getAPIKey(groupID uint64) (*ApiAuthentications, bool) {
	var auth ApiAuthentications
	exist := common.CheckGetQuery(tx.Where("group_id = ?", groupID).Get(&auth))
	logger.Debugf("Selected ApiAuthentications : %v", auth)

	return &auth, exist
}

func (tx *Tx) existAPIKey(apiKey string, zoneID uint64) bool {
	cnt, err := tx.Where("api_key = ?", apiKey).And("group_id = ?", zoneID).Count(&ApiAuthentications{})
	logger.Debugf("Selected API_KEY count : %d", cnt)

	if err != nil {
		panic(err)
	}

	return cnt > 0
}

func (tx *Tx) getLock(task string) (*TaskLock, bool) {
	var tl TaskLock
	exist := common.CheckGetQuery(tx.Where("task = ?", task).ForUpdate().Get(&tl))
	logger.Debugf("Selected TaskLock : %v", tl)

	return &tl, exist
}

func (tx *Tx) insertLock(tl *TaskLock) {
	cnt, err := tx.Insert(tl)

	if err != nil {
		panic(err)
	} else if cnt != 1 {
		common.PanicForUpdate("inserted", cnt, 1)
	}
}

func (tx *Tx) updateLock(tl *TaskLock) {
	cnt, err := tx.Where("task = ?", tl.Task).Update(tl)
	logger.Debugf("Updated TaskLock(%d) : %+v", cnt, tl)

	if err != nil {
		panic(err)
	} else if cnt != 1 {
		common.PanicForUpdate("updated", cnt, 1)
	}
}
