package common

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/go-pdf/fpdf"
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
		pdf, data, 20,
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
		pdf.SetTextColor(128, 128, 128) // Gray
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
	r.CellFormat(
		0, 7,
		"The following rules are registered to this class.",
		"", 2, "LT", false, 0, "",
	)

	generateSubSection := func(title, content string) {
		r.SetFontStyle("BI")
		r.CellFormat(
			0, 7,
			title,
			"LR", 1, "", false, 0, "",
		)
		r.SetFontStyle("")
		r.MultiCell(0, 6, content, "LRB", "LT", false)
	}

	// List rules.
	r.SetFillColor(128, 128, 128)
	rules := r.data.Rules
	for idx, rule := range rules {
		// Heading.
		r.SetFontStyle("BI")
		r.CellFormat(
			0, 7,
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
