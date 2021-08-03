package db

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

//User is intended to keep the personalized configurations of schedules
type User struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Profiles Profiles `json:"profiles,omitempty"`
}

type Users []User

var ErrUserNotFound = fmt.Errorf("user not found")

func hashPass(password string) (hash string, err error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(h), err
}

func Create(u User) (user User, err error) {

	u.ID = 0

	hashed, err := hashPass(u.Password)
	if err != nil {
		return
	}

	query := fmt.Sprintf("INSERT INTO public.user (name, password) VALUES ('%s', '%s')", u.Name, hashed)

	queryChan, outputChan, errorChan := make(chan string), make(chan pgx.Rows), make(chan error)
	go AcquireConn(queryChan, outputChan, errorChan)

	queryChan <- query
	close(queryChan)

	select {
	case err = <-errorChan:
		return
	case _ = <-outputChan:
	}

	user, err = Read(u)

	if err != nil {
		return
	}

	user.Profiles, err = CreateProfiles(user.ID, u.Profiles)

	if err != nil {
		return
	}

	return Read(user)
}

func Read(u User) (user User, err error) {
	query := "SELECT * FROM public.\"user\""
	if u.ID <= 0 {
		query += fmt.Sprintf(" WHERE name = '%s'", u.Name)
	} else {
		query += fmt.Sprintf(" WHERE id = %d", u.ID)
	}

	queryChan, outputChan, errChan := make(chan string), make(chan pgx.Rows), make(chan error)
	go AcquireConn(queryChan, outputChan, errChan)
	queryChan <- query
	close(queryChan)

	select {
	case err = <-errChan:
		if err != nil {
			return
		}
	case rows := <-outputChan:
		{
			var wtf Users
			var bugser User
			i := true

			for rows.Next() {
				if i {
					err = rows.Scan(&user.ID, &user.Name, &user.Password)
					if err != nil {
						return
					}
					i = false
				} else {
					err = rows.Scan(&bugser.ID, &bugser.Name, &bugser.Password)
					if err != nil {
						return
					}
					wtf = append(wtf, bugser)
				}
			}
			if len(wtf) > 0 {
				return user, fmt.Errorf("too many users found = %+v", wtf)
			}
		}
	}

	if user.ID == 0 {
		return User{}, ErrUserNotFound
	}

	user.Profiles, err = ReadProfiles(user.ID)
	return
}

func Update(u User) (err error) {
	hashed, err := hashPass(u.Password)
	if err != nil {
		return
	}
	query := fmt.Sprintf("UPDATE public.user SET name = '%s', password = '%s' WHERE id = %d", u.Name, hashed, u.ID)
	queryChan, outputChan, errChan := make(chan string), make(chan pgx.Rows), make(chan error)
	go AcquireConn(queryChan, outputChan, errChan)
	queryChan <- query
	close(queryChan)

	select {
	case err = <-errChan:
		if err != nil {
			return
		}
	case _ = <-outputChan:
	}

	return
}

func Delete(u User) (err error) {

	err = DeleteProfiles(u.ID)
	if err != nil {
		return
	}

	queryChan, outputChan, errorChan := make(chan string), make(chan pgx.Rows), make(chan error)
	go AcquireConn(queryChan, outputChan, errorChan)
	query := fmt.Sprintf("DELETE FROM public.\"user\" * WHERE id = %d", u.ID)
	queryChan <- query
	close(queryChan)
	select {
	case err = <-errorChan:
	case _ = <-outputChan:
	}
	return
}
