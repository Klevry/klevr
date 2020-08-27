package common

import "context"

var (
	// AppContext application context
	appContext = context.Background()
)

func GetAppContext() *context.Context {
	return &appContext
}
