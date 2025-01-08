package model

import (
	"time"

	"gorm.io/gorm"
)

// 迁移数据表，在没有数据结构变更的时候，建议注释不执行
// 只支持创建表，增加表中没有的字段和索引
// 为了保护数据，并不支持改变已有的字段类型或删除未被使用的字段
func MakeMigrate(db *gorm.DB) error {
	return nil
}

// 通用模型
type Model struct {
	ID       int       `gorm:"primary_key;auto_increment" json:"id"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}
