package common

import "sync"

type stateManager struct {
	*sync.Mutex
	state map[string]interface{}
}

var singleton *stateManager
var once sync.Once

func GetStateManager() *stateManager {
	once.Do(func() {
		singleton = &stateManager{Mutex: &sync.Mutex{}}
		singleton.state = make(map[string]interface{})
	})
	return singleton
}

func (sm *stateManager) GetState(key string) interface{} {
	sm.Lock()
	defer sm.Unlock()
	return sm.state[key]
}

func (sm *stateManager) SetState(key string, value interface{}) {
	sm.Lock()
	defer sm.Unlock()
	sm.state[key] = value
}
