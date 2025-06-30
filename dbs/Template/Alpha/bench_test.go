package Alpha

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"

	gormysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Alpha struct {
	bun.BaseModel `bun:"table:alpha"`
	Uuid          string `bun:"Uuid,pk" gorm:"column:Uuid;primaryKey"`
	FirstInsert   string `bun:"FirstInsert" gorm:"column:FirstInsert"`
	LastUpdate    string `bun:"LastUpdate" gorm:"column:LastUpdate"`
	Animal        string `bun:"Animal" gorm:"column:Animal"`
	BigNumber     string `bun:"BigNumber" gorm:"column:BigNumber"`
	TestField     string `bun:"test_field" gorm:"column:test_field"`
}

func (Alpha) TableName() string {
	return "alpha"
}

func BenchmarkEntityDBInsertRawSQL(b *testing.B) {
	_, err := DBTruncate()
	if err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	stmt, err := db.Prepare("INSERT INTO alpha (`Uuid`, `FirstInsert`, `LastUpdate`, `Animal`, `BigNumber`, `test_field`) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		b.Fatal(err)
	}
	defer stmt.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := stmt.Exec(
			uuid.NewString(),
			"2024-01-01 15:04:05.000000",
			"2024-01-01 15:04:05.000000",
			"Animal",
			"1234567890",
			"Test",
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEntityDBInsertMarGO(b *testing.B) {
	_, err := DBTruncate()
	if err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		e := Entity{
			Uuid:        uuid.NewString(),
			FirstInsert: "2024-01-01 15:04:05.000000",
			LastUpdate:  "2024-01-01 15:04:05.000000",
			Animal:      "Animal",
			BigNumber:   "1234567890",
			TestField:   "Test",
		}
		b.StartTimer()

		_, err := e.DBInsert(Fields)
		if err != nil {
			b.Fatalf("insert failed: %v", err)
		}
	}
}

func BenchmarkEntityDBInsertBun(b *testing.B) {
	_, err := DBTruncate()
	if err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	dbx := bun.NewDB(db, mysqldialect.New())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e := &Alpha{
			Uuid:        uuid.NewString(),
			FirstInsert: "2024-01-01 15:04:05.000000",
			LastUpdate:  "2024-01-01 15:04:05.000000",
			Animal:      "Animal",
			BigNumber:   "1234567890",
			TestField:   "Test",
		}
		if _, err := dbx.NewInsert().Model(e).Exec(context.Background()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEntityDBInsertGorm(b *testing.B) {
	_, err := DBTruncate()
	if err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	gdb, err := gorm.Open(gormysql.New(gormysql.Config{
		Conn: db, // reuse existing *sql.DB
	}), &gorm.Config{})
	if err != nil {
		b.Fatalf("gorm open failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e := Alpha{
			Uuid:        uuid.NewString(),
			FirstInsert: "2024-01-01 15:04:05.000000",
			LastUpdate:  "2024-01-01 15:04:05.000000",
			Animal:      "Animal",
			BigNumber:   "1234567890",
			TestField:   "Test",
		}
		if err := gdb.Create(&e).Error; err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEntityDBInsertXORM(b *testing.B) {
	_, err := DBTruncate()
	if err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	type Alpha struct {
		Uuid        string `xorm:"'Uuid' pk"`
		FirstInsert string `xorm:"'FirstInsert'"`
		LastUpdate  string `xorm:"'LastUpdate'"`
		Animal      string `xorm:"'Animal'"`
		BigNumber   string `xorm:"'BigNumber'"`
		TestField   string `xorm:"'test_field'"`
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e := &Alpha{
			Uuid:        uuid.NewString(),
			FirstInsert: "2024-01-01 15:04:05.000000",
			LastUpdate:  "2024-01-01 15:04:05.000000",
			Animal:      "Animal",
			BigNumber:   "1234567890",
			TestField:   "Test",
		}
		if _, err := cXorm.Insert(e); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEntityDBInsertEnt(b *testing.B) {
	_, err := DBTruncate()
	if err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := cEnt.Alpha.Create().
			SetUUID(uuid.NewString()).
			SetFirstInsert("2024-01-01 15:04:05.000000").
			SetLastUpdate("2024-01-01 15:04:05.000000").
			SetAnimal("Animal").
			SetBigNumber("1234567890").
			SetTestField("Test").
			Save(ctx)

		if err != nil {
			b.Fatal(err)
		}
	}
}
