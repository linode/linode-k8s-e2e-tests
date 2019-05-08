package framework

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/appscode/go/crypto/rand"
	jtypes "github.com/appscode/go/encoding/json/types"
	"github.com/appscode/go/types"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mrand "math/rand"
)

type KubedbTable struct {
	Id         int64
	Name       string
	IntValue   int64
	FloatValue float64
}

func (i *Invocation) Postgres(isLocalStorage bool) *api.Postgres {
	postgres := &api.Postgres{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix(api.ResourceSingularPostgres),
			Namespace: i.namespace,
			Labels: map[string]string{
				"app": i.app,
			},
		},
		Spec: api.PostgresSpec{
			Version:  jtypes.StrYo(DBCatalogName),
			Replicas: types.Int32P(3),
		},
	}
	if isLocalStorage {
		postgres.Spec.StorageType = "Ephemeral"
	} else {
		postgres.Spec.Storage = &core.PersistentVolumeClaimSpec{
			Resources: core.ResourceRequirements{
				Requests: core.ResourceList{
					core.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
			StorageClassName: types.StringP(StorageClass),
		}
	}

	return postgres
}

func (f *Framework) CreatePostgres(obj *api.Postgres) error {
	_, err := f.extClient.KubedbV1alpha1().Postgreses(obj.Namespace).Create(obj)
	return err
}

func (f *Framework) GetPostgres(meta metav1.ObjectMeta) (*api.Postgres, error) {
	return f.extClient.KubedbV1alpha1().Postgreses(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
}

func (f *Framework) DeletePostgres(meta metav1.ObjectMeta) error {
	return f.extClient.KubedbV1alpha1().Postgreses(meta.Namespace).Delete(meta.Name, &metav1.DeleteOptions{})
}

func (f *Framework) EventuallyPostgresRunning(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	return Eventually(
		func() bool {
			postgres, err := f.extClient.KubedbV1alpha1().Postgreses(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			return postgres.Status.Phase == api.DatabasePhaseRunning
		},
		time.Minute*15,
		time.Second*5,
	)
}

func (f *Framework) DeletePostgresPod(meta metav1.ObjectMeta) error {
	return f.kubeClient.CoreV1().Pods(meta.Namespace).Delete(meta.Name+"-0", &metav1.DeleteOptions{})
}

func (f *Framework) CheckPostgresPod(meta metav1.ObjectMeta) (*core.Pod, error) {
	return f.kubeClient.CoreV1().Pods(meta.Namespace).Get(meta.Name+"-0", metav1.GetOptions{})
}

func (f *Framework) EventuallyInsertRow(meta metav1.ObjectMeta, dbName, userName string, total int) GomegaAsyncAssertion {
	count := 0
	return Eventually(
		func() bool {
			clientPodName := f.GetPrimaryPodName(meta)
			tunnel, err := f.ForwardPort(meta, clientPodName)
			if err != nil {
				return false
			}
			defer tunnel.Close()

			en, err := f.GetPostgresClient(tunnel, dbName, userName)
			if err != nil {
				return false
			}
			defer en.Close()

			if err := en.Ping(); err != nil {
				return false
			}

			err = en.Sync(new(KubedbTable))
			if err != nil {
				fmt.Println("creation error", err)
				return false
			}

			for i := count; i < total; i++ {
				if _, err := en.Insert(&KubedbTable{
					Name:       fmt.Sprintf("KubedbName-%v", i),
					IntValue:   mrand.Int63(),
					FloatValue: mrand.Float64(),
				}); err != nil {
					return false
				}
				count++
			}
			return true
		},
		time.Minute*10,
		time.Second*10,
	)
}
