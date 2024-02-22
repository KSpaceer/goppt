package goppt

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/KSpaceer/goppt/internal/ioadapters"
)

const headerSize = 8

// recordType is an enumeration that specifies the record type of an atom record or a container record
// ([MS-PPT] 2.13.24 RecordType)
type recordType uint16

const (
	recordTypeUnspecified              recordType = 0
	recordTypeDocument                 recordType = 0x03E8
	recordTypeSlide                    recordType = 0x03EE
	recordTypeEnvironment              recordType = 0x03F2
	recordTypeSlidePersistAtom         recordType = 0x03F3
	recordTypeSlideShowSlideInfoAtom   recordType = 0x03F9
	recordTypeExternalObjectList       recordType = 0x0409
	recordTypeDrawingGroup             recordType = 0x040B
	recordTypeDrawing                  recordType = 0x040C
	recordTypeList                     recordType = 0x07D0
	recordTypeSoundCollection          recordType = 0x07E4
	recordTypeTextCharsAtom            recordType = 0x0FA0
	recordTypeTextBytesAtom            recordType = 0x0FA8
	recordTypeHeadersFooters           recordType = 0x0FD9
	recordTypeSlideListWithText        recordType = 0x0FF0
	recordTypeUserEditAtom             recordType = 0x0FF5
	recordTypePersistDirectoryAtom     recordType = 0x1772
	recordTypeRoundTripSlideSyncInfo12 recordType = 0x3714
)

// LowerPart returns lower byte of record type
func (r recordType) LowerPart() byte {
	const fullByte = 0xFF
	return byte(r & fullByte)
}

var errMismatchRecordType = errors.New("mismatch record type")

type record struct {
	header [headerSize]byte
	recordData
}

// Type returns recordType of record contained in it's header
func (r record) Type() recordType {
	return recordType(binary.LittleEndian.Uint16(r.header[2:4]))
}

// Length returns data length contained in record header
func (r record) Length() uint32 {
	return binary.LittleEndian.Uint32(r.header[4:8])
}

// Data returns all data from record except header
func (r record) Data() []byte {
	return r.recordData
}

type recordData []byte

// ReadAt copies bytes from record data at given offset into buffer p
func (rd recordData) ReadAt(p []byte, off int64) (n int, err error) {
	return ioadapters.BytesReadAt(rd, p, off)
}

// LongAt interprets 4 bytes of record data at given offset as uint32 value and returns it
func (rd recordData) LongAt(offset int) uint32 {
	return binary.LittleEndian.Uint32(rd[offset:])
}

// readRecord reads header and data of record. If wantedType is specified (not equals recordTypeUnspecified),
// also compares read type with the wanted one and returns an error is they are not equal
func readRecord(f io.ReaderAt, offset int64, wantedType recordType) (record, error) {
	r, err := readRecordHeaderOnly(f, offset, wantedType)
	if err != nil {
		return record{}, err
	}
	r.recordData = make([]byte, r.Length())
	_, err = f.ReadAt(r.recordData, offset+headerSize)
	if err != nil {
		return record{}, err
	}
	return r, nil
}

// readRecordHeaderOnly reads header of record. If wantedType is specified (not equals recordTypeUnspecified),
// also compares read type with the wanted one and returns an error is they are not equal
func readRecordHeaderOnly(f io.ReaderAt, offset int64, wantedType recordType) (record, error) {
	r := record{}
	_, err := f.ReadAt(r.header[:], offset)
	if err != nil {
		return record{}, err
	}
	if wantedType != recordTypeUnspecified && r.Type() != wantedType {
		return record{}, errMismatchRecordType
	}
	return r, nil
}
