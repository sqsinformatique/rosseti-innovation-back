package models

import "time"

type Message struct {
	ID        int       `bson:"id"`
	Text      string    `bson:"text"`
	Sender    int       `bson:"sender"`
	TimeStamp time.Time `bson:"timestamp"`
}

type ChatChannel struct {
	ID        int        `bson:"id"`
	Name      string     `bson:"name"`
	Messages  []*Message `bson:"messages"`
	LastMsgID int        `bson:"lastid"`
}
