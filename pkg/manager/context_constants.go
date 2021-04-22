package manager

import (
	"net/http"

	"github.com/gorilla/context"

	"github.com/Klevry/klevr/pkg/common"
)

// Context key constants
const (
	CtxPrefix = "CTX_"

	CtxRequestContext = CtxPrefix + "REQUEST"

	CtxServer    = CtxPrefix + "SERVER"
	CtxAPI       = CtxPrefix + "API"
	CtxDbConn    = CtxPrefix + "DB"
	CtxDbSession = CtxPrefix + "TX"
	CtxMqConn    = CtxPrefix + "MQConn"
	CtxMqChannel = CtxPrefix + "MQChannel"
	CtxMqQueue   = CtxPrefix + "MQQueue"
	CtxPrimary   = CtxPrefix + "PRIMARY"
)

// CtxGetServer get KlevrManager from context
func CtxGetServer(ctx *common.Context) *KlevrManager {
	return ctx.Get(CtxServer).(*KlevrManager)
}

// CtxGetDbConn get DB connection from context
func CtxGetDbConn(ctx *common.Context) *common.DB {
	return ctx.Get(CtxDbConn).(*common.DB)
}

// CtxGetDbSession get DB session from context
func CtxGetDbSession(ctx *common.Context) *Tx {
	return ctx.Get(CtxDbSession).(*Tx)
}

// CtxGetFromRequest get common.Context from request
func CtxGetFromRequest(r *http.Request) *common.Context {
	return context.Get(r, CtxRequestContext).(*common.Context)
}
