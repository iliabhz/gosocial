package mailer

import "embed"

const FromName = "GoSocial"
const UserInvitationTemplate = "user_invitation.tmpl"
const maxRetries = 3

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile string, username string, email string, data any, isSandbox bool) error
}
