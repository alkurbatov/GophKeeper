package gophtest

import "errors"

const (
	Username    = "admin"
	Password    = "1q2w3e"
	SecurityKey = "88bb5abaa61568b9f11ba091445d81772a3a264fb3f3054088f78baf7a091a9d"
	AccessToken = "SomeLongTokenInJWT"
	Secret      = "xxx"
)

var ErrUnexpected = errors.New("runtime error")
