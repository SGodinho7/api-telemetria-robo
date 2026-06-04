package util

import (
	"encoding/base64"
	"encoding/csv"
	"strings"
)

func DecodeBase64CSV(b64String string) ([][]string, error) {
	var (
		records [][]string
		err     error
	)

	decodedBytes, err := base64.StdEncoding.DecodeString(b64String)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(strings.NewReader(string(decodedBytes)))

	records, err = reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}
