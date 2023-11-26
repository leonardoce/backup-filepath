package injector

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
)

// BackupAdapterAnnotationName is the name of the annotation storing the configuration of the
// backup adapter
const BackupAdapterAnnotationName = "cnpg.io/backupAdapter"

// AdapterConfiguration contains the configuration for an external backup adapter
type AdapterConfiguration struct {
	// Id is the adapter ID, used by the injector
	ID string `json:"id"`

	// Parameters contains the configuration of the backup adapter
	// +optional
	Parameters map[string]string `json:"parameters,omitempty"`
}

// GetAdapterConfiguration returns the adapter configuration if stored into the Pod
func GetAdapterConfiguration(pod *corev1.Pod) (*AdapterConfiguration, error) {
	configurationString, ok := pod.Annotations[BackupAdapterAnnotationName]
	if !ok {
		return nil, nil
	}

	var result AdapterConfiguration
	if err := json.Unmarshal([]byte(configurationString), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
