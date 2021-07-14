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

func Create(u *User) (err error) {

	return
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

	return getProfiles(user)
}

func Update(u *User) (err error) {
	return
}

func Delete(u *User) (err error) {
	return
}
