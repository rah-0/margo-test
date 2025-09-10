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
			_, err := AllTypes.DBTruncate()
			if err != nil {
				return err
			}

			_, err = Alpha.DBTruncate()
			if err != nil {
				return err
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
	_, err := row.DBInsert(Alpha.Fields)
	if err != nil {
		t.Fatal("insert failed:", err)
	}

	results, err := QueryGetAllAnimals()
	if err != nil {
		t.Fatal("query failed:", err)
	}

	found := false
	for _, r := range results {
		if r.Animal == "cat" && r.BigNumber == "9000" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected row with Animal=cat and BigNumber=9000 not found in results: %+v", results)
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
	_, err := row.DBInsert(Alpha.Fields)
	if err != nil {
		t.Fatal("insert failed:", err)
	}

	results, err := QueryGetRecentCats()
	if err != nil {
		t.Fatal("query failed:", err)
	}

	found := false
	for _, r := range results {
		if r.Uuid == u {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected row with uuid %s not found in results: %+v", u, results)
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
	_, err := row.DBInsert(Alpha.Fields)
	if err != nil {
		t.Fatal("insert failed:", err)
	}

	// Run the query using uuid as argument
	r, err := QueryGetByUuid(u)
	if err != nil {
		t.Fatal("query failed:", err)
	}

	found := false
	if r.Animal == "dog" && r.TestField == "unique" {
		found = true
	}

	if !found {
		t.Errorf("expected row with Animal=dog and TestField=unique not found in results: %+v", r)
	}
}

func TestQueryCountNullBigNumbers(t *testing.T) {
	u := uuid.NewString()

	// Insert one row with BigNumber omitted (will be NULL)
	row := &Alpha.Entity{
		Uuid:        u,
		FirstInsert: "2025-06-30 16:00:00",
		LastUpdate:  "2025-06-30 16:00:00",
		Animal:      "nulltest",
		TestField:   "checknull",
	}
	_, err := row.DBInsert([]string{
		Alpha.FieldUuid,
		Alpha.FieldFirstInsert,
		Alpha.FieldLastUpdate,
		Alpha.FieldAnimal,
		Alpha.FieldTestField, // BigNumber is skipped = NULL
	})
	if err != nil {
		t.Fatal("insert failed:", err)
	}

	result, err := QueryCountNullBigNumbers()
	if err != nil {
		t.Fatal("query failed:", err)
	}

	count, err := strconv.Atoi(result.Count)
	if err != nil {
		t.Fatalf("invalid count returned: %v", result.Count)
	}

	if count < 1 {
		t.Errorf("expected at least 1 row with NULL BigNumber, got: %d", count)
	}
}

func TestExecInsertOne(t *testing.T) {
	u := uuid.NewString()
	_, err := ExecInsertOne(u, "hedgehog", "tf")
	if err != nil {
		t.Fatal("insert failed:", err)
	}

	r, err := QueryGetByUuid(u)
	if err != nil {
		t.Fatal("query failed:", err)
	}
	if r == nil || r.Animal != "hedgehog" || r.TestField != "tf" {
		t.Fatalf("row not inserted as expected: %+v", r)
	}
}

func TestExecInsertHardcoded(t *testing.T) {
	const hard = "11111111-1111-4111-8111-111111111111"

	// ensure a clean slate for this uuid
	_, _ = ExecDeleteByUuid(hard)

	if _, err := ExecInsertHardcoded(); err != nil {
		t.Fatal("insert hardcoded failed:", err)
	}

	r, err := QueryGetByUuid(hard)
	if err != nil {
		t.Fatal("query failed:", err)
	}
	if r == nil || r.Animal != "dog" {
		t.Fatalf("expected Animal=dog for hardcoded uuid, got: %+v", r)
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
	if _, err := row.DBInsert([]string{
		Alpha.FieldUuid, Alpha.FieldFirstInsert, Alpha.FieldLastUpdate, Alpha.FieldAnimal, Alpha.FieldTestField,
	}); err != nil {
		t.Fatal("seed insert failed:", err)
	}

	if _, err := ExecUpdateAnimalName("otter", u); err != nil {
		t.Fatal("update failed:", err)
	}

	r, err := QueryGetByUuid(u)
	if err != nil {
		t.Fatal("query failed:", err)
	}
	if r == nil || r.Animal != "otter" {
		t.Fatalf("expected Animal=otter after update, got: %+v", r)
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
	if _, err := row.DBInsert([]string{
		Alpha.FieldUuid, Alpha.FieldFirstInsert, Alpha.FieldLastUpdate, Alpha.FieldAnimal, Alpha.FieldTestField,
	}); err != nil {
		t.Fatal("seed insert failed:", err)
	}

	if _, err := ExecUpdateTestField(); err != nil {
		t.Fatal("update failed:", err)
	}

	r, err := QueryGetByUuid(u)
	if err != nil {
		t.Fatal("query failed:", err)
	}
	if r == nil || r.TestField != "updated" {
		t.Fatalf("expected test_field=updated after bulk update, got: %+v", r)
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
	if _, err := row.DBInsert([]string{
		Alpha.FieldUuid, Alpha.FieldFirstInsert, Alpha.FieldLastUpdate, Alpha.FieldAnimal, Alpha.FieldTestField,
	}); err != nil {
		t.Fatal("seed insert failed:", err)
	}

	if _, err := ExecDeleteByUuid(u); err != nil {
		t.Fatal("delete failed:", err)
	}

	r, err := QueryGetByUuid(u)
	if err != nil {
		t.Fatal("query failed:", err)
	}
	if r != nil {
		t.Fatalf("expected no row after delete, got: %+v", r)
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
	if _, err := oldRow.DBInsert([]string{
		Alpha.FieldUuid, Alpha.FieldFirstInsert, Alpha.FieldLastUpdate, Alpha.FieldAnimal, Alpha.FieldTestField,
	}); err != nil {
		t.Fatal("seed old insert failed:", err)
	}
	if _, err := newRow.DBInsert([]string{
		Alpha.FieldUuid, Alpha.FieldFirstInsert, Alpha.FieldLastUpdate, Alpha.FieldAnimal, Alpha.FieldTestField,
	}); err != nil {
		t.Fatal("seed new insert failed:", err)
	}

	if _, err := ExecDeleteOldRows(); err != nil {
		t.Fatal("delete old rows failed:", err)
	}

	ro, err := QueryGetByUuid(oldU)
	if err != nil {
		t.Fatal("query old failed:", err)
	}
	if ro != nil {
		t.Fatalf("expected old row to be deleted, got: %+v", ro)
	}

	rn, err := QueryGetByUuid(newU)
	if err != nil {
		t.Fatal("query new failed:", err)
	}
	if rn == nil || rn.Animal != "bee" {
		t.Fatalf("expected new row to remain, got: %+v", rn)
	}
}
