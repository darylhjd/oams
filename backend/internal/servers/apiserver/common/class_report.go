package common

import (
	"fmt"
	"strconv"

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
	report.setPageDefaults()

	pdf.SetHeaderFunc(func() {
		if pdf.PageNo() <= 1 {
			return
		}

		pdf.SetY(5)
		pdf.SetFont("Arial", "I", 10)
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
		pdf.SetFont("Arial", "I", 10)
		pdf.SetTextColor(128, 128, 128) // Gray
		pdf.CellFormat(
			0, 10,
			fmt.Sprintf("Page %d/{nb}", pdf.PageNo()),
			"", 0, "C", false, 0, "",
		)
	})

	return report
}

func (r *classReport) setPageDefaults() {
	r.pdf.SetFont("Arial", "", 13)
	r.pdf.SetTextColor(0, 0, 0)
}

func (r *classReport) generateTitlePage() {
	r.setPageDefaults()
	r.pdf.AddPage()

	_, height := r.pdf.GetPageSize()
	var middleOffset float64 = 15

	// OAMS Title
	r.pdf.SetFont("Arial", "BI", 50)
	r.pdf.SetTextColor(0, 191, 255) // Blue
	r.pdf.CellFormat(
		0, (height-2*r.margin)/2-middleOffset,
		"OAMS",
		"", 2, "CB", false, 0, "",
	)

	// Class Report Title
	r.pdf.SetFont("Arial", "B", 15)
	r.pdf.SetTextColor(0, 0, 0) // Black
	r.pdf.CellFormat(
		0, 20,
		"Class Report",
		"", 2, "CB", false, 0, "",
	)

	// Class Quick Info
	r.pdf.SetFont("Arial", "I", 13)
	r.pdf.CellFormat(
		0, 10,
		fmt.Sprintf("%s, %d/%s", r.reportData.Class.Code, r.reportData.Class.Year, r.reportData.Class.Semester),
		"", 0, "C", false, 0, "",
	)
}

func (r *classReport) fillData() {
	r.setPageDefaults()
	r.pdf.AddPage()
	for i := 0; i < 100; i++ {
		r.pdf.CellFormat(0, 10, strconv.Itoa(i), "1", 2, "C", false, 0, "")
	}
}

func (r *classReport) generateLastPage() {
	r.setPageDefaults()
	r.pdf.AddPage()

	_, height := r.pdf.GetPageSize()

	r.pdf.CellFormat(
		0, (height-2*r.margin)/2,
		"END OF REPORT",
		"", 0, "CB", false, 0, "",
	)
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
