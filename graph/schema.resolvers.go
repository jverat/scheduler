package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"scheduler/db"
	"scheduler/graph/generated"
	"scheduler/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (res *model.User, err error) {
	u, err := db.Create(db.User{Name: input.Name, Password: input.Password})
	if err != nil {
		return
	}

	res = &model.User{
		ID:       u.ID,
		Name:     u.Name,
		Password: nil,
	}

	return
}
func (r *mutationResolver) LogUser(ctx context.Context, input model.LoginUser) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateProfile(ctx context.Context, input *model.NewProfile) (*model.Profile, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, id int) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Profiles(ctx context.Context, userID int) ([]*model.Profile, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Profile(ctx context.Context, id int, userID int) (*model.Profile, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
