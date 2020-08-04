package manager

import (
	"time"

	"github.com/pkg/errors"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"xorm.io/xorm"
)

func getPrimaryAgent(conn *xorm.Session, zoneID uint64) *PrimaryAgents {
	var pa PrimaryAgents

	common.CheckGetQuery(conn.Where("group_id = ?", zoneID).Get(&pa))
	logger.Debugf("Selected PrimaryAgent : %v", pa)

	return &pa
}

func insertPrimaryAgent(conn *xorm.Session, pa *PrimaryAgents) (int64, error) {
	return conn.Insert(pa)
}

func deletePrimaryAgentIfOld(conn *xorm.Session, zoneID uint64, agentID uint64, time time.Duration) {
	sql := "delete p from PRIMARY_AGENTS p join AGENT a on a.ID = p.AGENT_ID where p.GROUP_ID = ? and p.AGENT_ID = ? and a.LAST_ACCESS_TIME < ?"
	res, err := conn.Exec(sql, zoneID, agentID, time)
	if err != nil {
		logger.Warningf("%+v", errors.Wrap(err, "sql error"))
	}

	logger.Debug(res)

}

func getAgentByAgentKey(conn *xorm.Session, agentKey string) *Agents {
	var a Agents

	common.CheckGetQuery(conn.Where("agent_key = ?", agentKey).Get(&a))
	logger.Debugf("Selected Agent : %v", a)

	return &a
}

func getAgentByID(conn *xorm.Session, id uint64) *Agents {
	var a Agents

	common.CheckGetQuery(conn.ID(id).Get(&a))
	logger.Debugf("Selected Agent : %v", a)

	return &a
}

func getAgentsForInactive(conn *xorm.Session, before time.Time) (int64, *[]Agents) {
	var agents []Agents

	cnt, err := conn.Where("IS_ACTIVE = ?", true).And("LAST_ALIVE_CHECK_TIME < ?", before).Cols("ID").FindAndCount(&agents)
	if err != nil {
		panic(err)
	}

	return cnt, &agents
}

func updateAgentStatus(conn *xorm.Session, ids []uint64) {
	cnt, err := conn.Table(new(Agents)).In("id", ids).Update(map[string]interface{}{"IS_ACTIVE": false})

	if err != nil {
		panic(err)
	} else if cnt != 1 {
		common.PanicForUpdate("updated", cnt, int64(len(ids)))
	}
}

func getAgentGroups(conn *xorm.Session) *[]AgentGroups {
	var ags []AgentGroups

	err := conn.Find(&ags)
	if err != nil {
		panic(err)
	}

	return &ags
}

func getAgentGroup(conn *xorm.Session, zoneID uint64) (*AgentGroups, bool) {
	var ag AgentGroups

	exist := common.CheckGetQuery(conn.ID(zoneID).Get(&ag))
	logger.Debugf("Selected AgentGroup : %v", ag)

	return &ag, exist
}

func addAgentGroup(conn *xorm.Session, ag *AgentGroups) {
	cnt, err := conn.Insert(ag)

	logger.Debugf("Inserted AgentGroup(%d) : %v", cnt, ag)

	if err != nil {
		panic(err)
	} else if cnt != 1 {
		common.PanicForUpdate("inserted", cnt, 1)
	}
}

func addAgent(conn *xorm.Session, a *Agents) {
	cnt, err := conn.Insert(a)
	logger.Debugf("Inserted Agent(%d) : %v", cnt, a)

	if err != nil {
		panic(err)
	} else if cnt != 1 {
		common.PanicForUpdate("inserted", cnt, 1)
	}
}

func updateAgent(conn *xorm.Session, a *Agents) {
	cnt, err := conn.Where("id = ?", a.Id).Update(a)
	logger.Debugf("Updated Agent(%d) : %v", cnt, a)

	if err != nil {
		panic(err)
	} else if cnt != 1 {
		common.PanicForUpdate("updated", cnt, 1)
	}
}

func addAPIKey(conn *xorm.Session, auth *ApiAuthentications) {
	cnt, err := conn.Insert(auth)
	if err != nil {
		panic(err)
	} else if cnt != 1 {
		common.PanicForUpdate("inserted", cnt, 1)
	}
}

func updateAPIKey(conn *xorm.Session, auth *ApiAuthentications) {
	cnt, err := conn.Where("group_id = ?", auth.GroupId).Update(auth)
	logger.Debugf("Updated APIKey(%d) : %v", cnt, auth)

	if err != nil {
		panic(err)
	} else if cnt != 1 {
		common.PanicForUpdate("updated", cnt, 1)
	}
}

func getAPIKey(conn *xorm.Session, groupID uint64) (*ApiAuthentications, bool) {
	var auth ApiAuthentications
	exist := common.CheckGetQuery(conn.Where("group_id = ?", groupID).Get(&auth))
	logger.Debugf("Selected ApiAuthentications : %v", auth)

	return &auth, exist
}

func existAPIKey(conn *xorm.Session, apiKey string, zoneID uint64) bool {
	cnt, err := conn.Where("api_key = ?", apiKey).And("group_id = ?", zoneID).Count(&ApiAuthentications{})
	logger.Debugf("Selected API_KEY count : %d", cnt)

	if err != nil {
		panic(err)
	}

	return cnt > 0
}

func getLock(conn *xorm.Session, task string) (*TaskLock, bool) {
	var tl TaskLock
	exist := common.CheckGetQuery(conn.Where("task = ?", task).ForUpdate().Get(&tl))
	logger.Debugf("Selected TaskLock : %v", tl)

	return &tl, exist
}

func insertLock(conn *xorm.Session, tl *TaskLock) {
	cnt, err := conn.Insert(tl)

	if err != nil {
		panic(err)
	} else if cnt != 1 {
		common.PanicForUpdate("inserted", cnt, 1)
	}
}

func updateLock(conn *xorm.Session, tl *TaskLock) {
	cnt, err := conn.Where("task = ?", tl.Task).Update(tl)
	logger.Debugf("Updated TaskLock(%d) : %+v", cnt, tl)

	if err != nil {
		panic(err)
	} else if cnt != 1 {
		common.PanicForUpdate("updated", cnt, 1)
	}
}
