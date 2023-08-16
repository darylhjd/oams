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

var Classes = newClassesTable("public", "classes", "class")

type classesTable struct {
	postgres.Table

	// Columns
	ID        postgres.ColumnInteger
	Code      postgres.ColumnString
	Year      postgres.ColumnInteger
	Semester  postgres.ColumnString
	Programme postgres.ColumnString
	Au        postgres.ColumnInteger
	CreatedAt postgres.ColumnTimestampz
	UpdatedAt postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ClassesTable struct {
	classesTable

	EXCLUDED classesTable
}

// AS creates new ClassesTable with assigned alias
func (a ClassesTable) AS(alias string) *ClassesTable {
	return newClassesTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ClassesTable with assigned schema name
func (a ClassesTable) FromSchema(schemaName string) *ClassesTable {
	return newClassesTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ClassesTable with assigned table prefix
func (a ClassesTable) WithPrefix(prefix string) *ClassesTable {
	return newClassesTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ClassesTable with assigned table suffix
func (a ClassesTable) WithSuffix(suffix string) *ClassesTable {
	return newClassesTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newClassesTable(schemaName, tableName, alias string) *ClassesTable {
	return &ClassesTable{
		classesTable: newClassesTableImpl(schemaName, tableName, alias),
		EXCLUDED:     newClassesTableImpl("", "excluded", ""),
	}
}

func newClassesTableImpl(schemaName, tableName, alias string) classesTable {
	var (
		IDColumn        = postgres.IntegerColumn("id")
		CodeColumn      = postgres.StringColumn("code")
		YearColumn      = postgres.IntegerColumn("year")
		SemesterColumn  = postgres.StringColumn("semester")
		ProgrammeColumn = postgres.StringColumn("programme")
		AuColumn        = postgres.IntegerColumn("au")
		CreatedAtColumn = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn = postgres.TimestampzColumn("updated_at")
		allColumns      = postgres.ColumnList{IDColumn, CodeColumn, YearColumn, SemesterColumn, ProgrammeColumn, AuColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns  = postgres.ColumnList{CodeColumn, YearColumn, SemesterColumn, ProgrammeColumn, AuColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return classesTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		Code:      CodeColumn,
		Year:      YearColumn,
		Semester:  SemesterColumn,
		Programme: ProgrammeColumn,
		Au:        AuColumn,
		CreatedAt: CreatedAtColumn,
		UpdatedAt: UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
