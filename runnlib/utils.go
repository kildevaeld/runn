package runnlib

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/kildevaeld/go-filecrypt"
)

func getInterpreterFromPath(ext string) []string {
	switch filepath.Ext(ext) {
	case ".sh":
		return []string{"sh", "-c"}
	case ".js":
		return []string{"node"}
	case ".py":
		return []string{"python"}
	default:
		return nil
	}
}

func UnarchiveToDir(path string, source io.Reader, size int64, key_raw []byte) error {

	key := filecrypt.Key(key_raw)

	out := bytes.NewBuffer(nil)

	if err := filecrypt.Decrypt(out, source, &key); err != nil {
		return fmt.Errorf("decrypt: %s", err)
	}

	arc, err := zip.NewReader(bytes.NewReader(out.Bytes()), size)
	if err != nil {
		return err
	}

	os.MkdirAll(path, 0766)

	for _, file := range arc.File {
		fullPath := filepath.Join(path, file.Name)
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(fullPath, file.Mode()); err != nil {
				return err
			}
			continue
		}

		fileReader, rerr := file.Open()
		if rerr != nil {
			return rerr
		}
		defer fileReader.Close()

		fileWriter, werr := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if werr != nil {
			return werr
		}
		defer fileWriter.Close()

		if _, err := io.Copy(fileWriter, fileReader); err != nil {
			return err
		}

	}

	return nil
}

func ArchieveDir(path string, name string, key_raw []byte) (io.Reader, error) {

	buf := bytes.NewBuffer(nil)
	arc := zip.NewWriter(buf)

	path, _ = filepath.Abs(path)

	//baseDir := filepath.Base(path)
	includeCurrentFolder := false

	var files []string

	err := addAll(path, path, includeCurrentFolder, func(info os.FileInfo, file io.Reader, entryName string) (err error) {

		// Create a header based off of the fileinfo
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// If it's a file, set the compression method to deflate (leave directories uncompressed)
		if !info.IsDir() {
			header.Method = zip.Deflate
		}

		// Set the header's name to what we want--it may not include the top folder
		header.Name = filepath.Join(strings.TrimPrefix(entryName, path))

		// Add a trailing slash if the entry is a directory
		if info.IsDir() {
			header.Name += "/"
		}

		// Get a writer in the archive based on our header
		if i := getInterpreterFromPath(header.Name); i != nil {
			header.SetMode(0777)
		}
		writer, err := arc.CreateHeader(header)
		if err != nil {
			return err
		}

		// If we have a file to write (i.e., not a directory) then pipe the file into the archive writer
		if file != nil {
			files = append(files, entryName)
			if _, err := io.Copy(writer, file); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	var bundle Bundle
	if err := GetBundleFromPath(path, &bundle); err != nil {
		if err != NotExistsError {
			return nil, err
		}
		bundle.Name = name

		for _, file := range files {
			interpreter := getInterpreterFromPath(file)
			if interpreter == nil {
				continue
			}

			base := filepath.Base(file)
			name := strings.TrimSuffix(base, filepath.Ext(base))
			filePath := strings.TrimPrefix(file, path)
			command := BundleCommand{
				Name: name,
				Command: CommandConfig{
					Cmd:         fmt.Sprintf("{{.WorkDir}}/%s", filePath),
					WorkDir:     fmt.Sprintf("{{.WorkDir}}"),
					Stderr:      "stderr",
					Stdout:      "stdout",
					Interpreter: interpreter,
				},
			}

			bundle.Commands = append(bundle.Commands, command)

		}

		bundleFile, _ := arc.Create("bundle.yaml")

		if b, e := yaml.Marshal(bundle); e != nil {
			return nil, e
		} else {

			if _, e := bundleFile.Write(b); e != nil {
				return nil, e
			}
			//bundleFile.Write(nil)
		}

	}

	if err := arc.Close(); err != nil {
		return nil, fmt.Errorf("zip-close: %s", err)
	}

	out := bytes.NewBuffer(nil)

	key := filecrypt.Key(key_raw)

	if _, err := filecrypt.Encrypt(out, bytes.NewReader(buf.Bytes()), &key); err != nil {
		return nil, fmt.Errorf("encrypt: %s", err)
	}

	return out, nil

}

type ArchiveWriteFunc func(info os.FileInfo, file io.Reader, entryName string) (err error)

func getSubDir(dir string, rootDir string, includeCurrentFolder bool) (subDir string) {

	subDir = strings.Replace(dir, rootDir, "", 1)
	// Remove leading slashes, since this is intentionally a subdirectory.
	if len(subDir) > 0 && subDir[0] == os.PathSeparator {
		subDir = subDir[1:]
	}
	subDir = path.Join(strings.Split(subDir, string(os.PathSeparator))...)

	if includeCurrentFolder {
		parts := strings.Split(rootDir, string(os.PathSeparator))
		subDir = path.Join(parts[len(parts)-1], subDir)
	}

	return
}

// addAll is used to recursively go down through directories and add each file and directory to an archive, based on an ArchiveWriteFunc given to it
func addAll(dir string, rootDir string, includeCurrentFolder bool, writerFunc ArchiveWriteFunc) error {

	// Get a list of all entries in the directory, as []os.FileInfo
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	// Loop through all entries
	for _, info := range fileInfos {

		full := filepath.Join(dir, info.Name())

		// If the entry is a file, get an io.Reader for it
		var file *os.File
		var reader io.Reader
		if !info.IsDir() {
			file, err = os.Open(full)
			if err != nil {
				return err
			}
			reader = file
		}

		// Write the entry into the archive
		subDir := getSubDir(dir, rootDir, includeCurrentFolder)
		entryName := path.Join(subDir, info.Name())
		if err := writerFunc(info, reader, entryName); err != nil {
			if file != nil {
				file.Close()
			}
			return err
		}

		if file != nil {
			if err := file.Close(); err != nil {
				return err
			}

		}

		// If the entry is a directory, recurse into it
		if info.IsDir() {
			addAll(full, rootDir, includeCurrentFolder, writerFunc)
		}
	}

	return nil
}
