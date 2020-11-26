package sessionv1

import (
	"github.com/sqsinformatique/rosseti-innovation-back/internal/db"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/utils"
	"github.com/sqsinformatique/rosseti-innovation-back/models"
)

func (s *SessionV1) CreateSession(id int) (*models.Session, error) {
	var request models.Session

	seq, err := utils.RuneSequence(100, utils.AlphaNum)
	if err != nil {
		return nil, err
	}
	request.ID = string(seq)
	request.UserID = id
	request.CreateTimestamp()

	result, err := s.orm.InsertInto("sessions", &request)
	if err != nil {
		return nil, err
	}

	return result.(*models.Session), nil
}

func (s *SessionV1) GetSession(id string) (data *models.Session, err error) {
	data = &models.Session{}

	conn := *s.db
	if conn == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	err = conn.Get(data, "select * from production.sessions where id=$1", id)
	if err != nil {
		return nil, err
	}

	s.log.Debug().Msgf("session %+v", data)

	return
}

func (s *SessionV1) DeleteSession(id string) (err error) {
	conn := *s.db
	if conn == nil {
		return db.ErrDBConnNotEstablished
	}

	_, err = conn.Exec(conn.Rebind("DELETE FROM production.sessions WHERE id=$1"), id)

	if err != nil {
		return err
	}

	return nil
}
