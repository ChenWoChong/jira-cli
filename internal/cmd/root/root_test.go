package root

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestHasJiraTokenFromConfig(t *testing.T) {
	t.Cleanup(viper.Reset)
	t.Setenv("JIRA_API_TOKEN", "")
	viper.Set("api_token", "config-token")

	assert.True(t, hasJiraToken("https://jira.example.com", "user"))
}

func TestHasJiraTokenWithoutCredentials(t *testing.T) {
	t.Cleanup(viper.Reset)
	t.Setenv("JIRA_API_TOKEN", "")
	viper.Set("api_token", "")

	assert.False(t, hasJiraToken("https://jira.example.com", "user"))
}
