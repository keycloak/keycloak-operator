FROM quay.io/keycloak/keycloak:11.0.0

USER root

RUN microdnf install -y net-tools procps findutils vim && microdnf clean all

USER 1000
