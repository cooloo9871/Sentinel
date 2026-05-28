package k8s

import (
	"context"
	"encoding/base64"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

// GetCredentials reads the sentinel-credentials Secret and returns a map of
// username → bcrypt hash. Secret data values are base64-encoded by the k8s API.
func GetCredentials(ctx context.Context, client dynamic.Interface, namespace string) (map[string]string, error) {
	obj, err := client.Resource(secretGVR).Namespace(namespace).Get(ctx, "sentinel-credentials", metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get sentinel-credentials secret: %w", err)
	}

	rawData, _, _ := unstructured.NestedMap(obj.Object, "data")
	result := make(map[string]string, len(rawData))
	for k, v := range rawData {
		str, ok := v.(string)
		if !ok {
			continue
		}
		decoded, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			return nil, fmt.Errorf("decode secret key %q: %w", k, err)
		}
		result[k] = string(decoded)
	}
	return result, nil
}
