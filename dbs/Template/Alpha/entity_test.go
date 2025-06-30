package Alpha

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/rah-0/testmark/testutil"
	"xorm.io/xorm"

	_ "github.com/go-sql-driver/mysql"

	"github.com/rah-0/margo-test/ent"
	"github.com/rah-0/margo-test/util"
)

var (
	c     *sql.DB
	cXorm *xorm.Engine
	cEnt  *ent.Client
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

			cXorm, err = xorm.NewEngine("mysql", dsn)
			if err != nil {
				return err
			}

			cEnt, err = ent.Open("mysql", dsn)
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

			err = cXorm.Close()
			if err != nil {
				return err
			}

			err = cEnt.Close()
			if err != nil {
				return err
			}

			return c.Close()
		},
	})
}

func TestEntityDBInsertWithUuid(t *testing.T) {
	e := Entity{
		Uuid:   uuid.New().String(),
		Animal: "Cat",
	}

	_, err := e.DBInsert([]string{FieldUuid, FieldAnimal})
	if err != nil {
		t.Fatal(err)
	}
}

func TestEntityDBInsertWithUuidAndDelete(t *testing.T) {
	e := Entity{
		Uuid:   uuid.New().String(),
		Animal: "Dog",
	}

	_, err := e.DBInsert([]string{FieldUuid, FieldAnimal})
	if err != nil {
		t.Fatal(err)
	}

	_, err = e.DBDeleteWhereAll([]string{FieldUuid})
	if err != nil {
		t.Fatal(err)
	}
}

func TestEntityDBSelectAll(t *testing.T) {
	u := uuid.New().String()
	entity := Entity{
		Uuid:   u,
		Animal: "Fox",
	}

	_, err := entity.DBInsert([]string{FieldUuid, FieldAnimal})
	if err != nil {
		t.Fatal(err)
	}

	entities, err := DBSelectAll()
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, e := range entities {
		if e.Uuid == u && e.Animal == "Fox" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("inserted entity not found in DBSelectAll results")
	}
}
