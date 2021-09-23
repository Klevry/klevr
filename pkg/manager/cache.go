package manager

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/go-redis/redis"
)

type Cache struct {
	client *redis.Client
}

func NewCache(address string, port int, password string) *Cache {
	cache := &Cache{}
	a := fmt.Sprintf("%s:%d", address, port)
	cache.client = redis.NewClient(&redis.Options{
		Addr:     a,
		Password: password,
		DB:       0,
	})

	err := cache.client.Ping().Err()
	if err != nil {
		logger.Error(err)
		return nil
	}

	return cache
}

func (c *Cache) Close() error {
	return c.client.Close()
}

func (c *Cache) UpdateAgent(ctx *common.Context, zoneID uint64, agent *Agents) error {
	mtx := ctx.Get(CtxCacheLock).(*sync.Mutex)
	mtx.Lock()
	defer mtx.Unlock()

	key := fmt.Sprintf("Agents.%d", zoneID)

	members, err := c.client.SMembers(key).Result()
	if err != nil {
		logger.Debug(err)
		return err
	}

	for _, member := range members {
		var buf Agents
		json.Unmarshal([]byte(member), &buf)
		if buf.AgentKey == agent.AgentKey {
			cnt, err := c.client.SRem(key, member).Result()
			if err != nil {
				logger.Debug(err)
			}
			if cnt == 0 {
				logger.Debug("there are no members of the deleted cache")
			}

		}
	}

	buf, err := json.Marshal(agent)
	if err != nil {
		logger.Debug(err)
	}

	logger.Debugf("##### UpdateAgent: SADD(%s): %s", key, string(buf))
	if err := c.client.SAdd(key, buf).Err(); err != nil {
		logger.Debug(err)
	}

	return nil
}

// UpdateZoneStatus에 전달된 agents의 정보는 Agents의 모든 필드를 채우고 있지 않음
// cache에 등록되어 있던 agent의 정보에 변경된 내역만 변경해서 다시 등록
func (c *Cache) UpdateZoneStatus(ctx *common.Context, zoneID uint64, agents []Agents) error {
	mtx := ctx.Get(CtxCacheLock).(*sync.Mutex)
	mtx.Lock()
	defer mtx.Unlock()

	/*
		logger.Debug("######### UpdatezoneStatus start ##########")
		defer logger.Debug("######### UpdateZoneStatus end ##########") //*/

	key := fmt.Sprintf("Agents.%d", zoneID)
	members, err := c.client.SMembers(key).Result()
	if err != nil {
		logger.Debug(err)
		return err
	}

	agentBuf := make([]Agents, 0)
	for _, member := range members {
		var buf Agents
		json.Unmarshal([]byte(member), &buf)
		agentBuf = append(agentBuf, buf)

		for _, agent := range agents {
			if buf.AgentKey == agent.AgentKey {
				cnt, err := c.client.SRem(key, member).Result()
				if err != nil {
					logger.Debug(err)
				}

				if cnt == 0 {
					logger.Debugf("there are no members of the deleted cache : %v", string(member))
				}
			}
		}
	}

	for _, agent := range agents {
		for _, a := range agentBuf {
			// 일치하는 agent를 찾아서 제거하고 새로운 agent를 추가한다.
			if a.AgentKey == agent.AgentKey {
				a.LastAliveCheckTime = agent.LastAliveCheckTime
				a.Cpu = agent.Cpu
				a.Memory = agent.Memory
				a.Disk = agent.Disk
				a.FreeMemory = agent.FreeMemory
				a.FreeDisk = agent.FreeDisk
				a.IsActive = agent.IsActive

				updateItem, _ := json.Marshal(a)
				if err := c.client.SAdd(key, updateItem).Err(); err != nil {
					logger.Debug(err)
				}

				break
			}
		}
	}

	return nil
}

func (c *Cache) DeleteAllAgent(ctx *common.Context, zoneID uint64) error {
	mtx := ctx.Get(CtxCacheLock).(*sync.Mutex)
	mtx.Lock()
	defer mtx.Unlock()

	/*
		logger.Debug("######### DeleteAllAgent start ##########")
		defer logger.Debug("######### DeleteAllAgent end ##########") //*/

	key := fmt.Sprintf("Agents.%d", zoneID)
	members, err := c.client.SMembers(key).Result()
	if err != nil {
		logger.Debugf("zoneID: %d", zoneID)
		logger.Debug(err)
		return err
	}

	for _, member := range members {
		cnt, err := c.client.SRem(key, member).Result()
		if err != nil {
			logger.Debug(err)
		}

		if cnt == 0 {
			logger.Debug("there are no members of the deleted cache")
		}
	}

	return nil
}

func (c *Cache) UpdateAccessAgent(ctx *common.Context, zoneID uint64, agentKey string, accessTime time.Time) error {
	mtx := ctx.Get(CtxCacheLock).(*sync.Mutex)
	mtx.Lock()
	defer mtx.Unlock()

	/*
		logger.Debug("######### UpdateAccessAgent start ##########")
		defer logger.Debug("######### UpdateAccessAgent end ##########") //*/

	key := fmt.Sprintf("Agents.%d", zoneID)
	members, err := c.client.SMembers(key).Result()
	if err != nil {
		logger.Debugf("agentKey: %s", agentKey)
		logger.Debug(err)
		return err
	}

	var findAgent Agents
	originLength := len(members)
	//logger.Debugf("### UpdateAccessAgent originLenght: %d", originLength)

	if originLength > 0 {
		for _, member := range members {
			var buf Agents
			json.Unmarshal([]byte(member), &buf)
			if buf.AgentKey == agentKey {
				//logger.Debugf("######### UpdateAccessAgent(SREM): %v", buf)
				findAgent = buf
				cnt, err := c.client.SRem(key, member).Result()
				if err != nil {
					logger.Debug(err)
				}

				if cnt == 0 {
					logger.Debugf("there are no members of the deleted cache: %v", string(member))
				}

				//logger.Debugf("######### UpdateAccessAgent(SREM cnt): %d", cnt)
				break
			}
		}
	}

	if findAgent.AgentKey == agentKey {
		findAgent.LastAccessTime = accessTime
		findAgent.IsActive = boolToByte(true)

		updateItem, _ := json.Marshal(findAgent)
		if err := c.client.SAdd(key, updateItem).Err(); err != nil {
			logger.Debug(err)
			return err
		}
	}

	return nil
}

