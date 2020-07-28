package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/context"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/gorilla/mux"
)

// CustomHeaderName custom header name
const CustomHeaderName = "CTX-CUSTOM-HEADER"

const (
	CHeaderApiKey         = "X-API-KEY"
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

	registURI(agent, PUT, "/handshake", api.receiveHandshake)
	registURI(agent, PUT, "/{agentKey}", api.receivePolling)
	registURI(agent, PUT, "/reports/{agentKey}", api.checkPrimaryInfo)

	// agent API 핸들러 추가
	agent.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ch := parseCustomHeader(r)

			// TODO: Support agent version 입력 추가

			// APIKey 인증
			if !authenticate(w, r, ch.ZoneID, ch.APIKey) {
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

func authenticate(w http.ResponseWriter, r *http.Request, zoneID uint, apiKey string) bool {
	logger.Debug(r.RequestURI)

	if !existAPIKey(GetDBConn(r), apiKey, zoneID) {
		common.WriteHTTPError(401, w, nil, "authentication failed")
		return false
	}

	return true
}

func parseCustomHeader(r *http.Request) *CustomHeader {
	zoneID, _ := strconv.ParseUint(strings.Join(r.Header.Values(CHeaderZoneID), ""), 10, 64)
	ts, _ := strconv.ParseInt(strings.Join(r.Header.Values(CHeaderTimestamp), ""), 10, 64)

	h := &CustomHeader{
		APIKey:         strings.Join(r.Header.Values(CHeaderApiKey), ""),
		AgentKey:       strings.Join(r.Header.Values(CHeaderAgentKey), ""),
		HashCode:       strings.Join(r.Header.Values(CHeaderHashCode), ""),
		ZoneID:         uint(zoneID),
		SupportVersion: strings.Join(r.Header.Values(CHeaderSupportVersion), ""),
		Timestamp:      ts,
	}

	context.Set(r, CustomHeaderName, h)

	return h
}

func (api *API) receiveHandshake(w http.ResponseWriter, r *http.Request) {
	ch := getCustomHeader(r)
	// var cr = &common.Request{r}
	var conn = GetDBConn(r)
	var paramAgent Agent

	err := json.NewDecoder(r.Body).Decode(&paramAgent)
	if err != nil {
		common.WriteHTTPError(500, w, err, "JSON parsing error")
		return
	}

	logger.Debug(fmt.Sprintf("CustomHeader : %v", ch))

	group := getAgentGroup(conn, ch.ZoneID)

	if group.Id == 0 {
		common.WriteHTTPError(400, w, nil, fmt.Sprintf("Does not exist zone for zoneId : %d", ch.ZoneID))
		return
	}

	agent := getAgentByAgentKey(conn, ch.AgentKey)

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

		addAgent(conn, agent)
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

		updateAgent(conn, agent)
	}

	rb := &Body{}
	groupPrimary := getPrimaryAgent(conn, ch.ZoneID)

	// Primary agent 정보 반환
	if groupPrimary.AgentId != 0 {
		primaryAgent := getAgentByID(conn, groupPrimary.AgentId)

		rp := &rb.Agent.Primary

		rp.IP = primaryAgent.Ip
		rp.Port = primaryAgent.Port
		rp.IsActive = primaryAgent.IsActive
		rp.LastAccessTime = primaryAgent.LastAccessTime.UTC().Unix()
	} else {
		// TODO: primary agent 선정
	}

	me := &rb.Me

	me.HmacKey = agent.HmacKey
	me.EncKey = agent.EncKey
	// me.CallCycle = 0 // seconds
	// me.LogLevel = "DEBUG"

	b, err := json.Marshal(rb)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

func (api *API) receivePolling(w http.ResponseWriter, r *http.Request) {

}

func (api *API) checkPrimaryInfo(w http.ResponseWriter, r *http.Request) {

}

func electPrimary() *Agents {
	return nil
}
