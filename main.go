package main

import (
	"os"

	"github.com/containers/libpod/v2/libpod/image"
	"github.com/containers/storage"
	"github.com/sirupsen/logrus"
)

// var images = []string{"quay.io/crio/redis@sha256:1780b5a5496189974b94eb2595d86731d7a0820e4beb8ea770974298a943ed55", "redis@sha256:1780b5a5496189974b94eb2595d86731d7a0820e4beb8ea770974298a943ed55"}

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
	ir := image.NewImageRuntimeFromStore(store)
	imgData, err := ir.NewFromLocal("redis@sha256:1780b5a5496189974b94eb2595d86731d7a0820e4beb8ea770974298a943ed55")
	if err != nil {
		return err
	}
	println(imgData.ID)
	println(imgData.Config.Env)

	images, err := ir.GetImages()
	if err != nil {
		return err
	}
	for _, img := range images {
		println(img.ID)
		println(img.Config.Env)
	}

	println()
	return retErr
}
