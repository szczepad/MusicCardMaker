package spotify

import (
	"bytes"
	"image/color"
	"log/slog"

	"github.com/go-pdf/fpdf"
	"github.com/skip2/go-qrcode"
)

func CreatePDF(tracks []Track) error {
	pdf := fpdf.New("P", "mm", "A4", "")

	pdf.SetFont("Helvetica", "B", 16)
	pdf.SetFillColor(52, 49, 45)
	pdf.SetTextColor(39, 154, 241)

	for i, track := range tracks {
		if i%4 == 0 {
			pdf.AddPage()
		}

		// Track Information
		pdf.Rect(5, float64((i%4)*70)+5, 65, 65, "FD")
		pdf.SetFontSize(13)
		pdf.SetXY(5, float64((i%4)*70))
		pdf.MultiCell(65, 20, track.Artist, "", "C", false)
		pdf.SetXY(5, float64((i%4)*70+55))
		pdf.MultiCell(65, 10, track.Name, "", "C", false)

		pdf.SetFontSize(32)
		pdf.SetXY(5, float64((i%4)*70+30))
		pdf.CellFormat(65, 15, track.ReleaseYear, "", 1, "C", false, 0, "")

		// Track QR-Code
		pdf.Rect(75, float64((i%4)*70)+5, 65, 65, "FD")

		qr, err := qrcode.New(track.Url, qrcode.High)
		if err != nil {
			slog.Error("Could not create QR-Code", "Error", err)
			return err
		}
		qr.ForegroundColor = color.RGBA{R: 39, G: 154, B: 241, A: 250}
		qr.BackgroundColor = color.RGBA{R: 52, G: 49, B: 45, A: 250}

		qrPng, err := qr.PNG(256)
		if err != nil {
			slog.Error("Could not create QR-Code", "Error", err)
			return err
		}

		reader := bytes.NewReader(qrPng)
		pdf.RegisterImageOptionsReader("qrcode", fpdf.ImageOptions{ImageType: "png"}, reader)
		pdf.ImageOptions(
			"qrcode",
			75,
			float64((i%4)*70)+5,
			65,
			65,
			false,
			fpdf.ImageOptions{},
			0,
			"",
		)
	}

	err := pdf.OutputFileAndClose("output.pdf")
	if err != nil {
		slog.Error("Could not create PDF", "Error", err)
		return err
	}
	return nil
}
