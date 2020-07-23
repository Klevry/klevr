package manager

import (
	"github.com/NexClipper/logger"
	"xorm.io/xorm"
)

func (api *API) getPrimaryAgent(zoneID uint) *PrimaryAgents {
	var pa PrimaryAgents

	// api.DB.Where(&PrimaryAgents{GroupId: zoneID}).First(&pa)

	return &pa
}

func (api *API) getAgent(agentKey string) *Agents {
	var a Agents

	// api.DB.Where(&Agents{AgentKey: agentKey}).First(&a)

	return &a
}

func getAgentGroup(conn *xorm.Session, zoneID uint) *AgentGroups {
	var m PrimaryAgents

	// api.DB.SetLogger(gorm.Logger{logger.})
	api.DB.Model(&PrimaryAgents{}).First(&m)
	api.DB.Debug().First(&m)

	logger.Debugf("%v", m)

	return &AgentGroups{}
}

func (api *API) addAgent(a *Agents) {
	api.DB.Create(a)
}
