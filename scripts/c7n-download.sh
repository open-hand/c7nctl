#!/usr/local/bin/bash

source ./c7n-offline.sh

pre_download
add_helm_repo
pull_image
pull_chart
