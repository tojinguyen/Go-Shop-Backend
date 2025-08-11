package seeder

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/toji-dev/go-shop/internal/pkg/converter"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/db/sqlc"
)

type Seeder struct {
	shopDB     *pgxpool.Pool
	userDB     *pgxpool.Pool
	queries    *sqlc.Queries
	ctx        context.Context
	usedEmails map[string]bool
}

func NewSeeder(shopDB, userDB *pgxpool.Pool) *Seeder {
	return &Seeder{
		shopDB:     shopDB,
		userDB:     userDB,
		queries:    sqlc.New(shopDB),
		ctx:        context.Background(),
		usedEmails: make(map[string]bool),
	}
}

// fetchAndPrepareSellers kiểm tra và "nâng cấp" người dùng thành seller nếu cần.
func (s *Seeder) fetchAndPrepareSellers(requiredCount int) ([]uuid.UUID, error) {
	log.Println("🔍 Checking for available sellers in user-service database...")

	// 1. Lấy danh sách sellers hiện có
	rows, err := s.userDB.Query(s.ctx, "SELECT id FROM user_accounts WHERE user_role = 'seller'")
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to query existing sellers: %w", err)
	}
	defer rows.Close()

	var existingSellerIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan seller ID: %w", err)
		}
		existingSellerIDs = append(existingSellerIDs, id)
	}
	log.Printf("Found %d existing sellers.", len(existingSellerIDs))

	// 2. Kiểm tra nếu đã đủ số lượng
	if len(existingSellerIDs) >= requiredCount {
		return existingSellerIDs, nil
	}

	// 3. Nếu chưa đủ, tìm và nâng cấp user `customer`
	needed := requiredCount - len(existingSellerIDs)
	log.Printf("⚠️ Not enough sellers. Need to promote %d more users from 'customer' to 'seller'.", needed)

	// Lấy ngẫu nhiên các user `customer`
	promoteRows, err := s.userDB.Query(s.ctx, "SELECT id FROM user_accounts WHERE user_role = 'customer' ORDER BY RANDOM() LIMIT $1", needed)
	if err != nil {
		return nil, fmt.Errorf("failed to find customers to promote: %w", err)
	}
	defer promoteRows.Close()

	var usersToPromoteIDs []uuid.UUID
	for promoteRows.Next() {
		var id uuid.UUID
		if err := promoteRows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan customer ID for promotion: %w", err)
		}
		usersToPromoteIDs = append(usersToPromoteIDs, id)
	}

	if len(usersToPromoteIDs) < needed {
		log.Printf("WARNING: Could only find %d customers to promote, but needed %d.", len(usersToPromoteIDs), needed)
		if len(usersToPromoteIDs) == 0 && len(existingSellerIDs) == 0 {
			return nil, fmt.Errorf("no customers available to promote to seller. Please run user-service seeder first")
		}
	}

	// Nâng cấp vai trò của họ thành 'seller'
	if len(usersToPromoteIDs) > 0 {
		_, err = s.userDB.Exec(s.ctx, "UPDATE user_accounts SET user_role = 'seller' WHERE id = ANY($1)", usersToPromoteIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to update user roles to seller: %w", err)
		}
		log.Printf("✅ Successfully promoted %d users to 'seller'.", len(usersToPromoteIDs))
	}

	// Gộp danh sách seller cũ và mới
	allSellerIDs := append(existingSellerIDs, usersToPromoteIDs...)
	return allSellerIDs, nil
}

// SeedShops tạo dữ liệu giả cho các cửa hàng
func (s *Seeder) SeedShops(count int) {
	// Sử dụng hàm mới để đảm bảo có đủ sellers
	sellerIDs, err := s.fetchAndPrepareSellers(count)
	if err != nil {
		log.Fatalf("❌ Could not prepare sellers: %v", err)
	}
	if len(sellerIDs) == 0 {
		log.Println("⚠️ No sellers available to create shops. Aborting.")
		return
	}

	log.Printf("🌱 Seeding %d shops...", count)

	for i := 0; i < count; i++ {
		// Bắt đầu transaction
		tx, err := s.shopDB.Begin(s.ctx)
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
			continue
		}
		qtx := s.queries.WithTx(tx)

		// Chuẩn bị dữ liệu
		shopID := uuid.New()
		addressID := uuid.New()
		ownerID := sellerIDs[rand.Intn(len(sellerIDs))] // Chọn ngẫu nhiên một seller

		// Tạo email duy nhất
		var shopEmail string
		for {
			shopEmail = faker.Email()
			if !s.usedEmails[shopEmail] {
				s.usedEmails[shopEmail] = true
				break
			}
		}

		address := faker.GetRealAddress().City
		// 1. Tạo bản ghi Address trước
		addressParams := sqlc.CreateShopAddressParams{
			ID:      converter.UUIDToPgUUID(addressID),
			ShopID:  converter.UUIDToPgUUID(shopID),
			Street:  faker.GetRealAddress().Address,
			City:    converter.StringToPgText(&address),
			Country: converter.StringToPgText(nil), // Mặc định là 'Vietnam' trong schema
		}
		_, err = qtx.CreateShopAddress(s.ctx, addressParams)
		if err != nil {
			log.Printf("Failed to create shop address: %v. Rolling back.", err)
			tx.Rollback(s.ctx)
			continue
		}

		// 2. Tạo bản ghi Shop
		shopDescription := faker.Sentence()
		shopParams := sqlc.CreateShopParams{
			ID:              converter.UUIDToPgUUID(shopID),
			OwnerID:         converter.UUIDToPgUUID(ownerID),
			ShopName:        faker.Sentence(),
			AvatarUrl:       "https://placehold.co/150x150/e8117f/ffffff.png?text=Shop",
			BannerUrl:       "https://placehold.co/800x200/333333/ffffff.png?text=Welcome",
			ShopDescription: converter.StringToPgText(&shopDescription),
			AddressID:       converter.UUIDToPgUUID(addressID),
			Phone:           faker.Phonenumber(),
			Email:           shopEmail,
		}

		_, err = qtx.CreateShop(s.ctx, shopParams)
		if err != nil {
			log.Printf("Failed to create shop: %v. Rolling back.", err)
			tx.Rollback(s.ctx)
			continue
		}

		// Commit transaction nếu mọi thứ thành công
		if err := tx.Commit(s.ctx); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
		} else {
			log.Printf("✅ Created shop '%s' for owner %s", shopParams.ShopName, ownerID)
		}
	}
	log.Println("🎉 Shop seeding complete.")
}
