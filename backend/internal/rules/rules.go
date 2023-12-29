package main

import (
	"fmt"
	"log"
	"time"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/expr-lang/expr"
)

// SimpleRule to check if student missed 3 consecutive classes.
var (
	SimpleRule = "all(enrollments[-consecutiveClasses:], {!.Attended})"
)

func main() {
	env := map[string]any{
		"enrollments": []model.SessionEnrollment{
			{
				ID:        1,
				SessionID: 0,
				UserID:    "",
				Attended:  false,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			{
				ID:        2,
				SessionID: 0,
				UserID:    "",
				Attended:  false,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			{
				ID:        3,
				SessionID: 0,
				UserID:    "",
				Attended:  false,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		"consecutiveClasses": 1,
	}

	programme, err := expr.Compile(SimpleRule, expr.Env(env))
	if err != nil {
		log.Fatal(err)
	}

	output, err := expr.Run(programme, env)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}
