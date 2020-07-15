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
	/*
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
	*/
	ir, err := image.NewImageRuntimeFromOptions(storeOptions)
	if err != nil {
		return err
	}
	images, err := ir.GetImages()
	if err != nil {
		return err
	}
	for _, img := range images {
		println(img.InputName)
		println(img.ID())
		for name := range img.Names() {
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
		println(img.ID())
		for name := range img.Names() {
			println(name)
		}
	}
	println()
	return retErr
}
