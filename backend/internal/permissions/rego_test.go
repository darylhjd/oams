package permissions

import (
	"context"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/stretchr/testify/assert"
)

func Test_RBAC(t *testing.T) {
	tts := []struct {
		name      string
		withInput RBACInput
		wantAllow bool
	}{
		{
			"allow with no permissions granted system admin role",
			RBACInput{
				model.UserRole_SystemAdmin,
				nil,
				nil,
			},
			true,
		},
		{
			"disallow with non-system admin role",
			RBACInput{
				model.UserRole_User,
				userRolePermissions,
				[]P{ClassCreate},
			},
			false,
		},
		{
			"allow with user role and granted permissions",
			RBACInput{
				model.UserRole_User,
				userRolePermissions,
				[]P{ClassRead},
			},
			true,
		},
	}

	for _, tt := range tts {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			ctx := context.Background()

			allow, err := RBAC(ctx, tt.withInput)
			if err != nil {
				t.Fatal(err)
			}

			a.Equal(tt.wantAllow, allow)
		})
	}
}
