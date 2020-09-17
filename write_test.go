package dicom

import (
	"encoding/binary"
	"testing"
	"io/ioutil"

	"github.com/stretchr/testify/assert"
	"github.com/suyashkumar/dicom/pkg/dicomio"
	"github.com/suyashkumar/dicom/pkg/tag"
	"github.com/suyashkumar/dicom/pkg/uid"
)

/*
FURTHER TESTING
	- Read written values back in and verify Datsets are the same
	- With 'wild' DICOMs with high variability, read in, write out, read in, and verify
*/

// TODO clean this function up big time
func TestWrite(t *testing.T) {
	location := "fullwrite.dcm"
	file, err := ioutil.TempFile("", location)
	assert.Nil(t, err)
	defer file.Close()

	mediaStorageSOPClassUID, err := newElement(tag.MediaStorageSOPClassUID, []string{"1.2.840.10008.5.1.4.1.1.1.2"})
	assert.Nil(t, err)
	mediaStorageSOPInstanceUID, err := newElement(tag.MediaStorageSOPInstanceUID, []string{"1.2.3.4.5.6.7"})
	assert.Nil(t, err)
	transferSyntax, err := newElement(tag.TransferSyntaxUID, []string{uid.ImplicitVRLittleEndian})
	assert.Nil(t, err)
	patientName, err := newElement(tag.PatientName, []string{"Robin Banks"})
	assert.Nil(t, err)

	elems := []*Element{
		mediaStorageSOPClassUID,
		mediaStorageSOPInstanceUID,
		transferSyntax,
		patientName,
	}

	ds := &Dataset{Elements: elems}

	err = Write(file, ds)
	assert.Nil(t, err)

	// TODO verify that the correct values are written
}

// TODO clean this function up big time
func TestWriteFileHeader(t *testing.T) {
	location := "fileheader.dcm"
	file, err := ioutil.TempFile("", location)
	assert.Nil(t, err)
	defer file.Close()

	w := dicomio.NewWriter(file, binary.LittleEndian, false)

	mediaStorageSOPClassUID, err := newElement(tag.MediaStorageSOPClassUID, []string{"1.2.840.10008.5.1.4.1.1.1.2"})
	assert.Nil(t, err)
	mediaStorageSOPInstanceUID, err := newElement(tag.MediaStorageSOPInstanceUID, []string{"1.2.3.4.5.6.7"})
	assert.Nil(t, err)
	transferSyntax, err := newElement(tag.TransferSyntaxUID, []string{uid.ImplicitVRLittleEndian})
	assert.Nil(t, err)
	metaElems := []*Element{
		mediaStorageSOPClassUID,
		mediaStorageSOPInstanceUID,
		transferSyntax,
	}
	ds := &Dataset{Elements: metaElems}

	err = writeFileHeader(w, ds, metaElems)
	assert.Nil(t, err)

	// TODO Verify the the corrrect things were written to the file header
}

func TestEncodeElementHeader(t *testing.T) {}

func TestWriteValue(t *testing.T) {}

func TestWriteTag(t *testing.T) {}

func TestWriteVRVL(t *testing.T) {}

func TestVerifyVR(t *testing.T) {
	tg := tag.Tag{ // FileMetaInformationGroupLength tag
		Group:   0x0002,
		Element: 0x0000,
	}

	// WRONG VR
	vr, err := verifyVR(tg, "OB")
	assert.Equal(t, "", vr)
	assert.NotNil(t, err)

	// NO VR
	vr, err = verifyVR(tg, "")
	assert.Nil(t, err)
	assert.Equal(t, "UL", vr)

	// MADE UP TAG
	tg = tag.Tag{
		Group:   0x9999,
		Element: 0x9999,
	}
	vr, err = verifyVR(tg, "")
	assert.Nil(t, err)
	assert.Equal(t, "UN", vr)
}

func TestVerifyValueType(t *testing.T) {
	tg := tag.Tag{ // FileMetaInformationGroupLength tag
		Group:   0x0002,
		Element: 0x0000,
	}

	// VALID
	value, err := NewValue([]int{128})
	assert.Nil(t, err)
	err = verifyValueType(tg, value, Ints, "UL")
	assert.Nil(t, err)

	// INVALID VR
	err = verifyValueType(tg, value, Ints, "NA")
	assert.NotNil(t, err)

	// WRONG VALUE TYPE
	err = verifyValueType(tg, value, Strings, "UL")
	assert.NotNil(t, err)
}