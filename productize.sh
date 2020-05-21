#!/bin/bash -ex

LABELS=$(cat <<-END
LABEL \
com.redhat.component="redhat-sso-7-sso74-rhel8-tech-preview-operator-container"  \
description="Red Hat Single Sign-On 7.4 Operator container image, based on the Red Hat Universal Base Image 8 Minimal container image" \
summary="Red Hat Single Sign-On 7.4 Operator container image, based on the Red Hat Universal Base Image 8 Minimal container image" \
version="7.4" \
io.k8s.description="Operator for Red Hat SSO" \
io.k8s.display-name="Red Hat SSO 7.4 Operator" \
io.openshift.tags="sso,sso74,keycloak,operator" \
name="rh-sso-7\/sso74-operator-rhel8-tech-preview" \
maintainer="Red Hat Single Sign-On Team"
END
)

sed -i \
    -e 's/registry.svc.ci.openshift.org\/openshift\/release:golang-1.13/openshift\/golang-builder:1.13/' \
    -e 's/FROM registry.access.redhat.com\/ubi8\/ubi-minimal:[0-9.]*/FROM ubi8-minimal:8-released/' \
    -e "s/##LABELS/$LABELS/g" \
    Dockerfile
