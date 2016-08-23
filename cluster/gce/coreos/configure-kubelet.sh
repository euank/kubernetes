#!/bin/bash

# Copyright 2015 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

MANIFESTS_DIR=/opt/kube-manifests/kubernetes

function load-master-components-images() {
  echo "Loading images for master components"
  export RKT_BIN=/opt/rkt/rkt
  export DOCKER2ACI_BIN=/opt/docker2aci/docker2aci
  ${SALT_DIR}/install.sh ${KUBE_BIN_TAR}
  ${SALT_DIR}/salt/kube-master-addons/kube-master-addons.sh

  # Get the image tag
  HYPERKUBE_DOCKER_TAG=$(cat ${KUBE_BIN_DIR}/hyperkube.docker_tag)
}

load-kubelet-components-images

echo "Configuring hostname"
hostnamectl set-hostname $(hostname | cut -f1 -d.)

echo "Configuring kubelet"
mkdir -p /var/lib/kubelet
mkdir -p /etc/kubernetes/manifests
src=${MANIFESTS_DIR}/kubelet-config.yaml
dst=/var/lib/kubelet/kubeconfig
cp ${src} ${dst}
sed -i 's/\"/\\\"/g' ${dst} # eval will remove the double quotes if they are not escaped
eval "echo \"$(< ${dst})\"" > ${dst}
