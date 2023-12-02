package rbac

import future.keywords.every
import future.keywords.if
import future.keywords.in

default allow := false

allow if input.user_role == `SYSTEM_ADMIN`

allow if {
	every _, required_permission in input.required_permissions {
		format_int(required_permission, 10), {} in input.has_permissions
	}
}
