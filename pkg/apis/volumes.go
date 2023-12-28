package apis

import "errors"

var (
	// ErrEmptyClaimName is raised when the claim name is empty
	ErrEmptyClaimName = errors.New("empty claim name")

	// ErrEmptyMountpath is raised when the mount path is empty
	ErrEmptyMountpath = errors.New("empty mount path")
)

// AdditionalVolume represent a new volume to be mounted into CNPG pods
type AdditionalVolume struct {
	// ClaimName is the name of the PVC to mount
	ClaimName string `yaml:"claimName"`

	// MountPath is the path where the volume should be mounted
	MountPath string `yaml:"mountPath"`

	// ReadOnly is true if the volume should be marked read only
	ReadOnly bool `yaml:"readOnly"`
}

// Validate returns an error when the required additional volume
func (volume AdditionalVolume) Validate() error {
	if len(volume.ClaimName) == 0 {
		return ErrEmptyClaimName
	}

	if len(volume.MountPath) == 0 {
		return ErrEmptyMountpath
	}

	return nil
}
