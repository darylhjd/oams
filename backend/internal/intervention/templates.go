package intervention

import (
	"embed"
	htmltemplate "html/template"
	texttemplate "text/template"

	"github.com/darylhjd/oams/backend/internal/database"
)

const (
	userEmailSubject        = "OAMS: Attendance Check Failure"
	ruleCreatorEmailSubject = "OAMS: Attendance Check Complete"
)

const (
	userTextTemplate        = "user_text_email.tmpl"
	userHtmlTemplate        = "user_html_email.tmpl"
	ruleCreatorTextTemplate = "rule_creator_text_email.tmpl"
	ruleCreatorHtmlTemplate = "rule_creator_html_email.tmpl"
)

var (
	//go:embed *.tmpl
	templates embed.FS

	userTextEmail *texttemplate.Template
	userHtmlEmail *htmltemplate.Template

	ruleCreatorTextEmail *texttemplate.Template
	ruleCreatorHtmlEmail *htmltemplate.Template
)

type userEmailArgs struct {
	UserInfo userKey
	Rules    []database.RuleInfo
}

type ruleCreatorEmailArgs struct {
	CreatorName  string
	CreatorEmail string
	RuleAndUsers map[ruleKey][]userKey
}

type ruleKey struct {
	ID            int64
	Title         string
	Description   string
	ClassCode     string
	ClassYear     int32
	ClassSemester string
}

func init() {
	var err error
	userTextEmail, err = texttemplate.ParseFS(templates, userTextTemplate)
	if err != nil {
		panic(err)
	}

	userHtmlEmail, err = htmltemplate.ParseFS(templates, userHtmlTemplate)
	if err != nil {
		panic(err)
	}

	ruleCreatorTextEmail, err = texttemplate.ParseFS(templates, ruleCreatorTextTemplate)
	if err != nil {
		panic(err)
	}

	ruleCreatorHtmlEmail, err = htmltemplate.ParseFS(templates, ruleCreatorHtmlTemplate)
	if err != nil {
		panic(err)
	}
}
