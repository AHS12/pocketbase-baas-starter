package seeders

import (
	"fmt"

	"ims-pocketbase-baas-starter/internal/database/factories"
	"ims-pocketbase-baas-starter/pkg/logger"

	"github.com/pocketbase/pocketbase/core"
)

// SeedUsersCLI seeds a specified number of fake users
func SeedUsersCLI(app core.App, count int) error {
	log := logger.FromApp(app)

	fmt.Printf("🌱 Seeding %d users...\\n", count)

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

	fmt.Printf("✅ Successfully seeded %d users\\n", count)

	log.Info("Successfully seeded users", "count", count)
	return nil
}

// SeedUsersWithRoleCLI seeds users and assigns them to a specific role
func SeedUsersWithRoleCLI(app core.App, count int, roleName string) error {
	log := logger.FromApp(app)

	fmt.Printf("🌱 Seeding %d users with role '%s'...\\n", count, roleName)

	log.Info("Seeding users with role", "count", count, "role", roleName)

	role, err := app.FindFirstRecordByFilter("roles", "name = {:name}", map[string]any{"name": roleName})
	if err != nil {
		return fmt.Errorf("failed to find role '%s': %w", roleName, err)
	}

	userFactory := factories.NewUserFactory(app)
	users, err := userFactory.GenerateMany(count)
	if err != nil {
		return fmt.Errorf("failed to generate users: %w", err)
	}

	for i, user := range users {
		user.Set("roles", []string{role.Id})

		if err := app.Save(user); err != nil {
			return fmt.Errorf("failed to save user %d: %w", i+1, err)
		}
		log.Info("Created user with role", "name", user.GetString("name"), "email", user.GetString("email"), "role", roleName)
	}

	fmt.Printf("✅ Successfully seeded %d users with role '%s'\\n", count, roleName)

	log.Info("Successfully seeded users with role", "count", count, "role", roleName)
	return nil
}
