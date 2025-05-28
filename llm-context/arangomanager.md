# ArangoManager - LLM Reference Documentation

This document provides comprehensive technical reference for the ArangoManager library, designed for Large Language Models (LLMs) to understand and work with this ArangoDB Go client library.

## Overview

ArangoManager is a Go library that provides high-level abstractions and utilities for working with ArangoDB, a multi-model NoSQL database. The library focuses on simplifying database operations, transaction management, and query building while maintaining type safety and functional programming patterns.

## Core Architecture

### Package Structure
```
arangomanager/
├── collection/          # Functional programming utilities for collections
├── command/flag/        # CLI flag definitions for ArangoDB connections
├── query/              # AQL query building and parsing utilities
├── testarango/         # Testing utilities for isolated test databases
├── database.go         # Main database operations
├── session.go          # Connection and session management
├── transaction.go      # Transaction handling with ACID compliance
├── result.go           # Single-row query result handling
├── resultset.go        # Multi-row query result handling
└── datasource.go       # Connection parameter definitions
```

### Key Types and Their Relationships

```
Session
├── Creates and manages → Database
└── Manages → Connection to ArangoDB

Database
├── Executes queries → Result/Resultset
├── Manages → Collections
├── Creates → TransactionHandler
└── Handles → Indexes and Graphs

TransactionHandler
├── Provides ACID → Transaction operations
├── Has context for → Query execution
└── Can be → Committed/Aborted

Result/Resultset
├── Wraps → ArangoDB cursor
├── Provides → Data reading interface
└── Manages → Resource cleanup
```

## Connection Management

### ConnectParams Structure
```go
type ConnectParams struct {
    User     string `validate:"required"`      // Database user
    Pass     string `validate:"required"`      // User password
    Database string `validate:"required"`      // Target database name
    Host     string `validate:"required"`      // ArangoDB host
    Port     int    `validate:"required"`      // ArangoDB port (default: 8529)
    Istls    bool                             // TLS encryption flag
}
```

### Session Management
The `Session` type manages the underlying ArangoDB client connection:
- `Connect()`: Creates a new session with connection parameters
- `NewSessionDb()`: Creates both session and database instances
- `CreateDB()`: Creates new databases
- `DB()`: Retrieves database instances
- `CreateUser()`: Manages user accounts
- `GrantDB()`: Handles database permissions

## Database Operations

### Core Database Interface
The `Database` type provides the main interface for ArangoDB operations:

#### Query Execution Methods
- `SearchRows(query, bindVars)`: Multi-row queries returning `*Resultset`
- `GetRow(query, bindVars)`: Single-row queries returning `*Result`
- `CountWithParams(query, bindVars)`: Count queries returning `int64`
- `Do(query, bindVars)`: Write operations with no return data
- `DoRun(query, bindVars)`: Write operations with return data

#### Collection Management
- `Collection(name)`: Retrieves existing collections
- `CreateCollection(name, opts)`: Creates new collections
- `FindOrCreateCollection(name, opts)`: Idempotent collection creation
- `Truncate(names...)`: Removes all data from collections

#### Index Management
- `EnsureGeoIndex()`: Creates/finds geo-spatial indexes
- `EnsureHashIndex()`: Creates/finds hash indexes
- `EnsurePersistentIndex()`: Creates/finds persistent indexes
- `EnsureSkipListIndex()`: Creates/finds skip-list indexes

#### Graph Operations
- `FindOrCreateGraph(name, defs)`: Manages named graphs

## Transaction System

### TransactionOptions Configuration
```go
type TransactionOptions struct {
    ReadCollections      []string  // Collections for read access
    WriteCollections     []string  // Collections for write access
    ExclusiveCollections []string  // Collections for exclusive access
    WaitForSync         bool      // Force disk sync before return
    AllowImplicit       bool      // Allow undeclared collection access
    LockTimeout         int       // Collection lock timeout (seconds)
    MaxTransactionSize  int       // Maximum transaction size (bytes)
}
```

### TransactionHandler Interface
The `TransactionHandler` provides ACID-compliant transaction operations:
- `BeginTransaction(ctx, opts)`: Starts a new transaction
- `Context()`: Returns transaction-bound context
- `Do(query, bindVars)`: Executes write operations in transaction
- `DoRun(query, bindVars)`: Executes queries returning data in transaction
- `Commit()`: Commits the transaction
- `Abort()`: Rolls back the transaction
- `Status()`: Retrieves transaction status information

### Transaction Usage Pattern
```go
// 1. Configure transaction options
txOpts := &arangomanager.TransactionOptions{
    WriteCollections: []string{"users", "orders"},
    ReadCollections:  []string{"products"},
}

// 2. Begin transaction
tx, err := db.BeginTransaction(context.Background(), txOpts)

// 3. Execute operations
err = tx.Do("INSERT {...} INTO users", bindVars)
result, err := tx.DoRun("FOR u IN users RETURN u", nil)

// 4. Commit or abort
err = tx.Commit() // or tx.Abort()
```

## Result Handling

