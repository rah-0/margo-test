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
	results, err := QueryGetByUuid(u)
	if err != nil {
		t.Fatal("query failed:", err)
	}

	found := false
	for _, r := range results {
		if r.Animal == "dog" && r.TestField == "unique" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected row with Animal=dog and TestField=unique not found in results: %+v", results)
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

	results, err := QueryCountNullBigNumbers()
	if err != nil {
		t.Fatal("query failed:", err)
	}

	if len(results) == 0 {
		t.Fatal("expected at least 1 result")
	}

	count, err := strconv.Atoi(results[0].Count)
	if err != nil {
		t.Fatalf("invalid count returned: %v", results[0].Count)
	}

	if count < 1 {
		t.Errorf("expected at least 1 row with NULL BigNumber, got: %d", count)
	}
}

func TestQueryInsertOne(t *testing.T) {
	u := uuid.NewString()

	args := []any{
		u,
		"duck",     // Animal
		"testmark", // test_field
	}

	rows, err := QueryInsertOne(args...)
	if err != nil {
		t.Fatal("insert query failed:", err)
	}
	rows.Close()

	// Confirm inserted row
	entity := &Alpha.Entity{Uuid: u}
	ok, err := entity.DBExists([]string{Alpha.FieldUuid})
	if err != nil {
		t.Fatal("existence check failed:", err)
	}
	if !ok {
		t.Errorf("expected row with uuid %s not found in DB", u)
	}
}

func TestQueryInsertHardcoded(t *testing.T) {
	// Run the query
	rows, err := QueryInsertHardcoded()
	if err != nil {
		t.Fatal("insert hardcoded query failed:", err)
	}
	rows.Close()

	// Confirm inserted row
	entity := &Alpha.Entity{Uuid: "11111111-1111-4111-8111-111111111111"}
	ok, err := entity.DBExists([]string{Alpha.FieldUuid})
	if err != nil {
		t.Fatal("existence check failed:", err)
	}
	if !ok {
		t.Error("expected hardcoded row not found in DB")
	}
}

func TestQueryUpdateAnimalName(t *testing.T) {
	u := uuid.NewString()

	// Insert initial row
	row := &Alpha.Entity{
		Uuid:        u,
		FirstInsert: "2025-06-30 17:00:00",
		LastUpdate:  "2025-06-30 17:00:00",
		Animal:      "bear",
		BigNumber:   "100",
		TestField:   "update-test",
	}
	_, err := row.DBInsert(Alpha.Fields)
	if err != nil {
		t.Fatal("insert failed:", err)
	}

	// Perform update
	_, err = QueryUpdateAnimalName("lion", u)
	if err != nil {
		t.Fatal("update query failed:", err)
	}

	// Verify update
	entity := &Alpha.Entity{Uuid: u}
	ok, err := entity.DBExists([]string{Alpha.FieldUuid})
	if err != nil {
		t.Fatal("existence check failed:", err)
	}
	if !ok {
		t.Fatal("expected row not found after update")
	}
	if entity.Animal != "lion" {
		t.Errorf("expected Animal = lion, got: %s", entity.Animal)
	}
}

func TestQueryUpdateTestField(t *testing.T) {
	u := uuid.NewString()

	// Insert a row with Animal = 'fox'
	row := &Alpha.Entity{
		Uuid:        u,
		FirstInsert: "2025-06-30 18:00:00",
		LastUpdate:  "2025-06-30 18:00:00",
		Animal:      "fox",
		BigNumber:   "321",
		TestField:   "oldval",
	}
	_, err := row.DBInsert(Alpha.Fields)
	if err != nil {
		t.Fatal("insert failed:", err)
	}

	// Run the update
	_, err = QueryUpdateTestField()
	if err != nil {
		t.Fatal("update query failed:", err)
	}

	// Verify the change
	entity := &Alpha.Entity{Uuid: u}
	ok, err := entity.DBExists([]string{Alpha.FieldUuid})
	if err != nil {
		t.Fatal("existence check failed:", err)
	}
	if !ok {
		t.Fatal("expected row not found after update")
	}
	if entity.TestField != "updated" {
		t.Errorf("expected test_field = updated, got: %s", entity.TestField)
	}
}

func TestQueryDeleteByUuid(t *testing.T) {
	u := uuid.NewString()

	// Insert a row to delete
	row := &Alpha.Entity{
		Uuid:        u,
		FirstInsert: "2025-06-30 19:00:00",
		LastUpdate:  "2025-06-30 19:00:00",
		Animal:      "delete-me",
		BigNumber:   "123",
		TestField:   "to-delete",
	}
	_, err := row.DBInsert(Alpha.Fields)
	if err != nil {
		t.Fatal("insert failed:", err)
	}

	// Ensure row exists
	found, err := row.DBExists([]string{Alpha.FieldUuid})
	if err != nil {
		t.Fatal("existence check failed:", err)
	}
	if !found {
		t.Fatal("row not found before delete")
	}

	// Run delete
	_, err = QueryDeleteByUuid(u)
	if err != nil {
		t.Fatal("delete query failed:", err)
	}

	// Verify deletion
	found, err = row.DBExists([]string{Alpha.FieldUuid})
	if err != nil {
		t.Fatal("existence check after delete failed:", err)
	}
	if found {
		t.Errorf("row with uuid %s should have been deleted", u)
	}
}

func TestQueryDeleteOldRows(t *testing.T) {
	u := uuid.NewString()

	// Insert a row with an old LastUpdate
	row := &Alpha.Entity{
		Uuid:        u,
		FirstInsert: "2022-12-31 12:00:00",
		LastUpdate:  "2022-12-31 12:00:00", // before the cutoff
		Animal:      "old",
		BigNumber:   "1",
		TestField:   "obsolete",
	}
	_, err := row.DBInsert(Alpha.Fields)
	if err != nil {
		t.Fatal("insert failed:", err)
	}

	// Confirm row is present before deletion
	found, err := row.DBExists([]string{Alpha.FieldUuid})
	if err != nil {
		t.Fatal("existence check failed:", err)
	}
	if !found {
		t.Fatal("row should exist before deletion")
	}

	// Run delete query
	_, err = QueryDeleteOldRows()
	if err != nil {
		t.Fatal("delete query failed:", err)
	}

	// Confirm row is gone
	found, err = row.DBExists([]string{Alpha.FieldUuid})
	if err != nil {
		t.Fatal("existence recheck failed:", err)
	}
	if found {
		t.Errorf("row with uuid %s should have been deleted", u)
	}
}
