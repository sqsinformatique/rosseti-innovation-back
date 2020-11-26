package centrifugov1

import (
	"context"
	"time"

	"github.com/sqsinformatique/rosseti-innovation-back/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *CentrifugoV1) chatsDB() *mongo.Collection {
	mongoconn := *c.mongoDB
	return mongoconn.Database(c.config.Mongo.ChatDB).Collection("chats")
}

func (c *CentrifugoV1) GetChat(id int) (*models.ChatChannel, error) {
	filter := bson.D{{"id", id}}

	var result models.ChatChannel
	err := c.chatsDB().FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		c.log.Debug().Err(err).Msg("failed GetChat")
		return nil, err
	}

	return &result, nil
}

func (c *CentrifugoV1) SaveToDB(chatID, userID int, name, message string) error {
	chatChannel, err := c.GetChat(chatID)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if err == mongo.ErrNoDocuments {
		c.log.Debug().Msgf("not found chat: %d", chatID)
		chatChannel = &models.ChatChannel{
			ID:       chatID,
			Name:     name,
			Messages: []*models.Message{},
		}
		_, err1 := c.chatsDB().InsertOne(context.TODO(), chatChannel)
		if err1 != nil {
			return err1
		}
	}

	chatChannel.LastMsgID += 1
	_, err = c.chatsDB().UpdateOne(
		context.TODO(),
		bson.M{"id": chatID},
		bson.M{"$push": bson.M{"messages": &models.Message{Sender: userID, Text: message, TimeStamp: time.Now(), ID: chatChannel.LastMsgID}}},
	)
	if err != nil {
		return err
	}

	_, err = c.chatsDB().UpdateOne(
		context.TODO(),
		bson.M{"id": chatID},
		bson.M{"$set": bson.M{"lastid": chatChannel.LastMsgID}},
	)
	if err != nil {
		return err
	}

	return nil
}
