package postgres

import (
	"database/sql"
	pb "ecoswap/genproto/users"
	"ecoswap/model"
	"fmt"
	"log"
	"math"
	"time"

	"golang.org/x/crypto/bcrypt"
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
					username, email, password_hash, full_name)
				VALUES
					($1, $2, $3, $4)`
	_, err := U.Db.Exec(query, req.Username, req.Email, req.Password, req.FullName)
	return err
}

func (U *UsersRepo) GetUserByEmail(email string) (model.InfoUser, error) {
	resp := model.InfoUser{}
	query := `
				SELECT 
					id, username, password_hash, full_name 
				FROM
					auth_service_users
				WHERE
					email = $1 AND deleted_at is null`
	err := U.Db.QueryRow(query, email).Scan(&resp.Id, &resp.Username, &resp.Password, &resp.FullName)
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
	tx, err := U.Db.Begin()
	if err != nil{
		log.Println(err)
		return nil, err
	}
	query := `
				UPDATE auth_service_users SET
					deleted_at = $1
				WHERE 
					id = $2 AND deleted_at is null`
	result, err := tx.Exec(query, time.Now(), userId.Id)
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
	query = `
				UPDATE refresh_token SET 
					deleted_at = $1
				WHERE
					id = $2`
	_, err = tx.Exec(query, time.Now(), userId.Id)
	if err != nil{
		return nil, err
	}
	err = tx.Commit()
	if err != nil{
		log.Println(err)
		return nil, err
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
		query += fmt.Sprintf(" limit %d", req.Limit)
	} else {
		req.Limit = total
	}
	if req.Offset > 0 {
		query += fmt.Sprintf(" offset %d", req.Offset)
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
		Page:  int32(math.Ceil(float64(total) / float64(req.Limit))),
		Limit: req.Limit,
	}, nil
}

func (U *UsersRepo) ResetPassword(email *pb.Email) (*pb.Status, error) {
	resp, err := U.GetUserByEmail(email.Email)
	if err != nil || resp.Id == "" {
		log.Println(err)
		return &pb.Status{
			Status:  false,
			Message: "Bazadan o'qishda xatolik yuz berdi",
		}, err
	}
	mail := gomail.NewMessage()
	mail.SetHeader("From", "nuriddinovhamidjon2@gmail.com")
	mail.SetHeader("To", email.Email)
	mail.SetHeader("Subject", "Kodni yangilash uchun link")

	mail.SetBody("URL", "localhost:7777/swagger/auth/updatePass/index.html#/auth/post_auth_updatePass")

	d := gomail.NewDialer("smtp.gmail.com", 587, "nuriddinovhamidjon2@gmail.com", "qkrj oxld lshb dgte")

	if err := d.DialAndSend(mail); err != nil {
		return &pb.Status{
			Status:  false,
			Message: "Link yuborilmadi",
		}, err
	}
	return &pb.Status{
		Status:  true,
		Message: "Parolingizni yangilash uchun emailingizga ko'rsatma yuborildi",
	}, nil
}

func (U *UsersRepo) UpdatePassword(req *model.ResetPassword) (*pb.Status, error) {
	resp, err := U.GetUserByEmail(req.Email)
	if err != nil {
		log.Println(err)
		return &pb.Status{
			Status:  false,
			Message: "Bunday foydalanuvchi mavjud emas",
		}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(resp.Password), []byte(req.OldPassword)); err != nil {
		log.Println(err)
		return &pb.Status{
			Status:  false,
			Message: "Parol xato",
		}, err
	}

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	query := `
		UPDATE auth_service_users SET
			password_hash = $1, updated_at = $2
		WHERE 
			email = $3`
	_, err = U.Db.Exec(query, string(hashpassword), time.Now(), req.Email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.Status{
		Status:  true,
		Message: "Parolingiz yangilandi",
	}, nil
}


func(U *UsersRepo) UpdateToken(token *model.RefreshToken)error{
	query := `
				UPDATE refresh_token SET
					token = $2, expires_at = $3
				WHERE
					user_id = $1 and deleted_at is null`

	_, err := U.Db.Exec(query, token.UserId, token.Token, token.ExpiresAt)
	if err != nil{
		log.Println(err)
		return err
	}
	return nil
}

func (U *UsersRepo) CancelToken(userId string)(error){
	query := `
				UPDATE refresh_token SET
					deleted_at = $1
				WHERE 
					user_id = $2`
	_, err := U.Db.Exec(query, time.Now(), userId)
	return err
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


func(h *UsersRepo) CreateEcoPointsByUser(req *pb.CreateEcoPoints)(*pb.InfoUserEcoPoints, error){
	user, err := h.GetEcoPointsByUser(&pb.UserId{Id: req.UserId})
	if err != nil{
		log.Println(err)
		return nil, err
	}
	query := `
				UPDATE auth_service_users SET
					eco_points = $1, updated_at = $2
				WHERE 
					id = $3 AND deleted_at is null`
	_, err = h.Db.Exec(query, user.EcoPoints + req.EcoPoints, time.Now(), req.UserId)
	if err != nil{
		log.Println(err)
		return nil, err
	}
	return &pb.InfoUserEcoPoints{
		UserId: req.UserId,
		EcoPoints: user.EcoPoints + req.EcoPoints,
		AddedPoints: req.EcoPoints,
		Reason: req.Reason,
		Date: time.Now().Format("16-07-2024"),
	}, err
}
