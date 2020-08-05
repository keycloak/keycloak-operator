FROM docker.io/jboss/keycloak:11.0.0

USER root

RUN microdnf update -y && \
    microdnf install -y net-tools procps findutils vim java-11-openjdk-devel.x86_64 && microdnf clean all

COPY jboss-cli-logging.properties /opt/jboss/keycloak/bin/jboss-cli-logging.properties
COPY change-database.sh /opt/jboss/tools/databases/change-database.sh
COPY standalone-configuration.cli /opt/jboss/tools/cli/standalone-configuration.cli

USER 1000
