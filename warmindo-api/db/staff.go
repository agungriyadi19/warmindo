package db

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type User struct {
	ID        int    `json:"id,omitempty"`
	Password  string `json:"password,omitempty"`
	Email     string `json:"email,omitempty"`
	Name      string `json:"name,omitempty"`
	Username  string `json:"username,omitempty"`
	RoleID    int    `json:"role_id,omitempty"`
	Phone     string `json:"phone,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

const (
	GetUserByEmailQuery = `SELECT id, email, name, username, role_id, phone, password, created_at, updated_at FROM staffs WHERE email = $1;`
)

func GetUserByEmail(dbConn *sql.DB, email string) (*User, error) {
	var user User
	err := dbConn.QueryRow(GetUserByEmailQuery, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.Username, &user.RoleID, &user.Phone, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (user *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) UserExists(dbConn *sql.DB) bool {
	rows, err := dbConn.Query(GetUserByEmailQuery, user.Email)
	if err != nil || !rows.Next() {
		return false
	}

	return true
}
