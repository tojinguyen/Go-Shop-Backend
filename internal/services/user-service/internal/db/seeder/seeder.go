package seeder

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/constant"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

// Seeder đóng gói kết nối database và các queries.
type Seeder struct {
	db  *pgxpool.Pool
	ctx context.Context
	// [SỬA LỖI] Sử dụng atomic counter để đảm bảo tính duy nhất khi chạy song song trong tương lai
	emailCounter *uint64
	phoneCounter *uint64
}

// NewSeeder tạo một Seeder instance mới.
func NewSeeder(db *pgxpool.Pool) *Seeder {
	var initialEmail uint64 = 1
	var initialPhone uint64 = 900000000 // Bắt đầu từ số điện thoại 0900000000
	return &Seeder{
		db:           db,
		ctx:          context.Background(),
		emailCounter: &initialEmail,
		phoneCounter: &initialPhone,
	}
}

// SeedAllUserTypes tạo users với phân bố theo tỉ lệ thực tế e-commerce
func (s *Seeder) SeedAllUserTypes(totalUsers int) {
	log.Printf("--- Starting to seed %d users with realistic distribution (Optimized Version) ---", totalUsers)

	customerCount := int(float64(totalUsers) * 0.87)
	sellerCount := int(float64(totalUsers) * 0.10)
	shipperCount := int(float64(totalUsers) * 0.025)
	adminCount := totalUsers - customerCount - sellerCount - shipperCount

	log.Printf("Distribution: Customers=%d, Sellers=%d, Shippers=%d, Admins=%d",
		customerCount, sellerCount, shipperCount, adminCount)

	// Seed tất cả các loại user
	s.seedUserBatch("customers", customerCount, constant.UserRoleCustomer, true, false)
	s.seedUserBatch("sellers", sellerCount, constant.UserRoleSeller, true, false)
	s.seedUserBatch("shippers", shipperCount, constant.UserRoleShipper, true, true)
	s.seedUserBatch("admins", adminCount, constant.UserRoleAdmin, true, false)

	log.Printf("🎉 Seeding complete! Total users created: %d", totalUsers)
	s.PrintSeedingStatistics(totalUsers, customerCount, sellerCount, shipperCount, adminCount)
}

// SeedAll chạy các hàm seeding (legacy mode)
func (s *Seeder) SeedAll(userCount, shipperCount int) {
	log.Println("--- Starting to seed all data (Legacy Mode) ---")
	s.seedUserBatch("customers", userCount, constant.UserRoleCustomer, true, false)
	s.seedUserBatch("shippers", shipperCount, constant.UserRoleShipper, true, true)
	log.Println("--- Seeding complete ---")
}

