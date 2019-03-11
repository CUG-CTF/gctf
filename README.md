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

### 单机部署

1. 你需要安装docker
2. 设置conf.json中datbase连接串
3. 启动database `docker run -d  -p 5432:5432 -e POSTGRES_PASSWORD=Ilovecug666  postgres `
4. 创建数据库`gctf`
5. 运行服务端

## some resource
<https://godoc.org/github.com/fsouza/go-dockerclient>  
<https://github.com/go-xorm/xorm>  
<http://www.xorm.io/docs/>  