package list

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func newJQLTestCommand() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().String("jql", "", "")
	return cmd
}

func TestApplyDefaultJQL(t *testing.T) {
	t.Cleanup(viper.Reset)

	t.Run("uses configured query when no explicit query is supplied", func(t *testing.T) {
		viper.Set("issue.list.jql", "assignee = currentUser()")
		cmd := newJQLTestCommand()

		applyDefaultJQL(cmd)

		assert.Equal(t, "assignee = currentUser()", cmd.Flag("jql").Value.String())
	})

	t.Run("does not override an explicit query", func(t *testing.T) {
		viper.Set("issue.list.jql", "assignee = currentUser()")
		cmd := newJQLTestCommand()
		assert.NoError(t, cmd.Flags().Set("jql", "status = Open"))

		applyDefaultJQL(cmd)

		assert.Equal(t, "status = Open", cmd.Flag("jql").Value.String())
	})

	t.Run("does nothing when no default is configured", func(t *testing.T) {
		viper.Set("issue.list.jql", "")
		cmd := newJQLTestCommand()

		applyDefaultJQL(cmd)

		assert.Equal(t, "", cmd.Flag("jql").Value.String())
	})
}

func TestApplyUnfinishedFilter(t *testing.T) {
	t.Run("adds the unfinished filter", func(t *testing.T) {
		cmd := newJQLTestCommand()
		cmd.Flags().Bool("unfinished", true, "")

		applyUnfinishedFilter(cmd)

		assert.Equal(t, "statusCategory != Done", cmd.Flag("jql").Value.String())
	})

	t.Run("combines with an existing query", func(t *testing.T) {
		cmd := newJQLTestCommand()
		cmd.Flags().Bool("unfinished", true, "")
		assert.NoError(t, cmd.Flags().Set("jql", "assignee = currentUser()"))

		applyUnfinishedFilter(cmd)

		assert.Equal(t, "(assignee = currentUser()) AND statusCategory != Done", cmd.Flag("jql").Value.String())
	})

	t.Run("does not change the query when disabled", func(t *testing.T) {
		cmd := newJQLTestCommand()
		cmd.Flags().Bool("unfinished", false, "")
		assert.NoError(t, cmd.Flags().Set("jql", "assignee = currentUser()"))

		applyUnfinishedFilter(cmd)

		assert.Equal(t, "assignee = currentUser()", cmd.Flag("jql").Value.String())
	})

	t.Run("is idempotent across refreshes", func(t *testing.T) {
		cmd := newJQLTestCommand()
		cmd.Flags().Bool("unfinished", true, "")

		applyUnfinishedFilter(cmd)
		applyUnfinishedFilter(cmd)

		assert.Equal(t, "statusCategory != Done", cmd.Flag("jql").Value.String())
	})
}
