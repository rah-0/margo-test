package Template

import (
	"database/sql"
	"strconv"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/rah-0/testmark/testutil"

	"github.com/rah-0/margo-test/dbs/Template/AllTypes"
	"github.com/rah-0/margo-test/dbs/Template/Alpha"
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

			AllTypes.SetDB(c)
			Alpha.SetDB(c)

			return SetDB(c)
		},
		UnloadResources: func() error {
			result1 := AllTypes.DBTruncate()
			if result1.Error != nil {
				return result1.Error
			}

			result2 := Alpha.DBTruncate()
			if result2.Error != nil {
				return result2.Error
			}

			return c.Close()
		},
	})
}

func TestQueryGetAllAnimals(t *testing.T) {
	u := uuid.NewString()

	// Insert uniquely identifiable row
	row := &Alpha.Entity{
		Uuid:        u,
		FirstInsert: "2025-06-30 12:00:00",
		LastUpdate:  "2025-06-30 12:00:00",
		Animal:      "cat",
		BigNumber:   "9000",
		TestField:   "test",
	}
	result := row.DBInsert(Alpha.NewQueryParams().WithInsert(Alpha.Fields...))
	if result.Error != nil {
		t.Fatal("insert failed:", result.Error)
	}

	qr := Alpha.QueryGetAllAnimals()
	if qr.Error != nil {
		t.Fatal("query failed:", qr.Error)
	}

	found := false
	for _, r := range qr.Entities {
		if r.Animal == "cat" && r.BigNumber == "9000" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected row with Animal=cat and BigNumber=9000 not found in results: %+v", qr.Entities)
	}
}

