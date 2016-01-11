# 3a-classic-server

[![Circle CI](https://circleci.com/gh/3a-classic/score-api-server.svg?style=shield&circle-token=05ab242168e17f5fd1b442f85002417e6e963a3a)](https://circleci.com/gh/3a-classic/score-api-server)
[![Apache2.0 License](http://img.shields.io/badge/license-Apache2.0-blue.svg?style=flat)](LICENSE)

## DESCRIPTION

This is api server for 3a-classic

## DEPENDENCY

* [3a-classic-ui](http://git.sadayuki-matsuno.com/3aclassic/3a-classic-ui)
* mongo

## REQUIRED CONFIGURATION


## USAGE

### MONGO

you can use official docker images of mongo

* docker mongo

```
$ docker run --name 3a-classic-mongo -d mongo
```

### 3A-CLASSIC-SERVER

* change directry to your work dir

```bash
$ cd /path/to/work/dir
```

* download this repositry

```bash
$ git clone http://git.sadayuki-matsuno.com/3aclassic/3a-classic-server.git
```

* build

```bash
$ gb build
```

* run docker

```bash
$ docker run -d -t \
  -v /home/matsuno/docker/3a-classic/server/3a-classic-server:/go/src \
  --name 3a-classic-server \
  --expose 80 \
  -p 8081:80 \
  -e "SLACK_INCOMING_HOOK_URL=https://hooks.slack.com/***********" \
  -e "SLACK_CHANNEL=#3a-classic" \
  -e "SLACK_USERNAME=3a-classic-error-log" \
  -e "MONGO_HOST=172.17.0.33" \
  -e "MONGO_PORT=27017" \
  -e "MONGO_DB=test" \
  -e "MONGO_LOG_COLLECTION=log" \
  -e "GIT_REMOTE_SERVICE=github" \
  -e "GIT_REMOTE_URL=https://github.com/sadayuki-matsuno/3aclassic_restful" \
  -e "AUTH_ADMIN_PASS_BASE=*********" \
  golang:latest \
  bash -c /go/src/bin/3aClassic
```

* access

```
http://localhost:8081
```

### Environment Variables

|Key|ExampleValue|Default|Explain|
|:--|:--|:--|:--|:--|
|SLACK_INCOMING_HOOK_URL|https://hooks.slack.com/***|""|Notify you of error logs by Slack|
|SLACK_CHANNEL|#3a-classic|Depend on [slackrus](https://github.com/johntdyer/slackrus)|Slack use this channel|
|SLACK_USERNAME|3a-classic-error-log|Depend on [slackrus](https://github.com/johntdyer/slackrus)|Slack use this username|
|MONGO_HOST|172.17.0.33|mongo|MongoDB hostname or IP|
|MONGO_PORT|27017|27017|MongoDB port|
|MONGO_DB|test|test|MongoDB database name|
|MONGO_LOG_COLLECTION|log|log|MongoDB log collection name|
|GIT_REMOTE_SERVICE|github/gitlab|github|for writing error code location in logs. you can also choose "gitlab"|
|GIT_REMOTE_URL|http://git.sadayuki-matsuno.com/3aclassic/3a-classic-server|http://git.sadayuki-matsuno.com/3aclassic/3a-classic-server|for writing error code location in logs.|

**Variable Name is not same as Docker one, because I heard Docker Link is "Legacy" somewhere.**

### CORS

if you access by chrome, you have to set CORS.
this is example to configure CORS in nginx.

```
$ vim /etc/nginx/conf.d/proxy.conf

add_header 'Access-Control-Allow-Origin' '*';
add_header 'Access-Control-Allow-Credentials' 'true';
add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';
add_header 'Access-Control-Allow-Headers' 'DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type';

if ($request_method = 'OPTIONS') {
        return 204;
}
```

if you use [nginx-proxy](https://github.com/jwilder/nginx-proxy), add file below.

```
$ vim /etc/nginx/vhost.d/default_location 

  add_header 'Access-Control-Allow-Origin' '$http_origin';
  add_header 'Access-Control-Allow-Credentials' 'true';
  add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS';
  add_header 'Access-Control-Allow-Headers' 'Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Mx-ReqToken,X-Requested-With';

if ($request_method = 'OPTIONS') {
        return 204;
}

if ($request_method = 'POST') {
  add_header Access-Control-Allow-Origin '';
}
```

### WEBSOCKET

if you use webdocket, you have to set below in nginx.

```
$ vim /etc/nginx/conf.d/proxy.conf

proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection upgrade;
```

## CAUTION

* change 3a-classic-server info in 3a-classic-ui
* 3a-classic-server have to be reachable from client node
* data in mongo-docker will disappear even if you save the docker image

## FOR DEVELOPER

if you edit and run this code, ready for development.

## GB

I use [gb](https://getgb.io/) as golang package manager

* install gb

```
$ go get github.com/constabulary/gb/...
```

* install bzr

* change directry to this repo

```
$ cd 3a-classic-server
```

* package install

```
$ gb vendor restore ./vendor/manifest
```

* build code

```
$ gb build
```

now you can run code compiled.

```
$ ./bin/3aClassic
```

## AUTHOR

 Sadayuki Matsuno
