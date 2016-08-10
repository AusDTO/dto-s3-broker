package catalog

import (
	"encoding/json"
	"os"

	"github.com/cloudfoundry-community/types-cf"
	"github.com/pkg/errors"
)

func LoadFromFile(path string) (*cf.Catalog, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open %s", path)
	}

	d := json.NewDecoder(f)
	var c cf.Catalog
	err = d.Decode(&c)
	return &c, errors.Wrap(err, "could not decode json")
}
