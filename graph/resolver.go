package graph

import "fmt"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

var ErrWrongPassword = fmt.Errorf("wrong password")
var ErrProfileNotCreated = fmt.Errorf("profile not created")
var ErrProfileNotFound = fmt.Errorf("profile not found")
var ErrNilProfilesArray = fmt.Errorf("the array of profiles is empty")

type Resolver struct{}
