package main

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var pool *pgxpool.Pool

func save(user *user) bool {
	commandTag, err := pool.Exec(context.Background(), "INSERT INTO users VALUES ($1::text, $2::text)",
		user.Username, user.Password)

	if err != nil {
		log.Error("Fail to save"+user.toString()+":", err)
		return false
	}

	return commandTag.Insert() && commandTag.RowsAffected() == 1
}

func delete(user *user) bool {
	commandTag, err := pool.Exec(context.Background(), "DELETE FROM users WHERE username = $1::text", user.Username)

	if err != nil {
		log.Error("Fail to delete "+user.toString()+":", err)
		return false
	}

	return commandTag.Delete() && commandTag.RowsAffected() == 1
}

func isNameExisted(user *user) bool {
	res := 0
	pool.QueryRow(context.Background(), "SELECT 1 FROM users WHERE username = $1::text", user.Username).Scan(&res)
	return res == 1
}

func getUser(userName string) *user {
	var user *user = &user{}
	r := pool.QueryRow(context.Background(), "SELECT * FROM users WHERE username = $1::text", userName)

	err := r.Scan(&user.Username, &user.Password)
	if err != nil {
		log.Error("Fail to find user with name "+userName, err)
		return nil
	}

	return user
}

func updatePassword(userName string, hashedNewPwd string) bool {
	log.Message(userName + " " + hashedNewPwd)
	commandTag, err := pool.Exec(context.Background(),
		"UPDATE users SET password = $1::text where username = $2::text", hashedNewPwd, userName)
	if err != nil {
		log.Error("Fail to change password:", err)
		return false
	}

	return commandTag.Update() && commandTag.RowsAffected() == 1
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func isUsernameAndPwdMatch(username string, rawPwd string) bool {
	var user *user
	user = getUser(username)
	if user == nil {
		return false
	}

	// Compare sent in password with saved user password hash (can use bcrypt package for pass hashing)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rawPwd))
	if err != nil {
		log.Error("Invalid password", err)
		return false
	}
	return true
}
