package keycloakclient

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	v13 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeycloakClientReconciler_Test_Creating_Client(t *testing.T) {
	// given
	keycloakCr := v1alpha1.Keycloak{}
	cr := &v1alpha1.KeycloakClient{
		ObjectMeta: v13.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: v1alpha1.KeycloakClientSpec{
			RealmSelector: &v13.LabelSelector{
				MatchLabels: map[string]string{"application": "sso"},
			},
			Client: &v1alpha1.KeycloakAPIClient{
				ClientID: "test",
				Secret:   "test",
			},
		},
	}

	currentState := &common.ClientState{
		Realm: &v1alpha1.KeycloakRealm{
			Spec: v1alpha1.KeycloakRealmSpec{
				Realm: &v1alpha1.KeycloakAPIRealm{
					Realm: "test",
				},
			},
		},
	}

	// when
	reconciler := NewKeycloakClientReconciler(keycloakCr)
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.IsType(t, common.PingAction{}, desiredState[0])
	assert.IsType(t, common.CreateClientAction{}, desiredState[1])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[2])
	assert.IsType(t, model.ClientSecret(cr), desiredState[2].(common.GenericCreateAction).Ref)
	assert.Equal(t, []byte("test"), model.ClientSecret(cr).Data[model.ClientSecretClientIDProperty])
	assert.Equal(t, []byte("test"), model.ClientSecret(cr).Data[model.ClientSecretClientSecretProperty])
}

func TestKeycloakClientReconciler_Test_PartialUpdate_Client(t *testing.T) {
	// given
	keycloakCr := v1alpha1.Keycloak{}
	cr := &v1alpha1.KeycloakClient{
		ObjectMeta: v13.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: v1alpha1.KeycloakClientSpec{
			RealmSelector: &v13.LabelSelector{
				MatchLabels: map[string]string{"application": "sso"},
			},
			Client: &v1alpha1.KeycloakAPIClient{
				ClientID: "test",
				Secret:   "test",
			},
		},
	}

	currentState := &common.ClientState{
		Realm: &v1alpha1.KeycloakRealm{
			Spec: v1alpha1.KeycloakRealmSpec{
				Realm: &v1alpha1.KeycloakAPIRealm{
					Realm: "test",
				},
			},
		},
		Client: &v1alpha1.KeycloakAPIClient{
			Name: "dummy",
		},
	}

	// when
	reconciler := NewKeycloakClientReconciler(keycloakCr)
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.IsType(t, common.PingAction{}, desiredState[0])
	assert.IsType(t, common.UpdateClientAction{}, desiredState[1])

	// client secret still needs to be created even if the client exists
	assert.IsType(t, common.GenericCreateAction{}, desiredState[2])
	assert.IsType(t, model.ClientSecret(cr), desiredState[2].(common.GenericCreateAction).Ref)
	assert.Equal(t, []byte("test"), model.ClientSecret(cr).Data[model.ClientSecretClientIDProperty])
	assert.Equal(t, []byte("test"), model.ClientSecret(cr).Data[model.ClientSecretClientSecretProperty])
}

func TestKeycloakClientReconciler_Test_Delete_Client(t *testing.T) {
	// given
	keycloakCr := v1alpha1.Keycloak{}
	cr := &v1alpha1.KeycloakClient{
		ObjectMeta: v13.ObjectMeta{
			Name:      "test",
			Namespace: "test",
			DeletionTimestamp: &v13.Time{
				Time: time.Now(),
			},
		},
		Spec: v1alpha1.KeycloakClientSpec{
			RealmSelector: &v13.LabelSelector{
				MatchLabels: map[string]string{"application": "sso"},
			},
			Client: &v1alpha1.KeycloakAPIClient{
				ClientID: "test",
				Secret:   "test",
			},
		},
	}

	currentState := &common.ClientState{
		Realm: &v1alpha1.KeycloakRealm{
			Spec: v1alpha1.KeycloakRealmSpec{
				Realm: &v1alpha1.KeycloakAPIRealm{
					Realm: "test",
				},
			},
		},
	}

	// when
	reconciler := NewKeycloakClientReconciler(keycloakCr)
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.IsType(t, common.PingAction{}, desiredState[0])
	assert.IsType(t, common.DeleteClientAction{}, desiredState[1])
}

