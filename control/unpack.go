package control

import (
	"archive/tar"
	"context"
	"io"
	"io/ioutil"
	"os"

	"github.com/containerd/containerd/archive"
	"github.com/containerd/containerd/archive/compression"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/platforms"
	"github.com/containerd/containerd/rootfs"
	"github.com/containerd/containerd/snapshots"
	bkidentity "github.com/moby/buildkit/identity"
	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/identity"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/tonistiigi/wasm-cli-plugin/util/singlemounter"
)

func unpackLayer(ctx context.Context, desc ocispec.Descriptor, cs content.Store, dest string) (*ocispec.Descriptor, error) {
	ra, err := cs.ReaderAt(ctx, desc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get reader from content store")
	}
	defer ra.Close()

	isCompressed, err := images.IsCompressedDiff(ctx, desc.MediaType)
	if err != nil {
		return nil, errors.Wrapf(err, "unsupported diff media type: %v", desc.MediaType)
	}

	r := content.NewReader(ra)
	if isCompressed {
		ds, err := compression.DecompressStream(r)
		if err != nil {
			return nil, err
		}
		defer ds.Close()
		r = ds
	}

	digester := digest.Canonical.Digester()
	rc := &readCounter{
		r: io.TeeReader(r, digester.Hash()),
	}

	uid := os.Getuid()
	gid := os.Getgid()

	if _, err := archive.Apply(ctx, dest, rc, archive.WithFilter(func(h *tar.Header) (bool, error) {
		h.Uid = uid
		h.Gid = gid
		return true, nil
	})); err != nil {
		return nil, err
	}

	// Read any trailing data
	if _, err := io.Copy(ioutil.Discard, rc); err != nil {
		return nil, err
	}

	ocidesc := ocispec.Descriptor{
		MediaType: ocispec.MediaTypeImageLayer,
		Size:      rc.c,
		Digest:    digester.Digest(),
	}

	return &ocidesc, nil
}

func unpack(ctx context.Context, desc ocispec.Descriptor, cs content.Store, sn snapshots.Snapshotter, platform platforms.MatchComparer) error {
	layers, err := getLayers(ctx, cs, desc, platform)
	if err != nil {
		return err
	}

	var chain []digest.Digest
	for _, layer := range layers {
		chain = append(chain, layer.Diff.Digest)
	}
	chainID := identity.ChainID(chain)

	if _, err := sn.Stat(ctx, chainID.String()); err == nil {
		return nil
	}

	key := bkidentity.NewID()
	mounts, err := sn.Prepare(ctx, key, "")
	if err != nil {
		return err
	}

	var noClean bool
	defer func() {
		if !noClean {
			sn.Remove(context.TODO(), key)
		}
	}()

	lm := singlemounter.SingleMounter(mounts)
	mp, err := lm.Mount()
	if err != nil {
		return err
	}
	defer lm.Unmount()

	var chain2 []digest.Digest
	for _, l := range layers {
		desc2, err := unpackLayer(ctx, l.Blob, cs, mp)
		if err != nil {
			return err
		}
		chain2 = append(chain2, desc2.Digest)
	}
	chainID2 := identity.ChainID(chain)

	if chainID != chainID2 {
		return errors.Errorf("chainid verification mismatch %s %s", chainID, chainID2)
	}

	if err := sn.Commit(ctx, chainID.String(), key); err != nil {
		return errors.Wrapf(err, "failed to commit snapshot %s", key)
	}

	cinfo := content.Info{
		Digest: desc.Digest,
		Labels: map[string]string{
			"containerd.io/gc.ref.snapshot.native": chainID.String(),
		},
	}

	if _, err := cs.Update(ctx, cinfo, "labels.containerd.io/gc.ref.snapshot.native"); err != nil {
		return err
	}

	return nil
}

type readCounter struct {
	r io.Reader
	c int64
}

func (rc *readCounter) Read(p []byte) (n int, err error) {
	n, err = rc.r.Read(p)
	rc.c += int64(n)
	return
}

func getLayers(ctx context.Context, provider content.Provider, desc ocispec.Descriptor, platform platforms.MatchComparer) ([]rootfs.Layer, error) {
	manifest, err := images.Manifest(ctx, provider, desc, platform)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	image := images.Image{Target: desc}
	diffIDs, err := image.RootFS(ctx, provider, platform)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve rootfs")
	}
	if len(diffIDs) != len(manifest.Layers) {
		return nil, errors.Errorf("mismatched image rootfs and manifest layers %+v %+v", diffIDs, manifest.Layers)
	}
	layers := make([]rootfs.Layer, len(diffIDs))
	for i := range diffIDs {
		layers[i].Diff = ocispec.Descriptor{
			// TODO: derive media type from compressed type
			MediaType: ocispec.MediaTypeImageLayer,
			Digest:    diffIDs[i],
		}
		layers[i].Blob = manifest.Layers[i]
	}
	return layers, nil
}
