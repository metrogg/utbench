#!/bin/bash

# 获取执行 UID
function _uid() {
  uuidgen | awk '{print tolower($0)}'
}

# 获取用户名
function _username() {
  local user_email="$(git config user.email)"
  local username="${user_email%@tencent.com}"

  if [[ "${username}" == "" ]]; then
    username="anonymous"
  fi

  echo "${username}"
}

# 获取仓库地址
function _repository() {
  local remote_url="$(git remote get-url origin)"
  local repo="${remote_url#git@}"
  repo="${repo#https://}"
  repo="${repo#http://}"
  repo="${repo%.git}"
  repo="${repo//://}"
  echo "${repo}"
}

# 获取当前命令
function _cmd() {
  if [[ "${CLAUDECODE}" -gt "0" ]]; then
    echo "claude"
  else
    echo "codebuddy"
  fi
}

if [[ "${1}" == "" ]]; then
  echo "no data to push"
  exit 0
fi

push_data='{
  "Version": 1,
  "UID": "'$(_uid)'",
  "Args": ["'$(_cmd)'", "/qta-gen-ut"],
  "Username": "'$(_username)'",
  "Repository": "'$(_repository)'",
  "Commit": "'$(git rev-parse HEAD)'",
  "Results": {
    "Generate": {
      "Batches": '${1}'
    }
  }
}'

echo "Request:"
echo "${push_data}"
echo "Response:"
curl http://qta.woa.com/api/qta/ut-agent/v1/runrecords -H 'Content-Type: application/json' --data "${push_data}"
