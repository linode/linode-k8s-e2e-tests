package framework

import (
	"fmt"
	"strings"
	"time"

	v1 "github.com/appscode/kutil/core/v1"
	"github.com/kubedb/postgres/pkg/controller"
	"github.com/kubedb/postgres/pkg/leader_election"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (f *Framework) GetPrimaryPodName(meta metav1.ObjectMeta) string {
	postgres, err := f.GetPostgres(meta)
	Expect(err).NotTo(HaveOccurred())
	Expect(postgres.Spec.Replicas).NotTo(BeNil())

	if *postgres.Spec.Replicas == 1 {
		return fmt.Sprintf("%v-0", postgres.Name)
	}

	pods, err := f.kubeClient.CoreV1().Pods(meta.Namespace).List(metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{
			MatchLabels: v1.UpsertMap(map[string]string{
				controller.NodeRole: leader_election.RolePrimary,
			}, postgres.OffshootSelectors()),
		}),
	})

	Expect(err).NotTo(HaveOccurred())
	Expect(len(pods.Items)).To(Equal(1))

	return pods.Items[0].Name
}

func (f *Framework) GetArbitraryStandbyPodName(meta metav1.ObjectMeta) string {
	postgres, err := f.GetPostgres(meta)
	Expect(err).NotTo(HaveOccurred())
	Expect(postgres.Spec.Replicas).NotTo(BeNil())

	if *postgres.Spec.Replicas == 1 {
		return ""
	}

	pods, err := f.kubeClient.CoreV1().Pods(meta.Namespace).List(metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{
			MatchLabels: v1.UpsertMap(map[string]string{
				controller.NodeRole: leader_election.RoleReplica,
			}, postgres.OffshootSelectors()),
		}),
	})
	Expect(err).NotTo(HaveOccurred())

	return pods.Items[0].Name
}

func (f *Framework) EventuallyStreamingReplication(meta metav1.ObjectMeta, clientPodName, dbName, userName string) GomegaAsyncAssertion {
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

			results, err := db.Query("select * from pg_stat_replication;")
			if err != nil {
				return -1
			}

			for _, result := range results {
				applicationName := string(result["application_name"])
				state := string(result["state"])
				if state != "streaming" || !strings.HasPrefix(applicationName, meta.Name) {
					return -1
				}
			}
			return len(results)
		},
		time.Minute*10,
		time.Second*5,
	)
}
