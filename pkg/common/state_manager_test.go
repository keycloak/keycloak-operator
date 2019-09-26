package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateManager_Test_(t *testing.T) {
	// given
	stateManager := GetStateManager()
	stateManagerTwo := GetStateManager()

	// when
	stateManager.SetState("Test", "string")

	// then
	assert.Nil(t, stateManager.GetState("NotSet"))
	assert.Equal(t, stateManager.GetState("Test"), "string")
	assert.Equal(t, stateManager, stateManagerTwo)
}
