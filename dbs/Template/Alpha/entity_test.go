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
			result := DBTruncate()
			if result.Error != nil {
				return result.Error
			}

			err := cXorm.Close()
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

	result := e.DBInsert(NewQueryParams().WithInsert(FieldUuid, FieldAnimal))
	if result.Error != nil {
		t.Fatal(result.Error)
	}
}

func TestEntityDBInsertWithUuidAndDelete(t *testing.T) {
	e := Entity{
		Uuid:   uuid.New().String(),
		Animal: "Dog",
	}

	result := e.DBInsert(NewQueryParams().WithInsert(FieldUuid, FieldAnimal))
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	result = e.DBDelete(NewQueryParams().WithWhere(FieldUuid))
	if result.Error != nil {
		t.Fatal(result.Error)
	}
}

func TestEntityDBSelectAll(t *testing.T) {
	u := uuid.New().String()
	entity := Entity{
		Uuid:   u,
		Animal: "Fox",
	}

	result := entity.DBInsert(NewQueryParams().WithInsert(FieldUuid, FieldAnimal))
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	result = DBSelectAll()
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	found := false
	for _, e := range result.Entities {
		if e.Uuid == u && e.Animal == "Fox" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("inserted entity not found in DBSelectAll results")
	}
}
