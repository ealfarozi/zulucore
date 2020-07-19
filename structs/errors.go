package structs

//ErrorMessage struct for general error
type ErrorMessage struct {
	Message    string `json:"message,omitempty"`
	Data       string `json:"data,omitempty"`
	SysMessage string `json:"system_message,omitempty"`
	Code       int    `json:"code,omitempty"`
}

//var Message and system message
var (
	ErrNotFound   = "The data can't be found"
	Success       = "success"
	Unauthorized  = "Unauthorized, please login"
	AuthNotFound  = "Auth header is not found"
	TokenInv      = "token is invalid"
	GenTokenErr   = "generate token is error"
	UserNotFound  = "user is not found"
	IncorrectPass = "password is not match"
	QueryErr      = "query error"
	PrepStmtErr   = "Prepared statement error"
	RowsAffErr    = "error while getting rows affected"
	LastIDErr     = "error whilte getting last inserted id"
	Validate      = "error when validating the data"
	NomorInd      = "nomor induk is exist for that institutions"
	Email         = "email is already exists"
	Family        = "Family ID is already exists"
	EmptyID       = "ID can't be empty"
)
