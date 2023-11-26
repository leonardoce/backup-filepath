package walmanager

import (
	"context"
	"path"

	"github.com/leonardoce/backup-adapter/pkg/adapter"
	"github.com/leonardoce/backup-filepath/internal/cp"
	"github.com/leonardoce/backup-filepath/internal/logging"
)

const (
	walsDirectory = "wals"
)

type WalManagerImplementation struct {
	adapter.WalManagerServer
	basePath string
}

func NewWalManagerImplementation(basePath string) *WalManagerImplementation {
	return &WalManagerImplementation{
		basePath: basePath,
	}
}

func (impl *WalManagerImplementation) getWalPrefix(walName string) string {
	return walName[0:8]
}

func (impl *WalManagerImplementation) getWalPath(clusterName string, walName string) string {
	return path.Join(impl.basePath, clusterName, walsDirectory, impl.getWalPrefix(walName), walName)
}

func (impl *WalManagerImplementation) ArchiveWal(ctx context.Context, request *adapter.ArchiveWalRequest) (*adapter.ArchiveWalResult, error) {
	walName := path.Base(request.SourceFileName)
	destinationPath := impl.getWalPath(request.ClusterName, walName)
	logging := logging.FromContext(ctx)

	err := cp.CopyFile(request.SourceFileName, destinationPath)
	logging.Info("Archived WAL File", "sourceFileName", request.SourceFileName, "destinationPath", destinationPath, "clusterName", request.ClusterName, "err", err)
	return &adapter.ArchiveWalResult{}, err
}

func (impl *WalManagerImplementation) RestoreWal(ctx context.Context, request *adapter.RestoreWalRequest) (*adapter.RestoreWalResult, error) {
	walPath := impl.getWalPath(request.ClusterName, request.SourceWalName)
	logging := logging.FromContext(ctx)

	err := cp.CopyFile(walPath, request.DestinationFileName)
	logging.Info("Restored WAL File", "walName", request.SourceWalName, "walPath", walPath, "destinationPath", request.DestinationFileName, "err", err)
	return &adapter.RestoreWalResult{}, err
}
