package Beta

import (
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
			result := DBTruncate()
			if result.Error != nil {
				return result.Error
			}

			return c.Close()
		},
	})
}

func TestEntityLastUpdateManualOverride(t *testing.T) {
	e := Entity{
		Uuid: uuid.New().String(),
		Name: "manual-update-test",
	}

	// Insert entity
	result := e.DBInsert(NewQueryParams().WithInsert(FieldUuid, FieldName))
	if result.Error != nil {
		t.Fatal("insert failed:", result.Error)
	}

	// Attempt to manually override last_update with a past timestamp
	expected := "2000-01-01 00:00:00.123456"
	e.LastUpdate = expected

	result = e.DBUpdate(NewQueryParams().WithUpdate(FieldLastUpdate).WithWhere(FieldUuid))
	if result.Error != nil {
		t.Fatal("update failed:", result.Error)
	}

	// Fetch back using DBExists
	var check Entity
	check.Uuid = e.Uuid
	result = check.DBExists(NewQueryParams().WithWhere(FieldUuid))
	if result.Error != nil {
		t.Fatal("DBExists failed:", result.Error)
	}
	if !result.Exists {
		t.Fatal("entity not found after update")
	}

	if check.LastUpdate != expected {
		t.Fatalf("expected last_update %s, got %s", expected, check.LastUpdate)
	}
}
