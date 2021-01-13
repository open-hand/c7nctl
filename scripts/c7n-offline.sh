#!/usr/local/bin/bash

declare -A c7nChart=(
  ["chartmuseum"]="2.15.0"
  ["minio"]="5.0.5"
  ["mysql"]="0.1.4"
  ["redis"]="0.2.5"
  ["gitlab-ha"]="0.2.2"
  ["harbor"]="1.2.3"
  ["mysql-client"]="0.1.0"
  ["gitlab-runner"]="0.2.4"
  ["sonatype-nexus"]="3.4.0"
  ["sonarqube"]="0.15.0-3"
  ["persistentvolumeclaim"]="0.1.0"
  ["choerodon-admin"]="0.24.0"
  ["choerodon-asgard"]="0.24.0"
  ["choerodon-file"]="0.24.0"
  ["choerodon-front-hzero"]="0.24.0"
  ["choerodon-front"]="0.24.0"
  ["choerodon-gateway"]="0.24.0"
  ["choerodon-iam"]="0.24.0"
  ["choerodon-message"]="0.24.0"
  ["choerodon-monitor"]="0.24.0"
  ["choerodon-oauth"]="0.24.0"
  ["choerodon-platform"]="0.24.0"
  ["choerodon-register"]="0.24.0"
  ["workflow-service"]="0.24.0"
  ["gitlab-service"]="0.24.0"
  ["devops-service"]="0.24.0"
  ["test-manager-service"]="0.24.0"
  ["knowledgebase-service"]="0.24.0"
  ["agile-service"]="0.24.0"
  ["elasticsearch-kb"]="0.24.0"
  ["code-repo-service"]="0.24.0"
  ["prod-repo-service"]="0.24.0"
)

# kubectl get po -n c7n-system -o yaml | grep -E "^      image:" |  awk '{print $2}' | sort | uniq | sed 's/:/ /g' | awk '{print "[\""$1"\"]=\"" $2"\""}'
declare -A dockerImage=(
    ["alpine"]="3.10.3"
    ["busybox"]="1.31"
    ["busybox"]="latest"
    ["goharbor/harbor-core"]="v1.9.3"
    ["goharbor/harbor-db"]="v1.9.3"
    ["goharbor/harbor-jobservice"]="v1.9.3"
    ["goharbor/harbor-portal"]="v1.9.3"
    ["goharbor/harbor-registryctl"]="v1.9.3"
    ["goharbor/redis-photon"]="v1.9.3"
    ["goharbor/registry-photon"]="v2.7.1-patch-2819-2553-v1.9.3"
    ["joosthofman/wget"]="1.0"
    ["minio/minio"]="RELEASE.2020-01-03T19-12-21Z"
    ["mysql"]="5.7.23"
    ["postgres"]="9.6.2"
    ["redis"]="4.0.11"
    ["sonatype/nexus3"]="3.19.1"
    ["sonarqube"]="7.6-community"
    ["registry.cn-shanghai.aliyuncs.com/c7n/code-repo-service"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/prod-repo-service"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/dbtool"]="0.7.2"
    ["registry.cn-shanghai.aliyuncs.com/c7n/dbtool"]="0.7.1"
    ["registry.cn-shanghai.aliyuncs.com/c7n/elasticsearch-kb"]="7.2-elasticsearch-kb"
    ["registry.cn-shanghai.aliyuncs.com/c7n/postgres"]="11.5"
    ["registry.cn-shanghai.aliyuncs.com/c7n/skywalking-agent"]="6.5.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/skywalking-agent"]="6.6.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/sonarqube"]="7.9.3-community"
    ["registry.cn-shanghai.aliyuncs.com/c7n/agile-service"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/chartmuseum"]="v0.11.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/knowledgebase-service"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/choerodon-admin"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/choerodon-asgard"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/choerodon-file"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/choerodon-front-hzero"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/choerodon-front"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/choerodon-gateway"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/choerodon-iam"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/choerodon-jmeter-agent"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/choerodon-message"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/choerodon-monitor"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/choerodon-oauth"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/choerodon-platform"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/choerodon-register"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/choerodon-swagger"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/devops-service"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/gitlab-service"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/test-manager-service"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/workflow-service"]="0.24.0"
    ["registry.cn-shanghai.aliyuncs.com/c7n/docker-nexus"]="3.25.1"
    ["registry.cn-shanghai.aliyuncs.com/c7n/postgresql"]="10-2"
    ["registry.cn-shanghai.aliyuncs.com/c7n/redis"]="4.0.9-2"
)

OFFLINE_PATH=_offline
CHARTS_PATH=${OFFLINE_PATH}/charts
IMAGES_PATH=${OFFLINE_PATH}/images

NEXUS_USERNAME="admin"
NEXUS_PASSWORD="admin123"
BASE_CHART="http://localhost:8081/repository/helm/"
BASE_IMAGE="localhost:8082/"

function add_helm_repo() {
    helm repo add c7n    https://openchart.choerodon.com.cn/choerodon/c7n/
    helm repo update
}

function pre_download() {
  mkdir -p ${CHARTS_PATH}
  mkdir -p ${IMAGES_PATH}
}

function pull_group_chart() {
    # shellcheck disable=SC2154
    # shellcheck disable=SC2068
    for key in ${!c7nChart[@]}
    do
        echo "    helm pull chart $key for version ${c7nChart[$key]}"
        helm pull c7n/"$key" --version "${c7nChart[$key]}"
    done
}

function pull_group_image() {
    # shellcheck disable=SC2068
    for i in ${!dockerImage[@]}
    do
        docker pull "$i":"${dockerImage[$i]}"

	      # shellcheck disable=SC2207
	      # shellcheck disable=SC2034
	      # shellcheck disable=SC2006
	      arr=(`echo "$i" | tr '/' ' '`)
        docker save "$i":"${dockerImage[$i]}" -o ${IMAGES_PATH}/"${arr[-1]}"-"${dockerImage[$i]}".tar
    done
}

function pull_chart() {
    echo "Staring pull charts"
    cd ${CHARTS_PATH} || exit
    pull_group_chart
    # shellcheck disable=SC2164
    cd ../..
}

function push_chart() {
  # shellcheck disable=SC2006
  # shellcheck disable=SC2034

  # shellcheck disable=SC2231
  for f in ${CHARTS_PATH}/*
  do
    curl -u "${NEXUS_USERNAME}":"${NEXUS_PASSWORD}" ${BASE_CHART} --upload-file "$f" -v
  done
}

function push_image() {
      # shellcheck disable=SC2068
    for i in ${!dockerImage[@]}
    do
	      # shellcheck disable=SC2207
	      # shellcheck disable=SC2034
	      # shellcheck disable=SC2006
	      arr=(`echo "$i" | tr '/' ' '`)
	      docker tag "$i":"${dockerImage[$i]}" ${BASE_IMAGE}"${arr[-1]}":"${dockerImage[$i]}"
        docker push ${BASE_IMAGE}"${arr[-1]}":"${dockerImage[$i]}"
    done
}

function list_version() {
    # shellcheck disable=SC2154
    # shellcheck disable=SC2068
    for key in ${!c7nChart[@]}
    do
        helm search repo c7n/"${key}" | grep c7n/"${key}"
    done
}
# pre_download
# add_helm_repo
# pull_group_image
# pull_chart
# push_image
# push_chart

list_version