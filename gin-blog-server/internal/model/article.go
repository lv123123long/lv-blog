package model
import (
	"time"

	"gorm.io/gorm"
)


const (
	STATUS_PUBLIC = iota + 1  // 公开
	STATUS_SECRET             // 私密
	STATUS_DRAFT              // 草稿
)

const (
	TYPE_ORIGINAL = iota + 1 //原创
	TYPE_REPRINT             // 转载
	TYPE_TRANSLATE           // 翻译
)

type Article struct {
	Model 

}