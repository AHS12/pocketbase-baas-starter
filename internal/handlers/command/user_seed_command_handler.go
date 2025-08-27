package command

import (
	"fmt"
	"strconv"

	"ims-pocketbase-baas-starter/internal/database/seeders"
	"ims-pocketbase-baas-starter/pkg/logger"

	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

// HandleSeedUsersCommand handles the 'seed-users' CLI command
func HandleSeedUsersCommand(app *pocketbase.PocketBase, cmd *cobra.Command, args []string) {
	log := logger.GetLogger(app)

	count := 10 // Default number of users to seed

	if len(args) > 0 {
		parsedCount, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error("Invalid count argument, using default count of 10", "error", err)
			fmt.Println("Invalid count argument, using default count of 10")
		} else {
			count = parsedCount
		}
	}

	log.Info("Starting user seeding process", "count", count)
	fmt.Printf("Starting user seeding process with %d users...\n", count)

	if err := seeders.SeedUsersCLI(app, count); err != nil {
		log.Error("❌ Error seeding users", "error", err)
		fmt.Printf("❌ Error seeding users: %v\n", err)
		return
	}

	log.Info("✅ User seeding completed successfully")
	fmt.Println("✅ User seeding completed successfully")
}

// HandleSeedUsersWithRoleCommand handles the 'seed-users-with-role' CLI command
func HandleSeedUsersWithRoleCommand(app *pocketbase.PocketBase, cmd *cobra.Command, args []string) {
	log := logger.GetLogger(app)

	if len(args) < 2 {
		log.Error("❌ Usage: seed-users-with-role <count> <role-name>")
		fmt.Println("❌ Usage: seed-users-with-role <count> <role-name>")
		return
	}

	count, err := strconv.Atoi(args[0])
	if err != nil {
		log.Error("❌ Invalid count argument", "error", err)
		fmt.Printf("❌ Invalid count argument: %v\n", err)
		return
	}

	roleName := args[1]

	log.Info("Starting user seeding process with role", "count", count, "role", roleName)
	fmt.Printf("Starting user seeding process with %d users and role '%s'...\n", count, roleName)

	if err := seeders.SeedUsersWithRoleCLI(app, count, roleName); err != nil {
		log.Error("❌ Error seeding users with role", "error", err)
		fmt.Printf("❌ Error seeding users with role: %v\n", err)
		return
	}

	log.Info("✅ User seeding with role completed successfully")
	fmt.Println("✅ User seeding with role completed successfully")
}
