package common

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB

type Authenticated struct {
	gorm.Model
	Snowflake    string
	Username     string
	AccessToken  string
	RefreshToken string
}

func OpenDatabase(dsn string) (err error) {
	database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return database.AutoMigrate(&Authenticated{})
}

func Find(conds ...interface{}) (authenticated *Authenticated, err error) {
	if err = database.Where(conds[0], conds[1]).First(&authenticated).Error; err != nil {
		return nil, err
	}
	return authenticated, nil
}

func Create(authenticated Authenticated) error {
	return database.Create(&authenticated).Error
}