func TestKeycloakClientReconciler_Test_Update_Client(t *testing.T) {
	// given
	keycloakCr := v1alpha1.Keycloak{}
	cr := &v1alpha1.KeycloakClient{
		ObjectMeta: v13.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: v1alpha1.KeycloakClientSpec{
			RealmSelector: &v13.LabelSelector{
				MatchLabels: map[string]string{"application": "sso"},
			},
			Client: &v1alpha1.KeycloakAPIClient{
				ClientID:                     "test",
				Secret:                       "test",
				AuthorizationServicesEnabled: true,
				DefaultClientScopes:          []string{"profile"},
				OptionalClientScopes:         []string{"email"},
			},
			Roles: []v1alpha1.RoleRepresentation{
				{ID: "delete_recreateID2", Name: "delete_recreate"},
				{ID: "renameID", Name: "rename_new"},
				{ID: "rename_recreateID", Name: "rename_recreate_new"},
				{Name: "update", Description: "update_description"},
				{Name: "rename_recreate"},
			},
			ScopeMappings: &v1alpha1.MappingsRepresentation{
				ClientMappings: map[string]v1alpha1.ClientMappingsRepresentation{"someclient": {Mappings: []v1alpha1.RoleRepresentation{{Name: "b"}, {Name: "c"}}}},
				RealmMappings:  []v1alpha1.RoleRepresentation{{Name: "rb"}, {Name: "rc"}},
			},
		},
	}

	currentState := &common.ClientState{
		Client:       &v1alpha1.KeycloakAPIClient{},
		ClientSecret: &v1.Secret{},
		Realm: &v1alpha1.KeycloakRealm{
			Spec: v1alpha1.KeycloakRealmSpec{
				Realm: &v1alpha1.KeycloakAPIRealm{
					Realm: "test",
				},
			},
		},
		Roles: []v1alpha1.RoleRepresentation{
			{ID: "deleteID", Name: "delete"},
			{ID: "delete_recreateID", Name: "delete_recreate"},
			{ID: "updateID", Name: "update"},
			{ID: "renameID", Name: "rename"},
			{ID: "rename_recreateID", Name: "rename_recreate"},
			{Name: umaRoleName},
		},
		ScopeMappings: &v1alpha1.MappingsRepresentation{
			ClientMappings: map[string]v1alpha1.ClientMappingsRepresentation{"someclient": {Mappings: []v1alpha1.RoleRepresentation{{Name: "a"}, {Name: "b"}}}},
			RealmMappings:  []v1alpha1.RoleRepresentation{{Name: "ra"}, {Name: "rb"}},
		},
		AvailableClientScopes: []v1alpha1.KeycloakClientScope{{Name: "address", ID: "222"}, {Name: "email", ID: "421"}, {Name: "profile", ID: "314"}},
		DefaultClientScopes:   []v1alpha1.KeycloakClientScope{},
		OptionalClientScopes:  []v1alpha1.KeycloakClientScope{{Name: "address", ID: "222"}},
	}

	// when
	reconciler := NewKeycloakClientReconciler(keycloakCr)
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.IsType(t, common.PingAction{}, desiredState[0])
	assert.IsType(t, common.UpdateClientAction{}, desiredState[1])
	assert.Equal(t, "test", desiredState[1].(common.UpdateClientAction).Realm)
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[2])
	assert.IsType(t, model.ClientSecretReconciled(cr, currentState.ClientSecret), desiredState[2].(common.GenericUpdateAction).Ref)
	assert.Equal(t, []byte("test"), model.ClientSecretReconciled(cr, currentState.ClientSecret).Data[model.ClientSecretClientIDProperty])
	assert.Equal(t, []byte("test"), model.ClientSecretReconciled(cr, currentState.ClientSecret).Data[model.ClientSecretClientSecretProperty])

	assert.IsType(t, common.DeleteClientRoleAction{}, desiredState[3])
	assert.Equal(t, "delete", desiredState[3].(common.DeleteClientRoleAction).Role.Name)
	assert.IsType(t, common.DeleteClientRoleAction{}, desiredState[4])
	assert.Equal(t, "delete_recreate", desiredState[4].(common.DeleteClientRoleAction).Role.Name)
	assert.IsType(t, common.UpdateClientRoleAction{}, desiredState[5])
	assert.Equal(t, "rename_new", desiredState[5].(common.UpdateClientRoleAction).Role.Name)
	assert.Equal(t, "rename", desiredState[5].(common.UpdateClientRoleAction).OldRole.Name)
	assert.IsType(t, common.UpdateClientRoleAction{}, desiredState[6])
	assert.Equal(t, "rename_recreate_new", desiredState[6].(common.UpdateClientRoleAction).Role.Name)
	assert.Equal(t, "rename_recreate", desiredState[6].(common.UpdateClientRoleAction).OldRole.Name)
	assert.IsType(t, common.UpdateClientRoleAction{}, desiredState[7])
	assert.Equal(t, "update", desiredState[7].(common.UpdateClientRoleAction).Role.Name)
	assert.Equal(t, "update_description", desiredState[7].(common.UpdateClientRoleAction).Role.Description)
	assert.IsType(t, common.CreateClientRoleAction{}, desiredState[8])
	assert.Equal(t, "rename_recreate", desiredState[8].(common.CreateClientRoleAction).Role.Name)
	assert.IsType(t, common.CreateClientRoleAction{}, desiredState[9])
	assert.Equal(t, "delete_recreate", desiredState[9].(common.CreateClientRoleAction).Role.Name)

	assert.IsType(t, common.CreateClientRealmScopeMappingsAction{}, desiredState[10])
	assert.IsType(t, common.CreateClientClientScopeMappingsAction{}, desiredState[11])
	assert.Equal(t, "someclient", desiredState[11].(common.CreateClientClientScopeMappingsAction).Mappings.Client)
	assert.IsType(t, common.DeleteClientRealmScopeMappingsAction{}, desiredState[12])
	assert.IsType(t, common.DeleteClientClientScopeMappingsAction{}, desiredState[13])

	assert.IsType(t, common.UpdateClientDefaultClientScopeAction{}, desiredState[14])
	assert.Equal(t, "314", desiredState[14].(common.UpdateClientDefaultClientScopeAction).ClientScope.ID)
	assert.IsType(t, common.UpdateClientOptionalClientScopeAction{}, desiredState[15])
	assert.Equal(t, "421", desiredState[15].(common.UpdateClientOptionalClientScopeAction).ClientScope.ID)
	assert.IsType(t, common.DeleteClientOptionalClientScopeAction{}, desiredState[16])
	assert.Equal(t, "222", desiredState[16].(common.DeleteClientOptionalClientScopeAction).ClientScope.ID)

	assert.Equal(t, 17, len(desiredState))
}

