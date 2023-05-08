package migrations

import (
	"database/sql"
	"goweb/app/models"
	"goweb/pkg/migrate"

	"gorm.io/gorm"
)

type Topic struct {
	models.BaseModel

	Title      string `gorm:"type:varchar(191);not null;index"`
	Body       string `gorm:"type:longtext;not null"`
	UserID     string `gorm:"type:bigint;not null;index"`
	CategoryID string `gorm:"type:bigint;not null;index"`

	// 会创建 user_id 和 category_id 外键的约束
	User     User
	Category Category

	models.CommonTimestampsField
}

// 指定表名，不指定的话，默认是结构体名字+s
func (Topic) TableName() string {
	return "topics"
}

func init() {
	type User struct {
		models.BaseModel
	}
	type Category struct {
		models.BaseModel
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.AutoMigrate(&Topic{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.DropTable(&Topic{})
	}

	migrate.Add("2022_07_26_140635_add_topics_table", up, down)
}