### Result (Single Row)
The `Result` type handles single-row query results:
- `IsEmpty()`: Checks if result contains data
- `Read(interface{})`: Unmarshals data into provided struct/variable

### Resultset (Multiple Rows)
The `Resultset` type handles multi-row query results:
- `IsEmpty()`: Checks if resultset contains data
- `Scan()`: Advances to next row (returns bool for iteration)
- `Read(interface{})`: Reads current row data
- `Close()`: Releases resources (important for cleanup)

### Result Iteration Pattern
```go
rs, err := db.SearchRows(query, bindVars)
if err != nil {
    return err
}
defer rs.Close() // Always close to free resources

if rs.IsEmpty() {
    return nil // Handle empty results
}

for rs.Scan() {
    var item MyStruct
    if err := rs.Read(&item); err != nil {
        return err
    }
    // Process item
}
```

## Query Package

### Filter Parsing and AQL Generation
The `query` package provides sophisticated filter parsing and AQL statement generation:

#### Filter Structure
```go
type Filter struct {
    Field    string  // Database field name
    Operator string  // Comparison operator
    Value    string  // Comparison value
    Logic    string  // Logical connector ("," for OR, ";" for AND)
}
```

#### Supported Operators
- **Standard**: `==`, `!=`, `>`, `<`, `>=`, `<=`, `=~`, `!~`
- **Date** (prefix `$`): `$==`, `$>`, `$<`, `$>=`, `$<=`
- **Array** (prefix `@`): `@==`, `@!=`, `@=~`, `@!~`

#### Key Functions
- `ParseFilterString(filterStr)`: Parses filter strings into `Filter` slices
- `GenQualifiedAQLFilterStatement(fieldMap, filters)`: Generates AQL with qualified field names
- `GenAQLFilterStatement(params)`: Generates AQL with statement parameters

#### Filter String Format
```
field operator value[logic]field operator value...
```
Where:
- `logic` is `,` for OR operations, `;` for AND operations
- Example: `status==active;created_at$>=2023-01-01,priority@==high`

### StatementParameters Structure
```go
type StatementParameters struct {
    Fmap    map[string]string  // Field name to database path mapping
    Filters []*Filter          // Filter conditions
    Doc     string            // Document variable name (FOR doc IN...)
    Vert    string            // Vertex variable name (graph queries)
}
```

## Collection Package (Functional Utilities)

### Core Functions
The `collection` package provides functional programming utilities:

#### Transformation Functions
- `Map[T1, T2](slice, func)`: Transforms slice elements
- `Filter[T](slice, predicate)`: Filters elements by predicate
- `Partition[T](slice, predicate)`: Splits slice into two based on predicate

#### Utility Functions
- `Include[T](slice, element)`: Checks element presence (sorted search)
- `RemoveStringItems(slice, items...)`: Removes specified string items
- `IsEmpty[T](slice)`: Checks if slice is empty

#### Functional Composition
- `Pipe2/3/4[T...](initial, f1, f2, ...)`: Creates function pipelines
- `CurriedMap/Filter/Partition[T](func)`: Returns curried functions

#### Tuple Types
```go
type Tuple2[T1, T2 any] struct {
    First  T1
    Second T2
}
```

### Iterator Support
- `MapSeq[T1, T2](seq, fn)`: Transforms iterator sequences
- Support for Go 1.23+ iterator patterns

## Testing with TestArango

### TestArango Structure
```go
type TestArango struct {
    *arangomanager.ConnectParams  // Connection parameters
    *arangomanager.Session        // Database session
}
```

### Key Functions
- `NewTestArangoFromEnv(isCreate)`: Creates test instance from environment
- `NewTestArango(user, pass, host, port, isCreate)`: Creates with explicit params
- `CreateTestDb(name, opts)`: Creates additional test databases

### Environment Variables
- `ARANGO_HOST`: ArangoDB server host
- `ARANGO_USER`: Database username
- `ARANGO_PASS`: Database password
- `ARANGO_PORT`: Server port (optional, defaults to 8529)

### Testing Pattern
```go
func TestMain(m *testing.M) {
    ta, err := testarango.NewTestArangoFromEnv(true)
    if err != nil {
        log.Fatal(err)
    }
    
    code := m.Run()
    
    // Cleanup
    db, _ := ta.DB(ta.Database)
    db.Drop()
    os.Exit(code)
}
```

## Command Line Integration

### Flag Package
The `command/flag` package provides pre-configured CLI flags for urfave/cli:

#### Available Flag Sets
- `ArangoFlags()`: Basic connection flags (host, port, user, pass, secure)
- `ArangodbFlags()`: Extended flags including database name

#### Environment Variable Support
All flags support environment variable overrides:
- `ARANGODB_PASS` → `--arangodb-pass`
- `ARANGODB_USER` → `--arangodb-user`
- `ARANGODB_DATABASE` → `--arangodb-database`
- `ARANGODB_SERVICE_HOST` → `--arangodb-host`
- `ARANGODB_SERVICE_PORT` → `--arangodb-port`

