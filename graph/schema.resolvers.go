package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"
	"scheduler/db"
	"scheduler/graph/generated"
	"scheduler/graph/model"

	"golang.org/x/crypto/bcrypt"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	u, err := db.Create(db.User{Name: input.Name, Password: input.Password})
	if err != nil {
		return nil, err
	}

	ctx.Value("IP")
	res := &model.User{
		ID:       u.ID,
		Name:     u.Name,
		Password: nil,
	}

	return res, nil
}

func (r *mutationResolver) LogUser(ctx context.Context, input model.LoginUser) (bool, error) {
	u, err := db.Read(db.User{Name: input.Name})
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, ErrWrongPassword
	} else if err != nil {
		return false, err
	}
	log.Print("user " + u.Name + " logged in")
	return true, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUser) (bool, error) {
	err := db.Update(db.User{
		ID:       input.ID,
		Name:     input.Name,
		Password: input.Password,
	})
	return err == nil, err
}

func (r *mutationResolver) DeleteUser(ctx context.Context, input model.DeleteUser) (bool, error) {
	err := db.Delete(db.User{ID: input.ID})

	return err == nil, err
}

func (r *mutationResolver) CreateProfile(ctx context.Context, input model.NewProfile) (*model.Profile, error) {
	pfs, err := db.CreateProfiles(input.UID, db.Profiles{{
		Name:                  input.Name,
		WorkblockDuration:     input.WorkblockDuration,
		RestblockDuration:     input.RestblockDuration,
		LongRestblockDuration: input.LongRestblockDuration,
		NWorkblocks:           input.NWorkblocks,
	}})

	if err != nil {
		return nil, err
	}
	var pf db.Profile
	for _, p := range pfs {
		if p.Name == input.Name {
			pf = p
		}
	}

	if pf.ID == 0 {
		return nil, ErrProfileNotCreated
	}

	res := &model.Profile{
		ID:                    pf.ID,
		Name:                  pf.Name,
		WorkblockDuration:     pf.WorkblockDuration,
		RestblockDuration:     pf.RestblockDuration,
		LongRestblockDuration: pf.LongRestblockDuration,
		NWorkblocks:           pf.NWorkblocks,
	}
	log.Print("?")
	return res, nil
}

func (r *mutationResolver) CreateProfiles(ctx context.Context, input []*model.NewProfile) ([]*model.Profile, error) {
	var profiles db.Profiles
	if len(input) == 0 {
		return nil, ErrNilProfilesArray
	}

	for _, in := range input {
		profiles = append(profiles, db.Profile{
			Name:                  in.Name,
			WorkblockDuration:     in.WorkblockDuration,
			RestblockDuration:     in.RestblockDuration,
			LongRestblockDuration: in.LongRestblockDuration,
			NWorkblocks:           in.NWorkblocks,
		})
	}

	profiles, err := db.CreateProfiles(input[0].UID, profiles)
	if err != nil {
		return nil, err
	}

	var res []*model.Profile
	var profilesNotCreated []interface{}
	for _, pf := range profiles {
		if pf.ID == 0 {
			profilesNotCreated = append(profilesNotCreated)
		} else {
			res = append(res, &model.Profile{
				ID:                    pf.ID,
				Name:                  pf.Name,
				WorkblockDuration:     pf.WorkblockDuration,
				RestblockDuration:     pf.RestblockDuration,
				LongRestblockDuration: pf.LongRestblockDuration,
				NWorkblocks:           pf.NWorkblocks,
			})
		}
	}
	if len(profilesNotCreated) != 0 {
		errS := "some profiles were not created: "
		for range profilesNotCreated {
			errS += "%s "

		}
		return nil, fmt.Errorf(errS, profilesNotCreated...)
	}
	return res, nil
}

func (r *mutationResolver) UpdateProfiles(ctx context.Context, input []*model.UpdateProfile) (bool, error) {
	_, err := db.Read(db.User{ID: input[0].UID})
	if err != nil {
		return false, err
	}
	if len(input) == 0 {
		return false, nil
	}
	pfs := db.Profiles{}
	for _, p := range input {
		pfs = append(pfs, db.Profile{
			ID:                    p.ID,
			Name:                  p.Name,
			WorkblockDuration:     p.WorkblockDuration,
			RestblockDuration:     p.RestblockDuration,
			LongRestblockDuration: p.LongRestblockDuration,
			NWorkblocks:           p.NWorkblocks,
		})
	}
	oldPfs, err := db.ReadProfiles(input[0].UID)
	if err != nil {
		return false, err
	}
	err = db.UpdateProfiles(input[0].UID, pfs, oldPfs)
	return err == nil, err
}

func (r *mutationResolver) DeleteProfile(ctx context.Context, input model.DeleteProfile) (bool, error) {
	_, err := db.Read(db.User{ID: input.UID})
	if err != nil {
		return false, err
	}
	pfs, err := db.ReadProfiles(input.UID)
	if err != nil {
		return false, err
	}
	for _, p := range pfs {
		if p.ID == input.ID {
			err = db.DeleteProfile(db.Profile{ID: input.ID})
			return err == nil, err
		}
	}

	return false, ErrProfileNotFound
}

func (r *queryResolver) User(ctx context.Context, id int) (*model.User, error) {
	u, err := db.Read(db.User{ID: id})
	if err != nil {
		return nil, err
	}

	res := &model.User{
		ID:   u.ID,
		Name: u.Name,
	}

	return res, nil
}

func (r *queryResolver) Profiles(ctx context.Context, userID int) ([]*model.Profile, error) {
	pfs, err := db.ReadProfiles(userID)
	if err != nil {
		return nil, err
	}
	var res []*model.Profile
	for _, pf := range pfs {
		res = append(res, &model.Profile{
			ID:                    pf.ID,
			Name:                  pf.Name,
			WorkblockDuration:     pf.WorkblockDuration,
			RestblockDuration:     pf.RestblockDuration,
			LongRestblockDuration: pf.LongRestblockDuration,
			NWorkblocks:           pf.NWorkblocks,
		})
	}
	return res, nil
}

func (r *queryResolver) Profile(ctx context.Context, id int, userID int) (*model.Profile, error) {
	_, err := db.Read(db.User{ID: userID})
	if err != nil {
		return nil, err
	}
	pfs, err := r.Profiles(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, pf := range pfs {
		if pf.ID == id {
			return pf, nil
		}
	}
	return nil, ErrProfileNotFound
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
