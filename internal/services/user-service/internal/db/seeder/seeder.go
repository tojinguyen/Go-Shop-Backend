// internal/services/user-service/internal/db/seeder/seeder.go
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

// Seeder đóng gói kết nối database và các queries.
type Seeder struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
	ctx     context.Context
}

// NewSeeder tạo một Seeder instance mới.
func NewSeeder(db *pgxpool.Pool) *Seeder {
	return &Seeder{
		db:      db,
		queries: sqlc.New(db),
		ctx:     context.Background(),
	}
}

// SeedAll chạy tất cả các hàm seeding.
func (s *Seeder) SeedAll(userCount, shipperCount int) {
	log.Println("--- Starting to seed all data ---")

	// Tạo người dùng customer trước
	createdCustomers, err := s.SeedUsers(userCount, constant.UserRoleCustomer)
	if err != nil {
		log.Fatalf("Failed to seed customers: %v", err)
	}
	log.Printf("Successfully seeded %d customers.", len(createdCustomers))

	// Với mỗi customer, tạo 1-3 địa chỉ
	s.SeedAddressesForUsers(createdCustomers)

	// Tạo người dùng shipper
	createdShippers, err := s.SeedUsers(shipperCount, constant.UserRoleShipper)
	if err != nil {
		log.Fatalf("Failed to seed shippers: %v", err)
	}
	log.Printf("Successfully seeded %d shippers.", len(createdShippers))

	// Với mỗi shipper, tạo hồ sơ shipper
	s.SeedShipperProfiles(createdShippers)

	log.Println("--- Seeding complete ---")
}

// SeedUsers tạo một số lượng người dùng với vai trò cụ thể.
// Trả về một slice chứa các user accounts đã được tạo.
func (s *Seeder) SeedUsers(count int, role constant.UserRole) ([]sqlc.CreateUserAccountRow, error) {
	log.Printf("Seeding %d users with role '%s'...", count, role)
	var createdUsers []sqlc.CreateUserAccountRow

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	for i := 0; i < count; i++ {
		// 1. Tạo User Account
		userAccountParams := sqlc.CreateUserAccountParams{
			Email:          faker.Email(),
			HashedPassword: string(hashedPassword),
			UserRole:       string(role),
		}

		createdAccount, err := s.queries.CreateUserAccount(s.ctx, userAccountParams)
		if err != nil {
			log.Printf("Could not create user account (might be a duplicate email): %v. Skipping...", err)
			continue // Bỏ qua nếu email trùng
		}

		// 2. Tạo User Profile
		genders := []constant.UserGender{constant.UserGenderMale, constant.UserGenderFemale, constant.UserGenderOther}
		randomGender := genders[rand.Intn(len(genders))]

		userProfileParams := sqlc.CreateUserProfileParams{
			UserID:    createdAccount.ID,
			Email:     createdAccount.Email,
			FullName:  faker.Name(),
			Birthday:  converter.StringToPgDate("1995-05-20"), // faker.Date() có thể dùng nhưng cần format
			Phone:     faker.Phonenumber(),
			UserRole:  createdAccount.UserRole,
			AvatarUrl: fmt.Sprintf("https://i.pravatar.cc/150?u=%s", createdAccount.Email),
			Gender:    string(randomGender),
		}

		_, err = s.queries.CreateUserProfile(s.ctx, userProfileParams)
		if err != nil {
			// Nếu tạo profile lỗi, lý tưởng là rollback, nhưng với seeder thì báo lỗi là đủ
			log.Fatalf("failed to create user profile for user %s: %v", createdAccount.ID, err)
		}

		createdUsers = append(createdUsers, createdAccount)
		log.Printf("Successfully created user and profile for %s", createdAccount.Email)
	}

	return createdUsers, nil
}

// SeedAddressesForUsers tạo địa chỉ ngẫu nhiên cho danh sách người dùng.
func (s *Seeder) SeedAddressesForUsers(users []sqlc.CreateUserAccountRow) {
	log.Printf("Seeding addresses for %d users...", len(users))
	for _, user := range users {
		// Tạo 1 đến 3 địa chỉ cho mỗi user
		numAddresses := rand.Intn(3) + 1
		city := faker.GetRealAddress().City
		for i := 0; i < numAddresses; i++ {
			addressParams := sqlc.CreateAddressParams{
				UserID:    user.ID,
				Street:    faker.GetRealAddress().Address,
				City:      converter.StringToPgText(&city),
				IsDefault: converter.BoolToPgBool(i == 0), // Địa chỉ đầu tiên là mặc định
			}
			_, err := s.queries.CreateAddress(s.ctx, addressParams)
			if err != nil {
				log.Printf("Failed to create address for user %s: %v", user.ID, err)
			}
		}
		log.Printf("Created %d addresses for user %s", numAddresses, user.Email)
	}
}

// SeedShipperProfiles tạo hồ sơ shipper cho danh sách người dùng.
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
			IdentifyCardUrl: converter.StringToPgText(nil), // Có thể thêm link ảnh giả
			VehicleImageUrl: converter.StringToPgText(nil), // Có thể thêm link ảnh giả
		}

		_, err := s.queries.CreateShipper(s.ctx, shipperParams)
		if err != nil {
			log.Printf("Failed to create shipper profile for user %s: %v", user.ID, err)
		} else {
			log.Printf("Created shipper profile for user %s", user.Email)
		}
	}
}
