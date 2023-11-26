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

const PluginId = "filepath.leonardoce.io"

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
	log.SetLogger(logger)

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

	logger = logger.WithValues("podName", pod.Name, "podNamespace", pod.Namespace)

	configuration, err := GetAdapterConfiguration(pod)
	if err != nil {
		return err
	}
	if configuration == nil {
		// This is not a CNPG Pod
		return nil
	}

	if configuration.ID != PluginId {
		// This s not the correct injector
		logger.Info(
			"This not is not applicable to this injector",
			"configuration", configuration,
			"pluginID", PluginId)
		return nil
	}

	// Inject sidecar
	if len(pod.Spec.Containers) > 0 {
		pod.Spec.Containers = append(
			pod.Spec.Containers,
			injector.getSidecarContainer())
	}

	// Inject backup volume
	if len(pod.Spec.Volumes) > 0 {
		pod.Spec.Volumes = append(
			pod.Spec.Volumes,
			injector.getBackupVolume(configuration))
	}

	// Inject annotations
	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	pod.Annotations["filepath-adapter.leonardoce.io/injected"] = "true"
	logger.Info("Sidecar injected")

	return nil
}
