package db

import (
	"fmt"
	"github.com/jackc/pgx/v4"
)

//User is intended to keep the personalized configurations of schedules
type User struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Profiles Profiles `json:"profiles,omitempty"`
}

type Users []User

func Create(u User) (user User, err error) {

	u.ID = 0

	query := fmt.Sprintf("INSERT INTO public.user (name, password) VALUES ('%s', '%s')", u.Name, u.Password)

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

	err = createProfiles(user.ID, u.Profiles)

	if err != nil {
		return
	}

	user.Profiles = u.Profiles

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

	user.Profiles, err = readProfiles(user.ID)
	return
}

func Update(u User) (err error) {
	query := fmt.Sprintf("UPDATE public.user SET name = '%s', password = '%s' WHERE id = %d", u.Name, u.Password, u.ID)
	conn, err := Connection.Acquire(ctx)
	if err != nil {
		return
	}

	_, err = conn.Query(ctx, query)
	conn.Release()
	if err != nil {
		return
	}

	oldU, err := Read(u)
	if err != nil {
		return
	}

	if identicals(u.Profiles, oldU.Profiles) {
		return
	}

	err = updateProfiles(u.ID, u.Profiles, oldU.Profiles)

	return
}

func Delete(u User) (err error) {

	err = deleteProfiles(u.ID)
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
