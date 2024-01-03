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
	userTextTemplate        = "template_user_text_email.tmpl"
	userHtmlTemplate        = "template_user_html_email.tmpl"
	ruleCreatorTextTemplate = "template_rule_creator_text_email.tmpl"
	ruleCreatorHtmlTemplate = "template_rule_creator_html_email.tmpl"
)

var (
	//go:embed *.tmpl
	templates embed.FS

	userTextEmail *texttemplate.Template
	userHtmlEmail *htmltemplate.Template

	ruleCreatorTextEmail *texttemplate.Template
	ruleCreatorHtmlEmail *htmltemplate.Template
)

type userKey struct {
	ID    string
	Name  string
	Email string
}

type userEmailArgs struct {
	UserInfo userKey
	Rules    []database.RuleInfo
}

type ruleCreatorEmailArgs struct {
	CreatorInfo  userKey
	RuleAndUsers []ruleAndFailedUsers
}

type ruleAndFailedUsers struct {
	Rule        database.RuleInfo
	FailedUsers []userKey
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
