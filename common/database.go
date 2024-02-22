package common

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
)

var database *gorm.DB

type Authenticated struct {
	gorm.Model
	Snowflake    string
	Username     string
	AccessToken  string
	RefreshToken string
}

func OpenDatabase(driver, dsn string) (err error) {
	var dialector gorm.Dialector
	if strings.EqualFold(driver, "postgres") {
		dialector = postgres.Open(dsn)
	} else if strings.EqualFold(driver, "sqlite") {
		dialector = sqlite.Open(dsn)
	} else {
		return fmt.Errorf("invalid driver")
	}

	database, err = gorm.Open(dialector, &gorm.Config{})
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
