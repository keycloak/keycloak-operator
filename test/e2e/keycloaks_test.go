package e2e

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"k8s.io/client-go/kubernetes"

	v1 "k8s.io/api/core/v1"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/model"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	v1apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewKeycloaksCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksCRWithExtension,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakDeploymentTest": {testFunction: keycloakDeploymentTest},
		},
	}
}

func NewKeycloaksWithLabelsCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksCRWithPodLabels,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakWithPodLabelsDeploymentTest": {testFunction: keycloakDeploymentWithLabelsTest},
		},
	}
}

func NewUnmanagedKeycloaksCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareUnmanagedKeycloaksCR,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakUnmanagedDeploymentTest": {testFunction: keycloakUnmanagedDeploymentTest},
		},
	}
}

func NewKeycloaksSSLTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksSSLWithDB,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakSSLDBTest": {testFunction: keycloakSSLDBTest},
		},
	}
}

func keycloakSSLDBTest(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	// get the Keycloak Statefulset
	keycloakStatefulset := v1apps.StatefulSet{}
	err := GetNamespacedObject(f, namespace, model.KeycloakDeploymentName, &keycloakStatefulset)
	if err != nil {
		return err
	}

	// check pod has the env var JDBC_PARAMS
	envExists := false
	for _, envvar := range keycloakStatefulset.Spec.Template.Spec.Containers[0].Env {
		if envvar.Name == "JDBC_PARAMS" && strings.Contains(envvar.Value, "sslmode") {
			envExists = true
			break
		}
	}
	if !envExists {
		return errors.Errorf("test Failed : Env var JDBC_PARAMS and value sslMode not found")
	}

	// check the volume to the crt exists too
	volumeExists := false
	for _, vol := range keycloakStatefulset.Spec.Template.Spec.Volumes {
		if vol.Name == model.DatabaseSecretSslCert+"-vol" {
			volumeExists = true
			break
		}
	}
	if !volumeExists {
		return errors.Errorf("test Failed : Volume to the secret not found")
	}
	return err
}

func prepareKeycloaksSSLWithDB(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	cr := getKeycloakCR(namespace)

	// create secret with crt
	secretWithSSLCertForPostgres := getSecretWithSSLCertForPostgres(namespace)
	err := f.Client.Create(context.TODO(), secretWithSSLCertForPostgres, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}

	// get secret with Postgres parameters
	secret := model.DatabaseSecret(cr)

	// deploy PostgreSQL
	postgresvc, err := deployPostgreSQLWithSSLon(secretWithSSLCertForPostgres, cr, secret, f, ctx)
	if err != nil {
		return err
	}

	// deploying the Keycloak CR
	cr.Spec.ExternalDatabase.Enabled = true
	secret.Data["SSLMODE"] = []byte("verify-ca")
	secret.Data["POSTGRES_EXTERNAL_ADDRESS"] = []byte(postgresvc.Name + "." + namespace + ".svc.cluster.local")

	err = f.Client.Create(context.TODO(), secret, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}

	err = deployKeycloaksCR(t, f, ctx, namespace, cr)

	return err
}

