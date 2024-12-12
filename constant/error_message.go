package constant

import "errors"

const Unauthorized = "Unauthorized"
const InternalServerError = "Internal Server Error"
const BadInput = "Format data not valid"

var BadRequest = errors.New("Format data not valid")
var ErrUnauthorized = errors.New("Unauthorized")
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
var ErrUpdateAvatar = errors.New("Failed to update avatar")

var ErrGenerateJWT = errors.New("failed to generate jwt token")

var ErrValidateJWT = errors.New("failed to validate jwt token")

var ErrHashPassword = errors.New("failed to hash password")

var ErrSizeFile = errors.New("file size exceeds limit")
var ErrContentTypeFile = errors.New("only image allowed")

// Users By Admin
var ErrUserDataEmpty = errors.New("User Data Empty")
var ErrGetUser = errors.New("Failed to get user")
var ErrUserNotFound = errors.New("User Not Found")
var ErrUserIDNotFound = errors.New("User ID Not Found")
var ErrEditUserByAdmin = errors.New("Error, Update at least one field")

var ErrCreateProduct = errors.New("failed to create product")
var ErrProductEmpty = errors.New("product is empty")
var ErrGetProduct = errors.New("failed to get product")
var ErrUpdateProduct = errors.New("failed to update product")
var ErrDeleteProduct = errors.New("failed to delete product")

var ErrImpactCategoryNotFound = errors.New("failed to get impact category")
var ErrCreateImpactCategory = errors.New("failed to create impact category")
var ErrDeleteImpactCategory = errors.New("failed to delete impact category")

var ErrFieldType = errors.New("field type error")

var ErrUserAlreadyParticipate = errors.New("User already participate")
var ErrCreateChallenge = errors.New("Failed to create challenge")
var ErrGetChallenge = errors.New("Failed to get challenge")
var ErrGetChallengeByID = errors.New("Failed to get challenge by ID")
var ErrEditChallenge = errors.New("Field Description Cannot be Empty")
var ErrChallengeNotFound = errors.New("Challenge Not Found")
var ErrUpdateChallenge = errors.New("Failed to update challenge")
var ErrDeleteChallenge = errors.New("Failed to delete challenge")
var ErrChallengeField = errors.New("Field cannot be empty")
var ErrChallengeFieldUpdate = errors.New("One or more fields for update cannot be empty")
var ErrChallengeFieldSwipe = errors.New("Field cannot be empty")
var ErrChallengeType = errors.New("Challenge type must be accept or decline")
var ErrChallengeFieldCreate = errors.New("Field cannot be empty")
var ErrInsertChallengeLog = errors.New("Failed to insert challenge log")
var ErrChallengeLogType = errors.New("Challenge log type must be accept or decline")
var ErrParticipateChallenge = errors.New("Failed to participate challenge")

var ErrInvalidDayNumber = errors.New("invalid day number")
var ErrTaskAlreadyExists = errors.New("task already exists for the given day")
var ErrCreateTask = errors.New("failed to create task")
var ErrFetchTasks = errors.New("failed to fetch tasks")
var ErrTaskNotFound = errors.New("task not found")
var ErrUpdateTask = errors.New("failed to update task")
var ErrDeleteTask = errors.New("failed to delete task")

var ErrCreateChallengeLog = errors.New("Failed to take challenge")
var ErrCreateChallengeConfirmation = errors.New("Failed to create challenge confirmation")
var ErrChallengeConfirmationNotFound = errors.New("challenge confirmation not found")
var ErrUpdateChallengeConfirmation = errors.New("failed to update challenge confirmation")
var ErrChallengeLogNotFound = errors.New("Challenge log not found")
var ErrUpdateChallengeLog = errors.New("Failed to update challenge log")
var ErrChallengeTaskNotFound = errors.New("Challenge Task not found")
var ErrRewardAlreadyClaimed = errors.New("Rewards already claimed")
var ErrChallengeAlreadyTaken = errors.New("challenge already taken by user")

var ErrTransactionEmpty = errors.New("transaction is empty")
