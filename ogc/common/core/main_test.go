package core

import (
	"net/url"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/PDOK/gokoala/engine"
	"golang.org/x/text/language"

	"github.com/stretchr/testify/assert"
)

func init() {
	// change working dir to root, to mimic behavior of 'go run' in order to resolve template files.
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../../")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestNewCommonCore(t *testing.T) {
	type args struct {
		e *engine.Engine
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test render templates with MINIMAL config",
			args: args{
				e: engine.NewEngineWithConfig(&engine.Config{
					Version:            "2.3.0",
					Title:              "Test API",
					Abstract:           "Test API description",
					AvailableLanguages: []language.Tag{language.Dutch},
					BaseURL:            engine.YAMLURL{URL: &url.URL{Scheme: "https", Host: "api.foobar.example", Path: "/"}},
					OgcAPI: engine.OgcAPI{
						GeoVolumes: nil,
						Tiles:      nil,
						Styles:     nil,
					},
				}, "", false, true),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			core := NewCommonCore(test.args.e)
			assert.NotEmpty(t, core.engine.Templates.RenderedTemplates)
		})
	}
}
