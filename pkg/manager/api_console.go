package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/gorilla/mux"
)

type ConsoleAPI struct{}

func (api *API) InitConsole(console *mux.Router) {
	logger.Debug("API InitConsole - init URI")

	tx := &Tx{api.DB.NewSession()}
	cnt, _ := tx.getConsoleMember("admin")
	if cnt == 0 {
		encPassword, err := common.Encrypt(api.Manager.Config.Server.EncryptionKey, "admin")
		if err == nil {
			p := &PageMembers{UserId: "admin", UserPassword: encPassword}
			tx.insertConsoleMember(p)
		} else {
			logger.Error(err)
		}
	}

	consoleAPI := &ConsoleAPI{}

	registURI(console, POST, "/signin", consoleAPI.SignIn)
	registURI(console, GET, "/signout", consoleAPI.SignOut)
	registURI(console, POST, "/changepassword", consoleAPI.ChangePassword)
	registURI(console, GET, "/activated/{id}", consoleAPI.Activated)
	registURI(console, DELETE, "/groups/{groupID}/agents/{agentKey}", consoleAPI.DeleteAgent)
}

// SignIn godoc
// @Summary SignIn
// @Description Klevr Console 사용자 SignIn.
// @Tags Console
// @Accept x-www-form-urlencoded
// @Produce json
// @Router /console/signin [post]
// @Param id formData string true "User ID"
// @Param pw formData string true "Current Password"
// @Success 200
func (api *ConsoleAPI) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	manager := CtxGetServer(ctx)

	id := r.FormValue("id")
	pw := r.FormValue("pw")

	if id != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cnt, pms := tx.getConsoleMember(id)
	if cnt == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	pm := (*pms)[0]
	decPassword, err := common.Decrypt(manager.Config.Server.EncryptionKey, pm.UserPassword)
	if err != nil || pw != decPassword {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	jwtHelper := common.NewJWTHelper([]byte(manager.Config.Console.Secret)).AddClaims("id", id).SetExpirationTime(expirationTime.Unix())
	tks, err := jwtHelper.GenToken()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp, err := json.Marshal(struct {
		Token string `json:"token"`
	}{
		tks,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "token", Value: tks, Expires: expirationTime})
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", resp)
}

// SignOut godoc
// @Summary Sign Out
// @Description Klevr Console 사용자 SignOut.
// @Tags Console
// @Accept json
// @Produce json
// @Router /console/signout [get]
// @Success 200
func (api *ConsoleAPI) SignOut(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now(),
		MaxAge:  -1,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(200)
}

// ChangePassword godoc
// @Summary Password 변경
// @Description Klevr Console 사용자의 패스워드를 변경한다.
// @Tags Console
// @Accept x-www-form-urlencoded
// @Produce json
// @Router /console/changepassword [post]
// @Param id formData string true "User ID"
// @Param pw formData string false "Current Password"
// @Param cpw formData string true "Confirmed Password"
// @Success 200
func (api *ConsoleAPI) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	manager := CtxGetServer(ctx)

	id := r.FormValue("id")
	pw := r.FormValue("pw")
	cpw := r.FormValue("cpw") // confirmed password

	if id != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cnt, pms := tx.getConsoleMember(id)
	if cnt == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	pm := (*pms)[0]
	if pm.Activated == true {
		decPassword, err := common.Decrypt(manager.Config.Server.EncryptionKey, pm.UserPassword)
		if err != nil || pw != decPassword {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	encPassword, err := common.Encrypt(manager.Config.Server.EncryptionKey, cpw)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	pm.UserPassword = encPassword
	pm.Activated = true
	tx.updateConsoleMember(&pm)

	w.WriteHeader(200)
}

// Activated godoc
// @Summary 사용자 활성화 상태
// @Description Klevr Console 사용자의 활성화 상태를 확인한다.
// @Tags Console
// @Accept json
// @Produce json
// @Router /console/activated/{id} [get]
// @Param id path string true "User ID"
// @Success 200 {object} string "{\"status\":activated/initialized}"
func (api *ConsoleAPI) Activated(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	vars := mux.Vars(r)
	userID := vars["id"]

	cnt, pms := tx.getConsoleMember(userID)
	if cnt == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	pm := (*pms)[0]
	var activatedStatus string
	if pm.Activated == true {
		activatedStatus = "activated"
	} else {
		activatedStatus = "initialized"
	}

	resp, err := json.Marshal(struct {
		Status string `json:"status"`
	}{
		activatedStatus,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", resp)
}

// DeleteAgent godoc
// @Summary Klevr Agent를 종료한다.
// @Description agentKey에 해당하는 Agent를 종료한다.
// @Tags Console
// @Accept json
// @Produce json
// @Router /console/groups/{groupID}/agents/{agentKey} [delete]
// @Param groupID path uint64 true "ZONE ID"
// @Param agentKey path string true "agent key"
// @Success 200 {object} string "{\"canceld\":true/false}"
func (api *ConsoleAPI) DeleteAgent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "{\"canceled\":%v}", true)

	/*ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	// groupID, agentKey
	qryGroupID := r.URL.Query()["groupID"]
	qryAgentKey := r.URL.Query()["agentKey"]

	groupID, err := strconv.ParseUint(qryGroupID[0], 0, 64)
	if err != nil {
		common.WriteHTTPError(400, w, err, fmt.Sprintf("invalid groupID - [%s]", qryGroupID[0]))
		return
	}

	var agentKey string
	if len(qryAgentKey) == 1 {
		agentKey = qryAgentKey[0]
	} else {
		common.WriteHTTPError(400, w, err, fmt.Sprint("invalid agentKey"))
		return
	}

	// agent 삭제를 위한 task를 생성
	t := common.KlevrTask{
		ZoneID:             groupID,
		Name:               "DeleteAgent",
		TaskType:           common.AtOnce,
		TotalStepCount:     1,
		Parameter:          "",
		AgentKey:           agentKey,
		ExeAgentChangeable: false,
		Steps: []*common.KlevrTaskStep{&common.KlevrTaskStep{
			Seq:         1,
			CommandName: "DeleteAgent",
			CommandType: common.RESERVED,
			Command:     "ForceShutdownAgent",
			IsRecover:   false,
		}},
		EventHookSendingType: common.EventHookWithAll,
	}

	// Task 상태 설정
	t = *common.TaskStatusAdd(&t)

	// DTO -> entity
	persistTask := *TaskDtoToPerist(&t)

	manager := ctx.Get(CtxServer).(*KlevrManager)

	// DB insert
	persistTask = *tx.insertTask(manager, &persistTask)

	task, _ := tx.getTask(manager, persistTask.Id)

	dto := TaskPersistToDto(task)

	b, err := json.Marshal(dto)
	if err != nil {
		panic(err)
	}

	logger.Debugf("response : [%s]", string(b))

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b) */
}