## Error Handling Patterns

### Validation
The library uses `github.com/go-playground/validator/v10` for parameter validation:
- Connection parameters are validated on session creation
- Query parameters support custom validation rules
- Filter operators are validated against supported operator maps

### Common Error Scenarios
1. **Connection Errors**: Invalid credentials, unreachable host
2. **Query Errors**: Syntax errors, invalid AQL
3. **Transaction Errors**: Deadlocks, timeout, invalid operations
4. **Resource Errors**: Collection not found, insufficient permissions

### Error Handling Best Practices
```go
// Always check errors from database operations
result, err := db.GetRow(query, bindVars)
if err != nil {
    return fmt.Errorf("query failed: %w", err)
}

// Handle empty results appropriately
if result.IsEmpty() {
    return ErrNotFound // or appropriate handling
}

// Always close resources
defer resultset.Close()
```

## Performance Considerations

### Connection Pooling
- Sessions maintain persistent connections
- Reuse session instances across operations
- Avoid creating new sessions for each operation

### Transaction Sizing
- Use `MaxTransactionSize` to limit memory usage
- Break large operations into smaller transactions
- Consider using `WaitForSync` for critical data

### Query Optimization
- Use bind parameters to prevent AQL injection
- Leverage indexes for filter operations
- Use the query package for consistent AQL generation

### Resource Management
- Always close `Resultset` instances
- Use context cancellation for long-running operations
- Monitor transaction lifetimes to prevent deadlocks

## Integration Examples

### Web Service Integration
```go
type UserService struct {
    db *arangomanager.Database
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    txOpts := &arangomanager.TransactionOptions{
        WriteCollections: []string{"users"},
    }
    
    tx, err := s.db.BeginTransaction(ctx, txOpts)
    if err != nil {
        return err
    }
    defer tx.Abort() // Safe to call multiple times
    
    query := "INSERT @user INTO users RETURN NEW"
    bindVars := map[string]interface{}{"user": user}
    
    result, err := tx.DoRun(query, bindVars)
    if err != nil {
        return err
    }
    
    if err := result.Read(user); err != nil {
        return err
    }
    
    return tx.Commit()
}
```

### Batch Processing
```go
func ProcessBatch(db *arangomanager.Database, items []Item) error {
    txOpts := &arangomanager.TransactionOptions{
        WriteCollections: []string{"items", "audit_log"},
        MaxTransactionSize: 50 * 1024 * 1024, // 50MB limit
    }
    
    tx, err := db.BeginTransaction(context.Background(), txOpts)
    if err != nil {
        return err
    }
    defer tx.Abort()
    
    for _, item := range items {
        query := "INSERT @item INTO items"
        bindVars := map[string]interface{}{"item": item}
        
        if err := tx.Do(query, bindVars); err != nil {
            return fmt.Errorf("failed to insert item %v: %w", item.ID, err)
        }
    }
    
    // Log batch completion
    logQuery := "INSERT {batch_size: @size, timestamp: @ts} INTO audit_log"
    logVars := map[string]interface{}{
        "size": len(items),
        "ts":   time.Now(),
    }
    
    if err := tx.Do(logQuery, logVars); err != nil {
        return err
    }
    
    return tx.Commit()
}
```

## Migration and Schema Management

While ArangoManager doesn't provide explicit migration utilities, the following patterns are recommended:

### Collection Setup
```go
func SetupCollections(db *arangomanager.Database) error {
    collections := []struct {
        name string
        opts *driver.CreateCollectionOptions
    }{
        {"users", &driver.CreateCollectionOptions{Type: driver.CollectionTypeDocument}},
        {"edges", &driver.CreateCollectionOptions{Type: driver.CollectionTypeEdge}},
    }
    
    for _, coll := range collections {
        _, err := db.FindOrCreateCollection(coll.name, coll.opts)
        if err != nil {
            return fmt.Errorf("failed to create collection %s: %w", coll.name, err)
        }
    }
    
    return nil
}
```

### Index Management
```go
func EnsureIndexes(db *arangomanager.Database) error {
    indexes := []struct {
        collection string
        fields     []string
        indexType  string
    }{
        {"users", []string{"email"}, "hash"},
        {"users", []string{"created_at"}, "skiplist"},
        {"locations", []string{"coordinates"}, "geo"},
    }
    
    for _, idx := range indexes {
        switch idx.indexType {
        case "hash":
            _, _, err := db.EnsureHashIndex(idx.collection, idx.fields, nil)
        case "skiplist":
            _, _, err := db.EnsureSkipListIndex(idx.collection, idx.fields, nil)
        case "geo":
            _, _, err := db.EnsureGeoIndex(idx.collection, idx.fields, nil)
        }
        
        if err != nil {
            return fmt.Errorf("failed to create %s index on %s: %w", 
                idx.indexType, idx.collection, err)
        }
    }
    
    return nil
}
```

This documentation provides a comprehensive technical reference for LLMs to understand and work effectively with the ArangoManager library, covering all major components, patterns, and best practices.