package permissions

import (
	"context"
	"embed"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/open-policy-agent/opa/loader"
	"github.com/open-policy-agent/opa/rego"
)

//go:embed policies/*.rego
var f embed.FS

const (
	policiesDir = "policies"
)

const (
	rbacQuery = "data.rbac.allow"
)

var (
	rbac rego.PreparedEvalQuery
)

func init() {
	ctx := context.Background()

	bundle, err := loader.NewFileLoader().WithFS(f).AsBundle(policiesDir)
	if err != nil {
		panic(err)
	}

	rbac, err = rego.New(
		rego.Query(rbacQuery),
		rego.ParsedBundle(policiesDir, bundle),
	).PrepareForEval(ctx)
	if err != nil {
		panic(err)
	}
}

type RBACInput struct {
	UserRole            model.UserRole `json:"user_role"`
	HasPermissions      PermissionMap  `json:"has_permissions"`
	RequiredPermissions []P            `json:"required_permissions"`
}

func RBAC(ctx context.Context, input RBACInput) (bool, error) {
	result, err := rbac.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return false, err
	}

	return result.Allowed(), nil
}