func deployPostgreSQLWithSSLon(secretWithSSLCertForPostgres *v1.Secret, cr *keycloakv1alpha1.Keycloak, secret *v1.Secret, f *framework.Framework, ctx *framework.Context) (*v1.Service, error) {
	// create postgre deployment
	// Create config map with config for Postgresql to start with SSL
	postgresqlConfFile, _ := ioutil.ReadFile("testdata/postgresql.conf")

	postgreSQLConfig := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgre-sql-config",
			Namespace: cr.Namespace,
			Labels:    CreateLabel(cr.Namespace),
		},
		Data: map[string]string{
			"custom.conf": string(postgresqlConfFile),
		},
	}
	err := f.Client.Create(context.TODO(), &postgreSQLConfig, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return nil, err
	}
	// add a volume+volumemount in postgresql to the secret with the crt
	modeCrt := int32(0444)
	modeKey := int32(0440)
	volume := v1.Volume{
		Name: model.DatabaseSecretSslCert + "-vol",
		VolumeSource: v1.VolumeSource{
			Projected: &v1.ProjectedVolumeSource{
				Sources: []v1.VolumeProjection{
					{
						Secret: &v1.SecretProjection{
							LocalObjectReference: v1.LocalObjectReference{
								Name: secretWithSSLCertForPostgres.Name,
							},
							Items: []v1.KeyToPath{
								{Key: "server.crt", Path: "server.crt", Mode: &modeCrt},
								{Key: "server.key", Path: "server.key", Mode: &modeKey},
							},
						},
					},
				},
			},
		},
	}
	volumeConfig := v1.Volume{
		Name: "postgre-sql-config-vol",
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "postgre-sql-config",
				},
			},
		},
	}
	volumeMount := v1.VolumeMount{
		Name:      model.DatabaseSecretSslCert + "-vol",
		MountPath: "/opt/app-root/src/certificates/",
	}
	volumeMountConfig := v1.VolumeMount{
		Name:      "postgre-sql-config-vol",
		MountPath: "/opt/app-root/src/postgresql-cfg/",
	}

	pvc := model.PostgresqlPersistentVolumeClaim(cr)
	// changing the name as apparently another one is created by the operator
	pvc.Name = externalPostgresClaim
	err = f.Client.Create(context.TODO(), pvc, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return nil, err
	}
	postgresql := model.PostgresqlDeployment(cr, false)
	//postgresql.Spec.Template.Spec.Containers[0].Image = "postgres:10.5-alpine"
	postgresql.Spec.Template.Spec.Volumes = append(postgresql.Spec.Template.Spec.Volumes, volume, volumeConfig)
	postgresql.Spec.Template.Spec.Containers[0].VolumeMounts = append(postgresql.Spec.Template.Spec.Containers[0].VolumeMounts, volumeMount, volumeMountConfig)
	runasuser := int64(26)
	fsgroup := int64(26)
	postgresql.Spec.Template.Spec.SecurityContext = &v1.PodSecurityContext{
		RunAsUser:          &runasuser,
		FSGroup:            &fsgroup,
		SupplementalGroups: []int64{999, 1000},
	}
	for _, vol := range postgresql.Spec.Template.Spec.Volumes {
		if vol.PersistentVolumeClaim != nil && vol.PersistentVolumeClaim.ClaimName == "keycloak-postgresql-claim" {
			vol.PersistentVolumeClaim.ClaimName = externalPostgresClaim
		}
	}

	err = f.Client.Create(context.TODO(), postgresql, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return nil, err
	}

	postgresvc := model.PostgresqlService(cr, secret, false)
	postgresvc.Name = "keycloak-postgresql-svc"
	err = f.Client.Create(context.TODO(), postgresvc, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return nil, err
	}
	return postgresvc, nil
}

func getSecretWithSSLCertForPostgres(namespace string) *v1.Secret {
	serverCrt, _ := ioutil.ReadFile("testdata/server.crt")
	serverKey, _ := ioutil.ReadFile("testdata/server.key")
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      model.DatabaseSecretSslCert,
			Namespace: namespace,
			Labels:    CreateLabel(namespace),
		},
		Data: map[string][]byte{
			"server.crt": serverCrt,
			"server.key": serverKey,
			"root.crt":   serverCrt,
		},
	}
}

func getKeycloakCR(namespace string) *keycloakv1alpha1.Keycloak {
	return &keycloakv1alpha1.Keycloak{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testKeycloakCRName,
			Namespace: namespace,
			Labels:    CreateLabel(namespace),
		},
		Spec: keycloakv1alpha1.KeycloakSpec{
			Instances:      1,
			ExternalAccess: keycloakv1alpha1.KeycloakExternalAccess{Enabled: true},
			Profile:        currentProfile(),
		},
	}
}

