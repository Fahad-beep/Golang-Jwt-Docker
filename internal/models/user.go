package models

import "database/sql"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) FetchUserByEmail(email string) (*User, error) {
	user := &User{}
	query := "SELECT id, email, password_hash FROM users WHERE email = $1"
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) CreateUser(user *User) error {
	query := `INSERT INTO users (email, password_hash, age) VALUES ($1, $2, $3) RETURNING ID`
	err := r.DB.QueryRow(query, user.Email, user.Password, user.Age).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil

}
