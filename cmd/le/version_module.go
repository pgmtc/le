package main

import "github.com/pgmtc/le/pkg/common"

type VersionModule struct{}

func (VersionModule) GetActions() map[string]common.Action {
	return map[string]common.Action{
		"default": &common.RawAction{
			Handler: func(ctx common.Context, args ...string) error {
				ctx.Log.Infof("Version: %s\n", VERSION)
				return nil
			},
		},
		"update": &updateCliAction,
	}
}
