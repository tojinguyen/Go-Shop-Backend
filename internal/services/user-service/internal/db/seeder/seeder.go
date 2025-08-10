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
	db           *pgxpool.Pool
	queries      *sqlc.Queries
	ctx          context.Context
	usedPhones   map[string]bool
	phoneCounter int64
}

// NewSeeder tạo một Seeder instance mới.
func NewSeeder(db *pgxpool.Pool) *Seeder {
	return &Seeder{
		db:           db,
		queries:      sqlc.New(db),
		ctx:          context.Background(),
		usedPhones:   make(map[string]bool),
		phoneCounter: 1000000000, // Bắt đầu từ số điện thoại 10 chữ số
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

// SeedAllUserTypes tạo 50,000 users với phân bố theo tỉ lệ thực tế e-commerce
func (s *Seeder) SeedAllUserTypes(totalUsers int) {
	log.Printf("--- Starting to seed %d users with realistic distribution ---", totalUsers)

	// Phân bố tỉ lệ theo e-commerce thực tế:
	// Customer: 87% (~43,500 users)
	// Seller: 10% (~5,000 users)
	// Shipper: 2.5% (~1,250 users)
	// Admin: 0.5% (~250 users)

	customerCount := int(float64(totalUsers) * 0.87)
	sellerCount := int(float64(totalUsers) * 0.10)
	shipperCount := int(float64(totalUsers) * 0.025)
	adminCount := totalUsers - customerCount - sellerCount - shipperCount // Đảm bảo tổng đúng

	log.Printf("Distribution: Customers=%d, Sellers=%d, Shippers=%d, Admins=%d",
		customerCount, sellerCount, shipperCount, adminCount)

	// 1. Seed Customers (87%)
	log.Println("🛒 Seeding customers...")
	createdCustomers, err := s.SeedUsers(customerCount, constant.UserRoleCustomer)
	if err != nil {
		log.Fatalf("Failed to seed customers: %v", err)
	}
	log.Printf("✅ Successfully seeded %d customers.", len(createdCustomers))

	// Tạo địa chỉ cho customers
	s.SeedAddressesForUsers(createdCustomers)

	// 2. Seed Sellers (10%)
	log.Println("🏪 Seeding sellers...")
	createdSellers, err := s.SeedUsers(sellerCount, constant.UserRoleSeller)
	if err != nil {
		log.Fatalf("Failed to seed sellers: %v", err)
	}
	log.Printf("✅ Successfully seeded %d sellers.", len(createdSellers))

	// Tạo địa chỉ cho sellers
	s.SeedAddressesForUsers(createdSellers)

	// 3. Seed Shippers (2.5%)
	log.Println("🚚 Seeding shippers...")
	createdShippers, err := s.SeedUsers(shipperCount, constant.UserRoleShipper)
	if err != nil {
		log.Fatalf("Failed to seed shippers: %v", err)
	}
	log.Printf("✅ Successfully seeded %d shippers.", len(createdShippers))

	// Tạo shipper profiles và địa chỉ
	s.SeedShipperProfiles(createdShippers)
	s.SeedAddressesForUsers(createdShippers)

	// 4. Seed Admins (0.5%)
	log.Println("👑 Seeding admins...")
	createdAdmins, err := s.SeedUsers(adminCount, constant.UserRoleAdmin)
	if err != nil {
		log.Fatalf("Failed to seed admins: %v", err)
	}
	log.Printf("✅ Successfully seeded %d admins.", len(createdAdmins))

	// Tạo địa chỉ cho admins
	s.SeedAddressesForUsers(createdAdmins)

	totalCreated := len(createdCustomers) + len(createdSellers) + len(createdShippers) + len(createdAdmins)
	log.Printf("🎉 Seeding complete! Total users created: %d", totalCreated)

	// In thống kê
	s.PrintSeedingStatistics(totalCreated, len(createdCustomers), len(createdSellers), len(createdShippers), len(createdAdmins))
}

// PrintSeedingStatistics in ra thống kê sau khi seed
func (s *Seeder) PrintSeedingStatistics(total, customers, sellers, shippers, admins int) {
	log.Println("--- SEEDING STATISTICS ---")
	log.Printf("Total Users: %d", total)
	log.Printf("Customers: %d (%.1f%%)", customers, float64(customers)/float64(total)*100)
	log.Printf("Sellers: %d (%.1f%%)", sellers, float64(sellers)/float64(total)*100)
	log.Printf("Shippers: %d (%.1f%%)", shippers, float64(shippers)/float64(total)*100)
	log.Printf("Admins: %d (%.1f%%)", admins, float64(admins)/float64(total)*100)
	log.Println("-------------------------")
}

// generateUniquePhone generates a unique phone number for seeding
func (s *Seeder) generateUniquePhone() string {
	for {
		// Generate a Vietnamese phone number format: +84xxxxxxxxx (11 digits total)
		var phone string

		// Try faker first for variety
		if rand.Float32() < 0.3 { // 30% chance to use faker
			phone = faker.Phonenumber()
			// Clean up the phone number - remove non-digits and ensure it's reasonable
			cleanPhone := ""
			for _, char := range phone {
				if char >= '0' && char <= '9' {
					cleanPhone += string(char)
				}
			}
			// Ensure it's a proper Vietnamese phone format
			if len(cleanPhone) >= 9 && len(cleanPhone) <= 12 {
				if len(cleanPhone) == 9 {
					phone = "+84" + cleanPhone
				} else if len(cleanPhone) == 10 && cleanPhone[0] == '0' {
					phone = "+84" + cleanPhone[1:]
				} else if len(cleanPhone) == 12 && cleanPhone[:2] == "84" {
					phone = "+" + cleanPhone
				} else {
					phone = fmt.Sprintf("+84%d", s.phoneCounter)
					s.phoneCounter++
				}
			} else {
				phone = fmt.Sprintf("+84%d", s.phoneCounter)
				s.phoneCounter++
			}
		} else {
			// Generate sequential Vietnamese phone number
			phone = fmt.Sprintf("+84%d", s.phoneCounter)
			s.phoneCounter++
		}

		// Check if this phone number is already used
		if !s.usedPhones[phone] {
			s.usedPhones[phone] = true
			return phone
		}

		// If faker generated a duplicate, fall back to sequential
		phone = fmt.Sprintf("+84%d", s.phoneCounter)
		s.phoneCounter++
		if !s.usedPhones[phone] {
			s.usedPhones[phone] = true
			return phone
		}
	}
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

	// Batch processing để tối ưu hiệu suất với số lượng lớn
	batchSize := 1000
	if count < batchSize {
		batchSize = count
	}

	for i := 0; i < count; i += batchSize {
		end := i + batchSize
		if end > count {
			end = count
		}

		log.Printf("Processing batch %d-%d for role %s...", i+1, end, role)

		for j := i; j < end; j++ {
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

			// Tạo birthday ngẫu nhiên từ 1970-2005
			birthYear := 1970 + rand.Intn(35)
			birthMonth := 1 + rand.Intn(12)
			birthDay := 1 + rand.Intn(28)
			birthday := fmt.Sprintf("%d-%02d-%02d", birthYear, birthMonth, birthDay)

			userProfileParams := sqlc.CreateUserProfileParams{
				UserID:    createdAccount.ID,
				Email:     createdAccount.Email,
				FullName:  faker.Name(),
				Birthday:  converter.StringToPgDate(birthday),
				Phone:     s.generateUniquePhone(),
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
		}

		// Log tiến độ
		if count > 1000 {
			progress := float64(end) / float64(count) * 100
			log.Printf("Progress for %s: %.1f%% (%d/%d)", role, progress, end, count)
		}
	}

	log.Printf("Successfully created %d users with role %s", len(createdUsers), role)
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
