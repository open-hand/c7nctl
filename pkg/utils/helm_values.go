package utils

import (
	"fmt"
	"github.com/ghodss/yaml"
	"helm.sh/helm/v3/pkg/strvals"
)

// Merges source and destination map, preferring values from the source map
func mergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			dest[k] = v
			continue
		}
		// Edge case: If the key exists in the destination, but isn't a map
		destMap, isMap := dest[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			dest[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		dest[k] = mergeValues(destMap, nextMap)
	}
	return dest
}

// vals merges values from files specified via -f/--values and
// directly via --set or --set-string or --set-file, marshaling them to YAML
func Vals(values []string, fileValues string) (map[string]interface{}, error) {
	base := map[string]interface{}{}

	// User specified a values files via -f/--values
	currentMap := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(fileValues), &currentMap); err != nil {
		return nil, fmt.Errorf("failed to parse %s", err)
	}
	// Merge with the previous map
	base = mergeValues(base, currentMap)

	// User specified a value via --set
	for _, value := range values {
		if err := strvals.ParseInto(value, base); err != nil {
			return nil, fmt.Errorf("failed parsing --set data: %s", err)
		}
	}

	return base, nil
}
