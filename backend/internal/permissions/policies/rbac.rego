package rbac

import future.keywords.if

default allow := false

allow if input.action == "hello"
