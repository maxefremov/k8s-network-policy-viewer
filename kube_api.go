package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"strings"
)

// don't allow errors to bubble up
// TODO: switch to byte array as that's what is used in the end
func getJsonData(buffer *string) {
	*buffer = "{}"
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("Can't create in-cluster configuration: %s\n", err.Error())
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Can't create clientset: %s\n", err.Error())
		return
	}

	// pods
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if errors.IsNotFound(err) {
		log.Printf("Pods not found\n")
		return
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error listing pods: %v\n", statusError.ErrStatus.Message)
		return
	} else if err != nil {
		log.Printf("API error: %s\n", err.Error())
		return
	}

	for index, _ := range pods.Items {
		pods.Items[index].Kind = "Pod"
	}

	podsJson, err := json.Marshal(pods.Items)
	if err != nil {
		log.Printf("Can't marshal pods: %s\n", err.Error())
		return
	}

	// namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if errors.IsNotFound(err) {
		log.Printf("Namespaces not found\n")
		return
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error listing namespaces %v\n", statusError.ErrStatus.Message)
		return
	} else if err != nil {
		log.Printf("API error: %s\n", err.Error())
		return
	}

	for index, _ := range namespaces.Items {
		pods.Items[index].Kind = "Namespace"
	}

	namespacesJson, err := json.Marshal(pods.Items)
	if err != nil {
		log.Printf("Can't marshal namespaces: %s\n", err.Error())
		return
	}

	// network policies
	networkPolicies, err := clientset.NetworkingV1().NetworkPolicies("").List(metav1.ListOptions{})
	if errors.IsNotFound(err) {
		log.Printf("Network policies not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error listing network policies: %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		log.Printf("API error: %s\n", err.Error())
		return
	}

	for index, _ := range networkPolicies.Items {
		networkPolicies.Items[index].Kind = "NetworkPolicy"
	}

	networkPoliciesJson, err := json.Marshal(networkPolicies.Items)
	if err != nil {
		log.Printf("Can't marshal network policies: %s\n", err.Error())
		return
	}

	// stringify, trim, assemble
	podsJsonString := string(podsJson)
	namespacesJsonString := string(namespacesJson)
	networkPoliciesJsonString := string(networkPoliciesJson)

	trimBrackets(&podsJsonString)
	trimBrackets(&namespacesJsonString)
	trimBrackets(&networkPoliciesJsonString)

	*buffer = fmt.Sprintf("{\"kind\":\"List\",\"apiVersion\":\"v1\",\"Items\":[%s,%s,%s]}", podsJsonString, namespacesJsonString, networkPoliciesJsonString)
}

func trimBrackets(s *string) {
	*s = strings.Trim(*s, "[]")
}
