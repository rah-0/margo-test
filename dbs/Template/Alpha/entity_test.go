package Alpha

import (
	"database/sql"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rah-0/testmark/testutil"
	"xorm.io/xorm"

	_ "github.com/go-sql-driver/mysql"

	"github.com/rah-0/margo-test/util"
)

var (
	c     *sql.DB
	cXorm *xorm.Engine
)

func TestMain(m *testing.M) {
	testutil.TestMainWrapper(testutil.TestConfig{
		M: m,
		LoadResources: func() error {
			var err error
			c, cXorm, err = util.GetConn()

			c.SetMaxIdleConns(runtime.NumCPU())
			c.SetConnMaxLifetime(time.Minute * 5)
			c.SetConnMaxIdleTime(time.Minute * 1)

			cXorm.SetMaxIdleConns(runtime.NumCPU())
			cXorm.SetConnMaxLifetime(time.Minute * 5)
			cXorm.SetConnMaxIdleTime(time.Minute * 1)
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
