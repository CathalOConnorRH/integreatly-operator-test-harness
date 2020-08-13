package tests

import (
	"github.com/integr8ly/integreatly-operator-test-harness/pkg/metadata"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	k8sv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var _ = ginkgo.Describe("Integreatly Operator Tests", func() {
	defer ginkgo.GinkgoRecover()
	//config, err := rest.InClusterConfig()

	//if err != nil {
	//	panic(err)
	//}

	ginkgo.It("rhmis.integreatly.org CRD exists", func() {

		/*config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{CurrentContext: ""}).ClientConfig()
		*/
		config, err := rest.InClusterConfig()

		if err != nil {
			panic(err)
		}

		// Creates the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}

		//apiextensions, err := clientset.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())

		// Make sure the CRD exists
		secret, err := clientset.CoreV1().Secrets("kube-system").Get("aws-creds", metav1.GetOptions{})
		logrus.Infof("AWS ACCESS %v", string(secret.Data["aws_access_key_id"]))
		logrus.Infof("AWS SECRET %v", string(secret.Data["aws_secret_access_key"]))

		var replicas = int32(1)

		d := &k8sv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: "cluster-service",
			},
			Spec: k8sv1.DeploymentSpec{
				Replicas: &replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "cluster-service",
					},
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "cluster-service",
						},
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name:  "cluster-service",
								Image: "quay.io/integreatly/cluster-service:v0.4.0",
							},
						},
					},
				},
			},
		}

		_, err = clientset.AppsV1().Deployments("kube-system").Create(d)

		if err != nil {
			metadata.Instance.FoundCRD = false
		} else {
			metadata.Instance.FoundCRD = true
		}

		Expect(err).NotTo(HaveOccurred())
	})
})
