package service

import (
	"context"
	"log"
	"net/http"
	"time"

	"myproject/app/helper"
	"myproject/app/model"
	"myproject/app/model/apperrors"

	"github.com/google/uuid"
)

// userService acts as a struct for injecting an implementation of UserRepository
// for use in service methods
type userService struct {
	UserRepository model.UserRepository
}

// USConfig will hold repositories that will eventually be injected into this
// this service layer
type USConfig struct {
	UserRepository model.UserRepository
}

// NewUserService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewUserService(c *USConfig) model.UserService {
	return &userService{
		UserRepository: c.UserRepository,
	}
}

// Get retrieves a user based on their uuid
func (s *userService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	u, err := s.UserRepository.FindByID(ctx, uid)

	return u, err
}

// Signup reaches our to a UserRepository to verify the
// email address is available and signs up the user if this is the case
func (s *userService) Signup(ctx context.Context, u *model.User) error {
	pw, err := hashPassword(u.Password)

	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", u.Email)
		return apperrors.NewInternal()
	}

	// now I realize why I originally used Signup(ctx, email, password)
	// then created a user. It's somewhat un-natural to mutate the user here
	u.Password = pw

	if err := s.UserRepository.Create(ctx, u); err != nil {
		return err
	}

	// If we get around to adding events, we'd Publish it here
	// err := s.EventsBroker.PublishUserUpdated(u, true)

	// if err != nil {
	// 	return nil, apperrors.NewInternal()
	// }

	return nil
}

// Signin reaches our to a UserRepository check if the user exists
// and then compares the supplied password with the provided password
// if a valid email/password combo is provided, u will hold all
// available user fields
func (s *userService) Signin(ctx context.Context, u *model.User) error {
	start := time.Now()
	uFetched, err := s.UserRepository.FindByEmail(ctx, u.Email)

	// Will return NotAuthorized to client to omit details of why
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        "POST-SIGNIN",
			ResponseTime: time.Since(start),
			Response:     err,
			Key:          u.Email,
		})
		return apperrors.NewAuthorization("Invalid email and password combination")
	}

	// verify password - we previously created this method
	match, err := comparePasswords(uFetched.Password, u.Password)

	if err != nil {
		return apperrors.NewInternal()
	}

	if !match {
		helper.StringLog("error", "Invalid email : "+u.Email+" and password : "+u.Password+" combination")
		return apperrors.NewAuthorization("Invalid email and password combination")
	}

	*u = *uFetched

	helper.LogE2E(&helper.Loge2e{
		Event:        "e2eSigninInvoice",
		StatusCode:   http.StatusOK,
		ResponseTime: time.Since(time.Now()),
		Method:       "POST",
		Request:      "signin-service",
		URL:          uFetched.ImageURL,
		Message:      "message from log e2e",
		Key:          uFetched.Email,
	}, "info", "business-telco")

	helper.StringLog("info", "Success create invoice orderId "+uFetched.Email)

	return nil
}

func (s *userService) UpdateDetails(ctx context.Context, u *model.User) error {
	// Update user in UserRepository
	err := s.UserRepository.Update(ctx, u)

	if err != nil {
		return err
	}

	// // Publish user updated
	// err = s.EventsBroker.PublishUserUpdated(u, false)
	// if err != nil {
	// 	return apperrors.NewInternal()
	// }

	return nil
}
