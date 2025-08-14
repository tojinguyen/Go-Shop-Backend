package services

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/toji-dev/go-shop/internal/pkg/constant"
	"github.com/toji-dev/go-shop/internal/pkg/converter"
	redis_infra "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/repository"
)

type UserService struct {
	userProfileRepo repository.UserProfileRepository
	redisService    redis_infra.RedisServiceInterface
	// In-memory cache để giảm Redis network calls
	memoryCache sync.Map
	// Mutex để tránh cache stampede
	cacheMutex sync.RWMutex
	// Cache metrics
	cacheHits   int64
	cacheMisses int64
}

// CachedProfile struct để lưu trữ profile với timestamp
type CachedProfile struct {
	Profile   domain.UserProfile
	CachedAt  time.Time
	ExpiresAt time.Time
}

func NewUserService(
	userProfileRepo repository.UserProfileRepository,
	redisService redis_infra.RedisServiceInterface,
) *UserService {
	return &UserService{
		userProfileRepo: userProfileRepo,
		redisService:    redisService,
		memoryCache:     sync.Map{},
		cacheMutex:      sync.RWMutex{},
		cacheHits:       0,
		cacheMisses:     0,
	}
}

func (s *UserService) CreateProfile(ctx *gin.Context, req dto.CreateUserProfileRequest) (domain.UserProfile, error) {
	userIDRaw, exists := ctx.Get(constant.ContextKeyUserID)
	if !exists {
		log.Printf("user ID not found in context")
		return domain.UserProfile{}, fmt.Errorf("user ID not found in context")
	}
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		log.Printf("user ID is not a string")
		return domain.UserProfile{}, fmt.Errorf("user ID is not a string")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("invalid user ID format: %v", err)
		return domain.UserProfile{}, fmt.Errorf("invalid user ID format: %w", err)
	}

	roleRaw, exists := ctx.Get(constant.ContextKeyUserRole)
	if !exists {
		log.Printf("user role not found in context")
		return domain.UserProfile{}, fmt.Errorf("user role not found in context")
	}
	role, ok := roleRaw.(string)
	if !ok {
		log.Printf("user role is not a string")
		return domain.UserProfile{}, fmt.Errorf("user role is not a string")
	}

	email, exists := ctx.Get(constant.ContextKeyUserEmail)
	if !exists {
		log.Printf("email not found in context")
		return domain.UserProfile{}, fmt.Errorf("email not found in context: %w", err)
	}
	emailStr, ok := email.(string)
	if !ok {
		log.Printf("email is not a string")
		return domain.UserProfile{}, fmt.Errorf("email is not a string")
	}

	params := sqlc.CreateUserProfileParams{
		UserID:           converter.UUIDToPgUUID(userID),
		Email:            emailStr,
		FullName:         req.FullName,
		Birthday:         converter.StringToPgDate(req.Birthday),
		Phone:            req.Phone,
		UserRole:         role,
		BannedAt:         converter.NullPgTime(),
		AvatarUrl:        req.AvatarURL,
		DefaultAddressID: converter.NullPgUUID(),
		Gender:           req.Gender,
	}

	profile, err := s.userProfileRepo.CreateUserProfile(ctx, params)
	if err != nil {
		log.Printf("failed to create user profile: %v", err)
		return domain.UserProfile{}, err
	}
	return *profile, nil
}

func (s *UserService) GetProfile(ctx *gin.Context, userID string) (domain.UserProfile, error) {
	// Validate UUID format trước
	_, errParse := uuid.Parse(userID)
	if errParse != nil {
		return domain.UserProfile{}, fmt.Errorf("invalid user ID format")
	}

	cacheKey := fmt.Sprintf("user_profile:%s", userID)
	// 1. Kiểm tra cache trước
	var cachedProfile domain.UserProfile
	err := s.redisService.GetJSON(cacheKey, &cachedProfile)
	if err == nil {
		log.Printf("Cache HIT for user ID: %s", userID)
		return cachedProfile, nil
	}

	log.Printf("Cache MISS for user ID: %s. Fetching from DB.", userID)

	profile, err := s.userProfileRepo.GetUserProfileByID(ctx, userID)
	if err != nil {
		return domain.UserProfile{}, err
	}

	// Cache với TTL 15 phút
	err = s.redisService.SetJSON(cacheKey, profile, 15*time.Minute)
	if err != nil {
		// Ghi log lỗi cache nhưng không làm hỏng request
		log.Printf("Warning: Failed to set cache for user ID %s: %v", userID, err)
	}
	return *profile, nil
}

func (s *UserService) UpdateProfile(ctx *gin.Context, req dto.UpdateUserProfileRequest) (domain.UserProfile, error) {
	userIDRaw, exists := ctx.Get(constant.ContextKeyUserID)
	if !exists {
		return domain.UserProfile{}, fmt.Errorf("user ID not found in context")
	}
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		return domain.UserProfile{}, fmt.Errorf("user ID is not a string")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Get current profile to preserve existing values
	currentProfile, err := s.userProfileRepo.GetUserProfileByID(ctx, userIDStr)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("failed to get current profile: %w", err)
	}

	// Build update parameters
	params := sqlc.UpdateUserProfileParams{
		UserID:           converter.UUIDToPgUUID(userID),
		Email:            currentProfile.Email, // Keep existing email
		FullName:         currentProfile.FullName,
		Birthday:         converter.StringToPgDate(currentProfile.Birthday),
		Phone:            currentProfile.Phone,
		UserRole:         currentProfile.Role,
		BannedAt:         converter.StringToPgTime(currentProfile.BannedAt),
		AvatarUrl:        currentProfile.AvatarURL,
		Gender:           currentProfile.Gender,
		DefaultAddressID: converter.StringToPgUUID(currentProfile.DefaultAddressID),
	}

	// Update only provided fields
	if req.FullName != "" {
		params.FullName = req.FullName
	}

	if req.Phone != "" {
		params.Phone = req.Phone
	}

	if req.AvatarURL != "" {
		params.AvatarUrl = req.AvatarURL
	}

	if req.Birthday != "" {
		params.Birthday = converter.StringToPgDate(req.Birthday)
	}

	profile, err := s.userProfileRepo.UpdateUserProfile(ctx, params)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("failed to update profile: %w", err)
	}

	// Invalidate cache sau khi update
	cacheKey := fmt.Sprintf("user_profile:%s", userIDStr)
	err = s.redisService.Delete(cacheKey)
	if err != nil {
		log.Printf("Warning: Failed to invalidate Redis cache for user ID %s: %v", userIDStr, err)
	}

	// Invalidate memory cache
	s.memoryCache.Delete(userIDStr)
	log.Printf("Invalidated both Redis and memory cache for user ID: %s", userIDStr)

	return *profile, nil
}

