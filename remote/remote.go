package remote

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rumpl/devoxx-docker/oci"
)

type ImagePuller struct {
	reference string
	client    *http.Client
}

func NewImagePuller(ref string) *ImagePuller {
	return &ImagePuller{
		reference: ref,
		client:    &http.Client{},
	}
}

func (p *ImagePuller) parseReference() (string, string) {
	parts := strings.Split(p.reference, "/")
	var repository, tag string

	if len(parts) == 1 {
		repository = "library/" + parts[0]
	} else if len(parts) == 2 {
		repository = strings.Join(parts[0:2], "/")
	} else {
		repository = strings.Join(parts[1:], "/")
	}

	// Handle tag
	repoParts := strings.Split(repository, ":")
	if len(repoParts) > 1 {
		repository = repoParts[0]
		tag = repoParts[1]
	} else {
		tag = "latest"
	}

	return repository, tag
}

func (p *ImagePuller) Pull() error {
	repository, tag := p.parseReference()

	// Get auth token
	tokenURL := fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull", repository)
	resp, err := p.client.Get(tokenURL)
	if err != nil {
		return fmt.Errorf("failed to get auth token: %v", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %v", err)
	}

	// First try to get the manifest list
	manifestURL := fmt.Sprintf("https://registry-1.docker.io/v2/%s/manifests/%s", repository, tag)
	req, err := http.NewRequest("GET", manifestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create manifest request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	// Request manifest list first
	req.Header.Set("Accept", "application/vnd.oci.image.index.v1+json")

	resp, err = p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get manifest list: %v", err)
	}
	defer resp.Body.Close()

	var manifestDigest string

	// Try to parse as manifest list first
	var idx oci.Index
	if err := json.NewDecoder(resp.Body).Decode(&idx); err != nil {
		return fmt.Errorf("failed to decode manifest list: %v", err)
	}

	// If we got a manifest list, find the right manifest for our platform
	if idx.MediaType == "application/vnd.oci.image.index.v1+json" {
		found := false
		for _, m := range idx.Manifests {
			if m.Platform.OS == runtime.GOOS && m.Platform.Architecture == runtime.GOARCH {
				manifestDigest = m.Digest
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("no manifest found for platform %s/%s", runtime.GOOS, runtime.GOARCH)
		}

		// Get the specific manifest
		manifestURL = fmt.Sprintf("https://registry-1.docker.io/v2/%s/manifests/%s", repository, manifestDigest)
		req, err = http.NewRequest("GET", manifestURL, nil)
		if err != nil {
			return fmt.Errorf("failed to create manifest request: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
		req.Header.Set("Accept", "application/vnd.oci.image.manifest.v1+json")

		resp, err = p.client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to get manifest: %v", err)
		}
		defer resp.Body.Close()
	} else {
		return fmt.Errorf("unexpected media type: %s", idx.MediaType)
	}

	var manifest oci.Manifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return fmt.Errorf("failed to decode manifest: %v", err)
	}

	// Create destination directory
	imageName := strings.Split(repository, "/")[len(strings.Split(repository, "/"))-1]
	destDir := filepath.Join("/fs", imageName)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Download and extract layers
	for _, layer := range manifest.Layers {
		layerURL := fmt.Sprintf("https://registry-1.docker.io/v2/%s/blobs/%s", repository, layer.Digest)
		req, err := http.NewRequest("GET", layerURL, nil)
		if err != nil {
			return fmt.Errorf("failed to create layer request: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+tokenResp.Token)

		resp, err = p.client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to download layer: %v", err)
		}
		defer resp.Body.Close()

		if err := extractTar(resp.Body, destDir); err != nil {
			return fmt.Errorf("failed to extract layer: %v", err)
		}
	}

	return nil
}

func extractTar(r io.Reader, dest string) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %v", err)
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", target, err)
			}
		case tar.TypeReg:
			dir := filepath.Dir(target)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dir, err)
			}

			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %v", target, err)
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return fmt.Errorf("failed to write to file %s: %v", target, err)
			}
			f.Close()
		case tar.TypeSymlink:
			dir := filepath.Dir(target)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dir, err)
			}
			if err := os.Symlink(header.Linkname, target); err != nil {
				return fmt.Errorf("failed to create symlink %s -> %s: %v", target, header.Linkname, err)
			}
		}
	}

	return nil
}
