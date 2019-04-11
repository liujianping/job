package config

import (
	"bufio"
	"bytes"
	"os"

	yaml "gopkg.in/yaml.v2"
)

const yamlSeparator = "\n---"

// splitYAMLDocument is a bufio.SplitFunc for splitting YAML streams into individual documents.
func splitYAMLDocument(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	sep := len([]byte(yamlSeparator))
	if i := bytes.Index(data, []byte(yamlSeparator)); i >= 0 {
		// We have a potential document terminator
		i += sep
		after := data[i:]
		if len(after) == 0 {
			// we can't read any more characters
			if atEOF {
				return len(data), data[:len(data)-sep], nil
			}
			return 0, nil, nil
		}
		if j := bytes.IndexByte(after, '\n'); j >= 0 {
			return i + j + 1, data[0 : i-sep], nil
		}
		return 0, nil, nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

//ParseJDs from config
func ParseJDs(cf string) ([]*JD, error) {
	fd, err := os.Open(cf)
	if err != nil {
		return nil, err
	}
	jds := []*JD{}
	scanner := bufio.NewScanner(fd)
	scanner.Split(splitYAMLDocument)
	for scanner.Scan() {
		var obj map[string]JD
		err := yaml.Unmarshal(scanner.Bytes(), &obj)
		if err != nil {
			return nil, err
		}
		for k, v := range obj {
			if k == "Job" {
				jds = append(jds, &v)
			}
		}
	}
	return jds, nil
}
