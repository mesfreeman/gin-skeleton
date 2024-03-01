#!/bin/bash
# 
# Gin skeleton project deployment script
# 
# Usage: 
# ./deploy.sh [server] [project_path]
# default server is tank.server.cn
# default project_path is /data/services/projects/gin-skeleton/code
#

set -e

DATE="$(date +%Y-%m-%d\ %H:%M:%S)"
INFO="\033[32m[$DATE] [INFO]\033[0m"
WARN="\033[33m[$DATE] [WARN]\033[0m"
SUCCESS="\033[32m[$DATE] [SUCCESS]\033[0m"
ERROR="\033[31m[$DATE] [ERROR]\033[0m"

SERVER="$1"
PROJECT_PATH="$2"
LATEST_VERSION=""

if [ -z "$SERVER" ]; then
  SERVER="tank.server.cn"
fi

if [ -z "$PROJECT_PATH" ]; then
  PROJECT_PATH="/data/services/projects/gin-skeleton/code"
fi

# Get the latest version number, and create a new tag
CURRENT_VERSION=`git tag | sort -V | tail -n 1`
if [ -z "$CURRENT_VERSION" ]; then
  CURRENT_VERSION="-"
fi
read -p "Please input new version, current version is $CURRENT_VERSION -> " LATEST_VERSION
if [ -z "$LATEST_VERSION" ]; then
  echo -e "$ERROR Latest version number can not be empty"
  exit 1
fi

# Create a new tag, and push it to the remote repository
echo -e "$INFO Create tag $LATEST_VERSION"
git reset --hard
git clean -df
git checkout master
git fetch origin master -p -t
git reset --hard origin/master
git tag -a "$LATEST_VERSION" -m "chore: version $LATEST_VERSION"
git push origin "$LATEST_VERSION"
echo -e "$SUCCESS Create tag $LATEST_VERSION success"

# Build the gin-skeleton project
echo -e "$INFO Build gin-skeleton project, version is $LATEST_VERSION"
go mod tidy
go vet ./...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X 'gin-skeleton/helper.Version=$LATEST_VERSION' -X 'gin-skeleton/helper.BuildTime=$(date +%Y-%m-%d\ %H:%M:%S)'" -o release/gin-skeleton main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X 'gin-skeleton/helper.Version=$LATEST_VERSION' -X 'gin-skeleton/helper.BuildTime=$(date +%Y-%m-%d\ %H:%M:%S)'" -o release/gin-cli cmd/main.go
echo -e "$SUCCESS Build gin-skeleton:$LATEST_VERSION success"

# Deploy the gin-skeleton project to the server
echo -e "$INFO Deploy gin-skeleton:$LATEST_VERSION to $SERVER"
rsync -raze "ssh -p 20022 -l root" --exclude="release/config" --exclude="release/storage" "release" "$SERVER:$PROJECT_PATH/"
rsync -raze "ssh -p 20022 -l root" --exclude="config/config.yaml" "config" "$SERVER:$PROJECT_PATH/"
rsync -raze "ssh -p 20022 -l root" --exclude="storage/logs/*.log" "storage" "$SERVER:$PROJECT_PATH/"
rsync -raze "ssh -p 20022 -l root" "pm2.json" "$SERVER:$PROJECT_PATH/"

# Initialize the gin-skeleton project
ssh -p 20022 root@"$SERVER" "chmod +x $PROJECT_PATH/release/gin-*"
ssh -p 20022 root@"$SERVER" "if [ ! -L $PROJECT_PATH/release/storage ]; then ln -s ../storage $PROJECT_PATH/release/storage; fi"
ssh -p 20022 root@"$SERVER" "if [ ! -L $PROJECT_PATH/release/config ]; then ln -s ../config $PROJECT_PATH/release/config; fi"
ssh -p 20022 root@"$SERVER" "chown -R root:root $PROJECT_PATH/"
ssh -p 20022 root@"$SERVER" "chown -R www-data:www-data $PROJECT_PATH/storage"

# Start the gin-skeleton service
if [ `ssh -p 20022 root@"$SERVER" "if [ ! -f $PROJECT_PATH/config/config.yaml ]; then echo 1; fi"` ]; then
  echo -e "$WARN Copy config file ./config/config.yaml.example to ./config/config.yaml"
  ssh -p 20022 root@"$SERVER" "cp $PROJECT_PATH/config/config.yaml.example $PROJECT_PATH/config/config.yaml"
  ssh -p 20022 root@"$SERVER" "chown root:root $PROJECT_PATH/config/config.yaml"
  echo -e "$WARN Please modify the config file $PROJECT_PATH/config/config.yaml"
  echo -e "$WARN Need to start the gin-skeleton service after modifying the config file, use the command: pm2 start $PROJECT_PATH/pm2.json"
  echo -e "$SUCCESS Deploy gin-skeleton:$LATEST_VERSION to $SERVER success"
else
  ssh -p 20022 root@"$SERVER" "pm2 startOrRestart $PROJECT_PATH/pm2.json"
  echo -e "$INFO Please wait 3 seconds" && sleep 3
  if [ `ssh -p 20022 root@"$SERVER" "pm2 list | grep gin-skeleton | grep online | wc -l"` -eq 1 ]; then
    echo -e "$SUCCESS Deploy gin-skeleton:$LATEST_VERSION to $SERVER success"
  else
    echo -e "$ERROR Deploy gin-skeleton:$LATEST_VERSION to $SERVER failed"
  fi
fi
