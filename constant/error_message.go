package constant

import "errors"

const Unauthorized = "Unauthorized"
const InternalServerError = "Internal Server Error"
const BadInput = "Format data not valid"

var BadRequest = errors.New("Format data not valid")
var ErrPageInvalid = errors.New("Page is invalid")
var ErrFieldData = errors.New("Error word types on the fields")

// Users Errors
var ErrInvalidUsername = errors.New("Username formating not valid")
var ErrInvalidPhone = errors.New("Phone formating not valid")
var ErrEmptyLogin = errors.New("Email or Password cannot be empty")
var UserNotFound = errors.New("User not found")
var ErrLoginIncorrectPassword = errors.New("Incorrect password")
var ErrEmptyEmailRegister = errors.New("Email cannot be empty")
var ErrEmptyPasswordRegister = errors.New("Password cannot be empty")
var ErrEmptyAddressRegister = errors.New("Address cannot be empty")
var ErrEmptyNameRegister = errors.New("Name cannot be empty")
var ErrEmptyGenderRegister = errors.New("Gender cannot be empty")
var ErrEmailAlreadyExist = errors.New("Email already exist")
var ErrUsernameAlreadyExist = errors.New("Username already exist")
var ErrEmptyPhoneRegister = errors.New("Phone cannot be empty")
var ErrRegister = errors.New("Failed to register user")
var ErrUpdateUser = errors.New("Failed to update user")
var ErrEmptyUpdate = errors.New("One or more fields for update cannot be empty")
var ErrEmailUsernameAlreadyExist = errors.New("Email or Username already exist")
var ErrEmptyEmail = errors.New("Email cannot be empty")
var ErrEmailNotFound = errors.New("Email not found")
var ErrForgotPassword = errors.New("Failed to forgot password")
var ErrOTPNotValid = errors.New("OTP not valid")
var ErrOTPExpired = errors.New("OTP expired")
var ErrEmptyOTP = errors.New("OTP cannot be empty")
var ErrResetPassword = errors.New("Failed to reset password")
var ErrDeleteUser = errors.New("Failed to delete user")
var ErrEmptyResetPassword = errors.New("Email, password and confirmation password cannot be empty")
var ErrPasswordNotMatch = errors.New("Password not match")
var ErrInvalidEmail = errors.New("Email is not valid")

var ErrGenerateJWT = errors.New("failed to generate jwt token")

var ErrValidateJWT = errors.New("failed to validate jwt token")

var ErrHashPassword = errors.New("failed to hash password")

// Users By Admin
var ErrUserDataEmpty = errors.New("User Data Empty")
var ErrGetUser = errors.New("Failed to get user")
var ErrUserNotFound = errors.New("User Not Found")
var ErrUserIDNotFound = errors.New("User ID Not Found")
var ErrEditUserByAdmin = errors.New("Error, Update at least one field")