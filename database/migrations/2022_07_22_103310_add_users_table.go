package migrations

import (
	"database/sql"
	"goweb/app/models"
	"goweb/pkg/migrate"

	"gorm.io/gorm"
)

type User struct {
	models.BaseModel

	Name     string `gorm:"type:varchar(191);not null;index"`
	Email    string `gorm:"type:varchar(191);index;default:null"`
	Phone    string `gorm:"type:varchar(20);index;default:null"`
	Password string `gorm:"type:varchar(191)"`

	models.CommonTimestampsField
}

// 指定表名，不指定的话，默认是结构体名字+s
func (User) TableName() string {
	return "users"
}

func init() {

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.AutoMigrate(&User{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.DropTable(&User{})
	}

	migrate.Add("2022_07_22_103310_add_users_table", up, down)
}