func getUnmanagedKeycloakCR(namespace string) *keycloakv1alpha1.Keycloak {
	keycloak := getKeycloakCR(namespace)
	keycloak.Name = testKeycloakUnmanagedCRName
	keycloak.Spec.Unmanaged = true
	return keycloak
}

func getExternalKeycloakCR(namespace string, url string) *keycloakv1alpha1.Keycloak {
	keycloak := getUnmanagedKeycloakCR(namespace)
	keycloak.Name = testKeycloakExternalCRName
	keycloak.Labels = CreateExternalLabel(namespace)
	keycloak.Spec.External.Enabled = true
	keycloak.Spec.External.URL = url
	return keycloak
}

func getDeployedKeycloakCR(f *framework.Framework, namespace string) keycloakv1alpha1.Keycloak {
	keycloakCR := keycloakv1alpha1.Keycloak{}
	_ = GetNamespacedObject(f, namespace, testKeycloakCRName, &keycloakCR)
	return keycloakCR
}

func getExternalKeycloakSecret(f *framework.Framework, namespace string) (*v1.Secret, error) {
	secret, err := f.KubeClient.CoreV1().Secrets(namespace).Get(context.TODO(), "credential-"+testKeycloakCRName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "credential-" + testKeycloakExternalCRName,
			Namespace: namespace,
		},
		Data:       secret.Data,
		StringData: secret.StringData,
		Type:       secret.Type,
	}, nil
}

func prepareKeycloaksCR(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	return deployKeycloaksCR(t, f, ctx, namespace, getKeycloakCR(namespace))
}

func prepareKeycloaksCRWithExtension(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	keycloakCR := getKeycloakCR(namespace)
	keycloakCR.Spec.Extensions = []string{"https://github.com/aerogear/keycloak-metrics-spi/releases/download/2.5.3/keycloak-metrics-spi-2.5.3.jar"}

	return deployKeycloaksCR(t, f, ctx, namespace, keycloakCR)
}

func prepareKeycloaksCRWithPodLabels(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	keycloakCR := getKeycloakCR(namespace)
	keycloakCR.Spec.KeycloakDeploymentSpec.PodLabels = map[string]string{"cr.first.label": "first.value", "cr.second.label": "second.value"}
	return deployKeycloaksCR(t, f, ctx, namespace, keycloakCR)
}

func deployKeycloaksCR(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string, keycloakCR *keycloakv1alpha1.Keycloak) error {
	err := doWorkaroundIfNecessary(f, ctx, namespace)
	if err != nil {
		return err
	}

	err = Create(f, keycloakCR, ctx)
	if err != nil {
		return err
	}

	err = WaitForStatefulSetReplicasReady(t, f.KubeClient, model.ApplicationName, namespace)
	if err != nil {
		return err
	}

	return err
}

func prepareUnmanagedKeycloaksCR(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	err := doWorkaroundIfNecessary(f, ctx, namespace)
	if err != nil {
		return err
	}

	keycloakCR := getUnmanagedKeycloakCR(namespace)
	err = Create(f, keycloakCR, ctx)
	if err != nil {
		return err
	}

	err = WaitForKeycloakToBeReady(t, f, namespace, testKeycloakUnmanagedCRName)
	if err != nil {
		return err
	}

	return err
}

func prepareExternalKeycloaksCR(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	keycloakCR := getDeployedKeycloakCR(f, namespace)
	keycloakURL := keycloakCR.Status.ExternalURL

	secret, err := getExternalKeycloakSecret(f, namespace)
	if err != nil {
		return err
	}

	err = Create(f, secret, ctx)
	if err != nil {
		return err
	}

	externalKeycloakCR := getExternalKeycloakCR(namespace, keycloakURL)
	err = Create(f, externalKeycloakCR, ctx)
	if err != nil {
		return err
	}

	err = WaitForKeycloakToBeReady(t, f, namespace, testKeycloakExternalCRName)
	if err != nil {
		return err
	}

	return err
}

