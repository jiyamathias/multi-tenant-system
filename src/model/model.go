// Package model holds all database schema's and universal constant variables
package model

type (
	ActorType string
)

const (
	// ActorTypeUser is an ActorType of user
	ActorTypeUser ActorType = "user"
	// ActorTypeTenant is an ActorType of tenant
	ActorTypeTenant ActorType = "tenant"

	// ActionSignup defined the action signup
	ActionSignup string = "signup"
)
