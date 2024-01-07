package common

import (
	"encoding/json"
	"fmt"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/go-pdf/fpdf"
)

type classReport struct {
	pdf        *fpdf.Fpdf
	reportData database.CoordinatingClassReportData

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
	r.pdf.AddPage()

	_, height := r.pdf.GetPageSize()

	// OAMS Title.
	r.pdf.SetFont("Times", "BI", 75)
	r.pdf.SetTextColor(0, 191, 255) // Blue
	r.pdf.CellFormat(
		0, (height-2*r.margin)/2-15,
		"OAMS",
		"", 2, "CB", false, 0, "",
	)

	// Class Report Subheading.
	r.pdf.SetFont("Times", "B", 15)
	r.pdf.SetTextColor(0, 0, 0) // Black
	r.pdf.CellFormat(
		0, 20,
		"Class Report",
		"", 2, "CB", false, 0, "",
	)

	// Class Quick Info.
	r.pdf.SetFont("Times", "I", 13)
	r.pdf.CellFormat(
		0, 10,
		fmt.Sprintf("%s, %d/%s", r.reportData.Class.Code, r.reportData.Class.Year, r.reportData.Class.Semester),
		"", 0, "C", false, 0, "",
	)
}

func (r *classReport) fillData() {
	r.fillRules()
	r.fillManagers()
}

func (r *classReport) fillRules() {
	r.setFontDefaults()
	r.pdf.AddPage()

	// Set Class Rules section title.
	r.drawSectionTitle("I. RULES")
	r.setFontDefaults()

	// Section description.
	r.pdf.CellFormat(
		0, 7,
		"The following rules are registered to this class.",
		"", 2, "LT", false, 0, "",
	)

	generateSubSection := func(title, content string) {
		r.pdf.SetFontStyle("BI")
		r.pdf.CellFormat(
			0, 7,
			title,
			"LR", 1, "", false, 0, "",
		)
		r.pdf.SetFontStyle("")
		r.pdf.MultiCell(0, 6, content, "LRB", "LT", false)
	}

	// List rules.
	r.pdf.SetFillColor(128, 128, 128)
	rules := r.reportData.Rules
	for idx, rule := range rules {
		// Heading.
		r.pdf.SetFontStyle("BI")
		r.pdf.CellFormat(
			0, 7,
			fmt.Sprintf("Rule %d", idx+1),
			"LTRB", 1, "", true, 0, "",
		)

		generateSubSection("Title:", rule.Title)
		generateSubSection("Description:", rule.Description)
		generateSubSection("Rule:", rule.Rule)

		e, err := json.MarshalIndent(rule.Environment, "", "    ")
		if err != nil {
			r.pdf.SetError(err)
		}

		generateSubSection("Environment:", string(e))
	}

	r.pdf.SetFillColor(0, 0, 0) // Reset the fill color.
}

func (r *classReport) fillManagers() {
	r.setFontDefaults()
	r.pdf.AddPage()

}

func (r *classReport) generateLastPage() {
	r.setFontDefaults()
	r.pdf.AddPage()

	_, height := r.pdf.GetPageSize()

	r.pdf.CellFormat(
		0, (height-2*r.margin)/2,
		"END OF REPORT",
		"", 0, "CB", false, 0, "",
	)
}

func (r *classReport) drawSectionTitle(title string) {
	// Set section title.
	r.pdf.SetFont("Times", "B", 12)
	r.pdf.CellFormat(
		0, 7,
		title,
		"", 2, "CM", false, 0, "",
	)

	// Add some margin from horizontal line break.
	r.pdf.SetY(r.pdf.GetY() + 2)
}

func (r *classReport) setFontDefaults() {
	r.pdf.SetFont("Times", "", 12)
	r.pdf.SetTextColor(0, 0, 0)
}

func (r *classReport) close() {
	r.pdf.AliasNbPages("")
}

func GenerateClassReport(data database.CoordinatingClassReportData) *fpdf.Fpdf {
	report := newClassReport(data)

	report.generateTitlePage()
	report.fillData()
	report.generateLastPage()

	report.close()
	return report.pdf
}
