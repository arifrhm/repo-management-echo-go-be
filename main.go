package main  

import (  
 "fmt"  
 "log"  
 "net/http"  
 "os"  
 "os/exec"  

 "github.com/joho/godotenv"  
 "github.com/labstack/echo/v4"  
 "github.com/labstack/echo/v4/middleware"  
)  

// RepositoryConfig represents the configuration for repository pull  
type RepositoryConfig struct {  
 Path   string `json:"path" validate:"required"`  
 Branch string `json:"branch,omitempty"`  
}  

// APIResponse standardizes API response structure  
type APIResponse struct {  
 Status  string `json:"status"`  
 Message string `json:"message"`  
 Branch  string `json:"branch,omitempty"`  
}  

// Logger configuration  
func setupLogger() *log.Logger {  
 return log.New(os.Stdout, "REPO_MANAGER: ", log.Ldate|log.Ltime|log.Lshortfile)  
}  

// validateAPIKey checks the provided API key  
func validateAPIKey(next echo.HandlerFunc) echo.HandlerFunc {  
 return func(c echo.Context) error {  
  apiKey := os.Getenv("REPO_MANAGEMENT_API_KEY")  
  providedKey := c.Request().Header.Get("X-API-Key")  

  if apiKey == "" || providedKey != apiKey {  
   return echo.NewHTTPError(http.StatusForbidden, "Invalid or missing API key")  
  }  

  return next(c)  
 }  
}  

// pullRepository handles the repository pull operation  
func pullRepository(c echo.Context) error {  
 // Initialize logger  
 logger := setupLogger()  

 // Parse request body  
 var repoConfig RepositoryConfig  
 if err := c.Bind(&repoConfig); err != nil {  
  logger.Printf("Binding error: %v", err)  
  return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")  
 }  

 // Validate repository path  
 if _, err := os.Stat(repoConfig.Path); os.IsNotExist(err) {  
  logger.Printf("Repository path not found: %s", repoConfig.Path)  
  return echo.NewHTTPError(http.StatusInternalServerError, "Repository path not found")  
 }  

 // Default to main branch if not specified  
 if repoConfig.Branch == "" {  
  repoConfig.Branch = "main"  
 }  

 // Perform git fetch  
 fetchCmd := exec.Command("git", "-C", repoConfig.Path, "fetch", "origin")  
 if fetchOutput, err := fetchCmd.CombinedOutput(); err != nil {  
  logger.Printf("Fetch failed: %v - Output: %s", err, string(fetchOutput))  
  return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Fetch failed: %v", err))  
 }  

 // Perform git reset  
 resetCmd := exec.Command("git", "-C", repoConfig.Path, "reset", "--hard", fmt.Sprintf("origin/%s", repoConfig.Branch))  
 if resetOutput, err := resetCmd.CombinedOutput(); err != nil {  
  logger.Printf("Reset failed: %v - Output: %s", err, string(resetOutput))  
  return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Reset failed: %v", err))  
 }  

 // Log successful operation  
 logger.Printf("Successfully pulled repository at %s on branch %s", repoConfig.Path, repoConfig.Branch)  

 // Return success response  
 return c.JSON(http.StatusOK, APIResponse{  
  Status:  "success",  
  Message: "Repository successfully updated",  
  Branch:  repoConfig.Branch,  
 })  
}  

func main() {  
 // Load environment variables  
 if err := godotenv.Load(); err != nil {  
  log.Println("No .env file found")  
 }  

 // Create Echo instance  
 e := echo.New()  

 // Middleware  
 e.Use(middleware.Recover())  
 e.Use(middleware.Logger())  
 e.Use(middleware.CORS())  

 // Routes  
 e.POST("/pull-repo", pullRepository, validateAPIKey)  

 // Configure server  
 port := os.Getenv("PORT")  
 if port == "" {  
  port = "8080"  
 }  

 // Start server  
 e.Logger.Fatal(e.Start(":" + port))  
} 