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
			_, err := DBTruncate()
			if err != nil {
				return err
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
	_, err := e.DBInsert([]string{FieldUuid, FieldName})
	if err != nil {
		t.Fatal("insert failed:", err)
	}

	// Attempt to manually override last_update with a past timestamp
	expected := "2000-01-01 00:00:00.123456"
	e.LastUpdate = expected

	_, err = e.DBUpdateWhereAll([]string{FieldLastUpdate}, []string{FieldUuid})
	if err != nil {
		t.Fatal("update failed:", err)
	}

	// Fetch back using DBExists
	var check Entity
	check.Uuid = e.Uuid
	found, err := check.DBExists([]string{FieldUuid})
	if err != nil {
		t.Fatal("DBExists failed:", err)
	}
	if !found {
		t.Fatal("entity not found after update")
	}

	if check.LastUpdate != expected {
		t.Fatalf("expected last_update %s, got %s", expected, check.LastUpdate)
	}
}
