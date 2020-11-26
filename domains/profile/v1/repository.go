package profilev1

import (
	"encoding/json"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/crypto"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/db"
	"github.com/sqsinformatique/rosseti-innovation-back/models"
)

func (o *ProfileV1) CreateProfile(request *models.Profile) (*models.Profile, error) {

	request.CreateTimestamp()

	sign, err := crypto.GenerateSign()
	if err != nil {
		return nil, err
	}

	request.PrivateKey, request.PublicKey = crypto.MarshalSign(sign)
	o.log.Debug().Msgf("length privateKey %d length publicKey %d", len(request.PrivateKey), len(request.PublicKey))

	result, err := o.orm.InsertInto("profiles", request)
	if err != nil {
		return nil, err
	}

	return result.(*models.Profile), nil
}

func (o *ProfileV1) GetProfileByID(id int64) (data *models.Profile, err error) {
	data = &models.Profile{}

	conn := *o.db
	if o.db == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	err = conn.Get(data, "select * from production.profiles where id=$1", id)
	if err != nil {
		return nil, err
	}

	o.log.Debug().Msgf("user %+v", data)

	return
}

func (p *ProfileV1) SearchProfileByLastName(value *models.Search) (data *ArrayOfProfileData, err error) {
	conn := *p.db
	if p.db == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	rows, err := conn.Queryx(conn.Rebind("select * from production.profiles where user_last_name ilike $1"), value.Value+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data = &ArrayOfProfileData{}

	for rows.Next() {
		var item models.Profile

		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		*data = append(*data, item)
	}

	return data, nil
}

func mergeProfileData(oldData *models.Profile, patch *[]byte) (newData *models.Profile, err error) {
	id := oldData.ID

	original, err := json.Marshal(oldData)
	if err != nil {
		return
	}

	merged, err := jsonpatch.MergePatch(original, *patch)
	if err != nil {
		return
	}

	err = json.Unmarshal(merged, &newData)
	if err != nil {
		return
	}

	// Protect ID from changes
	newData.ID = id

	newData.UpdatedAt.Time = time.Now()
	newData.UpdatedAt.Valid = true

	return newData, nil
}

func (u *ProfileV1) UpdateProfileByID(id int64, patch *[]byte) (writeData *models.Profile, err error) {
	data, err := u.GetProfileByID(id)
	if err != nil {
		return
	}

	writeData, err = mergeProfileData(data, patch)
	if err != nil {
		return
	}

	if u.db == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	_, err = u.orm.Update("profiles", writeData)
	if err != nil {
		return nil, err
	}

	return writeData, err
}

func (u *ProfileV1) SoftDeleteProfileByID(id int64) (err error) {
	data, err := u.GetProfileByID(id)
	if err != nil {
		return
	}

	if data.DeletedAt.Valid {
		return
	}

	data.DeletedAt.Time = time.Now()
	data.DeletedAt.Valid = true
	data.UpdatedAt.Time = time.Now()
	data.UpdatedAt.Valid = true

	if u.db == nil {
		return db.ErrDBConnNotEstablished
	}

	_, err = u.orm.Update("profiles", data)

	return
}

func (u *ProfileV1) HardDeleteProfileByID(id int64) (err error) {
	conn := *u.db
	if conn == nil {
		return db.ErrDBConnNotEstablished
	}

	_, err = conn.Exec(conn.Rebind("DELETE FROM production.profiles WHERE id=$1"), id)

	if err != nil {
		return err
	}

	return nil
}
