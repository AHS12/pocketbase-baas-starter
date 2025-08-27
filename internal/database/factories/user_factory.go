package factories

import (
	"fmt"

	"github.com/go-faker/faker/v4"
	"github.com/pocketbase/pocketbase/core"
)

// UserFactory generates fake user data
type UserFactory struct {
	app core.App
}

// NewUserFactory creates a new instance of UserFactory
func NewUserFactory(app core.App) *UserFactory {
	return &UserFactory{app: app}
}

// Generate creates a single fake user record without saving it
func (f *UserFactory) Generate() (*core.Record, error) {
	usersCollection, err := f.app.FindCollectionByNameOrId("users")
	if err != nil {
		return nil, fmt.Errorf("failed to find users collection: %w", err)
	}

	email := faker.Email()
	name := faker.Name()
	password := faker.Password()

	record := core.NewRecord(usersCollection)
	record.Set("email", email)
	record.Set("name", name)
	record.Set("verified", true)  // Default to verified for test data
	record.Set("is_active", true) // Default to active
	record.SetPassword(password)

	return record, nil
}

// GenerateMany creates multiple fake user records without saving them
func (f *UserFactory) GenerateMany(count int) ([]*core.Record, error) {
	records := make([]*core.Record, count)
	for i := 0; i < count; i++ {
		record, err := f.Generate()
		if err != nil {
			return nil, fmt.Errorf("failed to generate user %d: %w", i+1, err)
		}
		records[i] = record
	}
	return records, nil
}
