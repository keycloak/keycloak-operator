package keycloakclient

import (
	"fmt"

	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model"
)

const (
	umaRoleName = "uma_protection"
)

type Reconciler interface {
	Reconcile(cr *kc.KeycloakClient) error
}

type KeycloakClientReconciler struct { // nolint
	Keycloak kc.Keycloak
}

func NewKeycloakClientReconciler(keycloak kc.Keycloak) *KeycloakClientReconciler {
	return &KeycloakClientReconciler{
		Keycloak: keycloak,
	}
}

func (i *KeycloakClientReconciler) Reconcile(state *common.ClientState, cr *kc.KeycloakClient) common.DesiredClusterState {
	desired := common.DesiredClusterState{}

	desired.AddAction(i.pingKeycloak())
	if cr.DeletionTimestamp != nil {
		desired.AddAction(i.getDeletedClientState(state, cr))
		return desired
	}

	if state.Client == nil {
		desired.AddAction(i.getCreatedClientState(state, cr))
	} else {
		desired.AddAction(i.getUpdatedClientState(state, cr))
	}

	if state.ClientSecret == nil {
		desired.AddAction(i.getCreatedClientSecretState(state, cr))
	} else {
		desired.AddAction(i.getUpdatedClientSecretState(state, cr))
	}

	i.ReconcileRoles(state, cr, &desired)

	i.ReconcileScopeMappings(state, cr, &desired)

	i.ReconcileClientScopes(state, cr, &desired)

	return desired
}

func (i *KeycloakClientReconciler) ReconcileRoles(state *common.ClientState, cr *kc.KeycloakClient, desired *common.DesiredClusterState) {
	// delete existing roles for which no desired role is found that (matches by ID OR has no ID but matches by name)
	// this implies that specifying a role with matching name but different ID will result in deletion (and re-creation)
	rolesDeleted, _ := model.RoleDifferenceIntersection(state.Roles, cr.Spec.Roles)
	// Prevent uma_protection role from deletion when fine-grained authorization support is enabled but
	// not present in the CR. This role is automatically created by Keycloak for AuthZ - read more below:
	// https://www.keycloak.org/docs/latest/authorization_services/#_service_protection_whatis_obtain_pat
	// TODO: evaluate sync options (once available) for uma_protection role once implemented
	if cr.Spec.Client.AuthorizationServicesEnabled || cr.Spec.Client.AuthorizationSettings != nil {
		rolesDeleted = removeUMARole(rolesDeleted)
	}
	for _, role := range rolesDeleted {
		desired.AddAction(i.getDeletedClientRoleState(state, cr, role.DeepCopy()))
	}

	// update with desired roles that can be matched to existing roles and have an ID set, this includes all renames
	// note down all renames
	existingRoleByID := make(map[string]kc.RoleRepresentation)
	for _, role := range state.Roles {
		existingRoleByID[role.ID] = role
	}
	renamedRolesOldNames := make(map[string]bool)
	_, rolesMatching := model.RoleDifferenceIntersection(cr.Spec.Roles, state.Roles)
	for _, role := range rolesMatching {
		if role.ID != "" {
			oldRole := existingRoleByID[role.ID]
			desired.AddAction(i.getUpdatedClientRoleState(state, cr, role.DeepCopy(), oldRole.DeepCopy()))
			if role.Name != oldRole.Name {
				renamedRolesOldNames[oldRole.Name] = true
			}
		}
	}

	// seemingly matching roles without an ID can either be regular updates
	// or re-creations after renames (not deletions)
	// note that duplicate role names are impossible thanks to +listType=map
	for _, role := range rolesMatching {
		if role.ID == "" {
			if _, contains := renamedRolesOldNames[role.Name]; contains {
				desired.AddAction(i.getCreatedClientRoleState(state, cr, role.DeepCopy()))
			} else {
				desired.AddAction(i.getUpdatedClientRoleState(state, cr, role.DeepCopy(), role.DeepCopy()))
			}
		}
	}

	// always create roles that don't match any existing ones
	rolesNew, _ := model.RoleDifferenceIntersection(cr.Spec.Roles, state.Roles)
	for _, role := range rolesNew {
		desired.AddAction(i.getCreatedClientRoleState(state, cr, role.DeepCopy()))
	}
}

