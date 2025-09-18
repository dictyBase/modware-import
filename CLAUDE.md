# Go Coding Conventions

-- **Build, Test, and Lint Commands**
    - Run all tests: `gotestsum --format-hide-empty-pkg --format dots`
    - Run specific test: `gotestsum --format-hide-empty-pkg --format dots -- -run TestFindSimilar ./...`
    - Run tests with verbose output: `gotestum --format-hide-empty-pkg --format standard-verbose`
    - Format code: `gofumpt -w .`
    - Lint codebase: `golangci-lint run`

- **Code Style Guidelines**
    - Imports: Standard library first, then external packages, then internal packages
    - Prefer functional programming utilities from collection package where appropriate
    - Use `slices.DeleteFunc` for conditional element removal instead of manual for loops
      ```go
      import "slices"
      
      // Avoid: Manual loop for conditional removal
      var result []Item
      for _, item := range items {
          if !shouldRemove(item) {
              result = append(result, item)
          }
      }
      items = result
      
      // Preferred: Use slices.DeleteFunc
      items = slices.DeleteFunc(items, shouldRemove)
      
      // Example: Remove users with specific role
      users = slices.DeleteFunc(users, func(u User) bool {
          return u.Role == "banned"
      })
      
      // Example: Remove tags matching both name and value
      tags = slices.DeleteFunc(tags, func(tag TagProperty) bool {
          return tag.Name == targetName && tag.Value == targetValue
      })
      ```
      - Essential utility functions:
        ```go
        func Map[T, U any](ts []T, f func(T) U) []U {
            us := make([]U, len(ts))
            for i, t := range ts {
                us[i] = f(t)
            }
            return us
        }
        
        func Filter[T any](slice []T, predicate func(T) bool) []T {
            var result []T
            for _, item := range slice {
                if predicate(item) {
                    result = append(result, item)
                }
            }
            return result
        }
        
        func Find[T any](slice []T, predicate func(T) bool) (*T, bool) {
            for i := range slice {
                if predicate(slice[i]) {
                    return &slice[i], true
                }
            }
            return nil, false
        }
        
        func MapWithError[T, U any](ts []T, f func(T) (U, error)) ([]U, error) {
            us := make([]U, len(ts))
            for i, t := range ts {
                u, err := f(t)
                if err != nil {
                    return nil, err
                }
                us[i] = u
            }
            return us, nil
        }
        ```
      - Usage examples:
        ```go
        // Helper functions for transformations
        func stringToInt(s string) int {
            n, _ := strconv.Atoi(s)
            return n
        }
        
        func isEven(n int) bool {
            return n%2 == 0
        }
        
        func isBob(u string) bool {
            return u == "bob"
        }
        
        // Transform slice elements
        strings := []string{"1", "2", "3"}
        numbers := Map(strings, stringToInt)
        
        // Filter slice elements
        numbers := []int{1, 2, 3, 4, 5}
        evens := Filter(numbers, isEven)
        
        // Find first matching element
        users := []string{"alice", "bob", "charlie"}
        user, found := Find(users, isBob)
        
        // For complex operations with error handling
        func (c *Client) processArticleWithError(pmid string) (*Article, error) {
            article, err := c.GetArticle(pmid)
            if err != nil {
                return nil, fmt.Errorf("failed to process article %s: %w", pmid, err)
            }
            return article, nil
        }
        
        // Usage with MapWithError
        pmids := []string{"12345", "67890", "11111"}
        articles, err := MapWithError(pmids, c.processArticleWithError)
        ```
    - **Functional Programming with fp-go Library**
      - Use `github.com/IBM/fp-go` for monadic patterns, pipes, and functional composition
      - Import pattern: Use single-letter aliases for fp-go modules
      ```go
      import (
          A "github.com/IBM/fp-go/array"
          E "github.com/IBM/fp-go/either"
          F "github.com/IBM/fp-go/function"
          O "github.com/IBM/fp-go/option"
          T "github.com/IBM/fp-go/tuple"
          S "github.com/IBM/fp-go/string"
          J "github.com/IBM/fp-go/json"
          H "github.com/IBM/fp-go/context/readerioeither/http"
      )
      ```

      - **Either Monad for Error Handling**
        ```go
        // Use Either for operations that can fail
        func validateAndProcess(input string) E.Either[error, ProcessedData] {
            if input == "" {
                return E.Left[ProcessedData](fmt.Errorf("input cannot be empty"))
            }

            processed := ProcessedData{Value: strings.ToUpper(input)}
            return E.Right[error](processed)
        }

        // Chain Either operations with Map and Bind
        func processWorkflow(input string) E.Either[error, string] {
            return F.Pipe3(
                input,
                validateAndProcess,
                E.Map[error](func(data ProcessedData) ProcessedData {
                    data.Timestamp = time.Now()
                    return data
                }),
                E.Map[error](func(data ProcessedData) string {
                    return fmt.Sprintf("%s at %v", data.Value, data.Timestamp)
                }),
            )
        }

        // Handle Either results with Fold
        result := F.Pipe1(
            processWorkflow("hello"),
            E.Fold[error, string](
                func(err error) string {
                    return fmt.Sprintf("Error: %v", err)
                },
                func(success string) string {
                    return success
                },
            ),
        )
        ```

      - **Option Monad for Nullable Values**
        ```go
        // Use Option instead of pointers for optional values
        func findUser(id string) O.Option[User] {
            user, found := userDatabase[id]
            if found {
                return O.Some(user)
            }
            return O.None[User]()
        }

        // Chain Option operations
        func getUserEmail(id string) O.Option[string] {
            return F.Pipe1(
                findUser(id),
                O.Map(func(user User) string {
                    return user.Email
                }),
            )
        }

        // Option with fallback using Alt
        func getConfigValue(key string) O.Option[string] {
            return F.Pipe2(
                getFromEnv(key),
                O.Alt(func() O.Option[string] {
                    return getFromFile(key)
                }),
                O.Alt(func() O.Option[string] {
                    return getDefaultValue(key)
                }),
            )
        }

        // Extract Option values safely
        email := F.Pipe1(
            getUserEmail("user123"),
            O.GetOrElse(F.Constant("no-email@example.com")),
        )
        ```

      - **Function Composition with Pipe and Flow**
        ```go
        // Use F.Pipe for sequential composition (left to right)
        func processFileName(fileName string) string {
            return F.Pipe7(
                fileName,
                filepath.Base,
                strings.Split("."),
                A.Head[string],
                O.GetOrElse(F.Constant("")),
                strings.Split("_"),
                A.SliceRight[string](2),
                S.Join(":"),
            )
        }

        // Use F.Flow for creating reusable composed functions
        var processUserData = F.Flow3(
            validateInput,
            normalizeData,
            saveToDatabase,
        )

        // Complex pipeline with error handling
        func processFiles(files []string) E.Either[error, []ProcessedFile] {
            return F.Pipe2(
                E.TryCatchError(readFiles(files)),
                E.Map[error](func(fileContents []string) []ProcessedFile {
                    return F.Pipe3(
                        fileContents,
                        A.Filter(isValidContent),
                        A.Map(parseContent),
                        A.FilterMap(validateParsedContent),
                    )
                }),
                E.Fold[error, []ProcessedFile](
                    func(err error) E.Either[error, []ProcessedFile] {
                        return E.Left[[]ProcessedFile](err)
                    },
                    func(result []ProcessedFile) E.Either[error, []ProcessedFile] {
                        return E.Right[error](result)
                    },
                ),
            )
        }
        ```

      - **Curried Functions for Reusability**
        ```go
        // Create curried functions for partial application
        var HasField = F.Curry2(func(name string, field FieldType) bool {
            return field.Name == name
        })

        var SearchUser = F.Curry2(func(email string, workspace WorkspaceResp) int {
            return F.Pipe3(
                workspace.Users,
                A.FindFirst(HasUser(email)),
                O.Map(func(user WorkspaceUserResp) int { return user.Id }),
                O.GetOrElse(F.Constant(0)),
            )
        })

        // Usage of curried functions
        hasEmailField := HasField("email")
        if hasEmailField(userField) {
            // process email field
        }

        searchInWorkspace := SearchUser("user@example.com")
        userID := searchInWorkspace(workspace)

        // Curry custom functions
        var assignedByIdHandler = F.Curry2(
            func(aid int, loader *DataLoader) *DataLoader {
                if aid != 0 {
                    loader.Payload.AssignedBy = []AssignedBy{{Id: aid}}
                }
                return loader
            })
        ```

      - **Array/Slice Operations**
        ```go
        // Use fp-go array functions for slice operations
        func processUserList(users []User) []string {
            return F.Pipe3(
                users,
                A.Filter(func(u User) bool { return u.Active }),
                A.Map(func(u User) string { return u.Email }),
                A.FilterMap(func(email string) O.Option[string] {
                    if isValidEmail(email) {
                        return O.Some(email)
                    }
                    return O.None[string]()
                }),
            )
        }

        // Find operations with Option return
        func findFirstAdmin(users []User) O.Option[User] {
            return F.Pipe1(
                users,
                A.FindFirst(func(u User) bool { return u.Role == "admin" }),
            )
        }

        // Complex array transformations
        func extractTextFromCells(cells []Cell, startCol, endCol int) string {
            return F.Pipe3(
                createCellRange(startCol, endCol, len(cells)),
                A.FilterMap(extractTextFromCellAtIndex(cells)),
                A.Head[string],
                O.GetOrElse(F.Constant("")),
            )
        }
        ```

      - **Error Handling Patterns**
        ```go
        // Combine Either with TryCatchError for exception handling
        func safeOperation(input string) E.Either[error, Result] {
            return F.Pipe2(
                E.TryCatchError(riskyOperation(input)),
                E.Map[error](func(rawResult RawResult) Result {
                    return processResult(rawResult)
                }),
                E.MapLeft[Result](func(err error) error {
                    return fmt.Errorf("operation failed: %w", err)
                }),
            )
        }

        // Chain multiple Either operations
        func processWorkflowWithValidation(input Input) E.Either[error, Output] {
            return F.Pipe4(
                E.Right[error](input),
                E.Bind(validateInput),
                E.Bind(processData),
                E.Bind(formatOutput),
                E.MapLeft[Output](func(err error) error {
                    return fmt.Errorf("workflow failed: %w", err)
                }),
            )
        }
        ```

      - **Tuple Operations**
        ```go
        // Use tuples for functions returning multiple values
        func extractTextWithIndex(cells []Cell, index int) O.Option[T.Tuple2[string, int]] {
            return F.Pipe1(
                extractTextFromSingleCell(cells[index]),
                O.Map(func(text string) T.Tuple2[string, int] {
                    return T.MakeTuple2(text, index)
                }),
            )
        }

        // Extract tuple values
        result := extractTextWithIndex(cells, 0)
        text, index := F.Pipe1(
            result,
            O.GetOrElse(F.Constant(T.MakeTuple2("", -1))),
        ).F1, result.F2
        ```
    - Use options pattern for configurable components
      ```go
      // Option type for functional options
      type Option func(*Config)
      
      // Config struct holds the configuration
      type Config struct {
          timeout     time.Duration
          retries     int
          debug       bool
          maxConnections int
      }
      
      // Option functions
      func WithTimeout(timeout time.Duration) Option {
          return func(c *Config) {
              c.timeout = timeout
          }
      }
      
      func WithRetries(retries int) Option {
          return func(c *Config) {
              c.retries = retries
          }
      }
      
      func WithDebug(debug bool) Option {
          return func(c *Config) {
              c.debug = debug
          }
      }
      
      func WithMaxConnections(max int) Option {
          return func(c *Config) {
              c.maxConnections = max
          }
      }
      
      // Constructor with default values and options
      func NewClient(opts ...Option) *Client {
          cfg := &Config{
              timeout:        30 * time.Second,
              retries:        3,
              debug:          false,
              maxConnections: 10,
          }
          
          for _, opt := range opts {
              opt(cfg)
          }
          
          return &Client{config: cfg}
      }
      
      // Usage examples
      client1 := NewClient() // Uses all defaults
      
      client2 := NewClient(
          WithTimeout(60*time.Second),
          WithRetries(5),
          WithDebug(true),
      )
      
      client3 := NewClient(WithMaxConnections(20))
      ```
    - Document all exported functions, types, and constants with proper Go doc comments
    - Test coverage should be comprehensive with both unit and integration tests
    - Use `github.com/stretchr/testify/require` for all unit test assertions
      ```go
      import (
          "testing"
          "github.com/stretchr/testify/require"
      )

      // Basic assertions
      func TestBasicAssertions(t *testing.T) {
          value := "hello world"
          number := 42
          result := []string{"a", "b", "c"}

          // String assertions
          require.Equal(t, "hello world", value)
          require.NotEqual(t, "goodbye", value)
          require.Contains(t, value, "world")
          require.NotEmpty(t, value)

          // Numeric assertions
          require.Equal(t, 42, number)
          require.Greater(t, number, 40)
          require.GreaterOrEqual(t, number, 42)
          require.Less(t, number, 50)
          require.LessOrEqual(t, number, 42)

          // Collection assertions
          require.Len(t, result, 3)
          require.Contains(t, result, "b")
          require.NotContains(t, result, "d")
          require.ElementsMatch(t, []string{"c", "a", "b"}, result) // Order independent
      }

      // Error handling tests
      func TestErrorHandling(t *testing.T) {
          // Function that should return error
          _, err := someFunction("invalid input")
          require.Error(t, err)
          require.ErrorContains(t, err, "invalid input")

          // Function that should succeed
          result, err := someFunction("valid input")
          require.NoError(t, err)
          require.NotNil(t, result)
      }

      // Nil/Not Nil assertions
      func TestNilAssertions(t *testing.T) {
          var ptr *string
          value := "test"

          require.Nil(t, ptr)
          require.NotNil(t, &value)
      }

      // Boolean assertions
      func TestBooleanAssertions(t *testing.T) {
          active := true
          disabled := false

          require.True(t, active)
          require.False(t, disabled)
      }

      // Type assertions
      func TestTypeAssertions(t *testing.T) {
          var value interface{} = "hello"

          require.IsType(t, "", value)
          require.Implements(t, (*io.Reader)(nil), &bytes.Buffer{})
      }

      // Testing with subtests
      func TestUserValidation(t *testing.T) {
          tests := []struct {
              name     string
              input    User
              expected error
          }{
              {
                  name:     "valid user",
                  input:    User{Name: "John", Email: "john@example.com"},
                  expected: nil,
              },
              {
                  name:     "missing name",
                  input:    User{Email: "john@example.com"},
                  expected: ErrMissingName,
              },
              {
                  name:     "invalid email",
                  input:    User{Name: "John", Email: "invalid"},
                  expected: ErrInvalidEmail,
              },
          }

          for _, tt := range tests {
              t.Run(tt.name, func(t *testing.T) {
                  err := ValidateUser(tt.input)
                  if tt.expected == nil {
                      require.NoError(t, err)
                  } else {
                      require.ErrorIs(t, err, tt.expected)
                  }
              })
          }
      }

      // Testing collections and maps
      func TestCollections(t *testing.T) {
          users := []User{
              {ID: 1, Name: "Alice"},
              {ID: 2, Name: "Bob"},
          }
          userMap := map[string]int{
              "alice": 1,
              "bob":   2,
          }

          // Slice assertions
          require.Len(t, users, 2)
          require.Equal(t, "Alice", users[0].Name)

          // Map assertions
          require.Len(t, userMap, 2)
          require.Equal(t, 1, userMap["alice"])
          require.Contains(t, userMap, "bob")
          require.NotContains(t, userMap, "charlie")
      }

      // Testing with require.Eventually for async operations
      func TestAsyncOperation(t *testing.T) {
          service := NewAsyncService()
          service.Start()

          // Wait up to 5 seconds for condition to be true
          require.Eventually(t, func() bool {
              return service.IsReady()
          }, 5*time.Second, 100*time.Millisecond)

          require.True(t, service.IsReady())
      }

      // Testing with custom error messages
      func TestWithCustomMessages(t *testing.T) {
          result := calculateSum([]int{1, 2, 3})
          require.Equal(t, 6, result, "Sum calculation should equal 6")

          user := findUser("unknown")
          require.Nil(t, user, "User should not be found for unknown ID")
      }

      // Testing HTTP handlers
      func TestHTTPHandler(t *testing.T) {
          req := httptest.NewRequest("GET", "/users/123", nil)
          recorder := httptest.NewRecorder()

          handler := NewUserHandler()
          handler.ServeHTTP(recorder, req)

          require.Equal(t, http.StatusOK, recorder.Code)
          require.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
          require.Contains(t, recorder.Body.String(), "user_id")
      }

      // Setup and teardown patterns
      func TestWithSetup(t *testing.T) {
          // Setup
          db := setupTestDatabase(t)
          defer cleanupDatabase(db)

          user := &User{Name: "Test User", Email: "test@example.com"}
          err := db.CreateUser(user)
          require.NoError(t, err)
          require.NotZero(t, user.ID)

          // Verify user was created
          retrievedUser, err := db.GetUser(user.ID)
          require.NoError(t, err)
          require.Equal(t, user.Name, retrievedUser.Name)
          require.Equal(t, user.Email, retrievedUser.Email)
      }

      // Helper function for test setup
      func setupTestDatabase(t *testing.T) *Database {
          db, err := NewDatabase("test_connection_string")
          require.NoError(t, err, "Failed to create test database")
          return db
      }

      func cleanupDatabase(db *Database) {
          if db != nil {
              db.Close()
          }
      }
      ```
    - Use go-playground/validator for struct field and parameter validation
      ```go
      import (
          "github.com/go-playground/validator/v10"
      )
      
      // Struct with validation tags
      type CreateUserRequest struct {
          Name     string `validate:"required,min=2,max=50" json:"name"`
          Email    string `validate:"required,email" json:"email"`
          Age      int    `validate:"gte=18,lte=120" json:"age"`
          Password string `validate:"required,min=8" json:"password"`
          Role     string `validate:"required,oneof=admin user guest" json:"role"`
          Website  string `validate:"omitempty,url" json:"website"`
      }
      
      // Global validator instance (thread-safe singleton)
      var validate = validator.New()
      
      // Validate struct fields
      func CreateUser(req CreateUserRequest) error {
          if err := validate.Struct(req); err != nil {
              return fmt.Errorf("validation failed: %w", err)
          }
          
          // Process valid request
          return nil
      }
      
      // Validate individual parameters
      func UpdateUserEmail(userID string, newEmail string) error {
          if err := validate.Var(newEmail, "required,email"); err != nil {
              return fmt.Errorf("invalid email: %w", err)
          }
          
          if err := validate.Var(userID, "required,uuid"); err != nil {
              return fmt.Errorf("invalid user ID: %w", err)
          }
          
          // Process update
          return nil
      }
      
      // Complex validation with nested structs
      type Address struct {
          Street  string `validate:"required,min=5"`
          City    string `validate:"required"`
          Country string `validate:"required,iso3166_1_alpha2"`
          ZipCode string `validate:"required,postcode_iso3166_alpha2=US"`
      }
      
      type UserProfile struct {
          User    CreateUserRequest `validate:"required"`
          Address Address          `validate:"required"`
          Tags    []string         `validate:"dive,required,min=2"`
      }
      
      // Custom validation function
      func validateBusinessEmail(fl validator.FieldLevel) bool {
          email := fl.Field().String()
          // Business emails should not be from common free providers
          blockedDomains := []string{"gmail.com", "yahoo.com", "hotmail.com"}
          
          for _, domain := range blockedDomains {
              if strings.HasSuffix(email, "@"+domain) {
                  return false
              }
          }
          return true
      }
      
      // Register custom validator
      func init() {
          validate.RegisterValidation("business_email", validateBusinessEmail)
      }
      
      // Usage with custom validator
      type BusinessUser struct {
          Email string `validate:"required,email,business_email"`
      }
      
      // Helper function to format validation errors
      func FormatValidationError(err error) string {
          if validationErrors, ok := err.(validator.ValidationErrors); ok {
              var messages []string
              for _, e := range validationErrors {
                  switch e.Tag() {
                  case "required":
                      messages = append(messages, fmt.Sprintf("%s is required", e.Field()))
                  case "email":
                      messages = append(messages, fmt.Sprintf("%s must be a valid email", e.Field()))
                  case "min":
                      messages = append(messages, fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param()))
                  case "max":
                      messages = append(messages, fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param()))
                  default:
                      messages = append(messages, fmt.Sprintf("%s failed validation: %s", e.Field(), e.Tag()))
                  }
              }
              return strings.Join(messages, "; ")
          }
          return err.Error()
      }
      ```
    - Any function or method receiving more than three parameters should use a type struct
      ```go
      // Avoid: Too many parameters
      func CreateUser(name, email, phone, address string, age int, active bool) error {
          // implementation
      }
      
      // Preferred: Use a struct for parameters
      type CreateUserParams struct {
          Name    string
          Email   string
          Phone   string
          Address string
          Age     int
          Active  bool
      }
      
      func CreateUser(params CreateUserParams) error {
          // implementation
      }
      
      // Usage
      err := CreateUser(CreateUserParams{
          Name:    "John Doe",
          Email:   "john@example.com",
          Phone:   "555-0123",
          Address: "123 Main St",
          Age:     30,
          Active:  true,
      })
      ```