func keycloakDeploymentTest(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	keycloakCR := getDeployedKeycloakCR(f, namespace)
	assert.NotEmpty(t, keycloakCR.Status.InternalURL)
	assert.NotEmpty(t, keycloakCR.Status.ExternalURL)

	err := WaitForKeycloakToBeReady(t, f, namespace, testKeycloakCRName)
	if err != nil {
		return err
	}

	keycloakURL := keycloakCR.Status.ExternalURL

	// Skipping TLS verification is actually part of the test. In Kubernetes, if there's no signing
	// manager installed, Keycloak will generate its own, self-signed cert. Of course
	// we don't have a matching truststore for it, hence we need to skip TLS verification.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint

	err = WaitForSuccessResponse(t, f, keycloakURL+"/auth")
	if err != nil {
		return err
	}

	metricsBody, err := GetSuccessfulResponseBody(keycloakURL + "/auth/realms/master/metrics")
	if err != nil {
		return err
	}

	masterRealmBody, err := GetSuccessfulResponseBody(keycloakURL + "/auth/realms/master")
	if err != nil {
		return err
	}

	// there should be a redirect/rewrite from the metrics endpoint to master realm
	assert.Equal(t, masterRealmBody, metricsBody)

	return err
}
func keycloakDeploymentWithLabelsTest(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	// check that the creation labels are present
	keycloakPod := v1.Pod{}
	podName := "keycloak-0"
	_ = GetNamespacedObject(f, namespace, podName, &keycloakPod)
	assert.Contains(t, keycloakPod.Labels, "cr.first.label")
	assert.Contains(t, keycloakPod.Labels, "cr.second.label")

	//add runtime labels to the pod (as if it was the existing labels from previous installation)
	keycloakPod.ObjectMeta.Labels["pod.label.one"] = "value1"
	keycloakPod.ObjectMeta.Labels["pod.label.two"] = "value2"
	err := Update(f, &keycloakPod)
	if err != nil {
		return err
	}

	//modify the CR adding labels, to see ifthe reconcile process also adds the labels
	keycloakCR := getDeployedKeycloakCR(f, namespace)
	newlabels := map[string]string{"cr-reconc.label.one": "value1", "cr-reconc.label.two": "value1"}
	keycloakCR.Spec.KeycloakDeploymentSpec.PodLabels = model.AddPodLabels(&keycloakCR, newlabels)
	err = Update(f, &keycloakCR)
	if err != nil {
		return err
	}

	// we need to wait for the reconciliation
	err = WaitForPodHavingLabels(t, f.KubeClient, podName, namespace, keycloakCR.Spec.KeycloakDeploymentSpec.PodLabels)
	if err != nil {
		return err
	}

	// assert that runtime  labels added directly to the pod are still there
	// assert that new labels added to the CR are also present in the pod
	_ = GetNamespacedObject(f, namespace, podName, &keycloakPod)
	// Labels set in the CR on the creation
	assert.Contains(t, keycloakPod.Labels, "cr.first.label")
	assert.Contains(t, keycloakPod.Labels, "cr.second.label")
	// Labels in the pod set by the user
	assert.Contains(t, keycloakPod.Labels, "pod.label.one")
	assert.Contains(t, keycloakPod.Labels, "pod.label.two")
	// Labels added to the CR during runtime
	assert.Contains(t, keycloakPod.Labels, "cr-reconc.label.one")
	assert.Contains(t, keycloakPod.Labels, "cr-reconc.label.two")

	return nil
}

func keycloakUnmanagedDeploymentTest(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	keycloakCR := getDeployedKeycloakCR(f, namespace)
	assert.Empty(t, keycloakCR.Status.InternalURL)
	assert.Empty(t, keycloakCR.Status.ExternalURL)

	err := WaitForCondition(t, f.KubeClient, func(t *testing.T, c kubernetes.Interface) error {
		sts, err := f.KubeClient.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return errors.Errorf("list StatefulSet failed, ignoring for %v: %v", pollRetryInterval, err)
		}
		if len(sts.Items) == 0 {
			return nil
		}
		return errors.Errorf("found Statefulsets, this shouldn't be the case")
	})
	return err
}
