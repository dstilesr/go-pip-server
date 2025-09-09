package main

import (
	"crypto/sha256"
	"fmt"
	"go-pip-server/repository"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// PrepareFormData Extracts and validates form data from a multipart form for inserting a new project version.
func (p *PipServer) PrepareFormData(f *multipart.Form) (*repository.ProjectVersionInsert, error) {
	// Get required fields: name, version, filetype, and file content
	name, ok := f.Value["name"]
	if !ok || len(name) == 0 || name[0] == "" {
		return nil, fmt.Errorf("missing required field: name")
	}
	delete(f.Value, "name")

	version, ok := f.Value["version"]
	if !ok || len(version) == 0 || version[0] == "" {
		return nil, fmt.Errorf("missing required field: version")
	}
	delete(f.Value, "version")

	fType, ok := f.Value["filetype"]
	if !ok || len(fType) == 0 {
		return nil, fmt.Errorf("missing required field: filetype")
	} else if fType[0] != "source" && fType[0] != "bdist_wheel" {
		return nil, fmt.Errorf("invalid filetype: %s", fType[0])
	}
	delete(f.Value, "filetype")

	fileData, ok := f.File["content"]
	if !ok || len(fileData) == 0 || fileData[0] == nil {
		return nil, fmt.Errorf("missing required file: content")
	}

	// Compute SHA256 checksum of the uploaded file if not given
	dg, ok := f.Value["digest"]
	dt, ok2 := f.Value["digest_type"]
	var digest, digestType string

	fh, err := fileData[0].Open()
	if err != nil {
		return nil, fmt.Errorf("error opening uploaded file: %w", err)
	}
	defer fh.Close()

	if !(ok && ok2) || len(dg) == 0 || dg[0] == "" || len(dt) == 0 || dt[0] == "" {
		h, err := computeSHA256(fh)
		if err != nil {
			return nil, fmt.Errorf("error computing SHA256: %w", err)
		}
		digest = fmt.Sprintf("%x", h)
		digestType = "sha256"
	} else {
		digest, digestType = f.Value["digest"][0], f.Value["digest_type"][0]
	}

	// Save file to disk
	fp := filepath.Join(p.DataPath, name[0])
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		err := os.MkdirAll(fp, 0755)
		if err != nil {
			return nil, fmt.Errorf("error creating project directory: %w", err)
		}
	}
	fp = filepath.Join(fp, fileData[0].Filename)
	out, err := os.Create(fp)
	if err != nil {
		return nil, fmt.Errorf("error creating file on disk: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, fh)
	if err != nil {
		return nil, fmt.Errorf("error saving file to disk: %w", err)
	}

	vf := &repository.ProjectVersionInsert{
		ProjectName: name[0],
		Version:     version[0],
		Digest:      digest,
		DigestType:  digestType,
		FilePath:    fp,
		FileType:    fType[0],
	}
	return vf, nil
}

// computeSHA256 computes the SHA256 checksum of a multipart file.
func computeSHA256(file multipart.File) ([]byte, error) {
	h := sha256.New()
	_, err := file.Seek(0, 0)
	if err != nil {
		return []byte{}, err
	}
	_, err = io.Copy(h, file)
	if err != nil {
		return []byte{}, err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return []byte{}, err
	}
	return h.Sum(nil), nil
}
