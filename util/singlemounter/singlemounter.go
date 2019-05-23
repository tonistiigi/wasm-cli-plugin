package singlemounter

import (
	"github.com/containerd/containerd/mount"
	"github.com/pkg/errors"
)

type Mounter interface {
	Mount() (string, error)
	Unmount() error
}

// SingleMounter is a helper for mounting mountfactory that is a bind mount to
// a single directory
func SingleMounter(mounts []mount.Mount) Mounter {
	return &singleMounter{mounts: mounts}
}

type singleMounter struct {
	mounts []mount.Mount
}

func (sm *singleMounter) Mount() (string, error) {
	if len(sm.mounts) == 1 && (sm.mounts[0].Type == "bind" || sm.mounts[0].Type == "rbind") {
		ro := false
		for _, opt := range sm.mounts[0].Options {
			if opt == "ro" {
				ro = true
				break
			}
		}
		if !ro {
			return sm.mounts[0].Source, nil
		}
	}

	return "", errors.Errorf("mount not supported for singlemounter %v", sm)
}

func (sm *singleMounter) Unmount() error {
	return nil
}
