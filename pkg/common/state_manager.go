package common

import "sync"

type StateManager struct {
	*sync.Mutex
	state map[string]interface{}
}

var singleton *StateManager
var once sync.Once

func GetStateManager() *StateManager {
	once.Do(func() {
		singleton = &StateManager{Mutex: &sync.Mutex{}}
		singleton.state = make(map[string]interface{})
	})
	return singleton
}

func (sm *StateManager) GetState(key string) interface{} {
	sm.Lock()
	defer sm.Unlock()
	return sm.state[key]
}

func (sm *StateManager) SetState(key string, value interface{}) {
	sm.Lock()
	defer sm.Unlock()
	sm.state[key] = value
}

func (sm *StateManager) Clear() {
	sm.Lock()
	defer sm.Unlock()
	sm.state = make(map[string]interface{})
}
