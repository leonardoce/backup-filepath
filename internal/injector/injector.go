package injector

import (
	"context"
	"fmt"

	"github.com/leonardoce/backup-filepath/internal/logging"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

type Data struct {
	webhookPort     int
	image           string
	imagePullPolicy string
}

func New() *Data {
	return &Data{
		webhookPort:     443,
		image:           "filepath_adapter:latest",
		imagePullPolicy: "Never",
	}
}

func (injector *Data) Run(ctx context.Context) error {
	logger := logging.FromContext(ctx)
	log.SetLogger(*logger)

	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
		WebhookServer: webhook.NewServer(webhook.Options{
			Port: injector.webhookPort,
		}),
	})
	if err != nil {
		return err
	}

	err = builder.WebhookManagedBy(mgr).
		For(&corev1.Pod{}).
		WithDefaulter(injector).
		Complete()
	if err != nil {
		return err
	}

	logger.Info("starting manager", "injector", injector)
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		logger.Error(err, "unable to run manager")
		return err
	}

	return nil
}

func (injector *Data) Default(ctx context.Context, obj runtime.Object) error {
	logger := logging.FromContext(ctx)

	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("expected a Pod but got a %T", obj)
	}

	// Inject sidecar
	sidecarContainer := corev1.Container{
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
	if len(pod.Spec.Containers) > 0 {
		pod.Spec.Containers = append(pod.Spec.Containers, sidecarContainer)
	}

	// Inject backup volume
	backupVolume := corev1.Volume{
		Name: "backups",
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: "backups-pvc",
			},
		},
	}
	if len(pod.Spec.Volumes) > 0 {
		pod.Spec.Volumes = append(pod.Spec.Volumes, backupVolume)
	}

	// Inject annotations
	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	pod.Annotations["filepath-adapter.leonardoce.io"] = "injected"
	logger.Info("Injected sidecar into Pod", "podName", pod.Name, "namespace", pod.Namespace)

	return nil
}
