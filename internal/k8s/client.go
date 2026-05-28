package k8s

import (
	"os"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

var (
	tracingPolicyGVR = schema.GroupVersionResource{
		Group:    "cilium.io",
		Version:  "v1alpha1",
		Resource: "tracingpolicies",
	}
	tracingPolicyNamespacedGVR = schema.GroupVersionResource{
		Group:    "cilium.io",
		Version:  "v1alpha1",
		Resource: "tracingpoliciesnamespaced",
	}
	namespaceGVR = schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "namespaces",
	}
	secretGVR = schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "secrets",
	}
)

// NewDynamicClient creates a dynamic k8s client using in-cluster config.
func NewDynamicClient() (dynamic.Interface, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(cfg)
}

// CurrentNamespace returns the namespace the pod is running in.
// Falls back to "default" if the file is not found.
func CurrentNamespace() string {
	b, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return "default"
	}
	return string(b)
}
