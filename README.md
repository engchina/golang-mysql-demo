# golang-oracle-demo

### Start MySQL 8.0
```shell
mkdir /u01/data/mysql8; chmod 777 /u01/data/mysql8
docker run -d --name mysql8.0 --restart=always -v /u01/data/mysql8:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=oracle -e MYSQL_DATABASE=oracle -e MYSQL_USER=oracle -e MYSQL_PASSWORD=oracle -p 3308:3306 m
ysql:8.0.31

docker exec -it mysql8.0 mysql -uroot -poracle -e "alter user 'oracle'@'%' identified with mysql_native_password by 'oracle';"
```

### Usage

Open [http://localhost:3001](http://localhost:3001) in your browser.

### Reference:
- https://getbootstrap.com/docs/5.2/getting-started/introduction/
- https://pkg.go.dev/crypto/md5
- https://xorm.io/
- https://pkg.go.dev/github.com/godror/godror
- https://pkg.go.dev/github.com/gin-gonic/gin
- https://pkg.go.dev/github.com/spf13/viper
- https://pkg.go.dev/github.com/sirupsen/logrus
- https://pkg.go.dev/github.com/swaggo/gin-swagger
- https://pkg.go.dev/github.com/asaskevich/govalidator
- https://pkg.go.dev/github.com/gomodule/redigo