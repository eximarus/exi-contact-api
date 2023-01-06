package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.22

import (
	"context"
	"fmt"

	"github.com/eximarus/exi-contact-api/pkg/graph/model"
)

// CreateGuestbookEntry is the resolver for the createGuestbookEntry field.
func (r *mutationResolver) CreateGuestbookEntry(ctx context.Context, input model.CreateGuestbookEntryInput) (*model.GuestbookEntry, error) {
	return r.createGuestbookEntryImpl(ctx, &input)
}

// SubmitContactInfo is the resolver for the submitContactInfo field.
func (r *mutationResolver) SubmitContactInfo(ctx context.Context, input model.ContactInfoInput) (*bool, error) {
	return r.submitContactInfoImpl(ctx, &input)
}

// GetGuestbook is the resolver for the getGuestbook field.
func (r *queryResolver) GetGuestbook(ctx context.Context, cursor *string, limit *int) (*model.GetGuestbookOutput, error) {
	return r.getGuestbookImpl(ctx, cursor, limit)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) Guestbook(ctx context.Context) ([]*model.GuestbookEntry, error) {
	panic(fmt.Errorf("not implemented: Guestbook - guestbook"))
}