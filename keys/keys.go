package keys

import "time"

// common keys
const (
	FailedValidation        = "Failed validation"
	MissingUsername         = "username cannot be empty"
	MissingEmail            = "email cannot be empty"
	MissingPassword         = "password cannot be empty"
	ConfirmPasswordMismatch = "confirm password value should match password"
	Invalid                 = "invalid"
	InvalidEmail            = "not a valid email address"
	InvalidUsername         = "username can only contain letters and numbers"
	UserExists              = "username already exists"
	WrongCredentials        = "wrong credentials, please check username or email and try again"
	Missing                 = "missing"
	UsernameAlreadyExists   = "username already exists"
	InternalError           = "internal_error"

	BadJSON    = "Problem parsing json"
	BadBody    = "Problem reading request body"
	BadRequest = "bad request"
	Success    = "ok"

	ConfigDir = "/etc/bq"

	Session = "session"

	Forbidden = "forbidden"

	BadToken                   = "invalid token"
	BadQuey                    = "invalid query params"
	MaxTokenLife time.Duration = 24 * 7 * time.Hour
	IsNotExist                 = "does not exist"
)

const (

	//LoggerKey The key which stores logger instance in context
	LoggerKey = "logger"

	//JwtKey key used to sore *JWT instance in the context.
	JwtKey = "jwt"

	//DB stores *sql.DB instance.
	DB = "db"

	// RootURL default root url
	RootURL = "http://localhost:8000"

	// Home is the bq home directory.
	Home = "BQ_HOME"
)

//website related keys
const (
	WebsiteTitle       = "high performance qrcode service"
	WebsiteDescription = "create, scan, manage and integrate qrcode into your business workflow"
	WebsiteURL         = "https://bq.co.tz"

	// Warden this is the key used to store an instance of *ladon.Ladon in the
	// request context which is used to manage access of resources.
	Warden = "warden"
)

func Desc(key string) string {
	switch key {
	case Missing:
		return "this field must be provided"
	default:
		return key
	}
}
