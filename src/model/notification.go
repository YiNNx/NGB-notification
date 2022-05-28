package model

import (
	"time"
)

const (
	TypeMessage   = 1
	TypeComment   = 2
	TypeMentioned = 3
	TypeNewPost   = 4
)

type Notification struct {
	tableName struct{}

	Nid       int       `pg:",pk"`
	Time      time.Time `pg:"default:now()"`
	Uid       int
	Type      int
	ContentId int //私信为mid,关注人发帖和@为pid,评论为cid
	Status    int `pg:"default:0"` //0未读 1已读
}
