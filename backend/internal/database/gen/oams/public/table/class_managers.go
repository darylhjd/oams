//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var ClassManagers = newClassManagersTable("public", "class_managers", "class_manager")

type classManagersTable struct {
	postgres.Table

	// Columns
	ID           postgres.ColumnInteger
	UserID       postgres.ColumnString
	ClassID      postgres.ColumnInteger
	ManagingRole postgres.ColumnString
	CreatedAt    postgres.ColumnTimestampz
	UpdatedAt    postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ClassManagersTable struct {
	classManagersTable

	EXCLUDED classManagersTable
}

// AS creates new ClassManagersTable with assigned alias
func (a ClassManagersTable) AS(alias string) *ClassManagersTable {
	return newClassManagersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ClassManagersTable with assigned schema name
func (a ClassManagersTable) FromSchema(schemaName string) *ClassManagersTable {
	return newClassManagersTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ClassManagersTable with assigned table prefix
func (a ClassManagersTable) WithPrefix(prefix string) *ClassManagersTable {
	return newClassManagersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ClassManagersTable with assigned table suffix
func (a ClassManagersTable) WithSuffix(suffix string) *ClassManagersTable {
	return newClassManagersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newClassManagersTable(schemaName, tableName, alias string) *ClassManagersTable {
	return &ClassManagersTable{
		classManagersTable: newClassManagersTableImpl(schemaName, tableName, alias),
		EXCLUDED:           newClassManagersTableImpl("", "excluded", ""),
	}
}

func newClassManagersTableImpl(schemaName, tableName, alias string) classManagersTable {
	var (
		IDColumn           = postgres.IntegerColumn("id")
		UserIDColumn       = postgres.StringColumn("user_id")
		ClassIDColumn      = postgres.IntegerColumn("class_id")
		ManagingRoleColumn = postgres.StringColumn("managing_role")
		CreatedAtColumn    = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn    = postgres.TimestampzColumn("updated_at")
		allColumns         = postgres.ColumnList{IDColumn, UserIDColumn, ClassIDColumn, ManagingRoleColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns     = postgres.ColumnList{UserIDColumn, ClassIDColumn, ManagingRoleColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return classManagersTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:           IDColumn,
		UserID:       UserIDColumn,
		ClassID:      ClassIDColumn,
		ManagingRole: ManagingRoleColumn,
		CreatedAt:    CreatedAtColumn,
		UpdatedAt:    UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
