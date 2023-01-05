package facade

import (
	"github.com/engchina/golang-mysql-demo/db"
	"github.com/engchina/golang-mysql-demo/models"
	"github.com/engchina/golang-mysql-demo/service"
)

func GetMyUserList() (interface{}, error) {
	return db.MySQLClient.ReadOnlyTransaction(service.GetMyUserList)
}

func InsertOrUpdate(myUser *models.MyUser) (interface{}, error) {
	return db.MySQLClient.ReadWriteTransaction(service.InsertOrUpdate, myUser)
}

func UpdateWithOptimisticLock(myUser *models.MyUser) (interface{}, error) {
	return db.MySQLClient.ReadWriteTransaction(service.UpdateWithOptimisticLock, myUser)
}

func UpdateWithPessimisticLock(myUser *models.MyUser) (interface{}, error) {
	return db.MySQLClient.ReadWriteTransaction(service.UpdateWithPessimisticLock, myUser)
}
