package ddevapp_test

import (
	"github.com/drud/ddev/pkg/ddevapp"
	"github.com/drud/ddev/pkg/fileutil"
	"github.com/drud/ddev/pkg/globalconfig"
	"github.com/drud/ddev/pkg/testcommon"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	asrt "github.com/stretchr/testify/assert"
)

// TestComposer does trivial tests of the ddev composer command
// More tests are found in the cmd package
func TestComposer(t *testing.T) {
	assert := asrt.New(t)
	app := &ddevapp.DdevApp{}

	// Use drupal8 only for this test, just need a little composer action
	site := FullTestSites[1]
	// If running this with GOTEST_SHORT we have to create the directory, tarball etc.
	if site.Dir == "" || !fileutil.FileExists(site.Dir) {
		app := &ddevapp.DdevApp{Name: site.Name}
		_ = app.Stop(true, false)
		_ = globalconfig.RemoveProjectInfo(site.Name)

		err := site.Prepare()
		require.NoError(t, err)
		// nolint: errcheck
		defer os.RemoveAll(site.Dir)
	}

	testDir, _ := os.Getwd()
	// nolint: errcheck
	defer os.Chdir(testDir)
	_ = os.Chdir(site.Dir)

	testcommon.ClearDockerEnv()
	err := app.Init(site.Dir)
	assert.NoError(err)
	//nolint: errcheck
	defer app.Stop(true, false)
	app.Hooks = map[string][]ddevapp.YAMLTask{"post-composer": {{"exec-host": "touch hello-post-composer-" + app.Name}}, "pre-composer": {{"exec-host": "touch hello-pre-composer-" + app.Name}}}
	// Make sure we get rid of this for other uses
	defer func() {
		app.Hooks = nil
		_ = app.WriteConfig()
	}()
	err = app.Start()
	assert.NoError(err)
	_, _, err = app.Composer([]string{"install"})
	assert.NoError(err)
	assert.FileExists("hello-pre-composer-" + app.Name)
	assert.FileExists("hello-post-composer-" + app.Name)
	err = os.Remove("hello-pre-composer-" + app.Name)
	assert.NoError(err)
	err = os.Remove("hello-post-composer-" + app.Name)
	assert.NoError(err)
}
