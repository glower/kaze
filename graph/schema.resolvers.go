package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.41

// // CreatePowerPlant is the resolver for the createPowerPlant field.
// func (r *mutationResolver) CreatePowerPlant(ctx context.Context, input model.NewPowerPlantInput) (*model.PowerPlant, error) {
// 	panic(fmt.Errorf("not implemented: CreatePowerPlant - createPowerPlant"))
// }

// // UpdatePowerPlant is the resolver for the updatePowerPlant field.
// func (r *mutationResolver) UpdatePowerPlant(ctx context.Context, id string, input model.UpdatePowerPlantInput) (*model.PowerPlant, error) {
// 	panic(fmt.Errorf("not implemented: UpdatePowerPlant - updatePowerPlant"))
// }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }