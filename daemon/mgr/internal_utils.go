package mgr

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alibaba/pouch/hookplugins"
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

func updateDockerSpecPath(cid string) string {
	const dockerRuncRoot = "/run/runc"
	return filepath.Join(dockerRuncRoot, cid, "update.json")
}

func createUpdateSpec(cid string, specAnnotation map[string]string) (retErr error) {
	annotations := make(map[string]string)
	if len(specAnnotation) > 0 {
		// filter by the prefix of annotation key
		for k, v := range specAnnotation {
			if strings.HasPrefix(k, hookplugins.AnnotationPrefix) {
				annotations[k] = v
			}
		}
	}

	if len(annotations) == 0 {
		return nil
	}

	path := updateSpecPath(cid)
	_, err := os.Stat(filepath.Dir(path))
	if err != nil {
		if os.IsNotExist(err) {
			// if not exist, try to update in docker-runc
			path = updateDockerSpecPath(cid)
		} else {
			return err
		}
	}

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
		Annotations: annotations,
	}

	err = json.NewEncoder(f).Encode(spec)
	if err != nil {
		return fmt.Errorf("failed to encode annotations: %v", err)
	}

	return nil
}

//clearUpdateSpec clear update.json in docker and pouch runc root if file exists
func clearUpdateSpec(cid string) error {
	err1 := clearPath(updateSpecPath(cid))
	err2 := clearPath(updateDockerSpecPath(cid))

	if err1 != nil {
		return err1
	}

	return err2
}

func clearPath(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	return os.Remove(path)
}
