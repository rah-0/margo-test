package Alpha

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"

	gormysql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/rah-0/margo-test/ent/alpha"
)

type Alpha struct {
	bun.BaseModel `bun:"table:alpha"`
	Uuid          string `bun:"Uuid,pk" gorm:"column:Uuid;primaryKey" xorm:"'Uuid' pk"`
	FirstInsert   string `bun:"FirstInsert" gorm:"column:FirstInsert" xorm:"'FirstInsert'"`
	LastUpdate    string `bun:"LastUpdate" gorm:"column:LastUpdate" xorm:"'LastUpdate'"`
	Animal        string `bun:"Animal" gorm:"column:Animal" xorm:"'Animal'"`
	BigNumber     string `bun:"BigNumber" gorm:"column:BigNumber" xorm:"'BigNumber'"`
	TestField     string `bun:"test_field" gorm:"column:test_field" xorm:"'test_field'"`
}

func (Alpha) TableName() string {
	return "alpha"
}

func BenchmarkEntityDBInsertRawSQL(b *testing.B) {
	b.Skip()

	result := DBTruncate()
	if err := result.Error; err != nil {
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
	b.Skip()

	result := DBTruncate()
	if err := result.Error; err != nil {
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

		result := e.DBInsert(NewQueryParams().WithInsert(Fields...))
		if err := result.Error; err != nil {
			b.Fatalf("insert failed: %v", err)
		}
	}
}

func BenchmarkEntityDBInsertBun(b *testing.B) {
	b.Skip()

	result := DBTruncate()
	if err := result.Error; err != nil {
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
	b.Skip()

	result := DBTruncate()
	if err := result.Error; err != nil {
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

func BenchmarkEntityDBInsertEnt(b *testing.B) {
	b.Skip()

	result := DBTruncate()
	if err := result.Error; err != nil {
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

func BenchmarkEntityDBDeleteRawSQL(b *testing.B) {
	b.Skip()

	result := DBTruncate()
	if err := result.Error; err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	insertStmt, err := db.Prepare("INSERT INTO alpha (`Uuid`, `FirstInsert`, `LastUpdate`, `Animal`, `BigNumber`, `test_field`) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		b.Fatal(err)
	}
	defer insertStmt.Close()

	deleteStmt, err := db.Prepare("DELETE FROM alpha WHERE `Uuid` = ?")
	if err != nil {
		b.Fatal(err)
	}
	defer deleteStmt.Close()

	uuids := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		u := uuid.NewString()
		uuids[i] = u
		_, err := insertStmt.Exec(
			u,
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := deleteStmt.Exec(uuids[i]); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEntityDBDeleteMarGO(b *testing.B) {
	b.Skip()

	result := DBTruncate()
	if err := result.Error; err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	entities := make([]*Entity, b.N)
	for i := 0; i < b.N; i++ {
		e := &Entity{
			Uuid:        uuid.NewString(),
			FirstInsert: "2024-01-01 15:04:05.000000",
			LastUpdate:  "2024-01-01 15:04:05.000000",
			Animal:      "Animal",
			BigNumber:   "1234567890",
			TestField:   "Test",
		}
		result := e.DBInsert(NewQueryParams().WithInsert(Fields...))
		if err := result.Error; err != nil {
			b.Fatal(err)
		}
		entities[i] = e
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := entities[i].DBDelete(NewQueryParams().WithWhere(FieldUuid))
	if err := result.Error; err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEntityDBDeleteGorm(b *testing.B) {
	b.Skip()

	result := DBTruncate()
	if err := result.Error; err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	gdb, err := gorm.Open(gormysql.New(gormysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		b.Fatalf("gorm open failed: %v", err)
	}

	entities := make([]*Alpha, b.N)
	for i := 0; i < b.N; i++ {
		e := &Alpha{
			Uuid:        uuid.NewString(),
			FirstInsert: "2024-01-01 15:04:05.000000",
			LastUpdate:  "2024-01-01 15:04:05.000000",
			Animal:      "Animal",
			BigNumber:   "1234567890",
			TestField:   "Test",
		}
		if err := gdb.Create(e).Error; err != nil {
			b.Fatal(err)
		}
		entities[i] = e
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := gdb.Delete(entities[i]).Error; err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEntityDBDeleteBun(b *testing.B) {
	b.Skip()

	result := DBTruncate()
	if err := result.Error; err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	dbx := bun.NewDB(db, mysqldialect.New())

	entities := make([]*Alpha, b.N)
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
		entities[i] = e
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := dbx.NewDelete().Model(entities[i]).
			Where("`Uuid` = ?", entities[i].Uuid).
			Exec(context.Background()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEntityDBDeleteEnt(b *testing.B) {
	b.Skip()

	result := DBTruncate()
	if err := result.Error; err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	ctx := context.Background()

	uuids := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		u := uuid.NewString()
		_, err := cEnt.Alpha.Create().
			SetUUID(u).
			SetFirstInsert("2024-01-01 15:04:05.000000").
			SetLastUpdate("2024-01-01 15:04:05.000000").
			SetAnimal("Animal").
			SetBigNumber("1234567890").
			SetTestField("Test").
			Save(ctx)
		if err != nil {
			b.Fatal(err)
		}
		uuids[i] = u
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := cEnt.Alpha.Delete().
			Where(alpha.UUIDEQ(uuids[i])).
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEntityDBSelectRawSQL(b *testing.B) {
	b.Skip()

	result := DBTruncate()
	if err := result.Error; err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	stmtInsert, err := db.Prepare("INSERT INTO alpha (`Uuid`, `FirstInsert`, `LastUpdate`, `Animal`, `BigNumber`, `test_field`) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		b.Fatal(err)
	}
	defer stmtInsert.Close()

	uuids := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		u := uuid.NewString()
		uuids[i] = u
		_, err := stmtInsert.Exec(
			u,
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

	stmtSelect, err := db.Prepare("SELECT `Uuid`, `FirstInsert`, `LastUpdate`, `Animal`, `BigNumber`, `test_field` FROM alpha WHERE `Uuid` = ?")
	if err != nil {
		b.Fatal(err)
	}
	defer stmtSelect.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var e Entity
		var ptrUuid, ptrFirstInsert, ptrLastUpdate, ptrAnimal, ptrBigNumber, ptrTestField *string
		row := stmtSelect.QueryRow(uuids[i])
		if err := row.Scan(
			&ptrUuid,
			&ptrFirstInsert,
			&ptrLastUpdate,
			&ptrAnimal,
			&ptrBigNumber,
			&ptrTestField,
		); err != nil {
			b.Fatal(err)
		}
		_ = e // skip use
	}
}

func BenchmarkEntityDBSelectMarGO(b *testing.B) {
	b.Skip()

	result := DBTruncate()
	if err := result.Error; err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	entities := make([]*Entity, b.N)
	for i := 0; i < b.N; i++ {
		e := &Entity{
			Uuid:        uuid.NewString(),
			FirstInsert: "2024-01-01 15:04:05.000000",
			LastUpdate:  "2024-01-01 15:04:05.000000",
			Animal:      "Animal",
			BigNumber:   "1234567890",
			TestField:   "Test",
		}
		result := e.DBInsert(NewQueryParams().WithInsert(Fields...))
		if err := result.Error; err != nil {
			b.Fatal(err)
		}
		entities[i] = e
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e := &Entity{
			Uuid: entities[i].Uuid,
		}
		result := e.DBExists(NewQueryParams().WithWhere(FieldUuid))
	ok, err := result.Exists, result.Error
		if err != nil {
			b.Fatal(err)
		}
		if !ok {
			b.Fatal("record not found")
		}
	}
}

func BenchmarkEntityDBSelectBun(b *testing.B) {
	b.Skip()

	result := DBTruncate()
	if err := result.Error; err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	dbx := bun.NewDB(db, mysqldialect.New())

	entities := make([]*Alpha, b.N)
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
		entities[i] = e
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result Alpha
		err := dbx.NewSelect().
			Model(&result).
			Where("`Uuid` = ?", entities[i].Uuid).
			Limit(1).
			Scan(context.Background())
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEntityDBSelectGorm(b *testing.B) {
	b.Skip()

	result := DBTruncate()
	if err := result.Error; err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	gdb, err := gorm.Open(gormysql.New(gormysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		b.Fatalf("gorm open failed: %v", err)
	}

	uuids := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		u := uuid.NewString()
		e := Alpha{
			Uuid:        u,
			FirstInsert: "2024-01-01 15:04:05.000000",
			LastUpdate:  "2024-01-01 15:04:05.000000",
			Animal:      "Animal",
			BigNumber:   "1234567890",
			TestField:   "Test",
		}
		if err := gdb.Create(&e).Error; err != nil {
			b.Fatal(err)
		}
		uuids[i] = u
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var e Alpha
		if err := gdb.First(&e, "Uuid = ?", uuids[i]).Error; err != nil {
			b.Fatal(err)
		}
	}
}
