# top-coin

# Arch

The System is a composition of 3 GRPC services + GRPC Gateway. I choosed GRPC to speed up development and desgin process of these APIs and protobuffer 
is a very efficient dataformat to transfer almost any kind of data from service to service. 
Only the Gateway should be accessible from outside and all other service are not reachable from outside.

## Backend requirements

* [docker](https://www.docker.com/) - Build, Manage and Secure Your Apps Anywhere. Your Way.
* [docker-compose](https://docs.docker.com/compose/) - Compose is a tool for defining and running multi-container Docker applications. 
* [golang](https://golang.org/) - The Go Programming Language
* [golang mod](https://github.com/golang/go/wiki/Modules) - Go dependency management tool 

## Versions requirements
* golang **>=1.13.10**
* Docker **>=18.09.2**
* Docker-compose **>=1.21.0**

### Setup Linux

```bash
git clone git@github.com:donutloop/top-coin.git
cd ./top-coin
sudo make builddockerimages
sudo docker-compose up
```

### Example http call

```bash
http://localhost:8080/v1/topcoins?limit=100
```