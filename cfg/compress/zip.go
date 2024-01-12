package compress

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// 圧縮   targetを圧縮して、圧縮したファイルのパスを返す
func Compress(target string) string {
	var buffer bytes.Buffer
	zipWriter := zip.NewWriter(&buffer)

	isDir := isDirectory(target)

	if isDir {
		err := addZipFiles(zipWriter, target, "")
		cobra.CheckErr(err)
	} else {
		fileName := filepath.Base(target)
		addZipFile(zipWriter, target, fileName)
	}
	err := zipWriter.Close()
	cobra.CheckErr(err)

	dumpDir, err := os.MkdirTemp("", ".cli_tool")
	cobra.CheckErr(err)

	zipfile := filepath.Join(dumpDir, "test.zip")

	file, err := os.Create(zipfile)
	cobra.CheckErr(err)
	defer file.Close()

	_, err = file.Write(buffer.Bytes())
	cobra.CheckErr(err)

	return zipfile
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	cobra.CheckErr(err)
	return fileInfo.IsDir()
}

func addZipFiles(writer *zip.Writer, basePath, pathInZip string) error {
	fileInfoArray, err := os.ReadDir(basePath)
	cobra.CheckErr(err)

	basePath = complementPath(basePath)
	pathInZip = complementPath(pathInZip)

	for _, fileInfo := range fileInfoArray {
		newBasePath := basePath + fileInfo.Name()
		newPathInZip := pathInZip + fileInfo.Name()

		if fileInfo.IsDir() {
			addDirectory(writer, newBasePath)

			newBasePath = newBasePath + string(os.PathSeparator)
			newPathInZip = newPathInZip + string(os.PathSeparator)

			err = addZipFiles(writer, newBasePath, newPathInZip)
			cobra.CheckErr(err)
		} else {
			addZipFile(writer, newBasePath, newPathInZip)
		}
	}

	return nil
}

func addZipFile(writer *zip.Writer, targetFilePath, pathInZip string) {
	data, err := os.ReadFile(targetFilePath)
	cobra.CheckErr(err)

	fileInfo, err := os.Lstat(targetFilePath)
	cobra.CheckErr(err)

	header, err := zip.FileInfoHeader(fileInfo)
	cobra.CheckErr(err)

	header.Name = pathInZip
	header.Method = zip.Deflate

	w, err := writer.CreateHeader(header) // zipファイルに書き込む（圧縮する）
	cobra.CheckErr(err)

	_, err = w.Write(data)
	cobra.CheckErr(err)
}

func addDirectory(writer *zip.Writer, basePath string) {
	fileInfo, err := os.Lstat(basePath)
	cobra.CheckErr(err)

	header, err := zip.FileInfoHeader(fileInfo)
	cobra.CheckErr(err)

	_, err = writer.CreateHeader(header)
	cobra.CheckErr(err)
}

func complementPath(path string) string {
	l := len(path)
	if l == 0 {
		return path
	}

	lastChar := path[l-1 : l]
	if lastChar == "/" || lastChar == "\\" {
		return path
	} else {
		return path + string(os.PathSeparator)
	}
}

// 解凍
func Decompress(dest, target string) error {
	reader, err := zip.OpenReader(target) // readerがzipファイルの中身
	cobra.CheckErr(err)
	defer reader.Close()

	for _, zippedFile := range reader.File {
		path := filepath.Join(dest, zippedFile.Name)
		if zippedFile.FileInfo().IsDir() {
			err = os.MkdirAll(path, zippedFile.Mode())
			cobra.CheckErr(err)
		} else {
			createFileFromZipped(path, zippedFile)
		}
	}

	return nil
}

func createFileFromZipped(path string, zippedFile *zip.File) {
	reader, err := zippedFile.Open()
	cobra.CheckErr(err)
	defer reader.Close()

	destFile, err := os.Create(path)
	cobra.CheckErr(err)
	defer destFile.Close()

	_, err = io.Copy(destFile, reader)
	cobra.CheckErr(err)
}
