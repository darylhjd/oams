package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database/types"
	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	postgres2 "github.com/go-jet/jet/v2/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	host, err := strconv.Atoi(env.GetDatabasePort())
	if err != nil {
		log.Fatal(err)
	}

	genPath := "../gen"
	dbConn := postgres.DBConnection{
		Host:       env.GetDatabaseHost(),
		Port:       host,
		User:       env.GetDatabaseUser(),
		Password:   env.GetDatabasePassword(),
		SslMode:    env.GetDatabaseSSLMode(),
		DBName:     env.GetDatabaseName(),
		SchemaName: "public",
	}

	err = postgres.Generate(genPath, dbConn, template.Default(postgres2.Dialect).
		UseSchema(func(schemaMetaData metadata.Schema) template.Schema {
			return template.DefaultSchema(schemaMetaData).
				UseModel(template.DefaultModel().
					UseTable(func(table metadata.Table) template.TableModel {
						return template.DefaultTableModel(table).
							// Add json tags to model fields.
							UseField(func(columnMetaData metadata.Column) template.TableModelField {
								defaultTableModelField := template.DefaultTableModelField(columnMetaData)

								if table.Name == "class_attendance_rules" && columnMetaData.Name == "environment" {
									defaultTableModelField.Type = template.NewType(types.Environment{})
								}

								return defaultTableModelField.UseTags(
									fmt.Sprintf(`json:"%s"`, columnMetaData.Name))
							}).
							// Use singular of table plural name for model name.
							UseTypeName(toCamelCase(trimTableName(table.Name)))
					})).
				UseSQLBuilder(template.DefaultSQLBuilder().
					UseTable(func(table metadata.Table) template.TableSQLBuilder {
						// Set proper alias to custom type name.
						return template.DefaultTableSQLBuilder(table).UseDefaultAlias(trimTableName(table.Name))
					}))
		}))
	if err != nil {
		log.Fatal(err)
	}
}

// toCamelCase uses a snake case table name to generate a camel case name for the corresponding type.
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for idx := range parts {
		parts[idx] = cases.Title(language.English).String(parts[idx])
	}

	return strings.Join(parts, "")
}

// trimTableName generates the proper singular name from a table name.
func trimTableName(s string) string {
	trimmedName := strings.TrimSuffix(s, "s")
	if s == "classes" {
		trimmedName = "class"
	}

	return trimmedName
}
