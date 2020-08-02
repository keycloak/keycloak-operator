FROM quay.io/keycloak/keycloak:11.0.0

USER root

RUN microdnf install -y net-tools procps findutils vim java-11-openjdk-devel.x86_64 && microdnf clean all

COPY jboss-cli-logging.properties /opt/jboss/keycloak/bin/jboss-cli-logging.properties

USER 1000