// seedUserBatch là hàm chính để tạo user, profile, address, shipper profile theo lô
func (s *Seeder) seedUserBatch(typeName string, count int, role constant.UserRole, createAddress bool, createShipperProfile bool) {
	if count == 0 {
		return
	}
	log.Printf("🛒 Seeding %d %s...", count, typeName)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	const sampleSize = 200
	preGeneratedNames := make([]string, sampleSize)
	for i := 0; i < sampleSize; i++ {
		preGeneratedNames[i] = faker.Name()
	}

	const batchSize = 1000
	totalCreated := 0

	for i := 0; i < count; i += batchSize {
		batchEnd := i + batchSize
		if batchEnd > count {
			batchEnd = count
		}
		currentBatchSize := batchEnd - i

		accountsRows := make([][]interface{}, 0, currentBatchSize)
		profilesRows := make([][]interface{}, 0, currentBatchSize)
		addressesRows := make([][]interface{}, 0, currentBatchSize)
		shipperProfilesRows := make([][]interface{}, 0, currentBatchSize)

		for j := 0; j < currentBatchSize; j++ {
			userID := uuid.New()

			// [SỬA LỖI] Tạo email và phone tuần tự để đảm bảo UNIQUE
			emailNum := atomic.AddUint64(s.emailCounter, 1)
			phoneNum := atomic.AddUint64(s.phoneCounter, 1)
			email := fmt.Sprintf("user+%d@goshop.dev", emailNum)
			phone := fmt.Sprintf("0%d", phoneNum)

			fullName := preGeneratedNames[rand.Intn(sampleSize)]

			accountsRows = append(accountsRows, []interface{}{
				userID,
				email,
				string(hashedPassword),
				string(role),
			})

			birthYear := 1970 + rand.Intn(35)
			birthMonth := 1 + rand.Intn(12)
			birthDay := 1 + rand.Intn(28)
			birthday := time.Date(birthYear, time.Month(birthMonth), birthDay, 0, 0, 0, 0, time.UTC)
			genders := []string{"male", "female", "other"}

			profilesRows = append(profilesRows, []interface{}{
				userID,
				email,
				fullName,
				birthday,
				phone,
				string(role),
				nil, // banned_at
				fmt.Sprintf("https://i.pravatar.cc/150?u=%s", email),
				genders[rand.Intn(len(genders))],
			})

			if createAddress {
				numAddresses := rand.Intn(2) + 1
				for addrIdx := 0; addrIdx < numAddresses; addrIdx++ {
					addressesRows = append(addressesRows, []interface{}{
						uuid.New(),
						userID,
						addrIdx == 0,
						faker.GetRealAddress().Address,
						faker.GetRealAddress().State,
						faker.GetRealAddress().City,
					})
				}
			}

			if createShipperProfile {
				vehicleTypes := []string{"Motorbike", "Car"}
				shipperProfilesRows = append(shipperProfilesRows, []interface{}{
					userID,
					vehicleTypes[rand.Intn(len(vehicleTypes))],
					faker.CCNumber(),
				})
			}
		}

		tx, err := s.db.Begin(s.ctx)
		if err != nil {
			log.Printf("❌ Failed to begin transaction for batch %d-%d: %v", i+1, batchEnd, err)
			continue
		}

		_, err = tx.CopyFrom(s.ctx, pgx.Identifier{"user_accounts"}, []string{"id", "email", "hashed_password", "user_role"}, pgx.CopyFromRows(accountsRows))
		if err != nil {
			log.Printf("❌ Error inserting user_accounts batch %d-%d: %v", i+1, batchEnd, err)
			_ = tx.Rollback(s.ctx)
			continue
		}

		_, err = tx.CopyFrom(s.ctx, pgx.Identifier{"user_profiles"}, []string{"user_id", "email", "full_name", "birthday", "phone", "user_role", "banned_at", "avatar_url", "gender"}, pgx.CopyFromRows(profilesRows))
		if err != nil {
			log.Printf("❌ Error inserting user_profiles batch %d-%d: %v", i+1, batchEnd, err)
			_ = tx.Rollback(s.ctx)
			continue
		}

		if createAddress && len(addressesRows) > 0 {
			_, err = tx.CopyFrom(s.ctx, pgx.Identifier{"addresses"}, []string{"id", "user_id", "is_default", "street", "district", "city"}, pgx.CopyFromRows(addressesRows))
			if err != nil {
				log.Printf("❌ Error inserting addresses batch %d-%d: %v", i+1, batchEnd, err)
				_ = tx.Rollback(s.ctx)
				continue
			}
		}

		if createShipperProfile && len(shipperProfilesRows) > 0 {
			_, err = tx.CopyFrom(s.ctx, pgx.Identifier{"shipper_profiles"}, []string{"user_id", "vehicle_type", "license_plate"}, pgx.CopyFromRows(shipperProfilesRows))
			if err != nil {
				log.Printf("❌ Error inserting shipper_profiles batch %d-%d: %v", i+1, batchEnd, err)
				_ = tx.Rollback(s.ctx)
				continue
			}
		}

		if err := tx.Commit(s.ctx); err != nil {
			log.Printf("❌ Failed to commit transaction for batch %d-%d: %v", i+1, batchEnd, err)
			continue
		}

		totalCreated += currentBatchSize
		log.Printf("✅ Successfully seeded batch %d-%d for %s. Total seeded: %d/%d", i+1, batchEnd, typeName, totalCreated, count)
	}
}

// SeedUsers, SeedAddressesForUsers, SeedShipperProfiles giờ đây không còn cần thiết và có thể được xóa đi
// hoặc giữ lại như các wrapper rỗng để tránh lỗi biên dịch ở các chỗ gọi khác.

// SeedUsers (Hàm cũ) - Giờ đây chỉ là wrapper để tương thích.
func (s *Seeder) SeedUsers(count int, role constant.UserRole) ([]sqlc.CreateUserAccountRow, error) {
	s.seedUserBatch(string(role), count, role, true, role == constant.UserRoleShipper)
	return []sqlc.CreateUserAccountRow{}, nil
}

// SeedAddressesForUsers (Hàm cũ) - Không làm gì cả.
func (s *Seeder) SeedAddressesForUsers(users []sqlc.CreateUserAccountRow) {
	// Logic đã được tích hợp vào seedUserBatch.
}

// SeedShipperProfiles (Hàm cũ) - Không làm gì cả.
func (s *Seeder) SeedShipperProfiles(users []sqlc.CreateUserAccountRow) {
	// Logic đã được tích hợp vào seedUserBatch.
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