func TestQueryGetRecentCats(t *testing.T) {
	u := uuid.NewString()

	row := &Alpha.Entity{
		Uuid:        u,
		FirstInsert: "2025-06-30 12:00:00",
		LastUpdate:  "2025-06-30 13:00:00",
		Animal:      "cat",
		BigNumber:   "12345",
		TestField:   "recent",
	}
	result := row.DBInsert(Alpha.NewQueryParams().WithInsert(Alpha.Fields...))
	if result.Error != nil {
		t.Fatal("insert failed:", result.Error)
	}

	qr := QueryGetRecentCats()
	if qr.Error != nil {
		t.Fatal("query failed:", qr.Error)
	}

	found := false
	for _, r := range qr.Results {
		if r.Uuid == u {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected row with uuid %s not found in results: %+v", u, qr.Results)
	}
}

func TestQueryGetByUuid(t *testing.T) {
	u := uuid.NewString()

	// Insert a row with known values
	row := &Alpha.Entity{
		Uuid:        u,
		FirstInsert: "2025-06-30 15:00:00",
		LastUpdate:  "2025-06-30 15:00:00",
		Animal:      "dog",
		BigNumber:   "5555",
		TestField:   "unique",
	}
	result := row.DBInsert(Alpha.NewQueryParams().WithInsert(Alpha.Fields...))
	if result.Error != nil {
		t.Fatal("insert failed:", result.Error)
	}

	// Run the query using uuid as argument
	qr := QueryGetByUuid(NewQueryParams().WithParams(u))
	if qr.Error != nil {
		t.Fatal("query failed:", qr.Error)
	}

	found := false
	if len(qr.Results) > 0 {
		r := qr.Results[0]
		if r.Animal == "dog" && r.TestField == "unique" {
			found = true
		}
	}

	if !found {
		t.Errorf("expected row with Animal=dog and TestField=unique not found in results: %+v", qr.Results)
	}
}

func TestQueryCountBigNumbers(t *testing.T) {
	u := uuid.NewString()

	// Insert one row with BigNumber omitted (will be NULL)
	row := &Alpha.Entity{
		Uuid:        u,
		FirstInsert: "2025-06-30 16:00:00",
		LastUpdate:  "2025-06-30 16:00:00",
		Animal:      "nulltest",
		TestField:   "checknull",
	}
	result := row.DBInsert(Alpha.NewQueryParams().WithInsert(
		Alpha.FieldUuid,
		Alpha.FieldFirstInsert,
		Alpha.FieldLastUpdate,
		Alpha.FieldAnimal,
		Alpha.FieldTestField, // BigNumber is skipped = NULL
	))
	if result.Error != nil {
		t.Fatal("insert failed:", result.Error)
	}

	qr := QueryCountBigNumbers()
	if qr.Error != nil {
		t.Fatal("query failed:", qr.Error)
	}

	if len(qr.Results) == 0 {
		t.Fatal("no results returned")
	}

	count, err := strconv.Atoi(qr.Results[0].Count)
	if err != nil {
		t.Fatalf("invalid count returned: %v", qr.Results[0].Count)
	}

	if count < 1 {
		t.Errorf("expected at least 1 row with NULL BigNumber, got: %d", count)
	}
}

func TestExecInsertOne(t *testing.T) {
	u := uuid.NewString()
	qr := ExecInsertOne(NewQueryParams().WithParams(u, "hedgehog", "tf"))
	if qr.Error != nil {
		t.Fatal("insert failed:", qr.Error)
	}

	r := QueryGetByUuid(NewQueryParams().WithParams(u))
	if r.Error != nil {
		t.Fatal("query failed:", r.Error)
	}
	if len(r.Results) == 0 || r.Results[0].Animal != "hedgehog" || r.Results[0].TestField != "tf" {
		t.Fatalf("row not inserted as expected: %+v", r.Results)
	}
}

func TestExecInsertHardcoded(t *testing.T) {
	const hard = "11111111-1111-4111-8111-111111111111"

	// ensure a clean slate for this uuid
	_ = ExecDeleteByUuid(NewQueryParams().WithParams(hard))

	qr := ExecInsertHardcoded()
	if qr.Error != nil {
		t.Fatal("insert hardcoded failed:", qr.Error)
	}

	r := QueryGetByUuid(NewQueryParams().WithParams(hard))
	if r.Error != nil {
		t.Fatal("query failed:", r.Error)
	}
	if len(r.Results) == 0 || r.Results[0].Animal != "dog" {
		t.Fatalf("expected Animal=dog for hardcoded uuid, got: %+v", r.Results)
	}
}

func TestExecUpdateAnimalName(t *testing.T) {
	u := uuid.NewString()
	row := &Alpha.Entity{
		Uuid:        u,
		FirstInsert: "2025-06-30 10:00:00",
		LastUpdate:  "2025-06-30 10:00:00",
		Animal:      "cat",
		TestField:   "x",
	}
	result := row.DBInsert(Alpha.NewQueryParams().WithInsert(
		Alpha.FieldUuid, Alpha.FieldFirstInsert, Alpha.FieldLastUpdate, Alpha.FieldAnimal, Alpha.FieldTestField,
	))
	if result.Error != nil {
		t.Fatal("seed insert failed:", result.Error)
	}

	qr := ExecUpdateAnimalName(NewQueryParams().WithParams("otter", u))
	if qr.Error != nil {
		t.Fatal("update failed:", qr.Error)
	}

	r := QueryGetByUuid(NewQueryParams().WithParams(u))
	if r.Error != nil {
		t.Fatal("query failed:", r.Error)
	}
	if len(r.Results) == 0 || r.Results[0].Animal != "otter" {
		t.Fatalf("expected Animal=otter after update, got: %+v", r.Results)
	}
}

func TestExecUpdateTestField(t *testing.T) {
	u := uuid.NewString()
	row := &Alpha.Entity{
		Uuid:        u,
		FirstInsert: "2025-06-30 11:00:00",
		LastUpdate:  "2025-06-30 11:00:00",
		Animal:      "fox",
		TestField:   "old",
	}
	result := row.DBInsert(Alpha.NewQueryParams().WithInsert(
		Alpha.FieldUuid, Alpha.FieldFirstInsert, Alpha.FieldLastUpdate, Alpha.FieldAnimal, Alpha.FieldTestField,
	))
	if result.Error != nil {
		t.Fatal("seed insert failed:", result.Error)
	}

	qr := ExecUpdateTestField()
	if qr.Error != nil {
		t.Fatal("update failed:", qr.Error)
	}

	r := QueryGetByUuid(NewQueryParams().WithParams(u))
	if r.Error != nil {
		t.Fatal("query failed:", r.Error)
	}
	if len(r.Results) == 0 || r.Results[0].TestField != "updated" {
		t.Fatalf("expected test_field=updated after bulk update, got: %+v", r.Results)
	}
}

func TestExecDeleteByUuid(t *testing.T) {
	u := uuid.NewString()
	row := &Alpha.Entity{
		Uuid:        u,
		FirstInsert: "2025-06-30 12:00:00",
		LastUpdate:  "2025-06-30 12:00:00",
		Animal:      "toad",
		TestField:   "y",
	}
	result := row.DBInsert(Alpha.NewQueryParams().WithInsert(
		Alpha.FieldUuid, Alpha.FieldFirstInsert, Alpha.FieldLastUpdate, Alpha.FieldAnimal, Alpha.FieldTestField,
	))
	if result.Error != nil {
		t.Fatal("seed insert failed:", result.Error)
	}

	qr := ExecDeleteByUuid(NewQueryParams().WithParams(u))
	if qr.Error != nil {
		t.Fatal("delete failed:", qr.Error)
	}

	r := QueryGetByUuid(NewQueryParams().WithParams(u))
	if r.Error != nil {
		t.Fatal("query failed:", r.Error)
	}
	if len(r.Results) > 0 {
		t.Fatalf("expected no row after delete, got: %+v", r.Results)
	}
}

func TestExecDeleteOldRows(t *testing.T) {
	oldU := uuid.NewString()
	newU := uuid.NewString()

	oldRow := &Alpha.Entity{
		Uuid:        oldU,
		FirstInsert: "2022-12-31 23:59:59",
		LastUpdate:  "2022-12-31 23:59:59",
		Animal:      "ant",
		TestField:   "old",
	}
	newRow := &Alpha.Entity{
		Uuid:        newU,
		FirstInsert: "2025-01-01 00:00:01",
		LastUpdate:  "2025-01-01 00:00:01",
		Animal:      "bee",
		TestField:   "new",
	}
	result := oldRow.DBInsert(Alpha.NewQueryParams().WithInsert(
		Alpha.FieldUuid, Alpha.FieldFirstInsert, Alpha.FieldLastUpdate, Alpha.FieldAnimal, Alpha.FieldTestField,
	))
	if result.Error != nil {
		t.Fatal("seed old insert failed:", result.Error)
	}
	result = newRow.DBInsert(Alpha.NewQueryParams().WithInsert(
		Alpha.FieldUuid, Alpha.FieldFirstInsert, Alpha.FieldLastUpdate, Alpha.FieldAnimal, Alpha.FieldTestField,
	))
	if result.Error != nil {
		t.Fatal("seed new insert failed:", result.Error)
	}

	qr := ExecDeleteOldRows()
	if qr.Error != nil {
		t.Fatal("delete old rows failed:", qr.Error)
	}

	ro := QueryGetByUuid(NewQueryParams().WithParams(oldU))
	if ro.Error != nil {
		t.Fatal("query old failed:", ro.Error)
	}
	if len(ro.Results) > 0 {
		t.Fatalf("expected old row to be deleted, got: %+v", ro.Results)
	}

	rn := QueryGetByUuid(NewQueryParams().WithParams(newU))
	if rn.Error != nil {
		t.Fatal("query new failed:", rn.Error)
	}
	if len(rn.Results) == 0 || rn.Results[0].Animal != "bee" {
		t.Fatalf("expected new row to remain, got: %+v", rn.Results)
	}
}
