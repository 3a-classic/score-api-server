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
docker run --name 3a-classic-mongo -d mongo
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
$ docker run -d -t --link 3a-classic-mongo:mongo -v /home/matsuno/docker/3a-classic/server/3a-classic-server:/go/src --name 3a-classic-server --expose 80 -p 8081:80 golang bash -c /go/src/bin/3aClassic
```

* access

```
http://localhost:8081
```

## CAUTION

* change 3a-classic-server info in 3a-classic-ui
* 3a-classic-server have to be reachable from client node
* data in mongo-docker will disappear even if you save the docker image


## AUTHOR

 Sadayuki Matsuno

