package seeders

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

// CLISeederFunction represents a CLI seeder function with its metadata
type CLISeederFunction struct {
	Name        string
	Description string
	Function    func(core.App) error
}

// GetAllCLISeederFunctions returns a list of all CLI seeder functions
func GetAllCLISeederFunctions() []CLISeederFunction {
	return []CLISeederFunction{
		{
			Name:        "UserSeeder[10]",
			Description: "Seeds 10 fake users",
			Function: func(app core.App) error {
				return SeedUsersCLI(app, 10)
			},
		},
		{
			Name:        "UserWithRoleSeeder[Admin][5]",
			Description: "Seeds 5 fake users with Admin role",
			Function: func(app core.App) error {
				return SeedUsersWithRoleCLI(app, 5, "Admin")
			},
		},
		{
			Name:        "UserWithRoleSeeder[User][10]",
			Description: "Seeds 10 fake users with User role",
			Function: func(app core.App) error {
				return SeedUsersWithRoleCLI(app, 10, "User")
			},
		},
	}
}

// RunAllCLISeederFunctions runs all registered CLI seeder functions
func RunAllCLISeederFunctions(app core.App) error {
	functions := GetAllCLISeederFunctions()

	if len(functions) == 0 {
		fmt.Println("No CLI seeder functions registered")
		return nil
	}

	fmt.Printf("🌱 Running %d CLI seeder functions...\n", len(functions))

	for _, fn := range functions {
		fmt.Printf("▶️  Running %s...\n", fn.Name)
		fmt.Printf("   %s\n", fn.Description)

		if err := fn.Function(app); err != nil {
			return fmt.Errorf("failed to run seeder %s: %w", fn.Name, err)
		}

		fmt.Printf("✅ Completed %s\n", fn.Name)
	}

	fmt.Printf("✅ Successfully ran %d CLI seeder functions\n", len(functions))
	return nil
}
