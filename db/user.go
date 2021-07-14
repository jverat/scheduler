package db

import (
	"fmt"
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
	query := fmt.Sprintf("INSERT INTO public.user (name, password) VALUES ('%s', '%s')", u.Name, u.Password)
	_, err = Connection.Query(ctx, query)

	if err != nil {
		return
	}

	user, err = Read(u)

	if err != nil {
		return
	}

	if createProfiles(user.ID, u.Profiles) != nil {
		return
	}

	return Read(user)

}

func Read(u User) (user User, err error) {
	query := "SELECT * FROM public.user"
	if u.ID == 0 {
		query += fmt.Sprintf(" WHERE name = '%s'", u.Name)
	} else {
		query += fmt.Sprintf(" WHERE id = %d", u.ID)
	}

	row := Connection.QueryRow(ctx, query)
	err = row.Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		return
	}

	return readProfiles(user)
}

func Update(u User) (user User, err error) {
	return
}

func Delete(u User) (user User, err error) {
	return
}
