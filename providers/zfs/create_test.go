package zfs_test

import (
	"fmt"
	"path/filepath"

	"github.com/mistifyio/gozfs"
	"github.com/mistifyio/mistify/acomm"
	zfsp "github.com/mistifyio/mistify/providers/zfs"
)

func (s *zfs) TestCreate() {
	fs := gozfs.DatasetFilesystem
	vol := gozfs.DatasetVolume
	tests := []struct {
		args *zfsp.CreateArgs
		err  string
	}{
		{&zfsp.CreateArgs{Name: "", Type: fs}, "missing arg: name"},
		{&zfsp.CreateArgs{Name: "fs1", Type: fs}, eexist},
		{&zfsp.CreateArgs{Name: "fs1/~1", Type: fs}, einval},
		{&zfsp.CreateArgs{Name: "fsbadprop", Type: fs, Properties: map[string]interface{}{"foo": "bar"}}, einval},
		{&zfsp.CreateArgs{Name: "foobar/fs1", Type: fs}, enoent},
		{&zfsp.CreateArgs{Name: "fscreatebad", Type: "asdf"}, "missing or invalid arg: type"},
		{&zfsp.CreateArgs{Name: "fscreate1", Type: fs}, ""},
		{&zfsp.CreateArgs{Name: "fscreate2", Type: fs, Properties: map[string]interface{}{"foo:bar": "baz"}}, ""},

		{&zfsp.CreateArgs{Name: "", Type: vol, Properties: map[string]interface{}{"volsize": 8192}}, "missing arg: name"},
		{&zfsp.CreateArgs{Name: "vol1", Type: vol, Properties: nil}, "missing or invalid arg: volsize"},
		{&zfsp.CreateArgs{Name: "vol2", Type: vol, Volsize: 0, Properties: nil}, "missing or invalid arg: volsize"},
		{&zfsp.CreateArgs{Name: "foovol/vol1", Type: vol, Volsize: 8192, Properties: nil}, enoent},
		{&zfsp.CreateArgs{Name: "volbadprop", Type: vol, Volsize: 8192, Properties: map[string]interface{}{"foo": "bar"}}, einval},
		{&zfsp.CreateArgs{Name: "vol3", Type: vol, Volsize: 8192, Properties: nil}, ""},
		{&zfsp.CreateArgs{Name: "fs2/vol1", Type: vol, Volsize: 8192, Properties: nil}, eexist},
		{&zfsp.CreateArgs{Name: "vol4", Type: vol, Volsize: 1024, Properties: map[string]interface{}{"volblocksize": 1024}}, ""},
	}

	for _, test := range tests {
		if test.args.Name != "" {
			test.args.Name = filepath.Join(s.pool, test.args.Name)
		}
		argsS := fmt.Sprintf("%+v", test.args)

		req, err := acomm.NewRequest("zfs-create", "unix:///tmp/foobar", "", test.args, nil, nil)
		s.Require().NoError(err, argsS)

		res, streamURL, err := s.zfs.Create(req)
		s.Empty(streamURL, argsS)
		if test.err == "" {
			if !s.Nil(err, argsS) {
				continue
			}
			if !s.NotNil(res, argsS) {
				continue
			}

			result, ok := res.(*zfsp.DatasetResult)
			if !s.True(ok, argsS) {
				continue
			}
			ds := result.Dataset
			if !s.NotNil(ds, argsS) {
				continue
			}
			s.Equal(ds.Name, test.args.Name, argsS)
			s.Equal(ds.Properties.Type, test.args.Type, argsS)
			if _, ok := test.args.Properties["foo:bar"]; ok {
				s.Equal(ds.Properties.UserDefined["foo:bar"], "baz", argsS)
			}
			if test.args.Type == vol {

			}
		} else {
			s.EqualError(err, test.err, argsS)
		}
	}
}
