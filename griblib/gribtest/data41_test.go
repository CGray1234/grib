package gribtest

import (
	"testing"
)

func Test_read41_integrationtest_file(t *testing.T) {
	messages := openGrib(t, "../integrationtestdata/template5_41.grib2")

	t.Log(messages[0].Data()[:1000])

	// Section5.GetDataTemplate()
	// [0] = reference value
	// [1] = binary scale factor
	// [2] = decimal scale factor
	// [3] = reference value bit count
	// [4] = original field type

	// assert.Len(t, messages, 2, "should have exactly 2 message in testfile")

	// assert.Equal(t, uint16(0), messages[0].Section5.DataTemplateNumber, "Data template number should be 0")
}
