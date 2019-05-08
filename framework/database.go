package framework

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/appscode/kutil/tools/portforward"
	"github.com/go-xorm/xorm"
	"github.com/kubedb/postgres/pkg/controller"
	_ "github.com/lib/pq"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (f *Framework) ForwardPort(meta metav1.ObjectMeta, clientPodName string) (*portforward.Tunnel, error) {
	tunnel := portforward.NewTunnel(
		f.kubeClient.CoreV1().RESTClient(),
		f.restConfig,
		meta.Namespace,
		clientPodName,
		controller.PostgresPort,
	)
	if err := tunnel.ForwardPort(); err != nil {
		return nil, err
	}
	return tunnel, nil
}

func (f *Framework) GetPostgresClient(tunnel *portforward.Tunnel, dbName string, userName string) (*xorm.Engine, error) {
	cnnstr := fmt.Sprintf("user=%s host=127.0.0.1 port=%v dbname=%s sslmode=disable", userName, tunnel.Local, dbName)
	return xorm.NewEngine("postgres", cnnstr)
}

func (f *Framework) EventuallyCreateSchema(meta metav1.ObjectMeta, dbName, userName string) GomegaAsyncAssertion {
	sql := fmt.Sprintf(`
DROP SCHEMA IF EXISTS "data" CASCADE;
CREATE SCHEMA "data" AUTHORIZATION "%s";`, userName)
	return Eventually(
		func() bool {
			clientPodName := f.GetPrimaryPodName(meta)
			tunnel, err := f.ForwardPort(meta, clientPodName)
			if err != nil {
				return false
			}
			defer tunnel.Close()

			db, err := f.GetPostgresClient(tunnel, dbName, userName)
			if err != nil {
				return false
			}
			defer db.Close()

			if err := f.CheckPostgres(db); err != nil {
				return false
			}

			_, err = db.Exec(sql)
			if err != nil {
				return false
			}
			return true
		},
		time.Minute*5,
		time.Second*5,
	)
}

var randChars = []rune("abcdefghijklmnopqrstuvwxyzabcdef")

// Use this for generating random pat of a ID. Do not use this for generating short passwords or secrets.
func characters(len int) string {
	bytes := make([]byte, len)
	rand.Read(bytes)
	r := make([]rune, len)
	for i, b := range bytes {
		r[i] = randChars[b>>3]
	}
	return string(r)
}

func (f *Framework) EventuallyPingDatabase(meta metav1.ObjectMeta, clientPodName, dbName, userName string) GomegaAsyncAssertion {
	return Eventually(
		func() bool {
			tunnel, err := f.ForwardPort(meta, clientPodName)
			if err != nil {
				return false
			}
			defer tunnel.Close()

			db, err := f.GetPostgresClient(tunnel, dbName, userName)
			if err != nil {
				return false
			}
			defer db.Close()

			if err := f.CheckPostgres(db); err != nil {
				return false
			}

			return true
		},
		time.Minute*10,
		time.Second*5,
	)
}

func (f *Framework) EventuallyCreateTable(meta metav1.ObjectMeta, dbName, userName string, total int) GomegaAsyncAssertion {
	count := 0
	return Eventually(
		func() bool {
			clientPodName := f.GetPrimaryPodName(meta)

			tunnel, err := f.ForwardPort(meta, clientPodName)
			if err != nil {
				return false
			}
			defer tunnel.Close()

			db, err := f.GetPostgresClient(tunnel, dbName, userName)
			if err != nil {
				return false
			}
			defer db.Close()

			if err := f.CheckPostgres(db); err != nil {
				return false
			}

			for i := count; i < total; i++ {
				table := fmt.Sprintf("SET search_path TO \"data\"; CREATE TABLE %v ( id bigserial )", characters(5))
				_, err := db.Exec(table)
				if err != nil {
					return false
				}
				count++
			}
			return true
		},
		time.Minute*5,
		time.Second*5,
	)

	return nil
}

func (f *Framework) EventuallyCountTable(meta metav1.ObjectMeta, clientPodName, dbName, userName string) GomegaAsyncAssertion {
	return Eventually(
		func() int {
			tunnel, err := f.ForwardPort(meta, clientPodName)
			if err != nil {
				return -1
			}
			defer tunnel.Close()

			db, err := f.GetPostgresClient(tunnel, dbName, userName)
			if err != nil {
				return -1
			}
			defer db.Close()

			if err := f.CheckPostgres(db); err != nil {
				return -1
			}

			res, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='data'")
			if err != nil {
				return -1
			}

			return len(res)
		},
		time.Minute*10,
		time.Second*5,
	)
}

func (f *Framework) EventuallyCountTableFromPrimary(meta metav1.ObjectMeta, dbName, userName string) GomegaAsyncAssertion {
	return Eventually(
		func() int {
			clientPodName := f.GetPrimaryPodName(meta)

			tunnel, err := f.ForwardPort(meta, clientPodName)
			if err != nil {
				return -1
			}
			defer tunnel.Close()

			db, err := f.GetPostgresClient(tunnel, dbName, userName)
			if err != nil {
				return -1
			}
			defer db.Close()

			if err := f.CheckPostgres(db); err != nil {
				return -1
			}

			res, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='data'")
			if err != nil {
				return -1
			}
			return len(res)
		},
		time.Minute*10,
		time.Second*5,
	)
}

func (f *Framework) CheckPostgres(db *xorm.Engine) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}

type PgStatArchiver struct {
	ArchivedCount int
}