- **Project Structure**
    - Primary interface definitions in package root
    - Implementations in subdirectories by backing technology

- **Variable Name Length:**
    -  Favor variable names that are at least three characters long, except for
    loop indices (e.g., `i`, `j`), method receivers (e.g., `r` for `receiver`),
    and extremely common types (e.g., `r` for `io.Reader`, `w` for `io.Writer`).
    -  Prioritize clarity and readability.  Use the shortest name that
    effectively conveys the variable's purpose within its context.
    - Variable naming: camelCase, descriptive names, no abbreviations except for
    common ones

- **Naming Style:**
    - Use `camelCase` for variable and function names (e.g., `myVariableName`, `calculateTotal`).
    - Use `PascalCase` for exported (public) types, functions, and constants (e.g., `MyType`, `CalculateTotal`).
    - Avoid `snake_case` (e.g., `my_variable_name`) in most cases.

- **Clarity and Context:**
    - The further a variable is used from its declaration, the more descriptive
    its name should be
      ```go
      func processData() {
          // Short scope: single letter acceptable
          for i := 0; i < 10; i++ {
              fmt.Println(i)
          }
          
          // Medium scope: short but descriptive
          users := fetchUsers()
          for _, user := range users {
              processUser(user)
          }
      }
      
      func longRunningFunction() {
          // Long scope: highly descriptive names
          authenticatedUserRepository := NewUserRepository()
          configurationManager := NewConfigManager()
          emailNotificationService := NewEmailService()
          
          // These variables are used throughout the function
          for page := 1; page <= totalPages; page++ {
              paginatedUserResults := authenticatedUserRepository.GetUsersByPage(page)
              
              for _, individualUser := range paginatedUserResults {
                  userEmailAddress := individualUser.Email
                  notificationPreferences := configurationManager.GetPreferences(individualUser.ID)
                  
                  if notificationPreferences.EmailEnabled {
                      emailNotificationService.Send(userEmailAddress, "Welcome!")
                  }
              }
          }
      }
      
      // Function parameters and package-level variables: descriptive
      func CalculateMonthlySubscriptionRevenue(subscriptionDetails []Subscription, 
                                               discountCalculator DiscountService) decimal.Decimal {
          totalMonthlyRevenue := decimal.Zero
          
          for _, subscription := range subscriptionDetails {
              monthlyAmount := subscription.MonthlyPrice
              applicableDiscount := discountCalculator.Calculate(subscription)
              finalAmount := monthlyAmount.Sub(applicableDiscount)
              totalMonthlyRevenue = totalMonthlyRevenue.Add(finalAmount)
          }
          
          return totalMonthlyRevenue
      }
      ```
    - Choose names that clearly indicate the variable's purpose and the type of
    data it holds.

