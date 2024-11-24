package cmd

// import (
// 	"fmt"
// 	"os"
// 	"sync"

// 	"github.com/msetsma/RepoRover/core"
// 	"github.com/msetsma/RepoRover/models"
// 	"github.com/spf13/cobra"
// )

// var (
// 	org     string
// 	project string
// 	pat     string
// )

// var fetchCmd = &cobra.Command{
// 	Use:   "fetch",
// 	Short: "Fetch repositories from Azure DevOps",
// 	Run:   runFetchCommand,
// }

// func init() {
// 	// Register flags
// 	fetchCmd.Flags().StringVar(&org, "org", "", "Azure DevOps organization")
// 	fetchCmd.Flags().StringVar(&project, "project", "", "Azure DevOps project name")
// 	fetchCmd.Flags().StringVar(&pat, "pat", "", "Azure DevOps personal access token")

// 	// Add fetchCmd to rootCmd
// 	rootCmd.AddCommand(fetchCmd)
// }

// // runFetchCommand orchestrates the fetch process
// func runFetchCommand(cmd *cobra.Command, args []string) {
// 	dbName, dbConfig := loadAndValidateConfig()
// 	db := initializeDatabase(dbName)
// 	defer db.Close()

// 	repos := fetchRepositories(dbConfig)
// 	processRepositoryDetails(repos, db, dbConfig)
// }

// // loadAndValidateConfig loads the configuration, validates required fields, and merges with flags.
// func loadAndValidateConfig() (string, map[string]string) {
// 	// Load configuration
// 	config, err := core.LoadConfig()
// 	if err != nil {
// 		fmt.Printf("Error loading configuration: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// Determine the database name
// 	dbFlagValue := rootCmd.PersistentFlags().Lookup("db").Value.String()
// 	dbName := core.GetDatabaseName(dbFlagValue, "default")

// 	// Get database-specific config
// 	dbConfig := config.Databases[dbName]
// 	if dbConfig == nil {
// 		dbConfig = make(map[string]string)
// 	}

// 	// Merge flags with config
// 	if org != "" {
// 		dbConfig["azure_org"] = org
// 	}
// 	if project != "" {
// 		dbConfig["azure_project"] = project
// 	}
// 	if pat != "" {
// 		dbConfig["azure_pat"] = pat
// 	}

// 	// Validate required fields
// 	if dbConfig["azure_org"] == "" || dbConfig["azure_project"] == "" || dbConfig["azure_pat"] == "" {
// 		fmt.Println("Error: Missing required arguments --org, --project, or --pat.")
// 		fmt.Println("Please set them using `reporover config set` or provide them as flags.")
// 		os.Exit(1)
// 	}

// 	// Save updated config
// 	config.Databases[dbName] = dbConfig
// 	if err := core.SaveConfig(config); err != nil {
// 		fmt.Printf("Error saving configuration: %v\n", err)
// 		os.Exit(1)
// 	}

// 	return dbName, dbConfig
// }

// // initializeDatabase ensures the database is initialized and ready for use
// func initializeDatabase(dbName string) *core.Database {
// 	db, err := core.GetDatabaseInstance(dbName)
// 	if err != nil {
// 		fmt.Printf("Error initializing database: %v\n", err)
// 		os.Exit(1)
// 	}
// 	return db
// }

// // fetchRepositories fetches the list of repositories from Azure DevOps
// func fetchRepositories(config map[string]string) []models.Repository {
// 	fmt.Printf("Fetching repositories for organization: %s, project: %s\n", config["azure_org"], config["azure_project"])
// 	repos, err := core.FetchRepositories(config["azure_org"], config["azure_project"], config["azure_pat"])
// 	if err != nil {
// 		fmt.Printf("Error fetching repositories: %v\n", err)
// 		os.Exit(1)
// 	}
// 	return repos
// }

// // processRepositoryDetails fetches additional details and saves them to the database
// func processRepositoryDetails(repos []models.Repository, db *core.Database, config map[string]string) {
// 	var wg sync.WaitGroup
// 	repoChan := make(chan *models.FetchResult, len(repos))

// 	// Fetch additional data concurrently
// 	for _, repo := range repos {
// 		wg.Add(1)
// 		go func(repo models.Repository) {
// 			defer wg.Done()
// 			fullRepo, err := core.FetchAdditionalRepoData(config["azure_org"], config["azure_project"], config["azure_pat"], repo.ID)
// 			repoChan <- &models.FetchResult{
// 				Repo: fullRepo,
// 				Err:  err,
// 			}
// 		}(repo)
// 	}

// 	// Close the channel after all Goroutines finish
// 	go func() {
// 		wg.Wait()
// 		close(repoChan)
// 	}()

// 	// Process results
// 	for result := range repoChan {
// 		if result.Err != nil {
// 			fmt.Printf("Error fetching details for repo %s: %v\n", result.Repo.Name, result.Err)
// 			continue
// 		}

// 		err := db.SaveRepository(result.Repo)
// 		if err != nil {
// 			fmt.Printf("Error saving repository %s: %v\n", result.Repo.Name, err)
// 		} else {
// 			fmt.Printf("Saved repository: %s\n", result.Repo.Name)
// 		}
// 	}
// }