package postgres

import (
	"database/sql"
	pb "ecoswap/genproto/users"
	"ecoswap/model"
	"fmt"
	"strconv"
	"strings"
	"time"
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

func (U *UsersRepo) StoreRefreshToken(req *model.RefreshToken)error{
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
	query := `
				SELECT 
					id, username, full_name, eco_points
				FROM
					auth_service_users
				WHERE
					deleted_at is null`
	param := []string{}
	arr := []interface{}{}

	if len(req.FullName) > 0 {
		query += " AND full_name = :full_name"
		param = append(param, ":full_name")
		arr = append(arr, req.FullName)
	}
	if len(req.Limit) > 0 {
		query += fmt.Sprintf(" limit %s", req.Limit)
	}
	if len(req.Offset) > 0 {
		query += fmt.Sprintf(" offset %s", req.Offset)
	}

	for i, j := range param {
		query = strings.Replace(query, j, "$"+strconv.Itoa(i+1), 1)
	}

	users := []*pb.User{}
	rows, err := U.Db.Query(query, arr...)
	if err != nil {
		return &pb.Users{Users: users}, err
	}
	for rows.Next() {
		var user pb.User
		err := rows.Scan(&user.Id, &user.Username, &user.FullName, &user.EcoPoints)
		if err != nil {
			return &pb.Users{Users: users}, err
		}
		users = append(users, &user)
	}
	return &pb.Users{Users: users}, nil
}

func (U *UsersRepo) ResetPassword() {

}

func (U *UsersRepo) UpdateToken() {

}

func (U *UsersRepo) CancelToken(id *pb.UserId){
	
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
