package intervention

import (
	"embed"
	htmltemplate "html/template"
	texttemplate "text/template"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
)

const (
	textTemplate = "user_text_email.tmpl"
	htmlTemplate = "user_html_email.tmpl"

	userEmailSubject = "OAMS: Attendance Check Failure"
)

var (
	//go:embed *.tmpl
	templates embed.FS

	textEmail *texttemplate.Template
	htmlEmail *htmltemplate.Template
)

type EmailArgs struct {
	UserInfo userKey
	Rules    []model.ClassAttendanceRule
}

func init() {
	var err error
	textEmail, err = texttemplate.ParseFS(templates, textTemplate)
	if err != nil {
		panic(err)
	}

	htmlEmail, err = htmltemplate.ParseFS(templates, htmlTemplate)
	if err != nil {
		panic(err)
	}
}
