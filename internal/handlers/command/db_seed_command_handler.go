package command

import (
	"fmt"

	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
	"ims-pocketbase-baas-starter/internal/database/seeders"
	"ims-pocketbase-baas-starter/pkg/logger"
)

func HandleDBSeedCommand(app *pocketbase.PocketBase, cmd *cobra.Command, args []string) {
	log := logger.GetLogger(app)

	log.Info("Starting database seeding process")
	fmt.Println("🌱 Starting database seeding process...")

	// Run all registered CLI seeder functions
	if err := seeders.RunAllCLISeederFunctions(app); err != nil {
		log.Error("Failed to seed database", "error", err)
		fmt.Printf("❌ Error seeding database: %v\n", err)
		return
	}

	log.Info("Database seeding process completed successfully")
	fmt.Println("✅ Database seeding completed successfully")
}