- **Avoidance:**
    - Do not use spaces in variable names.
    - Variable names should start with a letter or underscore.
    - Do not use Go keywords as variable names.

- **Constants:**
    - Use `PascalCase` for constants. If a constant is unscoped, all letters in
    the constant should be capitalized. `const MAX_SIZE = 100`

- **Error Handling:**
    - When naming error variables, use `err` as the prefix:  `errMyCustomError`.
    - Always check errors and return meaningful wrapped errors
      ```go
      import (
          "fmt"
          "io"
          "os"
      )
      
      // Avoid: Ignoring errors
      func badExample() {
          file, _ := os.Open("config.txt")
          data, _ := io.ReadAll(file)
          fmt.Println(string(data))
      }
      
      // Preferred: Always check and wrap errors with context
      func goodExample() error {
          file, err := os.Open("config.txt")
          if err != nil {
              return fmt.Errorf("failed to open config file: %w", err)
          }
          defer file.Close()
          
          data, err := io.ReadAll(file)
          if err != nil {
              return fmt.Errorf("failed to read config file: %w", err)
          }
          
          fmt.Println(string(data))
          return nil
      }
      
      // Multiple operations: preserve error chain
      func processUserData(userID string) error {
          user, err := fetchUser(userID)
          if err != nil {
              return fmt.Errorf("failed to fetch user %s: %w", userID, err)
          }
          
          profile, err := loadProfile(user.ProfileID)
          if err != nil {
              return fmt.Errorf("failed to load profile for user %s: %w", userID, err)
          }
          
          err = validateProfile(profile)
          if err != nil {
              return fmt.Errorf("invalid profile for user %s: %w", userID, err)
          }
          
          err = saveProcessedData(user, profile)
          if err != nil {
              return fmt.Errorf("failed to save processed data for user %s: %w", userID, err)
          }
          
          return nil
      }
      
      // Custom error types for better error handling
      type ValidationError struct {
          Field   string
          Value   interface{}
          Message string
      }
      
      func (e ValidationError) Error() string {
          return fmt.Sprintf("validation failed for field '%s' with value '%v': %s", 
                           e.Field, e.Value, e.Message)
      }
      
      func validateEmail(email string) error {
          if email == "" {
              return ValidationError{
                  Field:   "email",
                  Value:   email,
                  Message: "email cannot be empty",
              }
          }
          
          if !strings.Contains(email, "@") {
              return fmt.Errorf("invalid email format: %w", ValidationError{
                  Field:   "email",
                  Value:   email,
                  Message: "must contain @ symbol",
              })
          }
          
          return nil
      }
      
      // Validation errors using go-playground/validator
      func CreateUserWithValidation(req CreateUserRequest) error {
          // First validate the struct
          if err := validate.Struct(req); err != nil {
              if validationErrors, ok := err.(validator.ValidationErrors); ok {
                  return fmt.Errorf("validation failed: %w", &UserValidationError{
                      Errors: validationErrors,
                  })
              }
              return fmt.Errorf("validation failed: %w", err)
          }
          
          // Additional business logic validation
          if strings.Contains(req.Email, "test") {
              return fmt.Errorf("test emails not allowed: %w", &BusinessLogicError{
                  Field:   "email",
                  Value:   req.Email,
                  Message: "production environment does not accept test emails",
              })
          }
          
          // Process valid request
          return nil
      }
      
      // Custom error types for validation
      type UserValidationError struct {
          Errors validator.ValidationErrors
      }
      
      func (e *UserValidationError) Error() string {
          var messages []string
          for _, err := range e.Errors {
              switch err.Tag() {
              case "required":
                  messages = append(messages, fmt.Sprintf("%s is required", err.Field()))
              case "email":
                  messages = append(messages, fmt.Sprintf("%s must be a valid email address", err.Field()))
              case "min":
                  messages = append(messages, fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param()))
              case "max":
                  messages = append(messages, fmt.Sprintf("%s must be at most %s characters long", err.Field(), err.Param()))
              case "gte":
                  messages = append(messages, fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param()))
              case "oneof":
                  messages = append(messages, fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param()))
              default:
                  messages = append(messages, fmt.Sprintf("%s failed validation rule: %s", err.Field(), err.Tag()))
              }
          }
          return strings.Join(messages, "; ")
      }
      
      type BusinessLogicError struct {
          Field   string
          Value   interface{}
          Message string
      }
      
      func (e *BusinessLogicError) Error() string {
          return fmt.Sprintf("business logic error for field '%s' with value '%v': %s", 
                           e.Field, e.Value, e.Message)
      }
      
      // Error checking with type assertions
      func HandleUserCreation(req CreateUserRequest) {
          err := CreateUserWithValidation(req)
          if err != nil {
              var validationErr *UserValidationError
              var businessErr *BusinessLogicError
              
              switch {
              case errors.As(err, &validationErr):
                  log.Printf("Validation errors: %s", validationErr.Error())
                  // Handle validation errors - return 400 Bad Request
              case errors.As(err, &businessErr):
                  log.Printf("Business logic error: %s", businessErr.Error())
                  // Handle business logic errors - return 422 Unprocessable Entity
              default:
                  log.Printf("Unknown error: %s", err.Error())
                  // Handle unknown errors - return 500 Internal Server Error
              }
          }
      }
      ```

- **Receivers:**
    - Use short, one or two-letter receiver names that reflect the type (e.g.,
    `r` for `io.Reader`, `f` for `*File`).


