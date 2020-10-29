package manager

import (
	"fmt"
	"strconv"
	"time"

	"xorm.io/builder"

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

func (tx *Tx) deletePrimaryAgent(zoneID uint64) {
	sql := "delete p from PRIMARY_AGENTS p where p.GROUP_ID = ?"
	res, err := tx.Exec(sql, zoneID)
	if err != nil {
		logger.Warningf("%+v", errors.Wrap(err, "sql error"))
		panic(err)
	}

	logger.Debug(res)

}

func (tx *Tx) getAgentByAgentKey(agentKey string, groupID uint64) *Agents {
	var a Agents

	common.CheckGetQuery(tx.Where("agent_key = ?", agentKey).And("group_id = ?", groupID).Get(&a))
	logger.Debugf("selected Agent - id : [%d], agentKey : [%s], isActive : [%v], lastAccessTime : [%v]", a.Id, a.AgentKey, a.IsActive, a.LastAccessTime)
	return &a
}

func (tx *Tx) getAgentByID(id uint64) *Agents {
	var a Agents

	common.CheckGetQuery(tx.ID(id).Get(&a))
	logger.Debugf("Selected Agent : %+v", a)

	return &a
}

func (tx *Tx) getAgentsByGroupId(groupID uint64) (int64, *[]Agents) {
	var agents []Agents

	cnt, err := tx.Where("GROUP_ID = ?", groupID).FindAndCount(&agents)
	if err != nil {
		panic(err)
	}

	return cnt, &agents
}

func (tx *Tx) getAgentsForInactive(before time.Time) (int64, *[]Agents) {
	var agents []Agents

	err := tx.Where("IS_ACTIVE = ?", true).And("LAST_ACCESS_TIME < ?", before).Cols("ID").Find(&agents)
	if err != nil {
		panic(err)
	}

	cnt := int64(len(agents))
	// cnt, err := tx.Table(&Agents{}).
	// 	Join("INNER", "PRIMARY_AGENTS", "AGENTS.ID = PRIMARY_AGENTS.AGENT_ID").
	// 	Where("AGENTS.IS_ACTIVE = ?", true).And("AGENTS.LAST_ACCESS_TIME < ?", before).
	// 	Cols("AGENTS.ID").FindAndCount(&agents)
	// if err != nil {
	// 	panic(err)
	// }

	logger.Debugf("getAgentsForInactive count : %d, agents : %+v", cnt, &agents)

	return cnt, &agents
}

func (tx *Tx) updateAgentStatus(ids []uint64) {
	cnt, err := tx.Table(new(Agents)).In("id", ids).Update(map[string]interface{}{"IS_ACTIVE": 0})
	logger.Debugf("Status updated Agent(%d) : [%+v]", cnt, ids)

	if err != nil {
		panic(err)
	} else if cnt != int64(len(ids)) {
		common.PanicForUpdate("updated", cnt, int64(len(ids)))
	}
}

func (tx *Tx) updateAccessAgent(id uint64, accessTime time.Time) {
	result, err := tx.Exec("UPDATE `AGENTS` SET `LAST_ACCESS_TIME` = ?, `IS_ACTIVE` = ? WHERE ID = ?",
		accessTime, 1, id)

	if err != nil {
		panic(err)
	}

	cnt, _ := result.RowsAffected()

	logger.Debugf("Access information updated Agent(%d) : [%+v]", cnt, id)
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
	cnt, err := tx.Table(new(Agents)).Where("id = ?", a.Id).Update(map[string]interface{}{
		"CPU":                   a.Cpu,
		"DISK":                  a.Disk,
		"ENC_KEY":               a.EncKey,
		"HMAC_KEY":              a.HmacKey,
		"IP":                    a.Ip,
		"IS_ACTIVE":             a.IsActive,
		"LAST_ACCESS_TIME":      a.LastAccessTime,
		"LAST_ALIVE_CHECK_TIME": a.LastAliveCheckTime,
		"MEMORY":                a.Memory,
		"PORT":                  a.Port,
		"VERSION":               a.Version,
		"UPDATED_AT":            time.Now().UTC(),
	})
	// cnt, err := tx.Where("id = ?", a.Id).Update(a)
	logger.Debugf("Updated Agent - id : [%d], agentKey : [%s], isActive : [%v], lastAccessTime : [%v]", a.Id, a.AgentKey, a.IsActive, a.LastAccessTime)
	// logger.Debugf("Updated Agent(%d) : %v", cnt, a)

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

func (tx *Tx) insertTask(t *Tasks) *Tasks {
	result, err := tx.Exec("INSERT INTO `TASKS` (`zone_id`,`name`,`task_type`,`schedule`,`agent_key`,`exe_agent_key`,`status`) VALUES (?,?,?,?,?,?,?)",
		t.ZoneId, t.Name, t.TaskType, t.Schedule, t.AgentKey, t.ExeAgentKey, t.Status)

	if err != nil {
		panic(err)
	}

	cnt, _ := result.RowsAffected()

	if cnt != 1 {
		common.PanicForUpdate("inserted", cnt, 1)
	}

	id, _ := result.LastInsertId()

	t.Id = uint64(id)

	logger.Debugf("Inserted task : %v", t)

	if t.TaskDetail != nil {
		t.TaskDetail.TaskId = t.Id

		cnt, err = tx.Insert(t.TaskDetail)
		if err != nil {
			panic(err)
		}

		if cnt != 1 {
			common.PanicForUpdate("inserted", cnt, 1)
		}
	}

	taskStepLen := int64(len(*t.TaskSteps))

	if t.TaskSteps != nil && taskStepLen > 0 {
		for i, ts := range *t.TaskSteps {
			(*t.TaskSteps)[i].TaskId = t.Id
			logger.Debugf("TaskStep %d : [%+v]", i, ts)
			logger.Debugf("%v, %v", ((*t.TaskSteps)[i]), ts)
		}

		cnt, err = tx.Insert(t.TaskSteps)
		if err != nil {
			panic(err)
		}

		if cnt != taskStepLen {
			common.PanicForUpdate("inserted", cnt, taskStepLen)
		}
	}

	return t
}

func (tx *Tx) updateHandoverTasks(ids []uint64) {
	cnt, err := tx.Table(new(Tasks)).In("ID", ids).Update(map[string]interface{}{"STATUS": common.HandOver})
	if err != nil {
		panic(err)
	} else if cnt > 1 {
		common.PanicForUpdate("updated", cnt, int64(len(ids)))
	}
}

func (tx *Tx) updateTask(t *Tasks) {
	cnt, err := tx.Where("ID = ?", t.Id).
		Cols("EXE_AGENT_KEY", "STATUS").
		Update(t)
	logger.Debugf("Updated TASK(%d) : %v", cnt, t)

	if err != nil {
		panic(err)
	}

	detail := t.TaskDetail
	if detail.Result != "" {
		cnt, err = tx.Where("TASK_ID = ?", t.Id).
			Cols("CURRENT_STEP", "RESULT", "FAILED_STEP", "IS_FAILED_RECOVER").
			Update(detail)
	} else {
		cnt, err = tx.Where("TASK_ID = ?", t.Id).
			Cols("CURRENT_STEP", "FAILED_STEP", "IS_FAILED_RECOVER").
			Update(detail)
	}

	if err != nil {
		panic(err)
	}

	if t.Logs.Logs != "" {
		logs := t.Logs

		_, err := tx.Exec("INSERT INTO `TASK_LOGS` (`TASK_ID`,`LOGS`) VALUES (?,?) ON DUPLICATE KEY UPDATE TASK_ID = ?, LOGS=CONCAT_WS('\\n', LOGS, ?)",
			t.Id, logs.Logs, t.Id, logs.Logs)

		if err != nil {
			panic(err)
		}
	}
}

func (tx *Tx) getTask(id uint64) (*Tasks, bool) {
	var rTask RetriveTask
	var task Tasks

	stmt := tx.Join("LEFT OUTER", "TASK_DETAIL", "TASK_DETAIL.TASK_ID = TASKS.ID")
	stmt = stmt.Join("LEFT OUTER", "TASK_LOGS", "TASK_LOGS.TASK_ID = TASKS.ID")

	exist := common.CheckGetQuery(stmt.Where("TASKS.ID = ?", id).Get(&rTask))
	logger.Debugf("Selected Task : %v", rTask)

	var steps []TaskSteps

	err := tx.Where("TASK_ID = ?", id).Find(&steps)
	if err != nil {
		panic(err)
	}

	task = rTask.Tasks
	task.TaskDetail = &rTask.TaskDetail
	task.Logs = &rTask.TaskLogs
	task.TaskSteps = &steps

	return &task, exist
}

func (tx *Tx) getTasksByIds(ids []uint64) (*[]Tasks, int64) {
	var rts []RetriveTask
	inqCnt := len(ids)

	stmt := tx.Join("LEFT OUTER", "TASK_DETAIL", "TASK_DETAIL.TASK_ID = TASKS.ID")
	stmt = stmt.Join("LEFT OUTER", "TASK_LOGS", "TASK_LOGS.TASK_ID = TASKS.ID")

	cnt, err := stmt.Where(builder.In("ID", ids)).FindAndCount(&rts)
	logger.Debugf("Selected Tasks : %+v", rts)

	if err != nil {
		panic(err)
	}

	if int(cnt) != inqCnt {
		panic("The number of inquired tasks does not match - inquired count : " +
			strconv.Itoa(inqCnt) +
			", selected count : " + strconv.Itoa(int(cnt)))
	}

	return toTasks(&rts), cnt
}

func (tx *Tx) getTasksWithSteps(groupID uint64, statuses []string) (*[]Tasks, int64) {
	var tasks *[]Tasks

	var rts []RetriveTask

	stmt := tx.Join("LEFT OUTER", "TASK_DETAIL", "TASK_DETAIL.TASK_ID = TASKS.ID")
	stmt = stmt.Join("LEFT OUTER", "TASK_LOGS", "TASK_LOGS.TASK_ID = TASKS.ID")

	cnt, err := stmt.Where("ZONE_ID = ?", groupID).And(builder.In("STATUS", statuses)).FindAndCount(&rts)

	if err != nil {
		panic(err)
	}

	logger.Debugf("selected retreive tasks : %d", cnt)

	tasks = toTasks(&rts)

	for i, t := range *tasks {
		var steps []TaskSteps

		cnt, err := tx.Where("TASK_ID = ?", t.Id).OrderBy("SEQ ASC").FindAndCount(&steps)
		if err != nil {
			panic(err)
		}

		logger.Debugf("select steps for %d - %d", t.Id, cnt)

		(*tasks)[i].TaskSteps = &steps
	}

	logger.Debugf("tasks : [%+v]", tasks)

	return tasks, cnt
}

func toTasks(rts *[]RetriveTask) *[]Tasks {
	var tasks = make([]Tasks, 0, len(*rts))

	for _, rt := range *rts {
		rt.Tasks.TaskDetail = &rt.TaskDetail
		rt.Tasks.Logs = &rt.TaskLogs
		tasks = append(tasks, rt.Tasks)
	}

	return &tasks
}

func (tx *Tx) getTasks(groupIDs []uint64, statuses []string, agentKeys []string, taskNames []string) (*[]Tasks, bool) {
	var tasks []Tasks

	stmt := tx.Where(builder.In("ZONE_ID", groupIDs))

	// condition 추가
	if statuses != nil {
		stmt = stmt.And(builder.In("STATUS", statuses))
	}
	if agentKeys != nil {
		stmt = stmt.And(builder.In("AGENT_KEY", agentKeys))
	}
	if taskNames != nil {
		stmt = stmt.And(builder.In("NAME", taskNames))
	}

	cnt, err := stmt.FindAndCount(&tasks)
	if err != nil {
		panic(err)
	}

	logger.Debugf("Selected Tasks : %d", cnt)

	return &tasks, cnt > 0
}

func (tx *Tx) updateScheduledTask() int64 {
	result, err := tx.Exec("UPDATE TASKS SET STATUS = ? WHERE STATUS = ? AND SCHEDULE < CURRENT_TIMESTAMP() AND DELETED_AT IS NULL",
		common.WaitPolling, common.Scheduled)
	if err != nil {
		panic(err)
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return cnt
}

func (tx *Tx) cancelTask(id uint64) bool {
	result, err := tx.Exec("UPDATE TASKS SET STATUS = ? WHERE ID = ? AND STATUS IN (?, ?)",
		common.Canceled, id, common.Scheduled, common.WaitPolling)
	if err != nil {
		panic(err)
	}

	cnt, _ := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	if cnt == 0 {
		return false
	}

	return true
}

// Task 부가 데이터 삭제
func (tx *Tx) purgeTask(id uint64) {
	cnt, err := tx.Where("TASK_ID = ?").Delete(&TaskDetail{})
	if err != nil {
		panic(err)
	}
	if cnt != 1 {
		common.PanicForUpdate("deleted", cnt, 1)
	}

	cnt, err = tx.Where("TASK_ID = ?").Delete(&TaskLogs{})
	if err != nil {
		panic(err)
	}

	cnt, err = tx.Where("TASK_ID = ?").Delete(&TaskSteps{})
	if err != nil {
		panic(err)
	}
}