// agent들을 inactive 상태로 변경
func (c *Cache) UpdateAgentDisabledStatus(ctx *common.Context, inactiveAgents *[]Agents) error {
	mtx := ctx.Get(CtxCacheLock).(*sync.Mutex)
	mtx.Lock()
	defer mtx.Unlock()

	/*
		logger.Debug("######### UpdateAgentDisabledStatus start ##########")
		defer logger.Debug("######### UpdateAgentDisabledStatus end ##########") //*/

	for _, a := range *inactiveAgents {
		key := fmt.Sprintf("Agents.%d", a.GroupId)
		members, err := c.client.SMembers(key).Result()
		if err != nil {
			logger.Debug(err)
			return err
		}

		for _, member := range members {
			var buf Agents
			json.Unmarshal([]byte(member), &buf)
			if buf.AgentKey == a.AgentKey {
				cnt, err := c.client.SRem(key, member).Result()
				if err != nil {
					logger.Debug(err)
				}

				if cnt == 0 {
					logger.Debug("there are no members of the deleted cache")
				}

				buf.IsActive = boolToByte(false)
				item, _ := json.Marshal(buf)
				c.client.SAdd(key, item)
			}

		}
	}
	return nil
}

// GetAgentsForInactive는 LastAccessTime이 before보다 이전에 접속했던 agent이거나
// IsActive가 false 상태인 agent를 찾아준다.
func (c *Cache) GetAgentsForInactive(ctx *common.Context, before time.Time) (int64, *[]Agents) {
	mtx := ctx.Get(CtxCacheLock).(*sync.Mutex)
	mtx.Lock()
	defer mtx.Unlock()

	/*
		logger.Debug("######### GetAgentsForInactive start ##########")
		defer logger.Debug("######### GetAgentsForInactive end ##########") //*/

	keys, err := c.client.Keys("Agents.*").Result()
	if err != nil {
		logger.Debug(err)
		return 0, nil
	}

	inactivedAgents := make([]Agents, 0)
	for _, k := range keys {
		members, err := c.client.SMembers(k).Result()
		if err != nil {
			logger.Debug(err)
			continue
		}

		for _, member := range members {
			var buf Agents
			json.Unmarshal([]byte(member), &buf)

			if byteToBool(buf.IsActive) == true {
				if res := buf.LastAccessTime.Before(before); res == true {
					inactivedAgents = append(inactivedAgents, buf)
				}
			} else {
				inactivedAgents = append(inactivedAgents, buf)
			}
		}
	}
	return int64(len(inactivedAgents)), &inactivedAgents
}

func (c *Cache) GetAgentsByZoneID(ctx *common.Context, zoneID uint64) (int64, *[]Agents) {
	mtx := ctx.Get(CtxCacheLock).(*sync.Mutex)
	mtx.Lock()
	defer mtx.Unlock()

	/*
		logger.Debug("######### GetAgentsByZoneID start ##########")
		defer logger.Debug("######### GetAgentsByZoneID end ##########")
		//*/

	key := fmt.Sprintf("Agents.%d", zoneID)
	members, err := c.client.SMembers(key).Result()
	if err != nil {
		logger.Debug(err)
		return 0, nil
	}

	agents := make([]Agents, 0)
	for _, member := range members {
		var buf Agents
		json.Unmarshal([]byte(member), &buf)
		agents = append(agents, buf)
	}

	return int64(len(agents)), &agents
}

func (c *Cache) GetAgentByID(ctx *common.Context, zoneID, agentID uint64) *Agents {
	mtx := ctx.Get(CtxCacheLock).(*sync.Mutex)
	mtx.Lock()
	defer mtx.Unlock()

	/*
		logger.Debug("######### GetAgentByID start ##########")
		defer logger.Debug("######### GetAgentByID end ##########") //*/

	key := fmt.Sprintf("Agents.%d", zoneID)
	members, err := c.client.SMembers(key).Result()
	if err != nil {
		logger.Debug(err)
		return nil
	}

	for _, member := range members {
		var buf Agents
		json.Unmarshal([]byte(member), &buf)
		if buf.Id == agentID {
			return &buf
		}
	}

	return nil
}

func (c *Cache) GetAgentByAgentKey(ctx *common.Context, zoneID uint64, agentKey string) *Agents {
	mtx := ctx.Get(CtxCacheLock).(*sync.Mutex)
	mtx.Lock()
	defer mtx.Unlock()

	/*
		logger.Debug("######### GetAgentByAgentKey start ##########")
		defer logger.Debug("######### GetAgentByAgentKey end ##########")
		//*/

	key := fmt.Sprintf("Agents.%d", zoneID)
	members, err := c.client.SMembers(key).Result()
	if err != nil {
		logger.Debug(err)
		return nil
	}

	for _, member := range members {
		var buf Agents
		json.Unmarshal([]byte(member), &buf)
		if buf.AgentKey == agentKey {
			return &buf
		}
	}

	return nil
}