func (i *KeycloakClientReconciler) ReconcileScopeMappings(state *common.ClientState, cr *kc.KeycloakClient, desired *common.DesiredClusterState) {
	if cr.Spec.ScopeMappings == nil {
		cr.Spec.ScopeMappings = &kc.MappingsRepresentation{}
	}
	for clientID, clientMappings := range cr.Spec.ScopeMappings.ClientMappings {
		clientMappings.Client = clientID
	}

	mappingsNew := scopeMappingDifference(cr.Spec.ScopeMappings, state.ScopeMappings)
	if mappingsNew.RealmMappings != nil {
		desired.AddAction(i.getCreatedClientRealmScopeMappingsState(state, cr, &mappingsNew.RealmMappings))
	}
	for _, clientMappings := range mappingsNew.ClientMappings {
		desired.AddAction(i.getCreatedClientClientScopeMappingsState(state, cr, clientMappings.DeepCopy()))
	}

	mappingsDeleted := scopeMappingDifference(state.ScopeMappings, cr.Spec.ScopeMappings)
	if mappingsDeleted.RealmMappings != nil {
		desired.AddAction(i.getDeletedClientRealmScopeMappingsState(state, cr, &mappingsDeleted.RealmMappings))
	}
	for _, clientMappings := range mappingsDeleted.ClientMappings {
		desired.AddAction(i.getDeletedClientClientScopeMappingsState(state, cr, clientMappings.DeepCopy()))
	}
}

func (i *KeycloakClientReconciler) ReconcileClientScopes(state *common.ClientState, cr *kc.KeycloakClient, desired *common.DesiredClusterState) {
	defaultClientScopes := model.FilterClientScopesByNames(state.AvailableClientScopes, cr.Spec.Client.DefaultClientScopes)

	defaultClientScopesNew, _ := model.ClientScopeDifferenceIntersection(defaultClientScopes, state.DefaultClientScopes)
	for _, clientScope := range defaultClientScopesNew {
		desired.AddAction(i.getCreatedClientDefaultClientScopeState(state, cr, clientScope.DeepCopy()))
	}

	defaultClientScopesDeleted, _ := model.ClientScopeDifferenceIntersection(state.DefaultClientScopes, defaultClientScopes)
	for _, clientScope := range defaultClientScopesDeleted {
		desired.AddAction(i.getDeletedClientDefaultClientScopeState(state, cr, clientScope.DeepCopy()))
	}

	optionalClientScopes := model.FilterClientScopesByNames(state.AvailableClientScopes, cr.Spec.Client.OptionalClientScopes)

	optionalClientScopesNew, _ := model.ClientScopeDifferenceIntersection(optionalClientScopes, state.OptionalClientScopes)
	for _, clientScope := range optionalClientScopesNew {
		desired.AddAction(i.getCreatedClientOptionalClientScopeState(state, cr, clientScope.DeepCopy()))
	}

	optionalClientScopesDeleted, _ := model.ClientScopeDifferenceIntersection(state.OptionalClientScopes, optionalClientScopes)
	for _, clientScope := range optionalClientScopesDeleted {
		desired.AddAction(i.getDeletedClientOptionalClientScopeState(state, cr, clientScope.DeepCopy()))
	}
}

// removeUMARole removes the uma_protection role from r if it is present
func removeUMARole(r []kc.RoleRepresentation) []kc.RoleRepresentation {
	filteredRoles, _ := model.RoleDifferenceIntersection(r, []kc.RoleRepresentation{{Name: umaRoleName}})
	return filteredRoles
}

// determine which scope mappings are present in a but not in b
// works on realm scope mappings and client scope mappings for each client separately
func scopeMappingDifference(a *kc.MappingsRepresentation, b *kc.MappingsRepresentation) (d *kc.MappingsRepresentation) {
	// no mappings = empty mappings
	if a == nil {
		a = &kc.MappingsRepresentation{}
	}
	if b == nil {
		b = &kc.MappingsRepresentation{}
	}
	// initialize empty result mappings
	d = &kc.MappingsRepresentation{ClientMappings: make(map[string]kc.ClientMappingsRepresentation)}

	// difference in realm scope mappings
	d.RealmMappings, _ = model.RoleDifferenceIntersection(a.RealmMappings, b.RealmMappings)

	// collect the client IDs (UUIDs) for all "clientID"s (unique names) giving priority to the IDs in a
	// this ensures the results will always contain the client ID if at all known
	// the ID is necessary for the path of the REST requests
	clientIDs := make(map[string]string)
	for clientID, mappings := range a.ClientMappings {
		clientIDs[clientID] = mappings.ID
	}
	for clientID, mappings := range b.ClientMappings {
		if id := clientIDs[clientID]; id == "" {
			clientIDs[clientID] = mappings.ID
		}
	}

	// calculate difference for each client separately
	for clientID, id := range clientIDs {
		rolesA := a.ClientMappings[clientID].Mappings
		rolesB := b.ClientMappings[clientID].Mappings
		rolesD, _ := model.RoleDifferenceIntersection(rolesA, rolesB)
		if len(rolesD) > 0 {
			d.ClientMappings[clientID] = kc.ClientMappingsRepresentation{
				ID:       id,
				Client:   clientID,
				Mappings: rolesD,
			}
		}
	}

	return d
}

