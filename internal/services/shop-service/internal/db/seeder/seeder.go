package seeder

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
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

// fetchSellerIDs lấy danh sách ID của các user có vai trò 'seller' từ DB của user-service
func (s *Seeder) fetchSellerIDs() ([]uuid.UUID, error) {
	log.Println("Fetching seller IDs from user-service database...")
	rows, err := s.userDB.Query(s.ctx, "SELECT id FROM user_accounts WHERE user_role = 'seller'")
	if err != nil {
		return nil, fmt.Errorf("failed to query seller IDs: %w", err)
	}
	defer rows.Close()

	var sellerIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan seller ID: %w", err)
		}
		sellerIDs = append(sellerIDs, id)
	}
	log.Printf("Found %d sellers.", len(sellerIDs))
	return sellerIDs, nil
}

// SeedShops tạo dữ liệu giả cho các cửa hàng
func (s *Seeder) SeedShops(count int) {
	sellerIDs, err := s.fetchSellerIDs()
	if err != nil {
		log.Fatalf("❌ Could not fetch sellers: %v", err)
	}
	if len(sellerIDs) == 0 {
		log.Println("⚠️ No sellers found in user-service DB. Please seed users with 'seller' role first.")
		log.Println("Run: make seed-users")
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
			Country: converter.StringToPgText(nil),
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
