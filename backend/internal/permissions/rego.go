package permissions

import (
	"context"
	"embed"
	"log"

	"github.com/open-policy-agent/opa/loader"
	"github.com/open-policy-agent/opa/rego"
)

//go:embed policies/*.rego
var f embed.FS

var (
	rbac rego.PreparedEvalQuery
)

func init() {
	bundle, err := loader.NewFileLoader().WithFS(f).AsBundle("policies")
	if err != nil {
		panic(err)
	}

	log.Println(bundle.Modules)

	rbac, err = rego.New(
		rego.Query("data.rbac.allow"),
		rego.ParsedBundle("", bundle),
	).PrepareForEval(context.Background())
}
