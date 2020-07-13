package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/containers/image/image"
	is "github.com/containers/image/storage"
	"github.com/containers/image/transports/alltransports"
	"github.com/containers/image/types"
	"github.com/containers/storage"
	pkgerrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	var debug bool
	debug = true
	println("test")
	logrus.SetLevel(logrus.ErrorLevel)
	if !debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if err := imageLookup(); err != nil {
		logrus.Error(err)
	}
}

// test product type/version aggregation from pod image lookup
func imageLookup() (retErr error) {
	images := []string{"quay.io/crio/redis:alpine", "quay.io/crio/redis@sha256:1780b5a5496189974b94eb2595d86731d7a0820e4beb8ea770974298a943ed55"}
	ctx := context.Background()
	var cancel context.CancelFunc = func() {}
	ctx, cancel = context.WithTimeout(ctx, time.Duration(10)*time.Second)
	defer cancel()
	storeOptions, err := storage.DefaultStoreOptions(false, 0)
	if err != nil {
		return err
	}
	storeOptions.GraphDriverName = "overlay"
	println(storeOptions.GraphDriverName)

	store, err := storage.GetStore(storeOptions)
	if err != nil {
		return err
	}
	errors := false
	defer func() {
		store.Shutdown(false)
		if errors {
			os.Exit(1)
		}
	}()

	imgs, err := store.Images()
	if err != nil {
		return err
	}
	for _, i := range imgs {
		println("digest =" + i.Digest.String())
	}
	ref, err := is.Transport.ParseStoreReference(store, "test")
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
	imgCtx := &types.SystemContext{
		OSChoice: "linux",
	}
	imagesCount := uniqueCount(images)
	for img, num := range imagesCount {
		fmt.Println(img)
		fmt.Printf("image is used %d time(s)\n", num)

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
		fmt.Println()
		fmt.Println("IMAGE LABELS -")
		for key, val := range config.Config.Labels {
			if key == "org.jboss.product" {
				fmt.Println(key + "=" + val)
			}
		}
	}
	fmt.Println()
	return retErr
}

func parseImageSource(ctx context.Context, imgCtx *types.SystemContext, name string) (types.ImageSource, error) {
	ref, err := alltransports.ParseImageName(name)
	if err != nil {
		return nil, err
	}
	return ref.NewImageSource(ctx, imgCtx)
}

func uniqueCount(intSlice []string) map[string]int {
	keys := make(map[string]bool)
	list := map[string]int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list[entry] = 1
		} else {
			list[entry] = list[entry] + 1
		}
	}
	return list
}
