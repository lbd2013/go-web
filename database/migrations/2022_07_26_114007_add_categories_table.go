package migrations

import (
	"database/sql"
	"goweb/app/models"
	"goweb/pkg/migrate"

	"gorm.io/gorm"
)

type Category struct {
	models.BaseModel

	Name        string `gorm:"type:varchar(191);not null;index"`
	Description string `gorm:"type:varchar(191);default:null"`

	models.CommonTimestampsField
}

// 指定表名，不指定的话，默认是结构体名字+s
func (Category) TableName() string {
	return "categorys"
}

func init() {

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.AutoMigrate(&Category{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.DropTable(&Category{})
	}

	migrate.Add("2022_07_26_114007_add_categories_table", up, down)
}
