package manager

import (
	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"xorm.io/xorm"
)

func getPrimaryAgent(conn *xorm.Session, zoneID uint) *PrimaryAgents {
	var pa PrimaryAgents

	common.CheckGetQuery(conn.Where("group_id = ?", zoneID).Get(&pa))
	logger.Debugf("Selected PrimaryAgent : %v", pa)

	return &pa
}

func getAgentByAgentKey(conn *xorm.Session, agentKey string) *Agents {
	var a Agents

	common.CheckGetQuery(conn.Where("agent_key = ?", agentKey).Get(&a))
	logger.Debugf("Selected Agent : %v", a)

	return &a
}

func getAgentByID(conn *xorm.Session, id uint) *Agents {
	var a Agents

	common.CheckGetQuery(conn.ID(id).Get(&a))
	logger.Debugf("Selected Agent : %v", a)

	return &a
}

func getAgentGroup(conn *xorm.Session, zoneID uint) *AgentGroups {
	var ag AgentGroups

	common.CheckGetQuery(conn.ID(zoneID).Get(&ag))
	logger.Debugf("Selected AgentGroup : %v", ag)

	return &ag
}

func addAgent(conn *xorm.Session, a *Agents) {
	cnt, err := conn.Insert(a)
	logger.Debugf("Inserted AgentGroup(%d) : %v", cnt, a)

	if cnt != 1 {
		common.PanicForUpdate("inserted", cnt, 1)
	} else if err != nil {
		panic(err)
	}
}

func updateAgent(conn *xorm.Session, a *Agents) {
	cnt, err := conn.Where("id = ?", a.Id).Update(a)
	logger.Debugf("Updated AgentGroup(%d) : %v", cnt, a)

	if cnt != 1 {
		common.PanicForUpdate("updated", cnt, 1)
	} else if err != nil {
		panic(err)
	}
}

func existAPIKey(conn *xorm.Session, apiKey string, zoneID uint) bool {
	cnt, err := conn.Where("api_key = ?", apiKey).And("group_id = ?", zoneID).Count(&ApiAuthentications{})
	logger.Debugf("Selected API_KEY count : %d", cnt)

	if err != nil {
		panic(err)
	}

	return cnt > 0
}

func getLock(conn *xorm.Session, task string) *TaskLock {
	var tl TaskLock
	common.CheckGetQuery(conn.Where("task = ?", task).ForUpdate().Get(&tl))
	logger.Debugf("Selected TaskLock : %v", tl)

	return &tl
}
