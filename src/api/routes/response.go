package routes

type LoginDetails struct {
	Username  string
	Password  string
	ReCaptcha string
}

type AuthResponse struct {
	Error   string
	Message string
	Token   string
}

var InvalidLogin AuthResponse = AuthResponse{
	Message: "Invalid login details",
}

var FailedCaptcha AuthResponse = AuthResponse{
	Message: "Failed reCaptcha challenge",
}

var MalformedInput AuthResponse = AuthResponse{
	Message: "Malformed input, detected invalid JSON",
}

var UserExists AuthResponse = AuthResponse{
	Message: "User already exists",
}

var InvalidRegister AuthResponse = AuthResponse{
	Message: "Invalid registration",
}

var InvalidAuth AuthResponse = AuthResponse{
	Error: "Invalid auth credentials",
}

var ExpiredAuth AuthResponse = AuthResponse{
	Error: "Expired auth credentials",
}

var ProfileNotFound AuthResponse = AuthResponse{
	Error: "Profile not found",
}

var InvalidProfileUpdate AuthResponse = AuthResponse{
	Error: "Invalild profile update",
}