func (i *KeycloakClientReconciler) pingKeycloak() common.ClusterAction {
	return common.PingAction{
		Msg: "check if keycloak is available",
	}
}

func (i *KeycloakClientReconciler) getDeletedClientState(state *common.ClientState, cr *kc.KeycloakClient) common.ClusterAction {
	return common.DeleteClientAction{
		Ref:   cr,
		Realm: state.Realm.Spec.Realm.Realm,
		Msg:   fmt.Sprintf("removing client %v/%v", cr.Namespace, cr.Spec.Client.ClientID),
	}
}

func (i *KeycloakClientReconciler) getCreatedClientState(state *common.ClientState, cr *kc.KeycloakClient) common.ClusterAction {
	return common.CreateClientAction{
		Ref:   cr,
		Realm: state.Realm.Spec.Realm.Realm,
		Msg:   fmt.Sprintf("create client %v/%v", cr.Namespace, cr.Spec.Client.ClientID),
	}
}

func (i *KeycloakClientReconciler) getUpdatedClientSecretState(state *common.ClientState, cr *kc.KeycloakClient) common.ClusterAction {
	return common.GenericUpdateAction{
		Ref: model.ClientSecretReconciled(cr, state.ClientSecret),
		Msg: fmt.Sprintf("update client secret %v/%v", cr.Namespace, cr.Spec.Client.ClientID),
	}
}

func (i *KeycloakClientReconciler) getUpdatedClientState(state *common.ClientState, cr *kc.KeycloakClient) common.ClusterAction {
	return common.UpdateClientAction{
		Ref:   cr,
		Realm: state.Realm.Spec.Realm.Realm,
		Msg:   fmt.Sprintf("update client %v/%v", cr.Namespace, cr.Spec.Client.ClientID),
	}
}

func (i *KeycloakClientReconciler) getCreatedClientSecretState(state *common.ClientState, cr *kc.KeycloakClient) common.ClusterAction {
	return common.GenericCreateAction{
		Ref: model.ClientSecret(cr),
		Msg: fmt.Sprintf("create client secret %v/%v", cr.Namespace, cr.Spec.Client.ClientID),
	}
}

func (i *KeycloakClientReconciler) getCreatedClientRoleState(state *common.ClientState, cr *kc.KeycloakClient, role *kc.RoleRepresentation) common.ClusterAction {
	return common.CreateClientRoleAction{
		Role:  role,
		Ref:   cr,
		Realm: state.Realm.Spec.Realm.Realm,
		Msg:   fmt.Sprintf("create client role %v/%v/%v", cr.Namespace, cr.Spec.Client.ClientID, role.Name),
	}
}

func (i *KeycloakClientReconciler) getUpdatedClientRoleState(state *common.ClientState, cr *kc.KeycloakClient, role, oldRole *kc.RoleRepresentation) common.ClusterAction {
	return common.UpdateClientRoleAction{
		Role:    role,
		OldRole: oldRole,
		Ref:     cr,
		Realm:   state.Realm.Spec.Realm.Realm,
		Msg:     fmt.Sprintf("update client role %v/%v/%v", cr.Namespace, cr.Spec.Client.ClientID, oldRole.Name),
	}
}

func (i *KeycloakClientReconciler) getDeletedClientRoleState(state *common.ClientState, cr *kc.KeycloakClient, role *kc.RoleRepresentation) common.ClusterAction {
	return common.DeleteClientRoleAction{
		Role:  role,
		Ref:   cr,
		Realm: state.Realm.Spec.Realm.Realm,
		Msg:   fmt.Sprintf("delete client role %v/%v/%v", cr.Namespace, cr.Spec.Client.ClientID, role.Name),
	}
}

