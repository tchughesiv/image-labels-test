package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/containers/image/v5/image"
	is "github.com/containers/image/v5/storage"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"github.com/containers/storage"
	pkgerrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var images = []string{"quay.io/crio/redis@sha256:1780b5a5496189974b94eb2595d86731d7a0820e4beb8ea770974298a943ed55", "redis@sha256:1780b5a5496189974b94eb2595d86731d7a0820e4beb8ea770974298a943ed55"}

func main() {
	debug := true
	logrus.SetLevel(logrus.ErrorLevel)
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if err := imageLookup(); err != nil {
		logrus.Error(err)
	}
}

// test product type/version aggregation from pod image lookup
func imageLookup() (retErr error) {
	ctx := context.Background()
	//var cancel context.CancelFunc = func() {}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(10)*time.Second)
	defer cancel()
	storeOptions, err := storage.DefaultStoreOptionsAutoDetectUID()
	if err != nil {
		return err
	}
	store, err := storage.GetStore(storeOptions)
	if err != nil {
		return err
	}
	defer func() {
		if _, err := store.Shutdown(false); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	}()
	/*
		imgs, err := store.Images()
		if err != nil {
			return err
		}
		for _, i := range imgs {
			println("digest = " + i.Digest.String())
		}
	*/
	for _, img := range images {
		ref, err := is.Transport.ParseStoreReference(store, img)
		if err != nil {
			return err
		}
		strRef := ref.StringWithinTransport()
		imgRef, err := is.Transport.ParseReference(strRef)
		if err != nil {
			return err
		}
		if imgRef == nil {
			return err
		}
	}
	imgCtx := &types.SystemContext{
		OSChoice: "linux",
	}
	for _, img := range images {
		println(img)

		imgSrc, err := parseImageSource(ctx, imgCtx, "containers-storage:"+img)
		if err != nil {
			return err
		}
		defer func() {
			if err = imgSrc.Close(); err != nil {
				retErr = pkgerrors.Wrapf(retErr, fmt.Sprintf("(could not close image: %v) ", err))
			}
		}()

		img, err := image.FromUnparsedImage(ctx, imgCtx, image.UnparsedInstance(imgSrc, nil))
		if err != nil {
			return fmt.Errorf("Error parsing manifest for image: %v", err)
		}

		config, err := img.OCIConfig(ctx)
		if err != nil {
			return fmt.Errorf("Error reading OCI-formatted configuration data: %v", err)
		}
		println()
		println("IMAGE LABELS -")
		for key, val := range config.Config.Labels {
			if key == "org.jboss.product" {
				println(key + "=" + val)
			}
		}
	}
	println()
	return retErr
}

func parseImageSource(ctx context.Context, imgCtx *types.SystemContext, name string) (types.ImageSource, error) {
	ref, err := alltransports.ParseImageName(name)
	if err != nil {
		return nil, err
	}
	return ref.NewImageSource(ctx, imgCtx)
}