func TestKeycloakClientReconciler_Test_Marshal_Client(t *testing.T) {
	// given
	cr := &v1alpha1.KeycloakClient{
		ObjectMeta: v13.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: v1alpha1.KeycloakClientSpec{
			RealmSelector: &v13.LabelSelector{
				MatchLabels: map[string]string{"application": "sso"},
			},
			Client: &v1alpha1.KeycloakAPIClient{
				ClientID:                  "test",
				Secret:                    "test",
				DirectAccessGrantsEnabled: false,
				StandardFlowEnabled:       false,
				PublicClient:              false,
			},
		},
	}

	// when
	b, err := json.Marshal(cr)
	s := string(b)

	// then
	assert.Nil(t, err, "Client couldn't be marshalled")
	assert.True(t, strings.Contains(s, "\"directAccessGrantsEnabled\":false"), "Element directAccessGrantsEnabled should not be omitted if false, as keycloaks default is true")
	assert.True(t, strings.Contains(s, "\"standardFlowEnabled\":false"), "Element standardFlowEnabled should not be omitted if false, as keycloaks default is true")
	assert.True(t, strings.Contains(s, "\"publicClient\":false"), "Element publicClient should not be omitted if false, as keycloaks default is true")
}

func TestKeycloakClientReconciler_Test_ScopeMapping_Difference(t *testing.T) {
	// given
	a := &v1alpha1.MappingsRepresentation{
		RealmMappings: []v1alpha1.RoleRepresentation{{Name: "realmRoleA"}, {Name: "realmRoleB"}},
		ClientMappings: map[string]v1alpha1.ClientMappingsRepresentation{
			"clientA":  {Mappings: []v1alpha1.RoleRepresentation{{Name: "allDeleted"}}},
			"clientB1": {Mappings: []v1alpha1.RoleRepresentation{{Name: "a"}, {Name: "b"}}, ID: "idB1"},
			"clientB2": {Mappings: []v1alpha1.RoleRepresentation{{Name: "a"}, {Name: "b"}}},
		},
	}
	b := &v1alpha1.MappingsRepresentation{
		RealmMappings: []v1alpha1.RoleRepresentation{{Name: "realmRoleB"}, {Name: "realmRoleC"}},
		ClientMappings: map[string]v1alpha1.ClientMappingsRepresentation{
			"clientA":  {Mappings: []v1alpha1.RoleRepresentation{{Name: "allDeleted"}}},
			"clientB1": {Mappings: []v1alpha1.RoleRepresentation{{Name: "b"}, {Name: "c"}}},
			"clientB2": {Mappings: []v1alpha1.RoleRepresentation{{Name: "b"}, {Name: "c"}}, ID: "idB2"},
			"clientC":  {Mappings: []v1alpha1.RoleRepresentation{{Name: "irrelevant"}}},
		},
	}

	// when
	d := scopeMappingDifference(a, b)

	// then
	expected := &v1alpha1.MappingsRepresentation{
		RealmMappings: []v1alpha1.RoleRepresentation{{Name: "realmRoleA"}},
		ClientMappings: map[string]v1alpha1.ClientMappingsRepresentation{
			"clientB1": {Client: "clientB1", ID: "idB1", Mappings: []v1alpha1.RoleRepresentation{{Name: "a"}}},
			"clientB2": {Client: "clientB2", ID: "idB2", Mappings: []v1alpha1.RoleRepresentation{{Name: "a"}}},
		},
	}
	assert.Equal(t, expected, d)

	_, ok := d.ClientMappings["clientA"]
	assert.False(t, ok)
}
