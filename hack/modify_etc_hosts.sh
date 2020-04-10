#!/bin/bash

# The main part of this script has been downloaded from: https://gist.github.com/jacobtomlinson/4b835d807ebcea73c6c8f602613803d4

set -x

INGRESSES=$1
MINIKUBE_IP=$2

if [ -z "$INGRESSES" ]; then
  echo "Ingress address not set"
  exit 1
fi

if [ -z "$MINIKUBE_IP" ]; then
  echo "Assuming Minikube running with the current user"
  MINIKUBE_IP=$(minikube ip || exit 1)
fi

HOSTS_ENTRY="$MINIKUBE_IP $INGRESSES"

if grep -Fq "$MINIKUBE_IP" /etc/hosts > /dev/null
then
    sudo sed -i "s/^$MINIKUBE_IP.*/$HOSTS_ENTRY/" /etc/hosts
    echo "Updated hosts entry"
else
    echo "$HOSTS_ENTRY" | sudo tee -a /etc/hosts
    echo "Added hosts entry"
fi