package manager

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/gorilla/mux"
)

// CustomHeader name constants
const (
	CHeaderAPIKey         = "X-API-KEY"
	CHeaderAgentKey       = "X-AGENT-KEY"
	CHeaderHashCode       = "X-HASH-CODE"
	CHeaderZoneID         = "X-ZONE-ID"
	CHeaderSupportVersion = "X-SUPPORT-AGENT-VERSION"
	CHeaderTimestamp      = "X-TIMESTAMP"
)

// InitAgent initialize agent API
func (api *API) InitAgent(agent *mux.Router) {
	logger.Debug("API InitAgent - init URI")

	// registURI(agent, PUT, "/handshake", api.receiveHandshake)
	// registURI(agent, PUT, "/:agentKey", api.receivePolling)

	registURI(agent, PUT, "/handshake", receiveHandshake)
	registURI(agent, PUT, "/{agentKey}", receivePolling)
	registURI(agent, GET, "/reports/{agentKey}", checkPrimaryInfo)
	registURI(agent, GET, "/commands/init", getInitCommand)
	registURI(agent, POST, "/zones/init", receiveInitResult)

	// agent API 핸들러 추가
	agent.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ch := parseCustomHeader(r)
			ctx := CtxGetFromRequest(r)

			// TODO: Support agent version 입력 추가

			// APIKey 인증
			logger.Debug(r.RequestURI)
			if !authenticate(ctx, ch.ZoneID, ch.APIKey) {
				return
			}

			// TODO: hash 검증 로직 추가

			next.ServeHTTP(w, r)

			// TODO: 전송구간 암호화 로직 추가

			// TODO: hash 생성 로직 추가

			// response header 설정
			h := w.Header()

			h.Set(CHeaderAgentKey, ch.AgentKey)
			h.Set(CHeaderHashCode, ch.HashCode)
			h.Set(CHeaderSupportVersion, ch.SupportVersion)
			h.Set(CHeaderTimestamp, string(time.Now().UTC().Unix()))
		})
	})
}

func authenticate(ctx *common.Context, zoneID uint64, apiKey string) bool {
	if !GetDBConn(ctx).existAPIKey(apiKey, zoneID) {
		panic(common.NewHTTPError(401, "authentication failed"))
	}

	return true
}

func parseCustomHeader(r *http.Request) *common.CustomHeader {
	zoneID, _ := strconv.ParseUint(strings.Join(r.Header.Values(CHeaderZoneID), ""), 10, 64)
	ts, _ := strconv.ParseInt(strings.Join(r.Header.Values(CHeaderTimestamp), ""), 10, 64)

	h := &common.CustomHeader{
		APIKey:         strings.Join(r.Header.Values(CHeaderAPIKey), ""),
		AgentKey:       strings.Join(r.Header.Values(CHeaderAgentKey), ""),
		HashCode:       strings.Join(r.Header.Values(CHeaderHashCode), ""),
		ZoneID:         uint64(zoneID),
		SupportVersion: strings.Join(r.Header.Values(CHeaderSupportVersion), ""),
		Timestamp:      ts,
	}

	ctx := *CtxGetFromRequest(r)
	ctx.Put(common.CustomHeaderName, h)

	return h
}

func receiveInitResult(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	ch := ctx.Get(common.CustomHeaderName).(*common.CustomHeader)
	tx := GetDBConn(ctx)

	var requestBody common.Body

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		common.WriteHTTPError(500, w, err, "JSON parsing error")
		return
	}

	agent := tx.getAgentByAgentKey(ch.AgentKey)

	tasks := requestBody.Task

	UpdateTaskStatus(tx, ch.ZoneID, &tasks)
	tx.Commit()

	// response 데이터 생성
	rb := &common.Body{}

	b, err := json.Marshal(rb)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)

	len := len(tasks)

	if len > 0 {
		task := tasks[0]
		result, _ := json.Marshal(task.Params)

		AddEvent(&KlevrEvent{
			EventType: PrimaryInit,
			AgentId:   agent.Id,
			GroupId:   agent.GroupId,
			EventTime: &JSONTime{time.Now().UTC()},
			Result:    string(result),
		})
	}
}

func UpdateTaskStatus(tx *Tx, zoneID uint64, tasks *[]common.Task) {
	len := len(*tasks)

	if len > 0 {
		for _, t := range *tasks {
			p, _ := json.Marshal(t.Params)

			param := &TaskParams{
				TaskId: t.ID,
				Params: string(p),
			}

			tx.updateTask(&Tasks{
				Id:          t.ID,
				ExeAgentKey: t.AgentKey,
				Status:      t.Status,
				Params:      param,
			})
		}
	}
}

