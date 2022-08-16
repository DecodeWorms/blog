package handlers

import (
	"blog/auth"
	"blog/errhandler"
	"blog/models"
	"blog/pkg"
	"blog/storage"
	"blog/util"
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService storage.UserServices
}

func NewUserHandler(use storage.UserServices) UserHandler {
	return UserHandler{
		userService: use,
	}

}

func (u UserHandler) AutoMigrate() error {
	err := u.userService.Automigrate()
	if err != nil {
		return errors.New("unable to migrate a table")
	}
	return nil

}

func (u UserHandler) Signup(ctx context.Context, data models.UserInput) *errhandler.UserError {
	//validate if the user already exists
	_, err := u.userService.UserByEmail(ctx, data.Email)
	if err == nil {
		return errhandler.EmailAlreadyExist
	}
	_, err = u.userService.UserByPhoneNumber(ctx, data.PhoneNumber)
	if err == nil {
		return errhandler.PhoneNumberAlreadyExist
	}
	if b := pkg.ValidateGender(data.Gender); !b {
		return errhandler.GenderNotMatched
	}

	//validate user datas before persisting it to the DB..
	val := pkg.InitValidator{Validate: pkg.NewInitValidator()}
	valErr := ValidateUserData(val, data)
	return errhandler.NewUserError(500, valErr)

	// hash password
	password, err := pkg.HashPassword(data.Password)
	if err != nil {
		return errhandler.IncorrectPassword
	}

	d := models.UserInput{
		Username:    data.Username,
		Email:       data.Email,
		Password:    password,
		Gender:      data.Gender,
		PhoneNumber: data.PhoneNumber,
	}

	if crErr := u.userService.Create(ctx, d); crErr != nil {
		return errhandler.UnableToCeateUser
	}
	return nil

}

func (u UserHandler) StoreOtherUserData(ctx *gin.Context, data models.StoreOtherUserDataInput) error {
	userId := fmt.Sprintf("%v", util.GetUserIdFromContext(ctx))
	fmt.Println("userId is", userId)
	user, err := u.userService.UserById(ctx, userId)
	if err != nil {
		return err
	}

	val := pkg.InitValidator{Validate: pkg.NewInitValidator()}
	if valErr := ValidateStoreOtherUserData(val, data.PersonalInfo); valErr != nil {
		return fmt.Errorf("validation err : %v", valErr)
	}

	val = pkg.InitValidator{Validate: pkg.NewInitValidator()}
	if verr := ValidateAddress(val, data.Address); verr != nil {
		return fmt.Errorf("validation err : %v", verr)
	}

	us := models.MoreInfo{
		UserId:        user.ID,
		MaritalStatus: data.PersonalInfo.MaritalStatus,
		Age:           data.PersonalInfo.Age,
		Title:         data.PersonalInfo.Title,
	}

	err = u.userService.StoreOtherUserData(ctx, us)
	if err != nil {
		return err
	}

	add := models.AddressInput{
		UserId:     user.ID,
		State:      data.Address.State,
		LGA:        data.Address.LGA,
		PostalCode: data.Address.PostalCode,
	}

	err = u.userService.CreateAddress(ctx, add)
	if err != nil {
		return err
	}

	return nil

}

func (u UserHandler) Login(ctx context.Context, data models.Login) (*models.LoginResponse, *errhandler.UserError) {
	//validate if the user already exists
	user, err := u.userService.UserByEmail(ctx, data.Email)
	if err != nil {
		return nil, errhandler.UserRecordNotFound
	}

	// validate login datas
	val := pkg.InitValidator{Validate: pkg.NewInitValidator()}
	valerr := ValidateLogin(val, data)
	return nil, errhandler.NewUserError(500, valerr)

	// compare password

	perr := pkg.ComparePassword(user.Password, data.Password)
	if perr != nil {
		return nil, errhandler.IncorrectPassword
	}

	//generate token
	tokenString, err := auth.GenerateJWT(user.Email, user.Username, user.ID)
	if err != nil {
		return nil, errhandler.UnableToGenerateToken
	}

	lgnResponse := &models.LoginResponse{
		Id:            user.ID,
		Email:         user.Email,
		Gender:        user.Gender,
		PhoneNumber:   user.PhoneNumber,
		MaritalStatus: user.MaritalStatus,
		TokenString:   tokenString,
		Title:         user.Title,
	}
	return lgnResponse, nil

}

func (u UserHandler) UserDetails(ctx context.Context, email string) (*models.UserDetails, error) {
	user, err := u.userService.UserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user records not found")
	}
	d := &models.UserDetails{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Gender:   user.Gender,
	}
	return d, nil
}

func (u UserHandler) UserById(ctx *gin.Context, id string) (*models.User, error) {
	user, err := u.userService.UserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func ValidateUserData(val pkg.InitValidator, data interface{}) []error {
	mpErr := make([]error, 0)

	err := val.Struct(data)
	if err != nil {
		for _, value := range err.(validator.ValidationErrors) {

			e := fmt.Errorf("field : %s and kind : %s", value.Field(), value.Kind())
			mpErr = append(mpErr, e)
		}
		return mpErr
	}
	return nil
}

func ValidateStoreOtherUserData(val pkg.InitValidator, data interface{}) []error {
	mpErr := make([]error, 0)

	err := val.Struct(data)
	if err != nil {
		for _, value := range err.(validator.ValidationErrors) {

			e := fmt.Errorf("field : %s and kind : %s", value.Field(), value.Kind())
			mpErr = append(mpErr, e)
		}
		return mpErr
	}
	return nil

}

func ValidateAddress(val pkg.InitValidator, data interface{}) []error {
	mpErr := make([]error, 0)

	err := val.Struct(data)
	if err != nil {
		for _, value := range err.(validator.ValidationErrors) {

			e := fmt.Errorf("field : %s and kind : %s", value.Field(), value.Kind())
			mpErr = append(mpErr, e)
		}
		return mpErr
	}
	return nil

}

func ValidateLogin(val pkg.InitValidator, data interface{}) []error {
	mpErr := make([]error, 0)

	err := val.Struct(data)
	if err != nil {
		for _, value := range err.(validator.ValidationErrors) {

			e := fmt.Errorf("field : %s and kind : %s", value.Field(), value.Kind())
			mpErr = append(mpErr, e)
		}
		return mpErr
	}
	return nil

}