func (i *KeycloakClientReconciler) getCreatedClientRealmScopeMappingsState(state *common.ClientState, cr *kc.KeycloakClient, mappings *[]kc.RoleRepresentation) common.ClusterAction {
	return common.CreateClientRealmScopeMappingsAction{
		Mappings: mappings,
		Ref:      cr,
		Realm:    state.Realm.Spec.Realm.Realm,
		Msg:      fmt.Sprintf("create client realm scope mappings for %v/%v", cr.Namespace, cr.Spec.Client.ClientID),
	}
}

func (i *KeycloakClientReconciler) getDeletedClientRealmScopeMappingsState(state *common.ClientState, cr *kc.KeycloakClient, mappings *[]kc.RoleRepresentation) common.ClusterAction {
	return common.DeleteClientRealmScopeMappingsAction{
		Mappings: mappings,
		Ref:      cr,
		Realm:    state.Realm.Spec.Realm.Realm,
		Msg:      fmt.Sprintf("delete client realm scope mappings for %v/%v", cr.Namespace, cr.Spec.Client.ClientID),
	}
}

func (i *KeycloakClientReconciler) getCreatedClientClientScopeMappingsState(state *common.ClientState, cr *kc.KeycloakClient, mappings *kc.ClientMappingsRepresentation) common.ClusterAction {
	return common.CreateClientClientScopeMappingsAction{
		Mappings: mappings,
		Ref:      cr,
		Realm:    state.Realm.Spec.Realm.Realm,
		Msg:      fmt.Sprintf("create client client scope mappings %v/%v => %v", cr.Namespace, cr.Spec.Client.ClientID, mappings.Client),
	}
}

func (i *KeycloakClientReconciler) getDeletedClientClientScopeMappingsState(state *common.ClientState, cr *kc.KeycloakClient, mappings *kc.ClientMappingsRepresentation) common.ClusterAction {
	return common.DeleteClientClientScopeMappingsAction{
		Mappings: mappings,
		Ref:      cr,
		Realm:    state.Realm.Spec.Realm.Realm,
		Msg:      fmt.Sprintf("delete client client scope mappings %v/%v => %v", cr.Namespace, cr.Spec.Client.ClientID, mappings.Client),
	}
}

func (i *KeycloakClientReconciler) getCreatedClientDefaultClientScopeState(state *common.ClientState, cr *kc.KeycloakClient, clientScope *kc.KeycloakClientScope) common.ClusterAction {
	return common.UpdateClientDefaultClientScopeAction{
		ClientScope: clientScope,
		Ref:         cr,
		Realm:       state.Realm.Spec.Realm.Realm,
		Msg:         fmt.Sprintf("create client default client scope %v/%v => %v", cr.Namespace, cr.Spec.Client.ClientID, clientScope.Name),
	}
}

func (i *KeycloakClientReconciler) getCreatedClientOptionalClientScopeState(state *common.ClientState, cr *kc.KeycloakClient, clientScope *kc.KeycloakClientScope) common.ClusterAction {
	return common.UpdateClientOptionalClientScopeAction{
		ClientScope: clientScope,
		Ref:         cr,
		Realm:       state.Realm.Spec.Realm.Realm,
		Msg:         fmt.Sprintf("create client optional client scope %v/%v => %v", cr.Namespace, cr.Spec.Client.ClientID, clientScope.Name),
	}
}

func (i *KeycloakClientReconciler) getDeletedClientDefaultClientScopeState(state *common.ClientState, cr *kc.KeycloakClient, clientScope *kc.KeycloakClientScope) common.ClusterAction {
	return common.DeleteClientDefaultClientScopeAction{
		ClientScope: clientScope,
		Ref:         cr,
		Realm:       state.Realm.Spec.Realm.Realm,
		Msg:         fmt.Sprintf("delete client default client scope %v/%v => %v", cr.Namespace, cr.Spec.Client.ClientID, clientScope.Name),
	}
}

func (i *KeycloakClientReconciler) getDeletedClientOptionalClientScopeState(state *common.ClientState, cr *kc.KeycloakClient, clientScope *kc.KeycloakClientScope) common.ClusterAction {
	return common.DeleteClientOptionalClientScopeAction{
		ClientScope: clientScope,
		Ref:         cr,
		Realm:       state.Realm.Spec.Realm.Realm,
		Msg:         fmt.Sprintf("delete client optional client scope %v/%v => %v", cr.Namespace, cr.Spec.Client.ClientID, clientScope.Name),
	}
}
