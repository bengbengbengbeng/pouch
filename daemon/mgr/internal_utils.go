package mgr

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultRuncRoot = "/run/default"
)

type internalUpdateSpec struct {
	// Annotations contains arbitrary metadata for the container.
	Annotations map[string]string `json:"annotations,omitempty"`
}

func updateSpecPath(cid string) string {
	return filepath.Join(defaultRuncRoot, cid, "update.json")
}

func createUpdateSpec(cid string, specAnnotation map[string]string) (retErr error) {
	path := updateSpecPath(cid)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", path, err)
	}

	defer func() {
		if retErr != nil {
			// if return error, remove updateSpec file
			os.Remove(path)
		}
	}()

	defer f.Close()

	spec := internalUpdateSpec{
		Annotations: specAnnotation,
	}

	err = json.NewEncoder(f).Encode(spec)
	if err != nil {
		return fmt.Errorf("failed to encode annotations: %v", err)
	}

	return nil
}

func clearUpdateSpec(cid string) error {
	path := updateSpecPath(cid)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	return os.Remove(path)
}
