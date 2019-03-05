# gctfs
A CTF platform base on docker,every problem is pack with docker
run init.sh to get depents.

still working...

## require
I develop with docker 18.09,and api version is 1.39,latest postgres
## database
1. start local db :`docker run -d  -p 5432:5432 -e POSTGRES_PASSWORD=Ilovecug666  postgres 
`
2. set env vars,if you are using goland,don't forget to del '\'
```sh 
GCTF_DB_DRIVER=postgres
GCTF_DB_STRING=postgres://postgres:Ilovecug666@localhost:5432/gctf?sslmode\=disable
```

## some resource

<https://www.cnblogs.com/franknihao/p/8490416.html>  
<https://www.cnblogs.com/aguncn/p/7058662.html>  
<https://docs.docker.com/develop/sdk/examples/>  
<https://godoc.org/github.com/docker/docker/client>  
<https://docs.docker.com/engine/api/v1.37/#>