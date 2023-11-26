package injector

import corev1 "k8s.io/api/core/v1"

func (injector *Data) getSidecarContainer() corev1.Container {
	return corev1.Container{
		Name: "filepath-adapter",
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "scratch-data",
				MountPath: "/controller",
			},
			{
				Name:      "backups",
				MountPath: "/backup",
			},
			{
				Name:      "pgdata",
				MountPath: "/var/lib/postgresql/data",
			},
		},
		Image:           injector.image,
		ImagePullPolicy: corev1.PullPolicy(injector.imagePullPolicy),
		Command: []string{
			"/app/bin/filepath_adapter",
			"server",
			"--listening-network",
			"unix",
			"--listening-address",
			"/controller/walmanager",
			"--base-path",
			"/backup",
		},
	}
}

func (injector *Data) getBackupVolume() corev1.Volume {
	return corev1.Volume{
		Name: "backups",
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: "backups-pvc",
			},
		},
	}
}
