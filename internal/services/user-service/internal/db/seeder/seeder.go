package seeder

import (
	"context"
	"fmt"
	"log"

	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/toji-dev/go-shop/internal/pkg/converter"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/constant"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

// SeedUsers creates a specified number of fake users and their profiles.
func SeedUsers(db *pgxpool.Pool, count int) {
	log.Printf("Seeding %d users...", count)
	queries := sqlc.New(db)
	ctx := context.Background()

	for i := 0; i < count; i++ {
		// 1. Create a User Account
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("failed to hash password: %v", err)
		}

		userAccountParams := sqlc.CreateUserAccountParams{
			Email:          faker.Email(),
			HashedPassword: string(hashedPassword),
			UserRole:       string(constant.UserRoleCustomer), // Default to customer
		}

		createdAccount, err := queries.CreateUserAccount(ctx, userAccountParams)
		if err != nil {
			// If email is a duplicate, just try again with a new fake email.
			log.Printf("Could not create user account (might be a duplicate email): %v. Retrying...", err)
			i-- // Decrement i to retry this iteration.
			continue
		}
		log.Printf("Created UserAccount ID: %s", createdAccount.ID)

		// 2. Create a User Profile linked to the account
		userProfileParams := sqlc.CreateUserProfileParams{
			UserID:    createdAccount.ID,
			Email:     createdAccount.Email,
			FullName:  faker.Name(),
			Birthday:  converter.StringToPgDate("1995-05-20"), // A static birthday for simplicity
			Phone:     faker.Phonenumber(),
			UserRole:  createdAccount.UserRole,
			AvatarUrl: fmt.Sprintf("https://i.pravatar.cc/150?u=%s", createdAccount.Email),
			Gender:    string(constant.UserGenderOther),
		}

		_, err = queries.CreateUserProfile(ctx, userProfileParams)
		if err != nil {
			log.Fatalf("failed to create user profile for user %s: %v", createdAccount.ID, err)
		}
		log.Printf("Created UserProfile for User ID: %s", createdAccount.ID)
	}

	log.Printf("Successfully seeded %d users.", count)
}
