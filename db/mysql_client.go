package db

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
		panic(fmt.Errorf("error in read config file %w", err))
	}
}

func InitGormEngine() {
	dataSourceName := viper.GetString("mysql.dataSourceName")
	GormEngine, errNewEngine = gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if errNewEngine != nil {
		panic(fmt.Errorf("error in init new engine %w", errNewEngine))
	}

	SqlDB, errSqlDB := GormEngine.DB()
	if errSqlDB != nil {
		panic(fmt.Errorf("error in init sql db %w", errSqlDB))
	}

	//defer func(SqlDB *sql.DB) {
	//	err := SqlDB.Close()
	//	if err != nil {
	//		panic(fmt.Errorf("error in close sql db %w", errSqlDB))
	//	}
	//}(SqlDB)

	errPing := SqlDB.Ping()
	if errPing != nil {
		panic(fmt.Errorf("error on ping db: %w", errPing))
	}

	SqlDB.SetMaxOpenConns(50)
	SqlDB.SetMaxIdleConns(20)
	SqlDB.SetConnMaxLifetime(time.Hour)
	SqlDB.SetConnMaxIdleTime(30 * time.Minute)

	// create table
	//errCreateTable := GormEngine.AutoMigrate(&models.MyUser{})
	//if errCreateTable != nil {
	//	panic(fmt.Errorf("error in create table %w", errCreateTable))
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
	result, err := f(tx, in)
	if err != nil {
		return result, err
	}
	tx.Commit()
	return result, nil
}

func (engine *MySQLClientEngine) ReadOnlyTransaction(f func(tx *gorm.DB) (interface{}, error)) (interface{}, error) {
	tx := engine.Begin()
	result, err := f(tx)
	if err != nil {
		return result, err
	}
	tx.Rollback()
	return result, nil
}
