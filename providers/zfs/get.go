package zfs

import (
	"errors"
	"net/url"

	"github.com/mistifyio/gozfs"
	"github.com/mistifyio/mistify/acomm"
)

// GetResult is the result data for the List handler.
type GetResult struct {
	Dataset *Dataset `json:"dataset"`
}

// Get returns information about a dataset.
func (z *ZFS) Get(req *acomm.Request) (interface{}, *url.URL, error) {
	var args CommonArgs
	if err := req.UnmarshalArgs(&args); err != nil {
		return nil, nil, err
	}

	if args.Name == "" {
		return nil, nil, errors.New("missing arg: name")
	}

	ds, err := gozfs.GetDataset(args.Name)
	if err != nil {
		return nil, nil, err
	}

	return &GetResult{newDataset(ds)}, nil, nil
}
