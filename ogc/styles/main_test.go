package styles

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
	dir := path.Join(path.Dir(filename), "../../")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestNewStyles(t *testing.T) {
	type args struct {
		e *engine.Engine
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test render templates with OGC Styles config",
			args: args{
				e: engine.NewEngineWithConfig(&engine.Config{
					Version:  "0.4.0",
					Title:    "Test API",
					Abstract: "Test API description",
					Resources: &engine.Resources{
						Directory: "/fakedirectory",
					},
					AvailableLanguages: []language.Tag{language.Dutch},
					BaseURL:            engine.YAMLURL{URL: &url.URL{Scheme: "https", Host: "api.foobar.example", Path: "/"}},
					OgcAPI: engine.OgcAPI{
						GeoVolumes: nil,
						Tiles: &engine.OgcAPITiles{
							TileServer: engine.YAMLURL{URL: &url.URL{Scheme: "https", Host: "tiles.foobar.example", Path: "/somedataset"}},
							Types:      []string{"vector"},
							SupportedSrs: []engine.SupportedSrs{
								{
									Srs: "EPSG:28992",
									ZoomLevelRange: engine.ZoomLevelRange{
										Start: 12,
										End:   12,
									},
								},
							},
						},
						Styles: &engine.OgcAPIStyles{
							Default: "foo",
							SupportedStyles: []engine.StyleMetadata{
								{
									ID:    "foo",
									Title: "bar",
								},
							},
						},
					},
				}, "", false, true),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			styles := NewStyles(test.args.e)
			assert.NotEmpty(t, styles.engine.Templates.RenderedTemplates)
		})
	}
}
