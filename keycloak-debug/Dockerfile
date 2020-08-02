FROM quay.io/keycloak/keycloak:11.0.0

USER root

RUN microdnf install -y net-tools procps findutils vim && microdnf clean all

COPY jboss-cli-logging.properties /opt/jboss/keycloak/bin/jboss-cli-logging.properties

USER 1000
