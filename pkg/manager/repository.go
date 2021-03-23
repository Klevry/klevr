package manager

import (
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

func (tx *Tx) getPrimaryAgent(zoneID uint64) (primary *PrimaryAgents, exist bool) {
	var pa PrimaryAgents

	ok := common.CheckGetQuery(tx.Where("group_id = ?", zoneID).Get(&pa))
	logger.Debugf("Selected PrimaryAgent : %v", pa)

	return &pa, ok
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

	err := tx.Where("IS_ACTIVE = ?", true).And("LAST_ALIVE_CHECK_TIME < ?", before).Cols("ID", "AGENT_KEY, GROUP_ID").Find(&agents)
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
	}
}

func (tx *Tx) updateAccessAgent(agentKey string, accessTime time.Time) int64 {
	result, err := tx.Exec("UPDATE `AGENTS` SET `LAST_ACCESS_TIME` = ?, `IS_ACTIVE` = 1 WHERE AGENT_KEY = ?",
		accessTime, agentKey)

	if err != nil {
		panic(err)
	}

	cnt, _ := result.RowsAffected()

	logger.Debugf("Access information updated Agent(%d) : [%+v]", cnt, agentKey)

	return cnt
}

func (tx *Tx) deleteAgent(zoneID uint64) {
	sql := "delete a from AGENTS a where a.GROUP_ID = ?"
	res, err := tx.Exec(sql, zoneID)
	if err != nil {
		logger.Warningf("%+v", errors.Wrap(err, "sql error"))
		panic(err)
	}

	logger.Debug(res)
}

