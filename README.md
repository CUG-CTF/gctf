# gctfs

run init.sh to get depents

still working...

## database
1. start local env :`docker run -d  -p 5432:5432 -e POSTGRES_PASSWORD=Ilovecug666  postgres 
`
2. set env vars
```sh 
GCTF_DB_DRIVER=postgres
GCTF_DB_STRING=postgres://postgres:Ilovecug666@dev.cugctf.top:5432/gctf?sslmode\=disable
```

## some resource

<https://docs.docker.com/engine/swarm/manage-nodes/>  
<https://www.cnblogs.com/franknihao/p/8490416.html>  
<https://www.cnblogs.com/aguncn/p/7058662.html>  
<https://docs.docker.com/develop/sdk/examples/>  
<https://godoc.org/github.com/docker/docker/client>  
<https://docs.docker.com/engine/api/v1.37/#>