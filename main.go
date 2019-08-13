package main

import (
	"crypto/rand"
	//"flag"
	"fmt"
	"github.com/speak2jc/k-op/pkg/apis/example/v1alpha1"
	keevaClientSet "github.com/speak2jc/k-op/pkg/generated/clientset/versioned"
	apixv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os/user"
	"path/filepath"
	"strings"
	//keevaRest "github.com/speak2jc/k-op/vendor/k8s.io/client-go/rest
	//keevaClientCmd "github.com/speak2jc/k-op/vendor/k8s.io/client-go
)

var (
	keevakindGVR = schema.GroupVersionResource{
		Group:    "example.keeva.com",
		Version:  "v1alpha1",
		Resource: "keevakinds",
	}
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	log.Print("Loading client config")
	//config, err := clientcmd.BuildConfigFromFlags("", userConfig())
	//errExit("Failed to load client conifg", err)

	log.Print("Loading dynamic client")
	//client, err := dynamic.NewForConfig(config)

	//errExit("Failed to create client", err)

	name := "keeva-" + RandomString(3)
	namespace := "james"

	//RegisterRuntimeClassCRD(config)

	//keevaConf, err := clientcmd.BuildConfigFromFlags("", "")
	//if err != nil {
	//	log.Fatalf("Error building kubeconfig: %s", err.Error())
	//}
	//var keevaClient *keevaClientSet.Clientset
	//var keevaConf *keevaRest.Config
	//flag.Parse()
	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, "/Users/jamez/.kube/config")
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}
	keevaClient, err := keevaClientSet.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Error building keeva clientset: %s", err.Error())
	}

	// Create
	CreateKeevakind(keevaClient, name, namespace)

	//// Get
	//existingKeevakind, err := GetKeevakind(client, name, namespace)
	//
	//if err == nil {
	//	log.Printf("Retrieved Keevakind %s", existingKeevakind)
	//
	//	//payload := existingKeevakind.Object["spec"].(map[string]interface{})
	//	//payload["group"] = "mygroup"
	//
	//	existingKeevakind.Spec.Group = "mygroup1"
	//	//log.Printf("Updating payload %s", payload)
	//
	//	// Update
	//	UpdateKeevakind(client, existingKeevakind)
	//
	//	// Get again
	//	existingKeevakind, _ = GetKeevakind(client, name, namespace)
	//	log.Printf("Retrieved Keevakind %s", existingKeevakind)
	//
	//}
}

func CreateKeevakind(keevaClient *keevaClientSet.Clientset, name string, namespace string) (*v1alpha1.Keevakind, error) {

	image := "Image-" + RandomString(5)
	group := "Group-" + RandomString(5)
	var port int32 = 8080
	var count int32 = 14

	log.Printf("Creating Keevakind %s", name)

	keevakind := v1alpha1.Keevakind{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec:       v1alpha1.KeevakindSpec{Count: count, Group: group, Port: port, Image: image},
		Status:     v1alpha1.KeevakindStatus{},
	}

	return keevaClient.ExampleV1alpha1().Keevakinds(namespace).Create(&keevakind)
}

func UpdateKeevakind(client dynamic.Interface, keevakind v1alpha1.Keevakind) {

	res := client.Resource(keevakindGVR)
	converted := mapToUnstructured(keevakind)
	rc, err := res.Update(converted, metav1.UpdateOptions{})
	errExit(fmt.Sprintf("Failed to update Keevakind %#v", rc), err)
}

func GetKeevakind2(client dynamic.Interface, name string, namespace string) (v1alpha1.Keevakind, error) {

	log.Printf("Getting Keevakind %s", name)

	var keevakind v1alpha1.Keevakind
	res := client.Resource(keevakindGVR)
	existingKeeva, err := res.Get(name, metav1.GetOptions{})
	errExit(fmt.Sprintf("Failed to Get Keevakind %s in namespace %s", name, namespace), err)

	if existingKeeva == nil {
		err := errors.NewNotFound(schema.GroupResource{"example.keeva.com", "keevakind"}, name)
		return keevakind, err
	}

	keevakind = mapToKeevakind(existingKeeva, namespace)
	return keevakind, nil
}

