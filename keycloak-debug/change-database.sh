#!/bin/bash -e

DB_VENDOR=$1

export KEYCLOAK_LOGLEVEL="TRACE"
export ROOT_LOGLEVEL="TRACE"

cd /opt/jboss/keycloak

echo "Setting loglevel"
bin/jboss-cli.sh --connect --file=/opt/jboss/tools/cli/loglevel.cli
echo "Starting standalone config"
bin/jboss-cli.sh --connect --file=/opt/jboss/tools/cli/databases/$DB_VENDOR/standalone-configuration.cli
rm -rf /opt/jboss/keycloak/standalone/configuration/standalone_xml_history
echo "Starting HA config"
bin/jboss-cli.sh --connect --file=/opt/jboss/tools/cli/databases/$DB_VENDOR/standalone-ha-configuration.cli
rm -rf standalone/configuration/standalone_xml_history/current/*
