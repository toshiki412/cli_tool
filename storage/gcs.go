package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
)

func Upload(target string, filename string, conf cfg.UploadGoogleStorageType) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	cobra.CheckErr(err)
	defer client.Close()

	f, err := os.Open(target)
	cobra.CheckErr(err)
	defer f.Close()

	var uploadpath = ""
	if conf.Dir == "" {
		uploadpath = filename
	} else {
		uploadpath = filepath.Join(conf.Dir, filename)
	}

	o := client.Bucket(conf.Bucket).Object(uploadpath)

	wc := o.NewWriter(ctx)
	_, err = io.Copy(wc, f)
	cobra.CheckErr(err)
	err = wc.Close()
	cobra.CheckErr(err)
}

func Download(target string, conf cfg.UploadGoogleStorageType) string {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	cobra.CheckErr(err)
	defer client.Close()

	var filePath = ""
	if conf.Dir == "" {
		filePath = target
	} else {
		filePath = filepath.Join(conf.Dir, target)
	}

	o := client.Bucket(conf.Bucket).Object(filePath)

	tmpDir, err := os.MkdirTemp("", ".cli_tool")
	cobra.CheckErr(err)

	tmpfile := filepath.Join(tmpDir, target)
	f, err := os.OpenFile(tmpfile, os.O_CREATE|os.O_WRONLY, 0644)
	cobra.CheckErr(err)
	defer f.Close()

	rc, err := o.NewReader(ctx)
	cobra.CheckErr(err)
	_, err = io.Copy(f, rc)
	cobra.CheckErr(err)

	err = rc.Close()
	cobra.CheckErr(err)

	return tmpfile
}

func IsExist(target string, conf cfg.UploadGoogleStorageType) bool {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	cobra.CheckErr(err)
	defer client.Close()

	var filePath = ""
	if conf.Dir == "" {
		filePath = target
	} else {
		filePath = filepath.Join(conf.Dir, target)
	}

	o := client.Bucket(conf.Bucket).Object(filePath)

	attrs, _ := o.Attrs(ctx)

	return attrs != nil
}
