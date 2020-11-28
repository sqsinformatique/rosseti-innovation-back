package centrifugov1

import (
	"context"
	"time"

	"github.com/sqsinformatique/rosseti-innovation-back/internal/db"
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

func (c *CentrifugoV1) CreateTheme(request *models.Theme) (*models.Theme, error) {

	request.CreateTimestamp()

	result, err := c.orm.InsertInto("theme", request)
	if err != nil {
		return nil, err
	}

	return result.(*models.Theme), nil
}

func (c *CentrifugoV1) SelectAllDirections() (data *ArrayOfDirectionData, err error) {
	conn := *c.db
	if c.db == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	rows, err := conn.Queryx("select * from production.direction")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data = &ArrayOfDirectionData{}

	for rows.Next() {
		var item models.Direction

		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		*data = append(*data, item)
	}

	return data, nil
}

func (c *CentrifugoV1) SelectThemesByDirection(id int) (data *ArrayOfThemesData, err error) {
	conn := *c.db
	if c.db == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	rows, err := conn.Queryx(conn.Rebind("select * from production.theme where direction=$1"), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data = &ArrayOfThemesData{}

	for rows.Next() {
		var item models.Theme

		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		*data = append(*data, item)
	}

	return data, nil
}

func (c *CentrifugoV1) SelectLastActiveThemes() (data *ArrayOfThemesData, err error) {
	conn := *c.db
	if c.db == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	query := "("
	i := 0
	for k := range c.lastActiveThemesMap {
		if i < len(c.lastActiveThemesMap)-1 {
			query = query + k + ","
		} else {
			query = query + k
		}
		i++
	}

	query = query + ")"

	c.log.Debug().Msgf("query: %s", "select * from production.theme where id in "+query)

	rows, err := conn.Queryx("select * from production.theme where id in " + query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data = &ArrayOfThemesData{}

	for rows.Next() {
		var item models.Theme

		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		*data = append(*data, item)
	}

	return data, nil
}
