package migrations

import (
	"database/sql"
	"goweb/app/models"
	"goweb/pkg/migrate"

	"gorm.io/gorm"
)

type Link struct {
	models.BaseModel

	Name string `gorm:"type:varchar(191);not null"`
	URL  string `gorm:"type:varchar(191);default:null"`

	models.CommonTimestampsField
}

// 指定表名，不指定的话，默认是结构体名字+s
func (Link) TableName() string {
	return "links"
}

func init() {

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.AutoMigrate(&Link{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.DropTable(&Link{})
	}

	migrate.Add("2022_07_27_102824_add_links_table", up, down)
}
