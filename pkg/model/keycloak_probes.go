package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	LivenessProbeImplementation = `#!/bin/bash
set -e

PASSWORD_FILE="/tmp/management-password"
PASSWORD="not set"
USERNAME="admin"

if [ -d "/opt/eap/bin" ]; then
    pushd /opt/eap/bin > /dev/null
else
    pushd /opt/jboss/keycloak/bin > /dev/null
fi

if [ -f "$PASSWORD_FILE" ]; then
    PASSWORD=$(cat $PASSWORD_FILE)
else
    PASSWORD=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
	./add-user.sh -m -u $USERNAME -p $PASSWORD > /dev/null
	echo $PASSWORD > $PASSWORD_FILE
fi

./jboss-cli.sh --connect --user=$USERNAME --password=$PASSWORD --command-timeout=10 --commands="/deployment=keycloak-server.war:read-attribute(name=status)" > /dev/null
`
	ReadinessProbeImplementation = `#!/bin/bash
set -e

PASSWORD_FILE="/tmp/management-password"
PASSWORD="not set"
USERNAME="admin"

DATASOURCE_POOL_TYPE="data-source"
DATASOURCE_POOL_NAME="KeycloakDS"

if [ -d "/opt/eap/bin" ]; then
    pushd /opt/eap/bin > /dev/null
	DATASOURCE_POOL_TYPE=xa-data-source
	DATASOURCE_POOL_NAME=keycloak_postgresql-DB
else
    pushd /opt/jboss/keycloak/bin > /dev/null
fi

if [ -f "$PASSWORD_FILE" ]; then
    PASSWORD=$(cat $PASSWORD_FILE)
else
    PASSWORD=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
	./add-user.sh -m -u $USERNAME -p $PASSWORD> /dev/null
	echo $PASSWORD > $PASSWORD_FILE
fi

./jboss-cli.sh --connect --user=$USERNAME --password=$PASSWORD --command-timeout=10 --commands="/subsystem=datasources/${DATASOURCE_POOL_TYPE}=${DATASOURCE_POOL_NAME}:test-connection-in-pool" > /dev/null

curl -s --max-time 10 http://$(hostname -i):8080/auth > /dev/null
`
)

func KeycloakProbes(cr *v1alpha1.Keycloak) *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: v12.ObjectMeta{
			Name:      KeycloakProbesName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app":           ApplicationName,
				ApplicationName: cr.Name,
			},
		},
		Data: map[string]string{
			LivenessProbeProperty:  LivenessProbeImplementation,
			ReadinessProbeProperty: ReadinessProbeImplementation,
		},
	}
}

func KeycloakProbesSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      KeycloakProbesName,
		Namespace: cr.Namespace,
	}
}
