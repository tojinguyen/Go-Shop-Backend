// internal/services/user-service/db/seeder/seeder.go
package seeder

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/toji-dev/go-shop/internal/pkg/converter"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/constant"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

// Seeder encapsulates the database connection and queries.
type Seeder struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
	ctx     context.Context
}

// NewSeeder creates a new Seeder instance.
func NewSeeder(db *pgxpool.Pool) *Seeder {
	return &Seeder{
		db:      db,
		queries: sqlc.New(db),
		ctx:     context.Background(),
	}
}

// SeedAll runs all seeding functions.
func (s *Seeder) SeedAll(userCount, shipperCount int) {
	log.Println("--- Starting to seed all data ---")

	// Create customer users first
	createdCustomers, err := s.SeedUsers(userCount, constant.UserRoleCustomer)
	if err != nil {
		log.Fatalf("Failed to seed customers: %v", err)
	}
	log.Printf("Successfully seeded %d customers.", len(createdCustomers))

	// For each customer, create 1-3 addresses
	s.SeedAddressesForUsers(createdCustomers)

	// Create shipper users
	createdShippers, err := s.SeedUsers(shipperCount, constant.UserRoleShipper)
	if err != nil {
		log.Fatalf("Failed to seed shippers: %v", err)
	}
	log.Printf("Successfully seeded %d shippers.", len(createdShippers))

	// For each shipper, create a shipper profile
	s.SeedShipperProfiles(createdShippers)

	log.Println("--- Seeding complete ---")
}

// SeedUsers creates a specified number of users with a given role.
// It returns a slice of the created user accounts.
func (s *Seeder) SeedUsers(count int, role constant.UserRole) ([]sqlc.CreateUserAccountRow, error) {
	log.Printf("Seeding %d users with role '%s'...", count, role)
	var createdUsers []sqlc.CreateUserAccountRow

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	for i := 0; i < count; i++ {
		// 1. Create a User Account
		userAccountParams := sqlc.CreateUserAccountParams{
			Email:          faker.Email(),
			HashedPassword: string(hashedPassword),
			UserRole:       string(role),
		}

		createdAccount, err := s.queries.CreateUserAccount(s.ctx, userAccountParams)
		if err != nil {
			log.Printf("Could not create user account (might be a duplicate email): %v. Skipping...", err)
			continue // Skip this user if creation fails
		}

		// 2. Create a User Profile
		genders := []constant.UserGender{constant.UserGenderMale, constant.UserGenderFemale, constant.UserGenderOther}
		randomGender := genders[rand.Intn(len(genders))]

		userProfileParams := sqlc.CreateUserProfileParams{
			UserID:    createdAccount.ID,
			Email:     createdAccount.Email,
			FullName:  faker.Name(),
			Birthday:  converter.StringToPgDate("1995-05-20"),
			Phone:     faker.Phonenumber(),
			UserRole:  createdAccount.UserRole,
			AvatarUrl: fmt.Sprintf("https://i.pravatar.cc/150?u=%s", createdAccount.Email),
			Gender:    string(randomGender),
		}

		_, err = s.queries.CreateUserProfile(s.ctx, userProfileParams)
		if err != nil {
			// If profile creation fails, we should ideally roll back the account creation.
			// For a seeder, we can just log a fatal error.
			log.Fatalf("failed to create user profile for user %s: %v", createdAccount.ID, err)
		}

		createdUsers = append(createdUsers, createdAccount)
		log.Printf("Successfully created user and profile for %s", createdAccount.Email)
	}

	return createdUsers, nil
}

// SeedAddressesForUsers creates random addresses for a given list of users.
func (s *Seeder) SeedAddressesForUsers(users []sqlc.CreateUserAccountRow) {
	log.Printf("Seeding addresses for %d users...", len(users))
	for _, user := range users {
		// Create 1 to 3 addresses for each user
		numAddresses := rand.Intn(3) + 1
		city := faker.GetRealAddress().City
		for i := 0; i < numAddresses; i++ {
			addressParams := sqlc.CreateAddressParams{
				UserID:    user.ID,
				Street:    faker.GetRealAddress().Address,
				City:      converter.StringToPgText(&city),
				IsDefault: converter.BoolToPgBool(i == 0), // Set the first address as default
			}
			_, err := s.queries.CreateAddress(s.ctx, addressParams)
			if err != nil {
				log.Printf("Failed to create address for user %s: %v", user.ID, err)
			}
		}
		log.Printf("Created %d addresses for user %s", numAddresses, user.Email)
	}
}

// SeedShipperProfiles creates shipper profiles for a given list of users.
func (s *Seeder) SeedShipperProfiles(users []sqlc.CreateUserAccountRow) {
	log.Printf("Seeding shipper profiles for %d users...", len(users))
	for _, user := range users {
		vehicleTypes := []string{"Motorbike", "Car", "Truck"}
		vehicleType := vehicleTypes[rand.Intn(len(vehicleTypes))]
		licensePlate := faker.CCNumber()
		shipperParams := sqlc.CreateShipperParams{
			UserID:          user.ID,
			VehicleType:     converter.StringToPgText(&vehicleType),
			LicensePlate:    converter.StringToPgText(&licensePlate),
			IdentifyCardUrl: converter.StringToPgText(nil),
			VehicleImageUrl: converter.StringToPgText(nil),
		}

		_, err := s.queries.CreateShipper(s.ctx, shipperParams)
		if err != nil {
			log.Printf("Failed to create shipper profile for user %s: %v", user.ID, err)
		} else {
			log.Printf("Created shipper profile for user %s", user.Email)
		}
	}
}
