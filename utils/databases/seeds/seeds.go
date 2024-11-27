package seeds

import (
	"fmt"
	"greenenvironment/features/admin"
	"greenenvironment/utils/databases/seed"

	"gorm.io/gorm"
)

func seeds() []seed.Seed {
	var seeds []seed.Seed = []seed.Seed{
		{
			Name: "CreateAdmin01",
			Run: func(db *gorm.DB) error {
				return CreateAdminLogin(db, admin.Admin{
					ID:       "16930c07-bdb5-49d2-8a81-32591833241b",
					Name:     "admin",
					Username: "admin",
					Email:    "admin@ecomate.store",
					Password: "admin",
				})
			},
		},
		{
			Name: "CreateAdmin02",
			Run: func(db *gorm.DB) error {
				return CreateAdminLogin(db, admin.Admin{
					ID:       "14adafd7-de6c-4586-a35e-3cf17ef3d351",
					Name:     "admin2",
					Username: "admin2",
					Email:    "admin2@ecomate.store",
					Password: "admin2",
				})
			},
		},
	}
	return seeds
}

func RunSeeds(db *gorm.DB) error {
	for _, s := range seeds() {
		fmt.Printf("Running seed: %s\n", s.Name)
		if err := s.Run(db); err != nil {
			fmt.Printf("Failed to run seed %s: %v\n", s.Name, err)
			return err
		}
		fmt.Printf("Seed %s completed successfully\n", s.Name)
	}
	return nil
}
