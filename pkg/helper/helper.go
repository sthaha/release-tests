package helper

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/openshift-pipelines/release-tests/pkg/client"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	secv1 "github.com/openshift/api/security/v1"
	secclient "github.com/openshift/client-go/security/clientset/versioned/typed/security/v1"
	op "github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	knativetest "knative.dev/pkg/test"
)

// AssertNoError confirms the error returned is nil
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func WaitForClusterCR(t *testing.T, cs *client.Clients, name string) *op.Config {
	t.Helper()

	objKey := types.NamespacedName{Name: name}
	cr := &op.Config{}

	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		err := cs.Client.Get(context.TODO(), objKey, cr)
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s cr\n", name)
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
	AssertNoError(t, err)
	return cr
}

// WaitForDeploymentDeletion checks to see if a given deployment is deleted
// the function returns an error if the given deployment is not deleted within the timeout
func WaitForDeploymentDeletion(t *testing.T, cs *client.Clients, namespace, name string) error {
	t.Helper()
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		kc := cs.KubeClient.Kube //test.Global.KubeClient
		_, err := kc.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{IncludeUninitialized: true})
		if err != nil {
			if apierrors.IsGone(err) || apierrors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		}
		t.Logf("Waiting for deletion of %s deployment\n", name)
		return false, nil
	})
	if err == nil {
		t.Logf("%s Deployment deleted\n", name)
	}
	return err
}

func WaitForServiceAccount(t *testing.T, cs *client.Clients, ns, targetSA string) *corev1.ServiceAccount {
	t.Helper()

	ret := &corev1.ServiceAccount{}

	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		saList, err := cs.KubeClient.Kube.CoreV1().ServiceAccounts(ns).List(metav1.ListOptions{})
		for _, sa := range saList.Items {
			if sa.Name == targetSA {
				ret = &sa
				return true, nil
			}
		}
		return false, err
	})

	AssertNoError(t, err)
	return ret
}

func DeleteClusterCR(t *testing.T, cs *client.Clients, name string) {
	t.Helper()

	// ensure object exists before deletion
	objKey := types.NamespacedName{Name: name}
	cr := &op.Config{}
	err := cs.Client.Get(context.TODO(), objKey, cr)
	if err != nil {
		t.Logf("Failed to find cluster CR: %s : %s\n", name, err)
	}
	AssertNoError(t, err)

	err = wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		err := cs.Client.Delete(context.TODO(), cr)
		if err != nil {
			t.Logf("Deletion of CR %s failed %s \n", name, err)
			return false, err
		}

		return true, nil
	})

	AssertNoError(t, err)
}

func ValidateSCCAdded(t *testing.T, cs *client.Clients, ns, sa string) {
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		privileged, err := GetPrivilegedSCC(cs)
		if err != nil {
			t.Logf("failed to get privileged scc: %s \n", err)
			return false, err
		}
		t.Logf("... looking at %v", privileged.Users)

		ctrlSA := fmt.Sprintf("system:serviceaccount:%s:%s", ns, sa)
		return inList(privileged.Users, ctrlSA), nil
	})
	AssertNoError(t, err)
}

func ValidateSCCRemoved(t *testing.T, cs *client.Clients, ns, sa string) {
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		privileged, err := GetPrivilegedSCC(cs)
		if err != nil {
			t.Logf("failed to get privileged scc: %s \n", err)
			return false, err
		}

		ctrlSA := fmt.Sprintf("system:serviceaccount:%s:%s", ns, sa)
		return !inList(privileged.Users, ctrlSA), nil
	})
	AssertNoError(t, err)
}

func inList(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func ValidateDeployments(t *testing.T, cs *client.Clients, ns string, deployments ...string) {
	t.Helper()

	kc := cs.KubeClient.Kube
	for _, d := range deployments {
		err := WaitForDeployment(t, kc, ns,
			d,
			1,
			config.APIRetry,
			config.APITimeout,
		)
		AssertNoError(t, err)
	}

}

func GetPrivilegedSCC(cs *client.Clients) (*secv1.SecurityContextConstraints, error) {
	sec, err := secclient.NewForConfig(cs.KubeConfig)
	if err != nil {
		return nil, err
	}
	return sec.SecurityContextConstraints().Get("privileged", metav1.GetOptions{})
}

func ValidateDeploymentDeletion(t *testing.T, cs *client.Clients, ns string, deployments ...string) {
	t.Helper()
	for _, d := range deployments {
		err := WaitForDeploymentDeletion(t, cs, ns, d)
		AssertNoError(t, err)
	}
}

func WaitForDeployment(t *testing.T, kc kubernetes.Interface, namespace, name string, replicas int, retryInterval, timeout time.Duration) error {

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		deployment, err := kc.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{IncludeUninitialized: true})
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s deployment\n", name)
				return false, nil
			}
			return false, err
		}

		if int(deployment.Status.AvailableReplicas) == replicas {
			return true, nil
		}
		t.Logf("Waiting for full availability of %s deployment (%d/%d)\n", name, deployment.Status.AvailableReplicas, replicas)
		return false, nil
	})
	if err != nil {
		return err
	}
	t.Logf("Deployment available (%d/%d)\n", replicas, replicas)
	return nil
}

func CreateNamespace(kc *knativetest.KubeClient, ns string) {
	log.Printf("Create namespace %s to deploy to", ns)
	_, err := kc.Kube.CoreV1().Namespaces().Create(
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: ns},
		})

	if err != nil {
		log.Printf("Failed to create namespace %s for tests: %s", ns, err)
	}
}

func DeleteNamespace(kc *knativetest.KubeClient, namespace string) {
	log.Printf("Deleting namespace %s", namespace)
	if err := kc.Kube.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{}); err != nil {
		log.Printf("Failed to delete namespace %s: %s", namespace, err)
	}
}

func VerifyServiceAccountExists(kc *knativetest.KubeClient, namespace string) {
	defaultSA := "pipeline"
	log.Printf("Verify SA %q is created in namespace %q", defaultSA, namespace)

	if err := wait.PollImmediate(config.APIRetry, config.APITimeout, func() (bool, error) {
		_, err := kc.Kube.CoreV1().ServiceAccounts(namespace).Get(defaultSA, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, nil
		}
		return true, err
	}); err != nil {
		log.Printf("Failed to get SA %q in namespace %q for tests: %s", defaultSA, namespace, err)
	}
}
