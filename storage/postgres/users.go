package postgres

import (
	"database/sql"
	pb "ecoswap/genproto/users"
	"ecoswap/model"
	"fmt"
	"log"
	"math"
	"time"

	"gopkg.in/gomail.v2"
)

type UsersRepo struct {
	Db *sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{Db: db}
}

func (U *UsersRepo) Register(req *pb.UserRegister) error {
	query := `
				INSERT INTO auth_service_users(
					username, email, password, full_name)
				VALUES
					($1, $2, $3, $4)`
	_, err := U.Db.Exec(query, req.Username, req.Email, req.Password, req.FullName)
	return err
}

func (U *UsersRepo) GetUserByEmail(email string) (model.InfoUser, error) {
	resp := model.InfoUser{}
	query := `
				SELECT 
					id, username, password, full_name 
				FROM
					auth_service_users
				WHERE
					email = $1 AND deleted_at is null`
	err := U.Db.QueryRow(query, email).Scan(&resp.Id, &resp.Username, &resp.Password, resp.FullName)
	return resp, err
}

func (U *UsersRepo) StoreRefreshToken(req *model.RefreshToken) error {
	query := `
				INSERT INTO refresh_token(
					user_id, token, expires_at)
				VALUES
					($1, $2, $3)`
	_, err := U.Db.Exec(query, req.UserId, req.Token, req.ExpiresAt)
	return err
}

func (U *UsersRepo) GetProfile(userId *pb.UserId) (*pb.UserInfo, error) {
	resp := pb.UserInfo{}
	query := `
				SELECT 
					id, username, email, full_name, eco_points, created_at, updated_at
				FROM 
					auth_service_users
				WHERE 
					id = $1 AND deleted_at is null`
	err := U.Db.QueryRow(query, userId.Id).Scan(&resp.Id, &resp.Username, &resp.Email, &resp.FullName,
		&resp.EcoPoints, &resp.CreatedAt, &resp.UpdatedAt)
	return &resp, err
}

func (U *UsersRepo) UpdateProfile(req *pb.ProfileUpdate) (*pb.UpdateResponse, error) {
	resp := pb.UpdateResponse{}
	query := `
				UPDATE auth_service_users SET 
					username = $1, bio = $2, full_name = $3, updated_at = $4
				WHERE
					id = $5 AND deleted_at is null
				RETURNING
					id, username, email, full_name, updated_at`
	err := U.Db.QueryRow(query, req.Username, req.Bio, req.FullName, time.Now(), req.Id).Scan(
		&resp.Id, &resp.Username, &resp.Email, &resp.FullName, &resp.UpdatedAt)
	return &resp, err
}

func (U *UsersRepo) DeleteProfile(userId *pb.UserId) (*pb.Status, error) {
	query := `
				UPDATE auth_service_users SET
					deleted_at = $1
				WHERE 
					id = $2 AND deleted_at is null`
	result, err := U.Db.Exec(query, time.Now(), userId.Id)
	if err != nil {
		return &pb.Status{
			Status:  false,
			Message: "Ma'lumotlar o'chirilmadi",
		}, err
	}
	del, err := result.RowsAffected()
	if err != nil || del == 0 {
		return &pb.Status{
			Status:  false,
			Message: "Ma'lumotlar o'chirilmadi",
		}, err
	}
	return &pb.Status{
		Status:  true,
		Message: "Ma'lumotlar muvaffaqiyatli o'chirildi.",
	}, err
}

func (U *UsersRepo) GetAllUsers(req *pb.FilterField) (*pb.Users, error) {
	var total int32
	query := `
				SELECT 
					id, count(id), username, full_name, eco_points
				FROM
					auth_service_users
				WHERE
					deleted_at is null`
	arr := []interface{}{}

	if len(req.FullName) > 0 {
		query += " AND full_name = $1"
		arr = append(arr, req.FullName)
	}

	err := U.Db.QueryRow(query, arr...).Scan(nil, &total, nil, nil, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if req.Limit > 0 {
		query += fmt.Sprintf(" limit %s", req.Limit)
	} else {
		req.Limit = total
	}
	if req.Offset > 0 {
		query += fmt.Sprintf(" offset %s", req.Offset)
	}

	users := []*pb.User{}
	rows, err := U.Db.Query(query, arr...)
	if err != nil {
		return &pb.Users{Users: users}, err
	}
	for rows.Next() {
		var user pb.User
		err := rows.Scan(&user.Id, nil, &user.Username, &user.FullName, &user.EcoPoints)
		if err != nil {
			return &pb.Users{Users: users}, err
		}
		users = append(users, &user)
	}
	return &pb.Users{
		Users: users,
		Total: total,
		Page:  int32(math.Ceil(float64(total / req.Limit))),
		Limit: req.Limit,
	}, nil
}

func (U *UsersRepo) ResetPassword(email *pb.Email) (*pb.Status, error) {
	mail := gomail.NewMessage()
	resp, err := U.GetUserByEmail(email.Email)
	if err != nil {
		log.Println(err)
		return &pb.Status{
			Status:  false,
			Message: "Bazadan o'qishda xatolik yuz berdi",
		}, err
	}
	mail.SetHeader("From", "nuriddinovhamidjon2@gmail.com")
	mail.SetHeader("To", email.Email)
	mail.SetHeader("Subject", "EcoSwap dasturiga kirish kodingiz")

	mail.SetBody("Password")
}

func (U *UsersRepo) UpdateToken() {

}

func (U *UsersRepo) CancelToken(id *pb.UserId) {

}

func (U *UsersRepo) GetEcoPointsByUser(userId *pb.UserId) (*pb.UserEcoPoints, error) {
	resp := pb.UserEcoPoints{}
	query := `
				SELECT 
					id, eco_points, updated_at
				FROM 
					auth_service_users
				WHERE 
					id = $1 AND deleted_at is null`
	err := U.Db.QueryRow(query, userId.Id).Scan(&resp.UserId, &resp.EcoPoints, &resp.LastUpdated)
	return &resp, err
}
