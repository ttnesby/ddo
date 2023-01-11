//go:build integration

package deploy_test

import (
	"ddo/deploy"
	"testing"
)

// dummy test
func TestDummy(t *testing.T) {
	t.Parallel()

	_ = deploy.New
}
