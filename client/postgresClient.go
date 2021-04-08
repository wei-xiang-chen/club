package client

import (
	"club/setting"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	DBEngine *gorm.DB
)

func NewDBEngine(databaseSetting *setting.DBSetting) (*gorm.DB, error) {
	db, err := gorm.Open(
		databaseSetting.DBtype,
		fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			databaseSetting.Host,
			databaseSetting.Part,
			databaseSetting.Username,
			databaseSetting.DBName,
			databaseSetting.Password,
		),
	)
	if err != nil {
		return nil, err
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)

	return db, nil
}
