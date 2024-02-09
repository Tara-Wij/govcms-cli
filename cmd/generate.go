/*
Copyright © 2024 Joseph Zhao pandaski@outlook.com.au

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/govcms-tests/govcms-cli/pkg/govcms"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:       "generate [resource]",
	Short:     "Creates a GovCMS distribution tailored for either SaaS or PaaS deployment",
	Long:      "Creates a GovCMS distribution tailored for either SaaS or PaaS deployment.",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"distribution", "saas", "paas"},
	Run: func(cmd *cobra.Command, args []string) {
		// Check if both flags are provided
		if cmd.Flags().Changed("pr") && cmd.Flags().Changed("branch") {
			fmt.Println("Error: Cannot specify both --pr and --branch flags together.")
			return
		}
		resource := args[0]
		prNumber, _ := cmd.Flags().GetInt("pr")
		branchName, _ := cmd.Flags().GetString("branch")

		// Call the generate function from the govcms package
		err := govcms.Generate(resource, prNumber, branchName)
		if err != nil {
			fmt.Printf("Error generating %s: %v\n", resource, err)
			return
		}

		pathErr := os.Mkdir("govcms", os.ModePerm)
		if pathErr != nil {
			fmt.Println("Invalid path")
			return
		}
		// Define the target folder where repositories will be cloned
		targetFolder := "govcms"
		// Clone the corresponding repository
		repoURL := map[string]string{
			"distribution": "govCMS/GovCMS",
			"saas":         "govCMS/scaffold",
			"paas":         "govCMS/scaffold",
		}[resource]
		repoPath := filepath.Join(targetFolder, resource)
		fmt.Printf("Cloning %s into %s\n", repoURL, repoPath)
		_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
			URL:      "https://github.com/" + repoURL + ".git",
			Progress: os.Stdout,
		})
		if err != nil && err != git.ErrRepositoryAlreadyExists {
			fmt.Printf("Error cloning repository: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().IntP("pr", "p", 0, "Github PR number")
	generateCmd.Flags().StringP("branch", "b", "", "Git branch name")
}