func (tx *Tx) updateZoneStatus(arrAgent *[]Agents) {
	for _, a := range *arrAgent {
		_, err := tx.Where("AGENT_KEY = ?", a.AgentKey).
			Cols("LAST_ALIVE_CHECK_TIME", "IS_ACTIVE", "CPU", "MEMORY", "DISK").
			Update(a)

		if err != nil {
			panic(err)
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
	}
}

func (tx *Tx) deleteAgentGroup(groupID uint64) {
	_, err := tx.Where("ID = ?", groupID).Delete(&AgentGroups{})
	if err != nil {
		panic(err)
	}
}

func (tx *Tx) addAgent(a *Agents) {
	cnt, err := tx.Insert(a)
	logger.Debugf("Inserted Agent(%d) : %v", cnt, a)

	if err != nil {
		panic(err)
	}
}

func (tx *Tx) updateAgent(a *Agents) {
	_, err := tx.Table(new(Agents)).Where("id = ?", a.Id).Update(map[string]interface{}{
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
	}
}

func (tx *Tx) addAPIKey(auth *ApiAuthentications) {
	_, err := tx.Insert(auth)
	if err != nil {
		panic(err)
	}
}

func (tx *Tx) updateAPIKey(auth *ApiAuthentications) {
	cnt, err := tx.Where("group_id = ?", auth.GroupId).Update(auth)
	logger.Debugf("Updated APIKey(%d) : %v", cnt, auth)

	if err != nil {
		panic(err)
	}
}

func (tx *Tx) getAPIKey(groupID uint64) (*ApiAuthentications, bool) {
	var auth ApiAuthentications
	exist := common.CheckGetQuery(tx.Where("group_id = ?", groupID).Get(&auth))
	qry, args := tx.LastSQL()
	logger.Debugf("query: %s, args: %v", qry, args)
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
	_, err := tx.Insert(tl)

	if err != nil {
		panic(err)
	}
}

func (tx *Tx) updateLock(tl *TaskLock) {
	cnt, err := tx.Where("task = ?", tl.Task).Update(tl)
	logger.Debugf("Updated TaskLock(%d) : %+v", cnt, tl)

	if err != nil {
		panic(err)
	}
}

func (tx *Tx) insertTask(manager *KlevrManager, t *Tasks) *Tasks {
	result, err := tx.Exec("INSERT INTO `TASKS` (`zone_id`,`name`,`task_type`,`schedule`,`agent_key`,`exe_agent_key`,`status`) VALUES (?,?,?,?,?,?,?)",
		t.ZoneId, t.Name, t.TaskType, t.Schedule, t.AgentKey, t.ExeAgentKey, t.Status)

	if err != nil {
		panic(err)
	}

	result.RowsAffected()

	id, _ := result.LastInsertId()

	t.Id = uint64(id)

	logger.Debugf("Inserted task : %v", t)

	if t.TaskDetail != nil {
		t.TaskDetail.TaskId = t.Id

		t.TaskDetail.Parameter = manager.encrypt(t.TaskDetail.Parameter)
		t.TaskDetail.Result = manager.encrypt(t.TaskDetail.Result)

		_, err = tx.Insert(t.TaskDetail)
		if err != nil {
			panic(err)
		}
	}

	taskStepLen := int64(len(*t.TaskSteps))

	if t.TaskSteps != nil && taskStepLen > 0 {
		steps := make([]*TaskSteps, 0)
		for i, ts := range *t.TaskSteps {
			(*t.TaskSteps)[i].TaskId = t.Id
			(*t.TaskSteps)[i].ReservedCommand = manager.encrypt(ts.ReservedCommand)
			(*t.TaskSteps)[i].InlineScript = manager.encrypt(ts.InlineScript)

			logger.Debugf("TaskStep %d : [%+v]", i, ts)
			logger.Debugf("%v, %v", ((*t.TaskSteps)[i]), ts)

			step := &TaskSteps{
				Id:              (*t.TaskSteps)[i].Id,
				Seq:             (*t.TaskSteps)[i].Seq,
				TaskId:          (*t.TaskSteps)[i].TaskId,
				CommandName:     (*t.TaskSteps)[i].CommandName,
				CommandType:     (*t.TaskSteps)[i].CommandType,
				ReservedCommand: (*t.TaskSteps)[i].ReservedCommand,
				InlineScript:    (*t.TaskSteps)[i].InlineScript,
				IsRecover:       (*t.TaskSteps)[i].IsRecover,
			}
			steps = append(steps, step)
		}

		_, err = tx.Insert(steps)
		if err != nil {
			panic(err)
		}

	}

	return t
}

func (tx *Tx) updateHandoverTasks(ids []uint64) {
	_, err := tx.Table(new(Tasks)).In("ID", ids).Update(map[string]interface{}{"STATUS": common.HandOver})
	if err != nil {
		panic(err)
	}
}

func (tx *Tx) updateTask(manager *KlevrManager, t *Tasks) {
	cnt, err := tx.Where("ID = ?", t.Id).
		Cols("EXE_AGENT_KEY", "STATUS").
		Update(t)
	logger.Debugf("Updated TASK(%d) : %v", cnt, t)

	if err != nil {
		panic(err)
	}

	detail := t.TaskDetail
	if detail.Result != "" {
		detail.Result = manager.encrypt(detail.Result)

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

		logs.Logs = manager.encrypt(logs.Logs)

		_, err := tx.Exec("INSERT INTO `TASK_LOGS` (`TASK_ID`,`LOGS`) VALUES (?,?) ON DUPLICATE KEY UPDATE TASK_ID = ?, LOGS= ?",
			t.Id, logs.Logs, t.Id, logs.Logs)

		if err != nil {
			panic(err)
		}
	}
}

func (tx *Tx) getTask(manager *KlevrManager, id uint64) (*Tasks, bool) {
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

	task = *rTask.Tasks
	task.TaskDetail = rTask.TaskDetail
	task.Logs = rTask.TaskLogs

	if task.TaskDetail.Parameter != "" {
		task.TaskDetail.Parameter = manager.decrypt(task.TaskDetail.Parameter)
	}
	if task.TaskDetail.Result != "" {
		task.TaskDetail.Result = manager.decrypt(task.TaskDetail.Result)
	}
	if task.Logs.Logs != "" {
		task.Logs.Logs = manager.decrypt(task.Logs.Logs)
	}

	for i := 0; i < len(steps); i++ {
		if steps[i].ReservedCommand != "" {
			steps[i].ReservedCommand = manager.decrypt(steps[i].ReservedCommand)
		}
		if steps[i].InlineScript != "" {
			steps[i].InlineScript = manager.decrypt(steps[i].InlineScript)
		}
	}
	task.TaskSteps = &steps

	return &task, exist
}

func (tx *Tx) getTasksByIds(manager *KlevrManager, ids []uint64) (*[]Tasks, int64) {
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

	return toTasks(manager, &rts), cnt
}

func (tx *Tx) getTasksWithSteps(manager *KlevrManager, groupID uint64, statuses []string) (*[]Tasks, int64) {
	var tasks *[]Tasks

	var rts []RetriveTask

	stmt := tx.Join("LEFT OUTER", "TASK_DETAIL", "TASK_DETAIL.TASK_ID = TASKS.ID")
	stmt = stmt.Join("LEFT OUTER", "TASK_LOGS", "TASK_LOGS.TASK_ID = TASKS.ID")

	cnt, err := stmt.Where("ZONE_ID = ?", groupID).And(builder.In("STATUS", statuses)).FindAndCount(&rts)

	if err != nil {
		panic(err)
	}

	logger.Debugf("selected retreive tasks : %d", cnt)

	tasks = toTasks(manager, &rts)

	for i, t := range *tasks {
		var steps []TaskSteps

		cnt, err := tx.Where("TASK_ID = ?", t.Id).OrderBy("SEQ ASC").FindAndCount(&steps)
		if err != nil {
			panic(err)
		}

		logger.Debugf("select steps for %d - %d", t.Id, cnt)

		for i := 0; i < len(steps); i++ {
			if steps[i].ReservedCommand != "" {
				steps[i].ReservedCommand = manager.decrypt(steps[i].ReservedCommand)
			}
			if steps[i].InlineScript != "" {
				steps[i].InlineScript = manager.decrypt(steps[i].InlineScript)
			}
		}

		(*tasks)[i].TaskSteps = &steps
	}

	logger.Debugf("tasks : [%+v]", tasks)

	return tasks, cnt
}

func toTasks(manager *KlevrManager, rts *[]RetriveTask) *[]Tasks {
	var tasks = make([]Tasks, 0, len(*rts))

	for _, rt := range *rts {
		rt.Tasks.TaskDetail = rt.TaskDetail
		rt.Tasks.Logs = rt.TaskLogs

		if rt.Tasks.TaskDetail.Parameter != "" {
			rt.Tasks.TaskDetail.Parameter = manager.decrypt(rt.Tasks.TaskDetail.Parameter)
		}
		if rt.Tasks.TaskDetail.Result != "" {
			rt.Tasks.TaskDetail.Result = manager.decrypt(rt.Tasks.TaskDetail.Result)
		}
		if rt.Tasks.Logs.Logs != "" {
			rt.Tasks.Logs.Logs = manager.decrypt(rt.Tasks.Logs.Logs)
		}

		tasks = append(tasks, *rt.Tasks)
	}

	return &tasks
}

func (tx *Tx) getTasks(groupIDs []uint64, statuses []string, agentKeys []string, taskNames []string) (*[]Tasks, bool) {
	var tasks []Tasks

	stmt := tx.Where(builder.In("ZONE_ID", groupIDs))

	// condition 추가
	if statuses != nil && !(len(statuses) == 1 && statuses[0] == "") {
		stmt = stmt.And(builder.In("STATUS", statuses))
	}
	if agentKeys != nil && !(len(statuses) == 1 && statuses[0] == "") {
		stmt = stmt.And(builder.In("AGENT_KEY", agentKeys))
	}
	if taskNames != nil && !(len(statuses) == 1 && statuses[0] == "") {
		stmt = stmt.And(builder.In("NAME", taskNames))
	}

	tx.Engine().ShowSQL(true)

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
	_, err := tx.Where("TASK_ID = ?").Delete(&TaskDetail{})
	if err != nil {
		panic(err)
	}

	_, err = tx.Where("TASK_ID = ?").Delete(&TaskLogs{})
	if err != nil {
		panic(err)
	}

	_, err = tx.Where("TASK_ID = ?").Delete(&TaskSteps{})
	if err != nil {
		panic(err)
	}
}

func (tx *Tx) getConsoleMember(userID string) (int64, *[]PageMembers) {
	var members []PageMembers

	cnt, err := tx.Where("USER_ID = ?", userID).FindAndCount(&members)
	if err != nil {
		panic(err)
	}

	return cnt, &members
}

func (tx *Tx) insertConsoleMember(p *PageMembers) *PageMembers {
	result, err := tx.Exec("INSERT INTO `PAGE_MEMBERS` (`user_id`,`user_password`, `activated`, `api_key`) VALUES (?,?,?,?)", p.UserId, p.UserPassword, p.Activated, p.ApiKey)
	if err != nil {
		panic(err)
	}

	result.RowsAffected()

	id, _ := result.LastInsertId()

	p.Id = uint64(id)

	logger.Debugf("Inserted task : %v", p)

	return p
}

func (tx *Tx) updateConsoleMember(p *PageMembers) {
	cnt, err := tx.Where("USER_ID = ?", p.UserId).Cols("ACTIVATED").Update(p)
	logger.Debugf("Updated PageMember(%d) : %v", cnt, p)

	if err != nil {
		panic(err)
	}
}

func (tx *Tx) deleteApiAuthentication(zoneID uint64) {
	sql := "delete a from API_AUTHENTICATIONS a where a.GROUP_ID = ?"
	res, err := tx.Exec(sql, zoneID)
	if err != nil {
		logger.Warningf("%+v", errors.Wrap(err, "sql error"))
		panic(err)
	}

	logger.Debug(res)
}

func (tx *Tx) getApiAuthenticationsByGroupId(groupID uint64) (int64, *[]ApiAuthentications) {
	var apiAuths []ApiAuthentications

	cnt, err := tx.Where("GROUP_ID = ?", groupID).FindAndCount(&apiAuths)
	if err != nil {
		panic(err)
	}

	return cnt, &apiAuths
}

func (tx *Tx) insertCredential(manager *KlevrManager, c *Credentials) *Credentials {
	c.Value = manager.encrypt(c.Value)

	result, err := tx.Exec("INSERT INTO `CREDENTIALS` (`zone_id`,`name`,`value`) VALUES (?,?,?)",
		c.ZoneId, c.Name, c.Value)

	if err != nil {
		panic(err)
	}

	result.RowsAffected()

	id, _ := result.LastInsertId()

	c.Id = uint64(id)

	logger.Debugf("Inserted credential : %v", c)

	return c
}

func (tx *Tx) getCredential(manager *KlevrManager, id uint64) (*Credentials, bool) {
	var credential Credentials

	exist := common.CheckGetQuery(tx.Where("CREDENTIALS.ID = ?", id).Get(&credential))
	logger.Debugf("Selected Credential : %v", credential)

	return &credential, exist
}

func (tx *Tx) getCredentialsByNames(groupIDs []uint64, credentialNames []string) (*[]Credentials, bool) {
	var credentials []Credentials

	stmt := tx.Where(builder.In("ZONE_ID", groupIDs))

	// condition 추가
	if credentialNames != nil {
		stmt = stmt.And(builder.In("NAME", credentialNames))
	}

	tx.Engine().ShowSQL(true)

	cnt, err := stmt.FindAndCount(&credentials)
	if err != nil {
		panic(err)
	}

	logger.Debugf("Selected Credentials : %d", cnt)

	return &credentials, cnt > 0
}

func (tx *Tx) getCredentials(groupID uint64) (*[]Credentials, int64) {
	var credentials []Credentials

	cnt, err := tx.Where("ZONE_ID = ?", groupID).FindAndCount(&credentials)
	if err != nil {
		panic(err)
	}

	return &credentials, cnt
}

func (tx *Tx) deleteCredential(id uint64) {
	_, err := tx.Where("ID = ?", id).Delete(&Credentials{})
	if err != nil {
		panic(err)
	}
}
