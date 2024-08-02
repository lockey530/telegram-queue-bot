package types

import (
	"time"
)

type QueueUser struct {
	QueueID    uint64    `db:"queue_id" note:"Postgres-generated identifier"`
	UserHandle string    `db:"user_handle" note:"Refers to the Telegram handle"`
	ChatID     int64     `db:"chat_id" note:"Refers to the ID used to identify message chats"`
	Joined_at  time.Time `db:"timezone" note:"Refers to the time the user joined a queue"`
}

type AdminUser struct {
	AdminID     uint64 `db:"admin_id" note:"serial count"`
	AdminHandle string `db:"admin_handle" note:"Refers to the Telegram handle"`
	Removable   bool   `db:"removable" note:"at least one admin should be un-removable to prevent loss of bot access."`
}
