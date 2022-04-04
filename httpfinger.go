package appfinger

import (
	"errors"
	"strings"
)

func New(source string) (*Database, error) {
	var data = &Database{database: []*FingerPrint{}}

	for _, line := range strings.Split(source, "\n") {
		err := data.Add(line)
		if err != nil {
			return nil, errors.New(err.Error() + line)
		}
	}
	return data, nil
}