func getInitCommand(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	ch := ctx.Get(common.CustomHeaderName).(*common.CustomHeader)
	tx := GetDBConn(ctx)

	url := "http://raw.githubusercontent.com/NexClipper/klevr_tasks/master/queue"
	logger.Debugf("%s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		common.WriteHTTPError(500, w, err, "Internel server error")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	res, err := (&http.Client{Transport: tr}).Do(req)
	if err != nil {
		common.WriteHTTPError(500, w, err, "Internel server error")
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	command := string(body)
	param := make(map[string]interface{})
	param["script"] = command

	// jsonParam, _ := json.Marshal(param)

	task := AddTask(tx, common.INLINE, "INIT", ch.ZoneID, ch.AgentKey, param)
	logger.Debug("created task : %v", task)

	// response 데이터 생성
	rb := &common.Body{}
	rb.Task = make([]common.Task, 1)
	rb.Task[0] = *common.NewTask(task.Id, common.TaskType(task.Type), command, task.AgentKey, task.Status, nil)

	b, err := json.Marshal(rb)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

func receiveHandshake(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	ch := ctx.Get(common.CustomHeaderName).(*common.CustomHeader)
	// var cr = &common.Request{r}

	tx := GetDBConn(ctx)
	var requestBody common.Body
	var paramAgent common.Me

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		common.WriteHTTPError(500, w, err, "JSON parsing error")
		return
	}

	paramAgent = requestBody.Me
	logger.Debug(fmt.Sprintf("CustomHeader : %v", ch))

	_, exist := tx.getAgentGroup(ch.ZoneID)

	if !exist {
		common.WriteHTTPError(400, w, nil, fmt.Sprintf("Does not exist zone for zoneId : %d", ch.ZoneID))
		return
	}

	agent := tx.getAgentByAgentKey(ch.AgentKey)

	// agent 생성 or 수정
	upsertAgent(tx, agent, ch, &paramAgent)

	tx.Commit()

	// response 데이터 생성
	rb := &common.Body{}

	// primary 조회
	rb.Agent.Primary = getPrimary(ctx, tx, ch.ZoneID, agent.Id)

	// 접속한 agent 정보
	me := &rb.Me

	me.HmacKey = agent.HmacKey
	me.EncKey = agent.EncKey
	// me.CallCycle = 0 // seconds
	// me.LogLevel = "DEBUG"

	// Primary agent인 경우 node 정보 추가
	if ch.AgentKey == rb.Agent.Primary.AgentKey {
		cnt, agents := tx.getAgentsByGroupId(ch.ZoneID)
		nodes := make([]common.Agent, cnt)

		if cnt > 0 {
			for i, a := range *agents {
				nodes[i] = common.Agent{
					AgentKey: a.AgentKey,
					IP:       a.Ip,
					Port:     a.Port,
					Version:  a.Version,
					Resource: &common.Resource{
						Core:   a.Cpu,
						Memory: a.Memory,
						Disk:   a.Disk,
					},
				}
			}

			rb.Agent.Nodes = nodes
		}
	}

	b, err := json.Marshal(rb)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)

	AddEvent(&KlevrEvent{
		EventType: AgentConnect,
		AgentId:   agent.Id,
		GroupId:   agent.GroupId,
		EventTime: &JSONTime{time.Now().UTC()},
	})
}

func receivePolling(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	ch := ctx.Get(common.CustomHeaderName).(*common.CustomHeader)
	// var cr = &common.Request{r}
	tx := GetDBConn(ctx)
	var param common.Body

	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		common.WriteHTTPError(500, w, err, "JSON parsing error")
		return
	}

	// response 데이터 생성
	rb := &common.Body{}

	// primary agent access 정보 갱신
	agent := updateAgentAccess(tx, ch.AgentKey)
	logger.Debugf("%v", agent)

	// TODO: primary agent 실행 파라미터 update
	// rb.Me.CallCycle =
	// rb.Me.LogLevel =

	// agent zone 상태 정보 업데이트
	nodes := param.Agent.Nodes
	len := len(nodes)
	arrAgent := make([]Agents, len)

	for i, agent := range nodes {
		arrAgent[i].AgentKey = agent.AgentKey
		arrAgent[i].LastAliveCheckTime = time.Unix(agent.LastAliveCheckTime, 0)
		arrAgent[i].IsActive = agent.IsActive
		arrAgent[i].Cpu = agent.Core
		arrAgent[i].Memory = agent.Memory
		arrAgent[i].Disk = agent.Disk
	}

	tx.updateZoneStatus(&arrAgent)
	tx.Commit()

	// TODO: 수행한 task 상태 정보 업데이트
	// for i, task := range param.Task {

	// }

	// 신규 task 할당

	b, err := json.Marshal(rb)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

func checkPrimaryInfo(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	ch := ctx.Get(common.CustomHeaderName).(*common.CustomHeader)
	// var cr = &common.Request{r}
	tx := GetDBConn(ctx)
	var param common.Body

	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		common.WriteHTTPError(500, w, err, "JSON parsing error")
		return
	}

	// response 데이터 생성
	rb := &common.Body{}

	// agent access 정보 갱신
	agent := updateAgentAccess(tx, ch.AgentKey)

	rb.Agent.Primary = getPrimary(ctx, tx, ch.ZoneID, agent.Id)

	b, err := json.Marshal(rb)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

func updateAgentAccess(tx *Tx, agentKey string) *Agents {
	agent := tx.getAgentByAgentKey(agentKey)

	// agent 접속 시간 갱신
	agent.LastAccessTime = time.Now().UTC()
	tx.updateAgent(agent)

	return agent
}

func getPrimary(ctx *common.Context, tx *Tx, zoneID uint64, agentID uint64) common.Primary {
	// primary agent 정보
	groupPrimary := tx.getPrimaryAgent(zoneID)
	var primaryAgent *Agents

	if groupPrimary.AgentId == 0 {
		primaryAgent = electPrimary(ctx, zoneID, agentID, false)
	} else {
		primaryAgent = tx.getAgentByID(groupPrimary.AgentId)

		logger.Debugf("primaryAgent : %+v", primaryAgent)

		if primaryAgent.Id == 0 || !primaryAgent.IsActive {
			primaryAgent = electPrimary(ctx, zoneID, agentID, true)

			logger.Debugf("changed primaryAgent : %+v", primaryAgent)
		}
	}

	return common.Primary{
		AgentKey:       primaryAgent.AgentKey,
		IP:             primaryAgent.Ip,
		Port:           primaryAgent.Port,
		IsActive:       primaryAgent.IsActive,
		LastAccessTime: primaryAgent.LastAccessTime.UTC().Unix(),
	}
}

// primary agent 선출
func electPrimary(ctx *common.Context, zoneID uint64, agentID uint64, oldDel bool) *Agents {
	logger.Debugf("electPrimary for %d", zoneID)

	var tx *Tx
	var agent *Agents

	common.Block{
		Try: func() {
			tx = &Tx{CtxGetDbConn(ctx).NewSession()}

			if oldDel {
				tx.deletePrimaryAgent(zoneID)
			}

			pa := &PrimaryAgents{
				GroupId: zoneID,
				AgentId: agentID,
			}

			cnt, err := tx.insertPrimaryAgent(pa)

			if err != nil {
				pa = tx.getPrimaryAgent(zoneID)
			} else if cnt != 1 {
				logger.Warning(fmt.Sprintf("insert primary agent cnt : %d", cnt))
				common.Throw(common.NewStandardError("elect primary failed."))
			}

			if pa.AgentId == 0 {
				logger.Debugf("invalid primary agent : %v", pa)
				common.Throw(common.NewStandardError("elect primary failed."))
			}

			agent = tx.getAgentByID(pa.AgentId)

			if agent.Id == 0 {
				logger.Warning(fmt.Sprintf("primary agent not exist for id : %d, [%v]", agent.Id, agent))
				common.Throw(common.NewStandardError("elect primary failed."))
			}

			tx.Commit()
		},
		Catch: func(e common.Exception) {
			if tx != nil {
				tx.Rollback()
			}

			logger.Warning(e)
			common.Throw(e)
		},
		Finally: func() {
			if tx != nil && !tx.IsClosed() {
				tx.Close()
			}
		},
	}.Do()

	return agent
}

func upsertAgent(tx *Tx, agent *Agents, ch *common.CustomHeader, paramAgent *common.Me) {
	if agent.AgentKey == "" { // 처음 접속하는 에이전트일 경우 신규 등록
		agent.AgentKey = ch.AgentKey
		agent.GroupId = ch.ZoneID
		agent.IsActive = true
		agent.LastAccessTime = time.Now().UTC()
		agent.Ip = paramAgent.IP
		agent.Port = paramAgent.Port
		agent.Cpu = paramAgent.Resource.Core
		agent.Memory = paramAgent.Resource.Memory
		agent.Disk = paramAgent.Resource.Disk
		agent.HmacKey = common.GetKey(16)
		agent.EncKey = common.GetKey(32)

		tx.addAgent(agent)
	} else { // 기존에 등록된 에이전트 재접속일 경우 접속 정보 업데이트
		agent.IsActive = true
		agent.LastAccessTime = time.Now().UTC()
		agent.Ip = paramAgent.IP
		agent.Port = paramAgent.Port
		agent.Cpu = paramAgent.Resource.Core
		agent.Memory = paramAgent.Resource.Memory
		agent.Disk = paramAgent.Resource.Disk
		agent.HmacKey = common.GetKey(16)
		agent.EncKey = common.GetKey(32)

		tx.updateAgent(agent)
	}
}
