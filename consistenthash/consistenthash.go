package consistenthash

import (
	"sync"
	"sort"
	"crypto/md5"
	jump "github.com/renstrom/go-jump-consistent-hash"
	"encoding/hex"
	"fmt"
)

type Map struct {
	Mutex    sync.RWMutex
	KeyAlive map[string]bool
	KeyMap   map[string]bool
}

func HashNew(nodes []string) *Map {
	m := &Map{
		Mutex:    sync.RWMutex{},
		KeyMap:   make(map[string]bool),
		KeyAlive:   make(map[string]bool),
	}
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	for _, key := range nodes {
		m.KeyMap[key] = false
	}
	return m
}

// keys 活跃节点.
func (m *Map) Add(keys ...string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	// 更新活跃状态
	for k, _ := range m.KeyMap {
		m.KeyMap[k] = false
		delete(m.KeyAlive, k)
	}

	for _, k := range keys {
		m.KeyMap[k] = true
		m.KeyAlive[k] = true
	}
}

func (m *Map) IsEmpty()  (bool){
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	if len(m.KeyMap) == 0 {
		return true
	}
	return false
}

func (m *Map) Get(key string) (host string, err error){
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	srcKeys := getMapKeys(m.KeyMap)
	idx := int32(0)
	if idx, host, err = jumpHash(key, srcKeys); err != nil {
		return host, err
	}

	if m.KeyMap[host] {
		return host,  nil
	}

	keyAlive := getMapKeys(m.KeyAlive)
	newKey := MD5(key)
	idxMove, moveHost, err := jumpHash(newKey, keyAlive)
	if err != nil {
		fmt.Errorf("err:%v srcHost:%v->moveHost:%v key:%v new_key:%v idx:%v idxMove:%v", err, host, moveHost, key, newKey, idx, idxMove)
		return moveHost, err
	}
	return moveHost, nil
}

func getMapKeys(srcMap map[string]bool) (keys []string) {
	for k, _ := range srcMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func jumpHash(key string, hosts []string) (idx int32, host string, err error){
	sort.Strings(hosts)
	hostsLen := len(hosts)
	if hostsLen <= 0 {
		return idx,"", fmt.Errorf("key: %v hosts:%v is %v", key,  hosts, hostsLen)
	}
	idx = jump.HashString(key, int32(len(hosts)), jump.CRC32)
	host = hosts[idx]
	return idx, host, nil
}

func (m *Map) Update(key string, stats bool) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	m.KeyMap[key] = stats
	if stats {
		m.KeyAlive[key] = true
	}else {
		delete(m.KeyAlive, key)
	}
}

func (m *Map) Exist(key string) (bool){
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	_, ok := m.KeyAlive[key]
	return ok
}

func (m *Map) GetStats() interface{} {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	statsInfo := make(map[string]interface{})
	statsInfo["KeyMap"] = m.KeyMap
	statsInfo["KeyAlive"] = m.KeyAlive
	return statsInfo
}

func MD5(text string) string{
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

