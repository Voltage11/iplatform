package repo

import (
	"context"
	"encoding/json"

	"fmt"
	"time"

	"github.com/Voltage11/iplatform/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CachedUserRepo struct {
	dbRepo *UserRepo
	rdb    *redis.Client
	ttl    time.Duration // Время жизни кэша
}

func NewCachedUserRepo(dbRepo *UserRepo, rdb *redis.Client, ttl time.Duration) *CachedUserRepo {
	return &CachedUserRepo{
		dbRepo: dbRepo,
		rdb:    rdb,
		ttl:    ttl,
	}
}

// Генерируем уникальный ключ для кэша пользователя
func userKey(id uuid.UUID) string {
	return fmt.Sprintf("user:%s", id.String())
}

func (r *CachedUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	key := userKey(id)

	// 1. Пытаемся получить данные из Redis
	val, err := r.rdb.Get(ctx, key).Result()
	if err == nil {
		// Нашли в кэше. Десериализуем JSON обратно в структуру
		var user domain.User
		if err := json.Unmarshal([]byte(val), &user); err == nil {
			return &user, nil
		}
	}

	// 2. Если в кэше нет (или ошибка), идем в основную БД Postgres
	user, err := r.dbRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 3. Сохраняем найденного пользователя в Redis в JSON
	if data, err := json.Marshal(user); err == nil {
		_ = r.rdb.Set(ctx, key, data, r.ttl).Err()
	}

	return user, nil
}

// При обновлении или удалении профиля кэш сбрасывать
func (r *CachedUserRepo) Update(ctx context.Context, user *domain.User) error {
	err := r.dbRepo.Update(ctx, user)
	if err != nil {
		return err
	}
	// Удаляем старый кэш, чтобы при следующем запросе загрузились новые данные
	_ = r.rdb.Del(ctx, userKey(user.ID)).Err()
	return nil
}

func (r *CachedUserRepo) SoftDelete(ctx context.Context, userID uuid.UUID) error {
	err := r.dbRepo.SoftDelete(ctx, userID)
	if err != nil {
		return err
	}
	_ = r.rdb.Del(ctx, userKey(userID)).Err()
	return nil
}

func (r *CachedUserRepo) Create(ctx context.Context, user *domain.User) error {
	return r.dbRepo.Create(ctx, user)
}
func (r *CachedUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.dbRepo.GetByEmail(ctx, email)
}
func (r *CachedUserRepo) SoftUnDelete(ctx context.Context, userID uuid.UUID) error {
	_ = r.rdb.Del(ctx, userKey(userID))
	return r.dbRepo.SoftUnDelete(ctx, userID)
}
