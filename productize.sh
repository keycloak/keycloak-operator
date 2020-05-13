#!/bin/bash -ex

LABELS=$(cat <<-END
LABEL \
com.redhat.component="redhat-sso-7-sso74-rhel8-operator-container"  \
description="Red Hat Single Sign-On 7.4 Operator container image, based on the Red Hat Universal Base Image 8 Minimal container image" \
summary="Red Hat Single Sign-On 7.4 Operator container image, based on the Red Hat Universal Base Image 8 Minimal container image" \
version="7.4" \
io.k8s.description="Operator for Red Hat SSO" \
io.k8s.display-name="Red Hat SSO 7.4 Operator" \
io.openshift.tags="sso,sso74,keycloak,operator" \
name="rh-sso-7\/sso74-operator-rhel8" \
maintainer="Red Hat Single Sign-On Team"
END
)

GIT_REPOSITORY="${1:-https://github.com/keycloak/keycloak-operator.git}"
GIT_BRANCH="${2:-master}"

GIT_LAST_COMMIT_HASH=$(git ls-remote "${GIT_REPOSITORY}" "${GIT_BRANCH}" | awk '{ print $1}')

# Ensure that commit ID is there to prevent returning cached build upon a new commit
GIT_COMMAND="git clone --single-branch --branch ${GIT_BRANCH} ${GIT_REPOSITORY} .  #${GIT_LAST_COMMIT_HASH}"

sed -i \
    -e 's/FROM registry.access.redhat.com\/ubi8\/ubi-minimal:[0-9.]*/FROM ubi8-minimal:8-released/' \
    -e "s/##LABELS/$LABELS/g" \
    -e "s,^## *RUN git clone .*\$,RUN $GIT_COMMAND," \
    Dockerfile
