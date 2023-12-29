package injector

import (
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"

	"github.com/cloudnative-pg/volume-injector/pkg/apis"
)

// getAdditionalVolumesConfiguration returns the additional volumes configuration if stored into the Pod
func getAdditionalVolumesConfiguration(pod *corev1.Pod) ([]apis.AdditionalVolume, error) {
	configurationString, ok := pod.Annotations[apis.AdditionalVolumesAnnotationName]
	if !ok {
		return nil, nil
	}

	var result []apis.AdditionalVolume
	if err := yaml.Unmarshal([]byte(configurationString), &result); err != nil {
		return nil, err
	}

	return result, nil
}

// createKubernetesVolumes creates the kubernetes volumes for the defined
// additional volumes
func createKubernetesVolumes(configuration []apis.AdditionalVolume) (result []corev1.Volume) {
	result = make([]corev1.Volume, len(configuration))
	for i, volumeConfiguration := range configuration {
		result[i] = corev1.Volume{
			Name: volumeConfiguration.ClaimName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: volumeConfiguration.ClaimName,
					ReadOnly:  volumeConfiguration.ReadOnly,
				},
			},
		}
	}

	return result
}

// createKubernetesVolumes creates the kubernetes volume mount points for the
// defined additional volumes
func createKubernetesVolumeMounts(configuration []apis.AdditionalVolume) (result []corev1.VolumeMount) {
	result = make([]corev1.VolumeMount, len(configuration))
	for i, volumeConfiguration := range configuration {
		result[i] = corev1.VolumeMount{
			Name:      volumeConfiguration.ClaimName,
			MountPath: volumeConfiguration.MountPath,
		}
	}

	return result
}
