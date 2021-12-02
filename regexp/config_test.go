package regexp

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config"
)

func TestConfig(t *testing.T) {
	testFile := path.Join(".", "testdata", "config.yaml")
	v, err := config.NewMapFromFile(testFile)
	require.NoError(t, err)

	actualConfigs := map[string]*Config{}
	require.NoErrorf(t, v.UnmarshalExact(&actualConfigs), "unable to unmarshal yaml from file %v", testFile)

	expectedConfigs := map[string]*Config{
		"regexp/default": {},
		"regexp/cachedisabledwithsize": {
			CacheEnabled:       false,
			CacheMaxNumEntries: 10,
		},
		"regexp/cacheenablednosize": {
			CacheEnabled: true,
		},
	}

	for testName, actualCfg := range actualConfigs {
		t.Run(testName, func(t *testing.T) {
			expCfg, ok := expectedConfigs[testName]
			assert.True(t, ok)
			assert.Equal(t, expCfg, actualCfg)

			fs, err := NewFilterSet([]string{}, actualCfg)
			assert.NoError(t, err)
			assert.NotNil(t, fs)
		})
	}
}
