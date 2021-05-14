package manager

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NexClipper/logger"
	"github.com/go-redis/redis"
)

type ICache interface {
	SyncAgent(zoneID uint64, agents *[]Agents) error
	UpdateZoneStatus(zoneID uint64, agents []Agents) error
	DeleteAllAgent(zoneID uint64) error
	UpdateAccessAgent(zoneID uint64, agentKey string, accessTime time.Time) error
	UpdateAgentDisabledStatus(inactiveAgents *[]Agents) error
	GetAgentsForInactive(before time.Time) (int64, *[]Agents)
	GetAgentsByZoneID(zoniID uint64) (int64, *[]Agents)
	GetAgentByID(zoneID, agentID uint64) *Agents
	GetAgentByAgentKey(zoneID uint64, agentKey string) *Agents
}

type Cache struct {
	client *redis.Client
}

func NewCache(address string, port int, password string) ICache {
	cache := &Cache{}
	a := fmt.Sprintf("%s:%d", address, port)
	cache.client = redis.NewClient(&redis.Options{
		Addr:     a,
		Password: password,
		DB:       0,
	})

	return cache
}

func (c *Cache) SyncAgent(zoneID uint64, agents *[]Agents) error {
	logger.Debugf("AddAgent: %v", agents)

	key := fmt.Sprintf("Agents.%d", zoneID)

	members, err := c.client.SMembers(key).Result()
	if err != nil {
		logger.Debug(err)
		return err
	}

	for _, agent := range *agents {
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
		if err := c.client.SAdd(key, buf).Err(); err != nil {
			logger.Debug(err)
		}
	}

	return nil
}

// UpdateZoneStatus에 전달된 agents의 정보는 Agents의 모든 필드를 채우고 있지 않음
// cache에 등록되어 있던 agent의 정보에 변경된 내역만 변경해서 다시 등록
func (c *Cache) UpdateZoneStatus(zoneID uint64, agents []Agents) error {
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
					logger.Debug("there are no members of the deleted cache")
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

func (c *Cache) DeleteAllAgent(zoneID uint64) error {
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

func (c *Cache) UpdateAccessAgent(zoneID uint64, agentKey string, accessTime time.Time) error {
	key := fmt.Sprintf("Agents.%d", zoneID)
	members, err := c.client.SMembers(key).Result()
	if err != nil {
		logger.Debugf("agentKey: %s", agentKey)
		logger.Debug(err)
		return err
	}

	agentBuf := make([]Agents, 0)
	var findAgent Agents
	for _, member := range members {
		var buf Agents
		json.Unmarshal([]byte(member), &buf)
		agentBuf = append(agentBuf, buf)
		if buf.AgentKey == agentKey {
			findAgent = buf
			cnt, err := c.client.SRem(key, member).Result()
			if err != nil {
				logger.Debug(err)
			}

			if cnt == 0 {
				logger.Debug("there are no members of the deleted cache")
			}

			break
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
func (c *Cache) UpdateAgentDisabledStatus(inactiveAgents *[]Agents) error {
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

			}

			buf.IsActive = boolToByte(false)
			item, _ := json.Marshal(buf)
			c.client.SAdd(key, item)
		}
	}
	return nil
}

// GetAgentsForInactive는 LastAccessTime이 before보다 이전에 접속했던 agent이거나
// IsActive가 false 상태인 agent를 찾아준다.
func (c *Cache) GetAgentsForInactive(before time.Time) (int64, *[]Agents) {
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

func (c *Cache) GetAgentsByZoneID(zoneID uint64) (int64, *[]Agents) {
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

func (c *Cache) GetAgentByID(zoneID, agentID uint64) *Agents {
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

func (c *Cache) GetAgentByAgentKey(zoneID uint64, agentKey string) *Agents {
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
