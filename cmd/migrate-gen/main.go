package main

import (
	"fmt"
	"ims-pocketbase-baas-starter/pkg/migration"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := ParseArgs(os.Args)
	if err != nil {
		ShowUsage()
		return err
	}

	migrations, err := migration.ScanExistingMigrations(config.OutputDir)
	if err != nil {
		return err
	}

	nextNumber := migration.GetNextMigrationNumber(migrations)
	template := CreateMigrationTemplate(nextNumber, config.MigrationName)
	content, err := GenerateMigrationContent(template)
	if err != nil {
		return err
	}

	migrationPath := migration.GenerateMigrationFilePath(nextNumber, config.MigrationName)
	schemaPath := migration.GenerateSchemaFilePath(nextNumber)

	if err := migration.WriteMigrationFile(migrationPath, content); err != nil {
		return err
	}

	fmt.Printf("✓ Generated migration file: %s\n", migrationPath)
	fmt.Printf("✓ Schema file expected at: %s\n", schemaPath)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Design your collections in PocketBase Admin UI")
	fmt.Printf("2. Export collections to %s\n", schemaPath)
	fmt.Println("3. Update the rollback function with collection names to delete")
	fmt.Println("4. Test the migration in development environment")

	return nil
}
