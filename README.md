# 3a-classic-ui

## DESCRIPTION

This is api server for 3a-classic

## DEPENDENCY

* [3a-classic-ui](http://git.sadayuki-matsuno.com/3aclassic/3a-classic-ui)
* mongo

## REQUIRED CONFIGURE


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

* run docker

```bash
$ docker run -d -t --link 3a-classic-mongo:mongo \
  -v /home/matsuno/docker/3a-classic/server/3a-classic-server:/go/src \
  --name 3a-classic-server \
  --expose 80 \
  -p 8081:80 \
  golang:latest \
  bash -c /go/src/bin/3aClassic
```

* access

```
http://localhost:8081
```


### CORS

if you access by chrome, you have to set CORS.
this is example to configure CORS in nginx.

```
$ vim /etc/nginx/conf.d/proxy.conf

add_header 'Access-Control-Allow-Origin' '*';
add_header 'Access-Control-Allow-Credentials' 'true';
add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';
add_header 'Access-Control-Allow-Headers' 'DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type';
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
$ gb biuld
```

now you can run code compiled.

```
$ ./bin/3aClassic
```

## AUTHOR

 Sadayuki Matsuno