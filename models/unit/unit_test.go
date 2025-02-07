// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package unit

import (
	"testing"

	"code.gitea.io/gitea/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestLoadUnitConfig(t *testing.T) {
	defer func(disabledRepoUnits, defaultRepoUnits, defaultForkRepoUnits []Type) {
		DisabledRepoUnits = disabledRepoUnits
		DefaultRepoUnits = defaultRepoUnits
		DefaultForkRepoUnits = defaultForkRepoUnits
	}(DisabledRepoUnits, DefaultRepoUnits, DefaultForkRepoUnits)
	defer func(disabledRepoUnits, defaultRepoUnits, defaultForkRepoUnits []string) {
		setting.Repository.DisabledRepoUnits = disabledRepoUnits
		setting.Repository.DefaultRepoUnits = defaultRepoUnits
		setting.Repository.DefaultForkRepoUnits = defaultForkRepoUnits
	}(setting.Repository.DisabledRepoUnits, setting.Repository.DefaultRepoUnits, setting.Repository.DefaultForkRepoUnits)

	t.Run("regular", func(t *testing.T) {
		setting.Repository.DisabledRepoUnits = []string{"repo.issues"}
		setting.Repository.DefaultRepoUnits = []string{"repo.code", "repo.releases", "repo.issues", "repo.pulls"}
		setting.Repository.DefaultForkRepoUnits = []string{"repo.releases"}
		LoadUnitConfig()
		assert.Equal(t, []Type{TypeIssues}, DisabledRepoUnits)
		assert.Equal(t, []Type{TypeCode, TypeReleases, TypePullRequests}, DefaultRepoUnits)
		assert.Equal(t, []Type{TypeCode, TypeReleases}, DefaultForkRepoUnits)
	})
	t.Run("invalid", func(t *testing.T) {
		setting.Repository.DisabledRepoUnits = []string{"repo.issues", "invalid.1"}
		setting.Repository.DefaultRepoUnits = []string{"repo.code", "invalid.2", "repo.releases", "repo.issues", "repo.pulls"}
		setting.Repository.DefaultForkRepoUnits = []string{"invalid.3", "repo.releases"}
		LoadUnitConfig()
		assert.Equal(t, []Type{TypeIssues}, DisabledRepoUnits)
		assert.Equal(t, []Type{TypeCode, TypeReleases, TypePullRequests}, DefaultRepoUnits)
		assert.Equal(t, []Type{TypeCode, TypeReleases}, DefaultForkRepoUnits)
	})
	t.Run("duplicate", func(t *testing.T) {
		setting.Repository.DisabledRepoUnits = []string{"repo.issues", "repo.issues"}
		setting.Repository.DefaultRepoUnits = []string{"repo.code", "repo.releases", "repo.issues", "repo.pulls", "repo.code"}
		setting.Repository.DefaultForkRepoUnits = []string{"repo.releases", "repo.releases"}
		LoadUnitConfig()
		assert.Equal(t, []Type{TypeIssues}, DisabledRepoUnits)
		assert.Equal(t, []Type{TypeCode, TypeReleases, TypePullRequests}, DefaultRepoUnits)
		assert.Equal(t, []Type{TypeCode, TypeReleases}, DefaultForkRepoUnits)
	})
}
