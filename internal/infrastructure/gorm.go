package infrastructure

import (
	"github.com/glebarez/sqlite"
	"go-clean-architecture/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var db *gorm.DB

func dbSetup() {
	var err error
	l := gormLogger.Default.LogMode(gormLogger.Silent)
	if cfg.Database.Driver == "mysql" {
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN: cfg.Database.DSN,
		}), &gorm.Config{
			Logger: l,
		})
	} else {
		db, err = gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{
			Logger: l,
		})
	}

	if err != nil {
		panic(err)
	}

	if cfg.IsDevelopment {
		if err := db.AutoMigrate(
			&domain.Author{},
			&domain.Article{},
		); err != nil {
			panic(err)
		}
	}

	var count int64
	if err := db.Model(&domain.Author{}).Count(&count).Error; err != nil {
		panic(err)
	}

	if count == 0 {
		author := domain.Author{
			Name: "John Doe",
		}
		if err := db.Create(&author).Error; err != nil {
			panic(err)
		}
	}

}
