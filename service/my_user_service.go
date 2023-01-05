package service

import (
	"github.com/engchina/golang-mysql-demo/models"
	"gorm.io/gorm"
	"time"
)

// GetMyUserList Get MyUser List
func GetMyUserList(tx *gorm.DB) (interface{}, error) {
	var allData []*models.MyUser
	allData, err := models.GetMyUserList(tx)
	if err != nil {
		return nil, err
	}
	return allData, nil
}

// InsertOrUpdate Insert or Update
func InsertOrUpdate(tx *gorm.DB, in interface{}) (interface{}, error) {
	myUser := in.(*models.MyUser)
	var myUserModel *models.MyUser
	myUserModel, has, err := models.GetMyUserInTxn(tx, myUser.UserId)
	if err != nil {
		return nil, err
	}

	var affected int64
	if !has {
		affected, err = myUser.InsertMyUserInTxn(tx)
	} else {
		myUserModel.Name = myUser.Name
		affected, err = myUserModel.UpdateMyUserInTxn(tx)
	}

	if err != nil {
		return -1, err
	}
	return affected, nil
}

// UpdateWithOptimisticLock Optimistic Lock
func UpdateWithOptimisticLock(tx *gorm.DB, in interface{}) (interface{}, error) {
	myUser := in.(*models.MyUser)
	var myUserModel *models.MyUser
	myUserModel, _, err := models.GetMyUserInTxn(tx, myUser.UserId)
	if err != nil {
		return nil, err
	}
	time.Sleep(5 * time.Second)

	myUserModel.Name = myUser.Name
	affected, err := myUserModel.UpdateMyUserInTxn(tx)
	if err != nil {
		return -1, err
	}
	return affected, nil
}

// UpdateWithPessimisticLock Pessimistic Lock
func UpdateWithPessimisticLock(tx *gorm.DB, in interface{}) (interface{}, error) {
	myUser := in.(*models.MyUser)
	var myUserModel *models.MyUser
	myUserModel, _, err := models.GetMyUserForUpdateInTxn(tx, myUser.UserId)
	if err != nil {
		return nil, err
	}
	time.Sleep(5 * time.Second)

	myUserModel.Name = myUser.Name
	affected, err := myUserModel.UpdateMyUserInTxn(tx)
	if err != nil {
		return -1, err
	}
	return affected, nil
}
