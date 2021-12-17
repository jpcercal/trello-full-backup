package internal

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

func TestFilesystem(t *testing.T) {
	t.Run("CreateDirectoryRecursively", func(t *testing.T) {
		t.Run("new directory is created", func(t *testing.T) {
			dir := filepath.Join(
				os.TempDir(),
				fmt.Sprintf("dir-%f", rand.Float64()),
			)
			CreateDirectoryRecursively(dir)
			require.DirExists(t, dir)
		})
		t.Run("new directory is created recursively including all parents", func(t *testing.T) {
			dir := filepath.Join(
				os.TempDir(),
				fmt.Sprintf("dir-%f", rand.Float64()),
				fmt.Sprintf("sub-dir-%f", rand.Float64()),
				fmt.Sprintf("sub-sub-dir-%f", rand.Float64()),
			)
			CreateDirectoryRecursively(dir)
			require.DirExists(t, dir)
		})
		t.Run("new directory could not be created", func(t *testing.T) {
			dir := "fake-dir"
			mkdir = func(path string, perm os.FileMode) error {
				return fmt.Errorf("error")
			}

			nullLogger, hook := test.NewNullLogger()
			nullLogger.ExitFunc = func(code int) {}
			logger = nullLogger

			CreateDirectoryRecursively(dir)
			require.NoDirExists(t, dir)

			require.Equal(t, logrus.FatalLevel, hook.LastEntry().Level)
			require.Equal(t, "failed to create directory", hook.LastEntry().Message)
			require.Equal(t, 1, len(hook.Entries))
		})
	})
	t.Run("SaveFile", func(t *testing.T) {
		t.Run("new file is created", func(t *testing.T) {
			filename := filepath.Join(
				os.TempDir(),
				fmt.Sprintf("filename-%f", rand.Float64()),
			)
			content := []byte("")

			SaveFile(filename, content)
			require.FileExists(t, filename)

			file, err := os.ReadFile(filename)
			require.NoError(t, err)
			require.Equal(t, content, file)
		})
		t.Run("directory does not exist, so the new file cannot be created", func(t *testing.T) {
			filename := filepath.Join(
				os.TempDir(),
				fmt.Sprintf("dir-%f", rand.Float64()),
				fmt.Sprintf("sub-dir-%f", rand.Float64()),
				fmt.Sprintf("sub-sub-dir-%f", rand.Float64()),
				fmt.Sprintf("filename-%f", rand.Float64()),
			)
			content := []byte("")

			nullLogger, hook := test.NewNullLogger()
			nullLogger.ExitFunc = func(code int) {}
			logger = nullLogger

			SaveFile(filename, content)
			require.NoFileExists(t, filename)

			require.Equal(t, logrus.FatalLevel, hook.LastEntry().Level)
			require.Equal(t, "failed to save file content", hook.LastEntry().Message)
			require.Equal(t, 1, len(hook.Entries))
		})
	})
	t.Run("Sanitize", func(t *testing.T) {
		type Test struct {
			input          string
			expectedOutput string
		}
		tests := []struct {
			name  string
			tests []Test
		}{
			{
				"does nothing as the string has nothing to be expectedOutput",
				[]Test{{
					input:          "test",
					expectedOutput: "test",
				}},
			},
			{
				"replaces special symbol by underscore",
				[]Test{
					{
						input:          "tes|t",
						expectedOutput: "tes_t",
					}, {
						input:          "tes/t",
						expectedOutput: "tes_t",
					}, {
						input:          "tes?t",
						expectedOutput: "tes_t",
					}, {
						input:          "tes*t",
						expectedOutput: "tes_t",
					}, {
						input:          "tes>t",
						expectedOutput: "tes_t",
					}, {
						input:          "tes<t",
						expectedOutput: "tes_t",
					}, {
						input:          "tes't",
						expectedOutput: "tes_t",
					}, {
						input:          "tes:t",
						expectedOutput: "tes_t",
					},
				},
			},
			{
				"replaces each occurrence of a special symbol by an underscore",
				[]Test{{
					input:          "a|b/c?d*e>f<g'h:i",
					expectedOutput: "a_b_c_d_e_f_g_h_i",
				}},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				for _, sanitizedString := range tt.tests {
					require.Equal(t, sanitizedString.expectedOutput, Sanitize(sanitizedString.input))
				}
			})
		}
	})
}
