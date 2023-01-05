package models

import (
	_ "github.com/godror/godror"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/optimisticlock"
	"time"
)

const (
	MyUserTableName = "my_user"
)

type MyUser struct {
	UserId     string                 `json:"userId"  gorm:"type:varchar(200);primaryKey"   form:"userId"`
	Name       string                 `json:"name"    gorm:"type:varchar(200);not null"     form:"name"`
	NumOfTried optimisticlock.Version `json:"version" gorm:"index"`
	Created    time.Time              `json:"created" gorm:"autoCreateTime"`
	Updated    time.Time              `json:"updated" gorm:"autoUpdateTime"`
	Deleted    gorm.DeletedAt         `json:"deleted" gorm:"index"`
}

func (*MyUser) TableName() string {
	return MyUserTableName
}

func GetMyUserList(tx *gorm.DB) ([]*MyUser, error) {
	allData := make([]*MyUser, 0)
	find := tx.Table(MyUserTableName).Order("user_id").Find(&allData)
	if find.Error != nil {
		return nil, find.Error
	}
	return allData, nil
}

func GetMyUserInTxn(tx *gorm.DB, userId string) (*MyUser, bool, error) {
	var has = false
	myUser := new(MyUser)
	first := tx.Table(MyUserTableName).First(myUser, userId)
	if first.Error != nil {
		if first.Error.Error() == "record not found" {
			return nil, false, nil
		} else {
			return nil, false, first.Error
		}
	}
	if first.RowsAffected == 1 {
		has = true
	}
	return myUser, has, nil
}

func GetMyUserForUpdateInTxn(tx *gorm.DB, userId string) (*MyUser, bool, error) {
	var has = false
	myUser := new(MyUser)
	first := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table(MyUserTableName).First(myUser, userId)
	if first.Error != nil {
		if first.Error.Error() == "record not found" {
			return nil, false, nil
		} else {
			return nil, false, first.Error
		}
	}
	if first.RowsAffected == 1 {
		has = true
	}
	return myUser, has, nil
}

func (myUser *MyUser) InsertMyUserInTxn(tx *gorm.DB) (int64, error) {
	create := tx.Table(MyUserTableName).Create(myUser)
	if create.Error != nil {
		return -1, create.Error
	}
	return create.RowsAffected, nil
}

func (myUser *MyUser) UpdateMyUserInTxn(tx *gorm.DB) (int64, error) {
	updates := tx.Table(MyUserTableName).Where("user_id = ?", myUser.UserId).Updates(myUser)
	if updates.Error != nil {
		return -1, updates.Error
	}
	return updates.RowsAffected, nil
}
