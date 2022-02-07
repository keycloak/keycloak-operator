package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateStatusSecondaryResources(t *testing.T) {
	// given
	var sd map[string][]string

	// when
	sd = UpdateStatusSecondaryResources(sd, "kind", "name-1")
	sd = UpdateStatusSecondaryResources(sd, "kind", "name-2")

	// then
	assert.Equal(t, map[string][]string{"kind": {"name-1", "name-2"}}, sd)
}

func TestDeleteFromStatusSecondaryResources(t *testing.T) {
	// given
	sd := map[string][]string{"kind": {"name-1", "name-2"}}

	// when
	DeleteFromStatusSecondaryResources(sd, "kind", "name-1")

	// then
	assert.Equal(t, map[string][]string{"kind": {"name-2"}}, sd)
}
