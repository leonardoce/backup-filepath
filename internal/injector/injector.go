package injector

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/cloudnative-pg/volume-injector/internal/logging"
	"github.com/cloudnative-pg/volume-injector/pkg/apis"
)

// Data is the configuration of the webhook server
type Data struct {
	webhookPort int
}

// New creates a new webhook runner
func New() *Data {
	return &Data{
		webhookPort: 443,
	}
}

// Run starts the webhooks web server
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

// Default implements the mutating webhook
func (injector *Data) Default(ctx context.Context, obj runtime.Object) error {
	logger := logging.FromContext(ctx)

	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("expected a Pod but got a %T", obj)
	}

	logger = logger.WithValues("podName", pod.Name, "podNamespace", pod.Namespace)

	additionalVolumesConfiguration, err := getAdditionalVolumesConfiguration(pod)
	if err != nil {
		return err
	}
	if additionalVolumesConfiguration == nil {
		// Nothing to see here, this pod is not interesting to us
		return nil
	}

	// Let's validate the annotation
	for i := range additionalVolumesConfiguration {
		if err := additionalVolumesConfiguration[i].Validate(); err != nil {
			logger.Info("Invalid volume configuration", "targetVolume", additionalVolumesConfiguration[i], "err", err)
			return err
		}
	}

	// Inject additional volumes
	pod.Spec.Volumes = append(
		pod.Spec.Volumes,
		createKubernetesVolumes(additionalVolumesConfiguration)...)

	// Inject additional volume mounts in every container
	for i := range pod.Spec.Containers {
		pod.Spec.Containers[i].VolumeMounts = append(
			pod.Spec.Containers[i].VolumeMounts,
			createKubernetesVolumeMounts(additionalVolumesConfiguration)...)
	}

	// Inject annotation
	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	pod.Annotations[apis.AdditionalVolumesInjectedAnnotationName] = "true"
	logger.Info("Volume injected", "configuration", additionalVolumesConfiguration)

	return nil
}
