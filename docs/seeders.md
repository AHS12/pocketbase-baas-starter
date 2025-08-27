# Database Seeders Guide

This document explains how to create and register new database seeders in the IMS PocketBase BaaS Starter.

## Overview

The database seeder system allows you to populate your database with test data or initial data for development and production environments. Seeders are organized into two categories:

1. **Migration Seeders** - Run automatically during database migrations (for essential data)
2. **CLI Seeders** - Run manually via CLI commands (for test data and development)

This guide focuses on **CLI Seeders** which are designed for development and testing purposes.

## CLI Seeder Architecture

CLI Seeders follow a simple manual registration pattern:

- **Seeder Functions** - Individual functions that perform the seeding logic
- **Registry** - Manual list of seeder functions in `internal/database/seeders/cli_seeder_registry.go`
- **CLI Commands** - Commands to run individual or all seeders

## Creating a New CLI Seeder

### Step 1: Create the Seeder Function

Add your seeder function to an existing file (like `user_seeder.go`) or create a new seeder file:

```go
// internal/database/seeders/user_seeder.go

// SeedUsersCLI seeds a specified number of fake users
func SeedUsersCLI(app core.App, count int) error {
    log := logger.FromApp(app)
    
    fmt.Printf("🌱 Seeding %d users...\n", count)
    
    log.Info("Seeding users", "count", count)

    userFactory := factories.NewUserFactory(app)

    users, err := userFactory.GenerateMany(count)
    if err != nil {
        return fmt.Errorf("failed to generate users: %w", err)
    }

    for i, user := range users {
        if err := app.Save(user); err != nil {
            return fmt.Errorf("failed to save user %d: %w", i+1, err)
        }
        log.Info("Created user", "name", user.GetString("name"), "email", user.GetString("email"))
    }

    fmt.Printf("✅ Successfully seeded %d users\n", count)
    
    log.Info("Successfully seeded users", "count", count)
    return nil
}
```

### Step 2: Register the Seeder Function

Add your seeder function to the registry in `internal/database/seeders/cli_seeder_registry.go`:

```go
// GetAllCLISeederFunctions returns a list of all CLI seeder functions
func GetAllCLISeederFunctions() []CLISeederFunction {
    return []CLISeederFunction{
        // Existing seeders...
        {
            Name:        "UserSeeder[10]",
            Description: "Seeds 10 fake users",
            Function: func(app core.App) error {
                return SeedUsersCLI(app, 10)
            },
        },
        // Add your new seeder here:
        {
            Name:        "CustomSeeder[50]",
            Description: "Seeds 50 custom records",
            Function: func(app core.App) error {
                return SeedCustomRecords(app, 50)
            },
        },
    }
}
```

### Step 3: (Optional) Create a Dedicated CLI Command

If you want to run your seeder individually, create a command handler in `internal/handlers/command/`:

```go
// internal/handlers/command/custom_seeder_command.go

// HandleCustomSeederCommand handles the 'seed-custom' CLI command
func HandleCustomSeederCommand(app *pocketbase.PocketBase, cmd *cobra.Command, args []string) {
    log := logger.GetLogger(app)
    
    log.Info("Starting custom seeder process")
    fmt.Println("🌱 Starting custom seeder process...")
    
    // Call the seeder function
    if err := seeders.SeedCustomRecords(app, 50); err != nil {
        log.Error("Failed to seed custom records", "error", err)
        fmt.Printf("❌ Error seeding custom records: %v\n", err)
        return
    }
    
    log.Info("Custom seeder process completed successfully")
    fmt.Println("✅ Custom seeder completed successfully")
}
```

Then register it in `internal/commands/commands.go`. For detailed instructions on creating custom CLI commands, see the [CLI Commands Guide](cli-commands.md).

## Running Seeders

### Run All Registered Seeders

```bash
./main db-seed
```

This command runs all functions registered in `GetAllCLISeederFunctions()`.

### Run Individual Seeders

If you created a dedicated command:

```bash
./main seed-custom
```

## Best Practices

1. **Use Descriptive Names**: Give your seeders clear, descriptive names
2. **Provide Useful Descriptions**: Help users understand what each seeder does
3. **Handle Errors Gracefully**: Return meaningful error messages
4. **Use the Logger**: Log important events with `logger.FromApp(app)`
5. **Console Output**: Provide clear console feedback with emojis and formatting
6. **Make Seeders Idempotent**: Design seeders so they can be run multiple times safely
7. **Use Factories**: Leverage existing factories for generating fake data

## Example: Complete Custom Seeder

Here's a complete example of a custom seeder:

```go
// internal/database/seeders/custom_seeder.go

// SeedProductsCLI seeds a specified number of fake products
func SeedProductsCLI(app core.App, count int) error {
    log := logger.FromApp(app)
    
    fmt.Printf("🌱 Seeding %d products...\n", count)
    
    log.Info("Seeding products", "count", count)

    // Get products collection
    productsCollection, err := app.FindCollectionByNameOrId("products")
    if err != nil {
        return fmt.Errorf("failed to find products collection: %w", err)
    }

    // Generate fake products
    for i := 0; i < count; i++ {
        product := core.NewRecord(productsCollection)
        product.Set("name", faker.Word())
        product.Set("price", faker.Price(10, 1000))
        product.Set("description", faker.Sentence())
        product.Set("in_stock", faker.Bool())
        
        if err := app.Save(product); err != nil {
            return fmt.Errorf("failed to save product %d: %w", i+1, err)
        }
        
        log.Info("Created product", "name", product.GetString("name"))
    }

    fmt.Printf("✅ Successfully seeded %d products\n", count)
    
    log.Info("Successfully seeded products", "count", count)
    return nil
}
```

Register it in the registry:

```go
{
    Name:        "ProductSeeder[25]",
    Description: "Seeds 25 fake products",
    Function: func(app core.App) error {
        return SeedProductsCLI(app, 25)
    },
},
```

Now when you run `./main db-seed`, your new product seeder will be executed along with all other registered seeders.