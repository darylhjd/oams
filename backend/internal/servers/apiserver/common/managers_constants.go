package common

const (
	userIdColumn         = 0
	classCodeColumn      = 1
	classYearColumn      = 2
	classSemesterColumn  = 3
	classGroupNameColumn = 4
	classTypeColumn      = 5
	managingRowColumn    = 6

	expectedManagersSheetCount     = 1
	expectedSanityCheckDataRows    = 1
	expectedSanityCheckDataColumns = 7
)

var managersColumnNames = []string{
	"user_id",
	"class_code",
	"class_year",
	"class_semester",
	"class_group_name",
	"class_type",
	"managing_role",
}
