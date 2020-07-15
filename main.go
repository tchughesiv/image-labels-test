package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/containers/libpod/v2/libpod/image"
	"github.com/containers/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	debug := true
	args := os.Args[1:]
	logrus.SetLevel(logrus.ErrorLevel)
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if err := imageLookup(args); err != nil {
		logrus.Error(err)
	}
}

// test product type/version aggregation from pod image lookup
func imageLookup(args []string) (retErr error) {
	// ERRO[0000] error opening "/var/lib/containers/storage/storage.lock": permission denied
	// https://github.com/containers/storage/blob/ed28f2457e2f57cb3d3f2f4029a85f72b35370ab/store.go#L656-L664
	storeOptions, err := storage.DefaultStoreOptionsAutoDetectUID()
	if err != nil {
		return err
	}
	println("graphroot is " + storeOptions.GraphRoot)
	ir, err := image.NewImageRuntimeFromOptions(storeOptions)
	if err != nil {
		return err
	}
	images, err := ir.GetImages()
	if err != nil {
		return err
	}
	for _, img := range images {
		println()
		println(img.InputName)
		for _, name := range img.Names() {
			println(name)
		}
	}
	for _, imgName := range args {
		img, err := ir.NewFromLocal(imgName)
		if err != nil {
			return err
		}
		println()
		println(img.InputName)
		for _, name := range img.Names() {
			println(name)
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)
		defer cancel()
		// get inspect image data
		imgData, err := img.Inspect(ctx)
		if err != nil {
			return err
		}
		println(imgData.ID)
		if imgData.Labels != nil {
			println("IMAGE LABELS:")
			for key, val := range imgData.Labels {
				if strings.Contains(key, "org.jboss.") {
					println(key + "=" + val)
				}
			}
			for key, val := range imgData.Labels {
				if strings.Contains(key, "com.redhat.") {
					println(key + "=" + val)
				}
			}
		}
	}
	println()
	return retErr
}