func GetKeevakind(client dynamic.Interface, name string, namespace string) (v1alpha1.Keevakind, error) {

	log.Printf("Getting Keevakind %s", name)

	var keevakind v1alpha1.Keevakind
	res := client.Resource(keevakindGVR)
	existingKeeva, err := res.Get(name, metav1.GetOptions{})
	errExit(fmt.Sprintf("Failed to Get Keevakind %s in namespace %s", name, namespace), err)

	if existingKeeva == nil {
		err := errors.NewNotFound(schema.GroupResource{"example.keeva.com", "keevakind"}, name)
		return keevakind, err
	}

	keevakind = mapToKeevakind(existingKeeva, namespace)
	return keevakind, nil
}

func mapToKeevakind(in *unstructured.Unstructured, namespace string) v1alpha1.Keevakind {

	objMap := in.Object

	//metadata := objMap["metadata"].(map[string]interface{})
	//apiVersion := objMap["apiVersion"].(string)
	//kind := objMap["kind"].(string)
	//name := metadata["name"].(string)
	//resourceVersion := metadata["resourceVersion"].(string)
	//	namespace := metadata["namespace"].(string) - TODO - add to CRD - currently does not seem to store it

	spec := objMap["spec"].(map[string]interface{})
	group := spec["group"].(string)
	image := spec["image"].(string)
	port := int32(spec["port"].(int64))
	count := int32(spec["count"].(int64))

	keevakind := v1alpha1.Keevakind{
		//TypeMeta:   metav1.TypeMeta{APIVersion: apiVersion, Kind: kind},
		//ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, ResourceVersion: resourceVersion},
		Spec:   v1alpha1.KeevakindSpec{Count: count, Group: group, Port: port, Image: image},
		Status: v1alpha1.KeevakindStatus{},
	}

	return keevakind
}

func mapToUnstructured(in v1alpha1.Keevakind) *unstructured.Unstructured {

	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Keevakind",
			"apiVersion": keevakindGVR.Group + "/v1alpha1",
			"metadata": map[string]interface{}{
				"resourceVersion": in.ObjectMeta.ResourceVersion,
				"name":            in.Name,
				"namespace":       in.Namespace,
			},
			"spec": v1alpha1.KeevakindSpec{
				Count: in.Spec.Count,
				Group: in.Spec.Group,
				Image: in.Spec.Image,
				Port:  in.Spec.Port,
			},
		},
	}

}

func errExit(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %#v", msg, err)
	}
}

func userConfig() string {
	usr, err := user.Current()
	errExit("Failed to get current user", err)
	return filepath.Join(usr.HomeDir, ".kube", "config")
}

func RandomString(len int) string {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	fmt.Println(s)

	return strings.ToLower(s)
}

func RegisterRuntimeClassCRD(config *rest.Config) {
	apixClient, err := apixv1beta1client.NewForConfig(config)
	errExit("Failed to load apiextensions client", err)

	crds := apixClient.CustomResourceDefinitions()

	runtimeClassCRD := &apixv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "keevakinds.example.keeva.com",
		},
		Spec: apixv1beta1.CustomResourceDefinitionSpec{
			Group:   "example.keeva.com",
			Version: "v1alpha1",
			Versions: []apixv1beta1.CustomResourceDefinitionVersion{{
				Name:    "v1alpha1",
				Served:  true,
				Storage: true,
			}},
			Names: apixv1beta1.CustomResourceDefinitionNames{
				Plural:   "keevakinds",
				Singular: "keevakind",
				Kind:     "Keevakind",
			},
			Scope: apixv1beta1.ClusterScoped,
			Validation: &apixv1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &apixv1beta1.JSONSchemaProps{
					Properties: map[string]apixv1beta1.JSONSchemaProps{
						"spec": {
							Properties: map[string]apixv1beta1.JSONSchemaProps{
								"runtimeHandler": {
									Type: "string",
									//Pattern: "abc",
								},
								"kind": {
									Type: "string",
								},
							},
						},
					},
				},
			},
		},
	}
	log.Print("Registering Keevakind CRD")
	_, err = crds.Create(runtimeClassCRD)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Print("Keevakind CRD already registered")
		} else {
			errExit("Failed to create Keevakind CRD", err)
		}
	}
}
