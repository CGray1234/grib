package griblib

import (
	"bytes"
	"image/color"
	"image/png"
	"io"
	"math"
)

// http://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_temp5-41.shtml
//
//	| Octet Number | Content
//	-----------------------------------------------------------------------------------------
//	| 12-15	     | Reference value (R) (IEEE 32-bit floating-point value)
//	| 16-17	     | Binary scale factor (E)
//	| 18-19	     | Decimal scale factor (D)
//	| 20	         | Number of bits used for each packed value for simple packing, or for each
//	|              | group reference value for complex packing or spatial differencing
//	| 21           | Type of original field values
//	|              |    - 0 : Floating point
//	|              |    - 1 : Integer
//	|              |    - 2-191 : reserved
//	|              |    - 192-254 : reserved for Local Use
//	|              |    - 255 : missing
type Data41 struct {
	Reference    float32 `json:"reference"`
	BinaryScale  uint16  `json:"binaryScale"`
	DecimalScale uint16  `json:"decimalScale"`
	Bits         uint8   `json:"bits"`
	Type         uint8   `json:"type"`
}

func (template Data41) getRefScale() (float64, float64) {
	bscale := math.Pow(2.0, float64(template.BinaryScale))
	dscale := math.Pow(10.0, -float64(template.DecimalScale))

	scale := bscale * dscale
	ref := dscale * float64(template.Reference)

	return ref, scale
}

func (template Data41) scaleFunc() func(uintValue int64) float64 {
	ref, scale := template.getRefScale()
	return func(value int64) float64 {
		signed := int64(value)
		return ref + float64(signed)*scale
	}
}

// ParseData0 parses data0 struct from the reader into the an array of floating-point values
func ParseData41(dataReader io.Reader, dataLength int, template *Data41) ([]float64, error) {

	fld := []float64{}

	if dataLength == 0 {
		return fld, nil
	}

	pngData := make([]byte, dataLength)
	_, err := io.ReadFull(dataReader, pngData)
	if err != nil {
		return fld, nil
	}

	img, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		return fld, nil
	}

	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			pixelColor := img.At(x, y).(color.Gray16)

			fld = append(fld, (float64(template.Reference)+float64(pixelColor.Y)*math.Pow(2, float64(template.BinaryScale)))/math.Pow(10, float64(template.DecimalScale)))
		}
	}

	return fld, nil
}
