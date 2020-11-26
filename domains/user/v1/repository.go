package userv1

import (
	"encoding/json"
	"errors"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/crypto"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/db"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/utils"
	"github.com/sqsinformatique/rosseti-innovation-back/models"
)

var (
	ErrMismatchedHashAndPassword = errors.New("hashedPassword is not the hash of the given password")
	ErrNewPasswordSameAsOld      = errors.New("new password same as old")
)

func (u *UserV1) CreateUser(request *models.NewCredentials) (*models.User, error) {
	// Normalize email
	email, err := utils.NormalizeEmail(request.Email)
	if err != nil {
		return nil, err
	}

	data := &models.User{
		Hash:  crypto.HashString(request.Password),
		Email: email,
		Phone: request.Phone,
		Role:  request.Role,
	}

	data.CreateTimestamp()

	result, err := u.orm.InsertInto("users", data)
	if err != nil {
		return nil, err
	}

	return result.(*models.User), nil
}

func (u *UserV1) GetUserByID(id int64) (data *models.User, err error) {
	data = &models.User{}

	conn := *u.db
	if conn == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	err = conn.Get(data, "select * from production.users where id=$1", id)
	if err != nil {
		return nil, err
	}

	u.log.Debug().Msgf("user %+v", data)

	return
}

func (u *UserV1) GetUserDataByCreds(c *models.Credentials) (data *models.User, err error) {
	data = &models.User{}

	conn := *u.db
	if conn == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	// Get user from DB
	if c.Phone != "" {
		err = conn.Get(data, "select * from production.users where user_phone=$1", c.Phone)
	} else if c.Email != "" {
		// Normalize email
		email, err1 := utils.NormalizeEmail(c.Email)
		if err1 != nil {
			return nil, err1
		}

		err = conn.Get(data, "select * from production.users where user_email=$1", email)
	}

	if err != nil {
		return nil, err
	}

	// Check password
	err = crypto.CompareHash(data.Hash, c.Password)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func mergeUserData(oldData *models.User, patch *[]byte) (newData *models.User, err error) {
	id := oldData.ID

	original, err := json.Marshal(oldData)
	if err != nil {
		return
	}

	merged, err := jsonpatch.MergePatch(original, *patch)
	if err != nil {
		return
	}

	// Save hash
	Hash := oldData.Hash

	err = json.Unmarshal(merged, &newData)
	if err != nil {
		return
	}

	// Normalize email
	email, err := utils.NormalizeEmail(newData.Email)
	if err != nil {
		return
	}
	newData.Email = email

	// Protect ID from changes
	newData.ID = id

	if newData.Hash == "" {
		newData.Hash = Hash
	}

	err = newData.Validate()
	if err != nil {
		return nil, err
	}

	newData.UpdatedAt.Time = time.Now()
	newData.UpdatedAt.Valid = true

	return newData, nil
}

func (u *UserV1) UpdateUserByID(id int64, patch *[]byte) (writeData *models.User, err error) {
	data, err := u.GetUserByID(id)
	if err != nil {
		return
	}

	writeData, err = mergeUserData(data, patch)
	if err != nil {
		return
	}

	if u.db == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	_, err = u.orm.Update("users", writeData)
	if err != nil {
		return nil, err
	}

	return writeData, err
}

func (u *UserV1) SoftDeleteUserByID(id int64) (err error) {
	data, err := u.GetUserByID(id)
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

	_, err = u.orm.Update("users", data)

	return
}

func (u *UserV1) HardDeleteUserByID(id int64) (err error) {
	conn := *u.db
	if conn == nil {
		return db.ErrDBConnNotEstablished
	}

	_, err = conn.Exec(conn.Rebind("DELETE FROM production.users WHERE user_id=$1"), id)

	if err != nil {
		return err
	}

	return nil
}

func (u *UserV1) UpdateUserCredsByID(id int64, c *models.UpdateCredentials) (data *models.User, err error) {
	data, err = u.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Check password
	if c.OldPassword != "" {
		err = crypto.CompareHash(data.Hash, c.OldPassword)
		if err != nil {
			return nil, err
		}
	}

	// Update password
	if c.Password != "" {
		// Checking that new password is not same as old password
		err = crypto.CompareHash(data.Hash, c.Password)
		if err != ErrMismatchedHashAndPassword {
			if err == nil {
				return nil, ErrNewPasswordSameAsOld
			}
			return nil, err
		}

		data.Hash = crypto.HashString(c.Password)
	}

	data.UpdatedAt.Time = time.Now()
	data.UpdatedAt.Valid = true

	if u.db == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	_, err = u.orm.Update("users", data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
