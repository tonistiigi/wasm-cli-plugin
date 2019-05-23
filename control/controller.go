package control

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/content/local"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/metadata"
	"github.com/containerd/containerd/platforms"
	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	"github.com/containerd/containerd/snapshots"
	"github.com/containerd/containerd/snapshots/native"
	"github.com/docker/docker/errdefs"
	"github.com/opencontainers/image-spec/identity"
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/sync/semaphore"
)

type Opt struct {
	Root string
}

type Controller struct {
	db   *bbolt.DB
	cs   content.Store
	is   images.Store
	mdb  *metadata.DB
	root string
}

func New(opt Opt) (*Controller, error) {
	root := opt.Root

	if err := os.MkdirAll(filepath.Join(root, "containers"), 0700); err != nil {
		return nil, err
	}

	db, err := bolt.Open(filepath.Join(root, "metadata.db"), 0644, nil)
	if err != nil {
		return nil, err
	}

	c, err := local.NewStore(filepath.Join(root, "content"))
	if err != nil {
		return nil, err
	}

	ctx := addNS(context.Background())

	snapshotter, err := native.NewSnapshotter(filepath.Join(root, "snapshots"))
	if err != nil {
		return nil, err
	}

	mdb := metadata.NewDB(db, c, map[string]snapshots.Snapshotter{
		"native": snapshotter,
	})
	if err := mdb.Init(ctx); err != nil {
		return nil, err
	}

	return &Controller{
		db:   db,
		cs:   mdb.ContentStore(),
		is:   metadata.NewImageStore(mdb),
		mdb:  mdb,
		root: root,
	}, nil
}

func (c *Controller) Close() error {
	_, err := c.mdb.GarbageCollect(context.TODO())
	if err != nil {
		return err
	}
	return c.db.Close()
}

func (c *Controller) Delete(ctx context.Context, ref string) error {
	ref, err := parseRef(ref)
	if err != nil {
		return err
	}
	ctx = addNS(ctx)
	return c.is.Delete(ctx, ref)
}

func (c *Controller) Images(ctx context.Context) ([]images.Image, error) {
	ctx = addNS(ctx)
	return c.is.List(ctx)
}

func (c *Controller) GetRuntimeImage(ctx context.Context, ref string, platform platforms.MatchComparer) (*images.Image, error) {
	ref, err := parseRef(ref)
	if err != nil {
		return nil, err
	}

	ctx = addNS(ctx)

	img, err := c.is.Get(ctx, ref)
	if err != nil {
		return nil, nil
	}

	chain, err := img.RootFS(ctx, c.cs, platform)
	if err != nil {
		return nil, nil
	}

	chainID := identity.ChainID(chain)

	sn := c.mdb.Snapshotter("native")

	if _, err := sn.Stat(ctx, chainID.String()); err == nil {
		return &img, nil
	}

	return nil, nil
}

func (c *Controller) Pull(ctx context.Context, ref string, platform platforms.MatchComparer) (*images.Image, error) {
	ref, err := parseRef(ref)
	if err != nil {
		return nil, err
	}

	ctx = addNS(ctx)

	var limiter *semaphore.Weighted

	r := docker.NewResolver(docker.ResolverOptions{}) // credentials

	name, desc, err := r.Resolve(ctx, ref)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve reference %q", ref)
	}

	fetcher, err := r.Fetcher(ctx, name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get fetcher for %q", name)
	}

	handler := images.ChildrenHandler(c.cs)
	handler = images.SetChildrenLabels(c.cs, handler)
	handler = remotes.FilterManifestByPlatformHandler(handler, platform)
	handler = images.LimitManifests(handler, platform, 1)

	handlers := images.Handlers(
		remotes.FetchHandler(c.cs, fetcher),
		handler,
	)

	if err := images.Dispatch(ctx, handlers, limiter, desc); err != nil {
		return nil, err
	}

	if err := unpack(ctx, desc, c.cs, c.mdb.Snapshotter("native"), platform); err != nil {
		return nil, err
	}

	img := images.Image{
		Name:   name,
		Target: desc,
	}

	for {
		if created, err := c.is.Create(ctx, img); err != nil {
			if !strings.Contains(errors.Cause(err).Error(), "already exists") {
				return nil, err
			}

			updated, err := c.is.Update(ctx, img)
			if err != nil {
				// if image was removed, try create again
				if errdefs.IsNotFound(err) {
					continue
				}
				return nil, err
			}

			img = updated
		} else {
			img = created
		}

		return &img, nil
	}
}
