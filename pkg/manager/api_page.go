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

type PageAPI struct{}

func (api *API) InitPage(page *mux.Router) {
	logger.Debug("API InitPage - init URI")

	tx := &Tx{api.DB.NewSession()}
	cnt, _ := tx.getPageMember("admin")
	if cnt == 0 {
		encPassword, err := common.Encrypt(api.Manager.Config.Server.EncryptionKey, "admin")
		if err == nil {
			p := &PageMembers{UserId: "admin", UserPassword: encPassword}
			tx.insertPageMember(p)
		} else {
			logger.Error(err)
		}
	}

	pageAPI := &PageAPI{}

	registURI(page, POST, "/signin", pageAPI.SignIn)
	registURI(page, GET, "/signout", pageAPI.SignOut)
	registURI(page, POST, "/changepassword", pageAPI.ChangePassword)
}

func (api *PageAPI) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	manager := CtxGetServer(ctx)

	id := r.FormValue("id")
	pw := r.FormValue("pw")

	if id != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cnt, pms := tx.getPageMember(id)
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
	jwtHelper := common.NewJWTHelper([]byte(manager.Config.Page.Secret)).AddClaims("id", id).SetExpirationTime(expirationTime.Unix())
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

func (api *PageAPI) SignOut(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now(),
		MaxAge:  -1,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(200)
}

func (api *PageAPI) ChangePassword(w http.ResponseWriter, r *http.Request) {
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

	cnt, pms := tx.getPageMember(id)
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
	tx.updatePageMember(&pm)

	w.WriteHeader(200)
}
