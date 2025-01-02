package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.62

import (
	"context"
	"fmt"
	"ggg/database"
	"ggg/graph/model"
	"log"
	"time"

	"github.com/patrickmn/go-cache"

	// "github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/google/uuid"
)

// All Queries
// # Create User
// # mutation {
// #   createUser(input: { name: "Faiz", email: "mans@gmail.com", age: 23 }) {
// #     id
// #     name
// #     email
// #     age
// #   }
// # }

// # remove User
// 	# mutation{
// 	# deleteUser(id:"d1e167bf-35e9-4ba4-83d3-22182b9f77ad") {
// 	# id
// 	# }
// 	# }

// # All users
// # query{
// #   users{
// #     id
// #     name
// #     email
// #     age
// #   }
// # }

// # Update User by id
// mutation {
//   updateUser(id: "9e3e5791-0fca-4b19-b6ef-649e2fb94d94", input: { name: "Updated Name", email: "updatedemail@example.com", age: 35 }) {
//     id
//     name
//     email
//     age
//   }
// }

// # single user
// # query GetUser {
// #   user(id: "ef75e9c0-382e-46a2-b680-8875adf1b219") {
// #     id
// #     name
// #     email
// #     age
// #   }
// # }

// 5 min expiration time and 10 min clean interval
var c = cache.New(5*time.Minute, 10*time.Minute)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	// Create a new user
	// Example query:
	// mutation {
	// 	createUser(input: { name: "Faiz", email: "mans@gmail.com", age: 23 }) {
	// 	  id
	// 	  name
	// 	  email
	// 	  age
	// 	}
	//   }
	// id is automatically generated using uuid.NewString()
	//  gorm.model doesnt work with this , it says the id is useless

	user := model.User{
		ID:    uuid.NewString(),
		Name:  input.Name,
		Email: input.Email,
		Age:   input.Age,
	}

	// Save user to the database
	result := database.DB.Create(&user)
	if result.Error != nil {
		log.Printf("ERROR: %v", result.Error)
		return nil, result.Error
	}

	return &user, nil
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.NewUser) (*model.User, error) {
	var user model.User
	// Check if the user exists
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update user fields if provided in input
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Age != 0 {
		user.Age = input.Age
	}

	// Save the updated user
	if err := database.DB.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("could not update user: %w", err)
	}

	return &user, nil

}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	// 	# remove User
	// 	mutation{
	//     deleteUser(id:"ef75e9c0-382e-46a2-b680-8875adf1b219") {
	//       id
	//     }
	//   }
	var user model.User
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		return false, fmt.Errorf("user not found: %w", err)
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return false, fmt.Errorf("could not delete user: %w", err)
	}

	return true, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	startTime := time.Now()

	//Example Query
	// query {
	// 	users {
	// 	  id
	// 	  name
	// 	  email
	// 	  age
	// 	}
	//   }
	var users []*model.User
	cachedUsers, found := c.Get("all_users")
	if found {
		users = cachedUsers.([]*model.User)
		elapsedTime := time.Since(startTime)
		log.Println("[CACHE HIT] Fetched all users from cache %v",elapsedTime)
		return users, nil
	}

	result := database.DB.Find(&users)
	if result.Error != nil {
		log.Printf("ERROR: %v", result.Error)
		return nil, result.Error
	}
	c.Set("all_users", users, cache.DefaultExpiration)
	elapsedTime := time.Since(startTime) // Calculate elapsed time
	log.Printf("[SUCCESS] Users fetched successfully in %v", elapsedTime)
	return users, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	//Example Query
	// query GetUser {
	// 	user(id: "ef75e9c0-382e-46a2-b680-8875adf1b219") {
	// 	  id
	// 	  name
	// 	  email
	// 	  age
	// 	}
	//   }
	startTime := time.Now()
	var user *model.User
	cachedUser, found := c.Get(id)
	if found {
		elapsedTime := time.Since(startTime)
		log.Printf("[CACHE HIT] User with ID: %s fetched from cache %v", id, elapsedTime)
		return cachedUser.(*model.User), nil
	}

	result := database.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		log.Printf("ERROR: %v", result.Error)
		return nil, result.Error
	}
	c.Set(id, user, cache.DefaultExpiration)

	elapsedTime := time.Since(startTime) // Calculate elapsed time
	log.Printf("[SUCCESS] User fetched successfully in %v", elapsedTime)
	return user, nil
	// panic(fmt.Errorf("not implemented: User - user"))
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
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented: CreateTodo - createTodo"))
}
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented: Todos - todos"))
}
*/
