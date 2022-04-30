package generator

import (
	"bytes"
	b64 "encoding/base64"
	"image"
	"unicode"
	"unicode/utf8"

	"github.com/jung-kurt/gofpdf"
)

// Contact contact a company informations
type Contact struct {
	Contractor string   `json:"contractor,omitempty"`
	Name       string   `json:"name,omitempty" validate:"required,min=1,max=256"`
	Logo       *[]byte  `json:"logo,omitempty"` // Logo byte array
	Address    *Address `json:"address,omitempty"`
	TaxId      string   `json:"tax_id,omitempty"`
}

func Chunks(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}

func (c *Contact) appendContactTODoc(x float64, y float64, fill bool, logoAlign string, pdf *gofpdf.Fpdf) float64 {
	pdf.SetXY(x, y)

	// Logo
	if c.Logo != nil {
		// Create filename
		fileName := b64.StdEncoding.EncodeToString([]byte(c.Name))
		// Create reader from logo bytes
		ioReader := bytes.NewReader(*c.Logo)
		// Get image format
		_, format, _ := image.DecodeConfig(bytes.NewReader(*c.Logo))
		// Register image in pdf
		imageInfo := pdf.RegisterImageOptionsReader(fileName, gofpdf.ImageOptions{
			ImageType: format,
		}, ioReader)

		if imageInfo != nil {
			var imageOpt gofpdf.ImageOptions
			imageOpt.ImageType = format

			pdf.ImageOptions(fileName, pdf.GetX(), y, 0, 30, false, imageOpt, 0, "")

			pdf.SetY(y + 30)
		}
	}

	// Name
	if fill {
		pdf.SetFillColor(GreyBgColor[0], GreyBgColor[1], GreyBgColor[2])
	} else {
		pdf.SetFillColor(255, 255, 255)
	}

	// Reset x
	pdf.SetX(x)
	var heightOfname int = 0
	/* 	if len(c.Name) < 37 {
		pdf.Rect(x, pdf.GetY(), 70, 8, "F")

		// Set name
		pdf.SetFont("NotoSerif", "B", 10)
		pdf.Cell(40, 8, c.Name)
	} else { */

	// Name rect
	// Set name
	pdf.SetFont("NotoSerif", "B", 10)
	if len(c.Name) < 58 {
		pdf.Rect(x, pdf.GetY(), 109, 8, "F")

		pdf.Cell(40, 8, c.Name)
	} else {
		heightOfname = 3

		chunks := Chunks(c.Name, 58)
		if len(chunks) > 2 {
			heightOfname += len(chunks)
		}
		pdf.Rect(x, pdf.GetY(), 109, 8+float64(heightOfname), "F")

		height := 8

		for jj := 0; jj < len(chunks); jj++ {
			r := []rune(chunks[jj])
			var nextArrFirstR rune
			if len(chunks) != jj+1 {
				nextArrFirstR = []rune(chunks[jj+1])[0]
			} else {
				nextArrFirstR = r[len(r)-1]
			}
			if !unicode.IsSpace(r[len(r)-1]) && len(chunks[jj]) >= 58 && !unicode.IsSpace(nextArrFirstR) {
				pdf.Cell(40, float64(height), trimLastChar(chunks[jj])+"-")
			} else {
				if jj > 0 && !unicode.IsSpace(r[len(r)-1]) && !unicode.IsSpace(rune(chunks[jj][0])) {
					pdf.Cell(40, float64(height), chunks[jj-1][len(chunks[jj-1])-1:]+chunks[jj])
				} else {
					pdf.Cell(40, float64(height), chunks[jj])
				}
			}
			height += 8
			pdf.SetX(x)
		}
		//}
	}
	pdf.SetFont("NotoSerif", "", 10)

	if c.Address != nil {
		// Address rect
		var addrRectHeight float64 = 17

		if len(c.Contractor) > 0 {
			addrRectHeight = addrRectHeight + 5
		}

		if len(c.TaxId) > 0 {
			addrRectHeight = addrRectHeight + 5
		}

		if len(c.Address.Address2) > 0 {
			addrRectHeight = addrRectHeight + 5
		}

		if len(c.Address.Country) == 0 {
			addrRectHeight = addrRectHeight - 5
		}

		pdf.Rect(x, pdf.GetY()+9+float64(heightOfname), 109, addrRectHeight, "F")

		// Set address
		pdf.SetFont("NotoSerif", "", 10)
		pdf.SetXY(x, pdf.GetY()+10+float64(heightOfname))

		if len(c.TaxId) > 0 {
			pdf.Cell(70, 5, c.TaxId)
			pdf.SetXY(x, pdf.GetY()+5)
		}

		if len(c.Contractor) > 0 {
			pdf.Cell(70, 5, "c/o "+c.Contractor)
			pdf.SetXY(x, pdf.GetY()+5)
		}

		pdf.MultiCell(70, 5, c.Address.ToString(), "0", "L", false)
	}

	return pdf.GetY()
}

func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func (c *Contact) appendCompanyContactToDoc(pdf *gofpdf.Fpdf) float64 {
	x, y, _, _ := pdf.GetMargins()
	return c.appendContactTODoc(x, y+40, true, "L", pdf)
}

func (c *Contact) appendCustomerContactToDoc(pdf *gofpdf.Fpdf) float64 {
	x, y, _, _ := pdf.GetMargins()
	// return c.appendContactTODoc(130, BaseMarginTop+25, true, "R", pdf)
	return c.appendContactTODoc(x, y, true, "L", pdf)
}
