package common

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/go-pdf/fpdf"
)

const (
	classReportPageMargin       = 20
	classReportGreyFillColor    = 192
	classReportGreyTextColor    = 128
	classReportNormalLineHeight = 7

	classReportRulesLineHeight = 6

	classReportManagerUserIDWidth       = 42.5
	classReportManagerUserNameWidth     = 42.5
	classReportManagerClassGroupWidth   = 30
	classReportManagerManagingRoleWidth = 55
)

type classReport struct {
	*fpdf.Fpdf
	data database.CoordinatingClassReportData

	// PDF settings
	margin float64
}

func newClassReport(data database.CoordinatingClassReportData) *classReport {
	pdf := fpdf.New(fpdf.OrientationPortrait, fpdf.UnitMillimeter, fpdf.PageSizeA4, "")
	report := &classReport{
		pdf, data, classReportPageMargin,
	}

	pdf.SetMargins(report.margin, report.margin, report.margin)
	pdf.SetAutoPageBreak(true, report.margin)
	report.setFontDefaults()

	pdf.SetHeaderFunc(func() {
		if pdf.PageNo() <= 1 {
			return
		}

		pdf.SetY(5)
		pdf.SetFont("Times", "I", 10)
		pdf.SetTextColor(classReportGreyTextColor, classReportGreyTextColor, classReportGreyTextColor) // Gray
		pdf.CellFormat(
			0, 10,
			fmt.Sprintf("Class Report %s, %d/%s", data.Class.Code, data.Class.Year, data.Class.Semester),
			"", 0, "C", false, 0, "",
		)
		pdf.Ln(20)
	})
	pdf.SetFooterFunc(func() {
		if pdf.PageNo() <= 1 {
			return
		}

		pdf.SetY(-15)
		pdf.SetFont("Times", "I", 10)
		pdf.SetTextColor(128, 128, 128) // Gray
		pdf.CellFormat(
			0, 10,
			fmt.Sprintf("Page %d/{nb}", pdf.PageNo()),
			"", 0, "C", false, 0, "",
		)
	})

	return report
}

func (r *classReport) generateTitlePage() {
	r.setFontDefaults()
	r.AddPage()

	_, height := r.GetPageSize()

	// OAMS Title.
	r.SetFont("Times", "BI", 75)
	r.SetTextColor(0, 191, 255) // Blue
	r.CellFormat(
		0, (height-2*r.margin)/2-15,
		"OAMS",
		"", 2, "CB", false, 0, "",
	)

	// Class Report Subheading.
	r.SetFont("Times", "B", 15)
	r.SetTextColor(0, 0, 0) // Black
	r.CellFormat(
		0, 20,
		"Class Report",
		"", 2, "CB", false, 0, "",
	)

	// Class Quick Info.
	r.SetFont("Times", "I", 13)
	r.CellFormat(
		0, 10,
		fmt.Sprintf("%s, %d/%s", r.data.Class.Code, r.data.Class.Year, r.data.Class.Semester),
		"", 0, "C", false, 0, "",
	)
}

func (r *classReport) fillData() {
	r.fillRules()
	r.fillManagers()
}

func (r *classReport) fillRules() {
	r.setFontDefaults()
	r.AddPage()

	// Set Class Rules section title.
	r.drawSectionTitle("I. RULES")
	r.setFontDefaults()

	// Section description.
	rules := r.data.Rules
	r.MultiCell(
		0, classReportNormalLineHeight,
		fmt.Sprintf("There are currently %d rules registered to this class. Each rule has a title and a description."+
			" OAMS suggests using informative titles and descriptions as these are used to provide students with details"+
			" during rule checking. For more management options, please visit Class Management Menu > Attendance Rules.",
			len(rules),
		),
		"", "LT", false,
	)

	generateSubSection := func(title, content string) {
		r.SetFontStyle("BI")
		r.CellFormat(
			0, classReportNormalLineHeight,
			title,
			"LR", 1, "", false, 0, "",
		)
		r.SetFontStyle("")
		r.MultiCell(0, classReportRulesLineHeight, content, "LRB", "LT", false)
	}

	// List rules.
	r.SetFillColor(classReportGreyFillColor, classReportGreyFillColor, classReportGreyFillColor)
	for idx, rule := range rules {
		// Heading.
		r.SetFontStyle("BI")
		r.CellFormat(
			0, classReportNormalLineHeight,
			fmt.Sprintf("Rule %d", idx+1),
			"LTRB", 1, "", true, 0, "",
		)

		generateSubSection("Title:", rule.Title)
		generateSubSection("Description:", rule.Description)
		generateSubSection("Rule:", rule.Rule)

		e, err := json.MarshalIndent(rule.Environment, "", strings.Repeat(" ", 4))
		if err != nil {
			r.SetError(err)
		}

		generateSubSection("Environment:", string(e))
	}

	r.SetFillColor(0, 0, 0) // Reset the fill color.
}

func (r *classReport) fillManagers() {
	r.setFontDefaults()
	r.AddPage()

	// Set Class Rules section title.
	r.drawSectionTitle("II. MANAGERS")
	r.setFontDefaults()

	// Section description.
	r.MultiCell(
		0, classReportNormalLineHeight,
		"The following users are managers of this class. Course Coordinators have full access to class data,"+
			" while Teaching Assistants are only allowed to manage attendance for the class.",
		"", "LT", false,
	)

	// Table Header.
	headers := []string{"User ID", "User Name", "Class Group", "Managing Role"}

	r.SetFillColor(classReportGreyFillColor, classReportGreyFillColor, classReportGreyFillColor)
	columnWidths := []float64{
		classReportManagerUserIDWidth,
		classReportManagerUserNameWidth,
		classReportManagerClassGroupWidth,
		classReportManagerManagingRoleWidth,
	}

	for idx, head := range headers {
		r.CellFormat(
			columnWidths[idx], classReportNormalLineHeight,
			head,
			"LRTB", 0, "", true, 0, "",
		)
	}

	// List managers.
	r.Ln(classReportNormalLineHeight) // Add newline from header.
	for _, manager := range r.data.Managers {
		data := []string{
			manager.UserID,
			manager.UserName,
			manager.ClassGroupName,
			manager.ManagingRole.String(),
		}
		for idx, d := range data {
			r.CellFormat(
				columnWidths[idx], classReportNormalLineHeight,
				d,
				"LRB", 0, "", false, 0, "",
			)
		}
		r.Ln(classReportNormalLineHeight)
	}

	r.SetFillColor(0, 0, 0) // Reset fill color
}

func (r *classReport) generateLastPage() {
	r.setFontDefaults()
	r.AddPage()

	_, height := r.GetPageSize()

	r.CellFormat(
		0, (height-2*r.margin)/2,
		"END OF REPORT",
		"", 0, "CB", false, 0, "",
	)
}

func (r *classReport) drawSectionTitle(title string) {
	// Set section title.
	r.SetFont("Times", "B", 12)
	r.CellFormat(
		0, 7,
		title,
		"", 2, "CM", false, 0, "",
	)

	// Add some margin from horizontal line break.
	r.SetY(r.GetY() + 2)
}

func (r *classReport) setFontDefaults() {
	r.SetFont("Times", "", 12)
	r.SetTextColor(0, 0, 0)
}

func (r *classReport) close() {
	r.AliasNbPages("")
}

func GenerateClassReport(data database.CoordinatingClassReportData) *fpdf.Fpdf {
	report := newClassReport(data)

	report.generateTitlePage()
	report.fillData()
	report.generateLastPage()

	report.close()
	return report.Fpdf
}
