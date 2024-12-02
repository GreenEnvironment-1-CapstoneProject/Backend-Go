package google

import (
	"greenenvironment/configs"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var conf = configs.InitConfig().Google

var GoogleOauthConfig = &oauth2.Config{
	ClientID:     conf.ClientID,
	ClientSecret: conf.ClientKey,
	RedirectURL:  conf.CallbackURL,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
	},
	Endpoint: google.Endpoint,
}