func (s *UserService) DeleteProfile(ctx *gin.Context, userID string) error {
	// Verify that the user exists before attempting to delete
	_, err := s.userProfileRepo.GetUserProfileByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Perform soft delete
	err = s.userProfileRepo.DeleteProfile(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}
	return nil
}

func (s *UserService) GetProfileByID(ctx *gin.Context, userID string) (domain.UserProfile, error) {
	// Validate UUID format trước để tránh cache pollution
	_, errParse := uuid.Parse(userID)
	if errParse != nil {
		return domain.UserProfile{}, fmt.Errorf("invalid user ID format")
	}

	// 1. Kiểm tra in-memory cache trước (nhanh nhất)
	if cached, ok := s.memoryCache.Load(userID); ok {
		cachedProfile := cached.(CachedProfile)
		if time.Now().Before(cachedProfile.ExpiresAt) {
			s.cacheHits++
			log.Printf("Memory Cache HIT for user ID: %s", userID)
			return cachedProfile.Profile, nil
		}
		// Expired, remove from memory cache
		s.memoryCache.Delete(userID)
	}

	// 2. Sử dụng mutex để tránh cache stampede
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	// Double-check memory cache sau khi có lock
	if cached, ok := s.memoryCache.Load(userID); ok {
		cachedProfile := cached.(CachedProfile)
		if time.Now().Before(cachedProfile.ExpiresAt) {
			s.cacheHits++
			log.Printf("Memory Cache HIT (double-check) for user ID: %s", userID)
			return cachedProfile.Profile, nil
		}
		s.memoryCache.Delete(userID)
	}

	cacheKey := fmt.Sprintf("user_profile:%s", userID)

	// 3. Kiểm tra Redis cache
	var cachedProfile domain.UserProfile
	err := s.redisService.GetJSON(cacheKey, &cachedProfile)
	if err == nil {
		s.cacheHits++
		log.Printf("Redis Cache HIT for user ID: %s", userID)

		// Lưu vào memory cache để lần sau nhanh hơn
		memCache := CachedProfile{
			Profile:   cachedProfile,
			CachedAt:  time.Now(),
			ExpiresAt: time.Now().Add(5 * time.Minute), // Memory cache 5 phút
		}
		s.memoryCache.Store(userID, memCache)

		return cachedProfile, nil
	}

	s.cacheMisses++
	if err == redis.Nil {
		log.Printf("Cache MISS for user ID: %s. Fetching from DB.", userID)
	} else {
		log.Printf("Warning: Redis error on GetJSON for key %s: %v", cacheKey, err)
	}

	// 4. Fetch from database
	profile, err := s.userProfileRepo.GetUserProfileByID(ctx, userID)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("user not found")
	}

	// 5. Cache trong Redis với TTL 15 phút
	err = s.redisService.SetJSON(cacheKey, profile, 15*time.Minute)
	if err != nil {
		log.Printf("Warning: Failed to set Redis cache for user ID %s: %v", userID, err)
	}

	// 6. Cache trong memory với TTL 5 phút
	memCache := CachedProfile{
		Profile:   *profile,
		CachedAt:  time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	s.memoryCache.Store(userID, memCache)

	// 7. Log cache statistics
	totalRequests := s.cacheHits + s.cacheMisses
	if totalRequests%100 == 0 { // Log every 100 requests
		hitRatio := float64(s.cacheHits) / float64(totalRequests) * 100
		log.Printf("Cache Stats: Hits=%d, Misses=%d, Hit Ratio=%.2f%%",
			s.cacheHits, s.cacheMisses, hitRatio)
	}

	return *profile, nil
}

// GetCacheStats trả về thống kê cache performance
func (s *UserService) GetCacheStats() map[string]interface{} {
	totalRequests := s.cacheHits + s.cacheMisses
	hitRatio := float64(0)
	if totalRequests > 0 {
		hitRatio = float64(s.cacheHits) / float64(totalRequests) * 100
	}

	// Count memory cache entries
	memoryCacheSize := 0
	s.memoryCache.Range(func(key, value interface{}) bool {
		memoryCacheSize++
		return true
	})

	return map[string]interface{}{
		"cache_hits":        s.cacheHits,
		"cache_misses":      s.cacheMisses,
		"total_requests":    totalRequests,
		"hit_ratio_percent": hitRatio,
		"memory_cache_size": memoryCacheSize,
	}
}

// ClearMemoryCache xóa tất cả memory cache
func (s *UserService) ClearMemoryCache() {
	s.memoryCache = sync.Map{}
	log.Printf("Memory cache cleared")
}
