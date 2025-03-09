package internal

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type CredentialAttribute struct {
	Key   string
	Value string
}

func ReadCredentialAttributes(r io.Reader) ([]CredentialAttribute, error) {
	s := bufio.NewScanner(r)
	a := make([]CredentialAttribute, 0)

	for s.Scan() {
		line := s.Text()

		if line == "" {
			break
		}

		key, value, ok := strings.Cut(line, "=")

		if !ok {
			return nil, fmt.Errorf("missing '=' character in git credential input")
		}

		a = append(a, CredentialAttribute{key, value})
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return a, nil
}

func WriteCredentialAttributes(w io.Writer, attrs []CredentialAttribute) error {
	for _, attr := range attrs {
		if _, err := fmt.Fprintf(w, "%s=%s\n", attr.Key, attr.Value); err != nil {
			return err
		}
	}

	_, err := fmt.Fprintln(w)

	return err
}
