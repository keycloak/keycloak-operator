package common

import "sync"

type stateManager struct {
	state map[string]interface{}
}

var singleton *stateManager
var once sync.Once

func GetStateManager() *stateManager {
	once.Do(func() {
		singleton = &stateManager{}
		singleton.state = make(map[string]interface{})
	})
	return singleton
}
func (sm *stateManager) GetState(key string) interface{} {
	return sm.state[key]
}
func (sm *stateManager) SetState(key string, value interface{}) {
	sm.state[key] = value
}
