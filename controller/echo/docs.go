package echo

import _ "StoryGoAPI/docs"

//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	GPL 3
//	@license.url	https://www.gnu.org/licenses/gpl-3.0.html

//	@securitydefinitions.apikey	JwtAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
//
//	@securitydefinitions.apikey	GuestAuth
//	@in							header
//	@name						X-Guest-Token
//	@description				Just put the guest api key in header.
//
//	@tag.name					user
//	@tag.description			all user actions
//	@tag.name					story
//	@tag.description			all story actions for both guests and users
//	@tag.name					guest
//	@tag.name					all guest action
