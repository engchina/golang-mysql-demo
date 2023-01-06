package db

import (
	"database/sql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type MySQLClientEngine struct {
	*gorm.DB
}

var (
	GormEngine   *gorm.DB
	errNewEngine error
	SqlDB        *sql.DB
	MySQLClient  MySQLClientEngine
)

func InitConfig() {
	viper.SetConfigName("application")
	viper.AddConfigPath("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("error in read config file", err)
	}
}

func InitGormEngine() {
	dataSourceName := viper.GetString("mysql.dataSourceName")
	GormEngine, errNewEngine = gorm.Open(mysql.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if errNewEngine != nil {
		log.Fatalln("error in init database connection", errNewEngine)
	}

	SqlDB, errSqlDB := GormEngine.DB()
	if errSqlDB != nil {
		log.Fatalln("error in init sql db", errSqlDB)
	}

	SqlDB.SetMaxOpenConns(50)
	SqlDB.SetMaxIdleConns(20)
	SqlDB.SetConnMaxLifetime(time.Hour)
	SqlDB.SetConnMaxIdleTime(30 * time.Minute)

	//create table
	//errCreateTable := GormEngine.AutoMigrate(&models.MyUser{})
	//if errCreateTable != nil {
	//	log.Fatalln("error in create table", errCreateTable)
	//}
}

func InitMySQLClient() {
	MySQLClient.DB = GormEngine
}

func init() {
	InitConfig()
	InitGormEngine()
	InitMySQLClient()
}

func (engine *MySQLClientEngine) ReadWriteTransaction(f func(tx *gorm.DB, in interface{}) (interface{}, error), in interface{}) (interface{}, error) {
	tx := engine.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	result, err := f(tx, in)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return result, nil
}

func (engine *MySQLClientEngine) ReadOnlyTransaction(f func(tx *gorm.DB) (interface{}, error)) (interface{}, error) {
	tx := engine.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	result, err := f(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Rollback()
	return result, nil
}
