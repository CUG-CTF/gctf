# gctfs
A CTF platform base on docker,every problem is pack with docker
run init.sh to get depents.

still working...

## require
I develop with docker 18.09,and api version is 1.39,latest postgres
## database
1. start local db :`docker run -d  -p 5432:5432 -e POSTGRES_PASSWORD=Ilovecug666  postgres 
`
2. see `conf.json`

## some resource
<https://godoc.org/github.com/fsouza/go-dockerclient>
<https://github.com/go-xorm/xorm>
<http://www.xorm.io/docs/>