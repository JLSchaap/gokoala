package engine

import (
	"net/url"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newOpenAPI(t *testing.T) {
	openAPIViaCliArgument, err := filepath.Abs("engine/testdata/ogcapi-tiles-1.bundled.json")
	if err != nil {
		t.Fatalf("can't locate testdata %v", err)
	}

	type args struct {
		config      *Config
		openAPIFile string
	}
	tests := []struct {
		name                         string
		args                         args
		expectedStringsInOpenAPISpec []string
	}{
		{
			name: "Test render OpenAPI spec with MINIMAL config",
			args: args{
				config: &Config{
					Version:  "2.3.0",
					Title:    "Test API",
					Abstract: "Test API description",
					BaseURL:  YAMLURL{&url.URL{Scheme: "https", Host: "api.foobar.example", Path: "/"}},
					OgcAPI: OgcAPI{
						GeoVolumes: nil,
						Tiles:      nil,
						Styles:     nil,
					},
				},
			},
			expectedStringsInOpenAPISpec: []string{
				"Landing page",
				"/conformance",
				"/api",
			},
		},
		{
			name: "Test render OpenAPI spec with OGC Tiles config",
			args: args{
				config: &Config{
					Version:  "2.3.0",
					Title:    "Test API",
					Abstract: "Test API description",
					BaseURL:  YAMLURL{&url.URL{Scheme: "https", Host: "api.foobar.example", Path: "/"}},
					OgcAPI: OgcAPI{
						Tiles: &OgcAPITiles{
							TileServer: YAMLURL{&url.URL{Scheme: "https", Host: "tiles.foobar.example", Path: "/somedataset"}},
						},
					},
				},
			},
			expectedStringsInOpenAPISpec: []string{
				"Landing page",
				"/conformance",
				"/api",
				"Vector Tiles",
				"/tiles/{tileMatrixSetId}",
			},
		},
		{
			name: "Test render OpenAPI spec with OGC Styles config",
			args: args{
				config: &Config{
					Version:  "2.3.0",
					Title:    "Test API",
					Abstract: "Test API description",
					BaseURL:  YAMLURL{&url.URL{Scheme: "https", Host: "api.foobar.example", Path: "/"}},
					OgcAPI: OgcAPI{
						Styles: &OgcAPIStyles{},
					},
				},
			},
			expectedStringsInOpenAPISpec: []string{
				"Landing page",
				"/conformance",
				"/api",
				"/styles",
				"/styles/{styleId}",
				"/resources",
			},
		},
		{
			name: "Test render OpenAPI spec with OGC GeoVolumes config",
			args: args{
				config: &Config{
					Version:  "2.3.0",
					Title:    "Test API",
					Abstract: "Test API description",
					BaseURL:  YAMLURL{&url.URL{Scheme: "https", Host: "api.foobar.example", Path: "/"}},
					OgcAPI: OgcAPI{
						GeoVolumes: &OgcAPI3dGeoVolumes{
							TileServer: YAMLURL{&url.URL{Scheme: "https", Host: "api.foobar.example", Path: "/"}},
							Collections: GeoSpatialCollections{
								GeoSpatialCollection{ID: "feature1"},
								GeoSpatialCollection{ID: "feature2"},
							},
						},
					},
				},
			},
			expectedStringsInOpenAPISpec: []string{
				"Landing page",
				"/conformance",
				"/api",
				"Collections",
				"/collections",
			},
		},
		{
			name: "Test render OpenAPI spec with OGC Tiles and extra spec provided through CLI for overwrite",
			args: args{
				config: &Config{
					Version:  "2.3.0",
					Title:    "Test API",
					Abstract: "Test API description",
					BaseURL:  YAMLURL{&url.URL{Scheme: "https", Host: "api.foobar.example", Path: "/"}},
					OgcAPI: OgcAPI{
						Tiles: &OgcAPITiles{
							TileServer: YAMLURL{&url.URL{Scheme: "https", Host: "tiles.foobar.example", Path: "/somedataset"}},
						},
					},
				},
				openAPIFile: openAPIViaCliArgument,
			},
			expectedStringsInOpenAPISpec: []string{
				"Test API",
				"/conformance",
				"/api",
				"Vector Tiles",
				"/tiles/{tileMatrixSetId}",
				"Map Tiles",                             // extra from given spec through CLI
				"/collections/{collectionId}/map/tiles", // extra from given spec through CLI
			},
		},
		{
			name: "Test render OpenAPI spec with ALL OGC APIs (common, tiles, styles, features, geovolumes)",
			args: args{
				config: &Config{
					Version:  "2.3.0",
					Title:    "Test API",
					Abstract: "Test API description",
					BaseURL:  YAMLURL{&url.URL{Scheme: "https", Host: "api.foobar.example", Path: "/"}},
					OgcAPI: OgcAPI{
						GeoVolumes: &OgcAPI3dGeoVolumes{
							TileServer: YAMLURL{&url.URL{Scheme: "https", Host: "api.foobar.example", Path: "/"}},
							Collections: GeoSpatialCollections{
								GeoSpatialCollection{ID: "feature1"},
								GeoSpatialCollection{ID: "feature2"},
							},
						},
						Tiles: &OgcAPITiles{
							TileServer: YAMLURL{&url.URL{Scheme: "https", Host: "tiles.foobar.example", Path: "/somedataset"}},
						},
						Styles: &OgcAPIStyles{},
						Features: &OgcAPIFeatures{
							Limit: Limit{
								Default: 20,
								Max:     2000,
							},
							Collections: []GeoSpatialCollection{
								{
									ID: "foobar",
									Features: &CollectionEntryFeatures{
										Datasources: &Datasources{
											DefaultWGS84: Datasource{
												GeoPackage: &GeoPackage{
													Local: &GeoPackageLocal{
														File: "./examples/resources/addresses-crs84.gpkg",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedStringsInOpenAPISpec: []string{
				"Landing page",
				"/conformance",
				"/api",
				"Vector Tiles",
				"Features",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			openAPI := newOpenAPI(test.args.config, []string{test.args.openAPIFile}, nil)
			assert.NotNil(t, openAPI)

			// verify resulting OpenAPI spec contains expected strings (keywords, paths, etc)
			for _, expectedStr := range test.expectedStringsInOpenAPISpec {
				assert.Contains(t, string(openAPI.SpecJSON), expectedStr,
					"\"%s\" not found in spec: %s", expectedStr, string(openAPI.SpecJSON))
			}
		})
	}
}
