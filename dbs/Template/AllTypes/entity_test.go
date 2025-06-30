package AllTypes

import (
	"bytes"
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/rah-0/testmark/testutil"

	"github.com/rah-0/margo-test/util"
)

var (
	c *sql.DB
)

func TestMain(m *testing.M) {
	testutil.TestMainWrapper(testutil.TestConfig{
		M: m,
		LoadResources: func() error {
			dsn := util.GetDsn()
			var err error

			c, err = sql.Open("mysql", dsn)
			if err != nil {
				return err
			}

			SetDB(c)
			return err
		},
		UnloadResources: func() error {
			_, err := DBTruncate()
			if err != nil {
				return err
			}
			return c.Close()
		},
	})
}

func TestAllFieldsRoundtrip(t *testing.T) {
	u := uuid.New().String()
	e := Entity{
		Id:              "1",
		TinySigned:      "42",
		TinyUnsigned:    "42",
		SmallSigned:     "42",
		SmallUnsigned:   "42",
		MediumSigned:    "42",
		MediumUnsigned:  "42",
		IntSigned:       "42",
		IntUnsigned:     "42",
		BigSigned:       "42",
		BigUnsigned:     "42",
		FloatField:      "1.23",
		DoubleField:     "3.14159",
		RealField:       "2.71828",
		DecimalField:    "1234567890.1234567890",
		DecField:        "12345.12345",
		NumericField:    "999.9999999",
		FixedField:      "9999.999999",
		Bit1:            "\x01",
		Bit8:            "\x7F",
		Bit64:           "\x00\x00\x00\x00\x00\x00\x00\x01",
		BoolField:       "1",
		BooleanField:    "0",
		CharField:       "char10___",
		VarcharField:    "varchar test",
		TextField:       "some long text",
		TinytextField:   "tinytext",
		MediumtextField: "mediumtext content",
		LongtextField:   "longtext content",
		EnumField:       "two",
		SetField:        "a,b",
		BinaryField:     string(append([]byte{0x01, 0x02, 0x03}, make([]byte, 13)...)),
		VarbinaryField:  string([]byte{0x04, 0x05, 0x06}),
		BlobField:       "blob_data",
		TinyblobField:   "tinyblob",
		MediumblobField: string(bytes.Repeat([]byte("M"), 128)),
		LongblobField:   string(bytes.Repeat([]byte("L"), 256)),
		DateField:       "2025-06-29",
		TimeField:       "12:34:56",
		YearField:       "2025",
		DatetimeField:   "2025-06-29 12:34:56.000000",
		TimestampField:  "2025-06-29 12:34:56",
		UuidField:       u,
	}

	_, err := e.DBInsert(Fields)
	if err != nil {
		t.Fatal(err)
	}

	list, err := DBSelectAll()
	if err != nil {
		t.Fatal(err)
	}

	var found *Entity
	for i := range list {
		if list[i].UuidField == u {
			found = &list[i]
			break
		}
	}
	if found == nil {
		t.Fatal("inserted entity not found")
	}

	if found.Id != e.Id {
		t.Errorf("Id mismatch: got %v, want %v", found.Id, e.Id)
	}
	if found.TinySigned != e.TinySigned {
		t.Errorf("TinySigned mismatch: got %v, want %v", found.TinySigned, e.TinySigned)
	}
	if found.TinyUnsigned != e.TinyUnsigned {
		t.Errorf("TinyUnsigned mismatch: got %v, want %v", found.TinyUnsigned, e.TinyUnsigned)
	}
	if found.SmallSigned != e.SmallSigned {
		t.Errorf("SmallSigned mismatch: got %v, want %v", found.SmallSigned, e.SmallSigned)
	}
	if found.SmallUnsigned != e.SmallUnsigned {
		t.Errorf("SmallUnsigned mismatch: got %v, want %v", found.SmallUnsigned, e.SmallUnsigned)
	}
	if found.MediumSigned != e.MediumSigned {
		t.Errorf("MediumSigned mismatch: got %v, want %v", found.MediumSigned, e.MediumSigned)
	}
	if found.MediumUnsigned != e.MediumUnsigned {
		t.Errorf("MediumUnsigned mismatch: got %v, want %v", found.MediumUnsigned, e.MediumUnsigned)
	}
	if found.IntSigned != e.IntSigned {
		t.Errorf("IntSigned mismatch: got %v, want %v", found.IntSigned, e.IntSigned)
	}
	if found.IntUnsigned != e.IntUnsigned {
		t.Errorf("IntUnsigned mismatch: got %v, want %v", found.IntUnsigned, e.IntUnsigned)
	}
	if found.BigSigned != e.BigSigned {
		t.Errorf("BigSigned mismatch: got %v, want %v", found.BigSigned, e.BigSigned)
	}
	if found.BigUnsigned != e.BigUnsigned {
		t.Errorf("BigUnsigned mismatch: got %v, want %v", found.BigUnsigned, e.BigUnsigned)
	}
	if found.FloatField != e.FloatField {
		t.Errorf("FloatField mismatch: got %v, want %v", found.FloatField, e.FloatField)
	}
	if found.DoubleField != e.DoubleField {
		t.Errorf("DoubleField mismatch: got %v, want %v", found.DoubleField, e.DoubleField)
	}
	if found.RealField != e.RealField {
		t.Errorf("RealField mismatch: got %v, want %v", found.RealField, e.RealField)
	}
	if found.DecimalField != e.DecimalField {
		t.Errorf("DecimalField mismatch: got %v, want %v", found.DecimalField, e.DecimalField)
	}
	if found.DecField != e.DecField {
		t.Errorf("DecField mismatch: got %v, want %v", found.DecField, e.DecField)
	}
	if found.NumericField != e.NumericField {
		t.Errorf("NumericField mismatch: got %v, want %v", found.NumericField, e.NumericField)
	}
	if found.FixedField != e.FixedField {
		t.Errorf("FixedField mismatch: got %v, want %v", found.FixedField, e.FixedField)
	}
	if found.Bit1 != e.Bit1 {
		t.Errorf("Bit1 mismatch: got %v, want %v", []byte(found.Bit1), []byte(e.Bit1))
	}
	if found.Bit8 != e.Bit8 {
		t.Errorf("Bit8 mismatch: got %v, want %v", []byte(found.Bit8), []byte(e.Bit8))
	}
	if found.Bit64 != e.Bit64 {
		t.Errorf("Bit64 mismatch: got %v, want %v", []byte(found.Bit64), []byte(e.Bit64))
	}
	if found.BoolField != e.BoolField {
		t.Errorf("BoolField mismatch: got %v, want %v", found.BoolField, e.BoolField)
	}
	if found.BooleanField != e.BooleanField {
		t.Errorf("BooleanField mismatch: got %v, want %v", found.BooleanField, e.BooleanField)
	}
	if found.CharField != e.CharField {
		t.Errorf("CharField mismatch: got %v, want %v", found.CharField, e.CharField)
	}
	if found.VarcharField != e.VarcharField {
		t.Errorf("VarcharField mismatch: got %v, want %v", found.VarcharField, e.VarcharField)
	}
	if found.TextField != e.TextField {
		t.Errorf("TextField mismatch: got %v, want %v", found.TextField, e.TextField)
	}
	if found.TinytextField != e.TinytextField {
		t.Errorf("TinytextField mismatch: got %v, want %v", found.TinytextField, e.TinytextField)
	}
	if found.MediumtextField != e.MediumtextField {
		t.Errorf("MediumtextField mismatch: got %v, want %v", found.MediumtextField, e.MediumtextField)
	}
	if found.LongtextField != e.LongtextField {
		t.Errorf("LongtextField mismatch: got %v, want %v", found.LongtextField, e.LongtextField)
	}
	if found.EnumField != e.EnumField {
		t.Errorf("EnumField mismatch: got %v, want %v", found.EnumField, e.EnumField)
	}
	if found.SetField != e.SetField {
		t.Errorf("SetField mismatch: got %v, want %v", found.SetField, e.SetField)
	}
	if found.BinaryField != e.BinaryField {
		t.Errorf("BinaryField mismatch: got %v, want %v", []byte(found.BinaryField), []byte(e.BinaryField))
	}
	if found.VarbinaryField != e.VarbinaryField {
		t.Errorf("VarbinaryField mismatch: got %v, want %v", []byte(found.VarbinaryField), []byte(e.VarbinaryField))
	}
	if found.BlobField != e.BlobField {
		t.Errorf("BlobField mismatch: got %v, want %v", []byte(found.BlobField), []byte(e.BlobField))
	}
	if found.TinyblobField != e.TinyblobField {
		t.Errorf("TinyblobField mismatch: got %v, want %v", []byte(found.TinyblobField), []byte(e.TinyblobField))
	}
	if found.MediumblobField != e.MediumblobField {
		t.Errorf("MediumblobField mismatch: got %v, want %v", []byte(found.MediumblobField), []byte(e.MediumblobField))
	}
	if found.LongblobField != e.LongblobField {
		t.Errorf("LongblobField mismatch: got %v, want %v", []byte(found.LongblobField), []byte(e.LongblobField))
	}
	if found.DateField != e.DateField {
		t.Errorf("DateField mismatch: got %v, want %v", found.DateField, e.DateField)
	}
	if found.TimeField != e.TimeField {
		t.Errorf("TimeField mismatch: got %v, want %v", found.TimeField, e.TimeField)
	}
	if found.YearField != e.YearField {
		t.Errorf("YearField mismatch: got %v, want %v", found.YearField, e.YearField)
	}
	if found.DatetimeField != e.DatetimeField {
		t.Errorf("DatetimeField mismatch: got %v, want %v", found.DatetimeField, e.DatetimeField)
	}
	if found.TimestampField != e.TimestampField {
		t.Errorf("TimestampField mismatch: got %v, want %v", found.TimestampField, e.TimestampField)
	}
	if found.UuidField != e.UuidField {
		t.Errorf("UuidField mismatch: got %v, want %v", found.UuidField, e.UuidField)
	}
}
