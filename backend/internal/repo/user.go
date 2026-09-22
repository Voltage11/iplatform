package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Voltage11/iplatform/internal/db"
	"github.com/Voltage11/iplatform/internal/domain"
	"github.com/Voltage11/iplatform/internal/types/apperr"
	"github.com/Voltage11/iplatform/internal/types/filterbool"
	"github.com/Voltage11/iplatform/internal/types/pagination"
)

// UserRepo репозиторий работы с пользователями
type UserRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepo создаёт репозиторий пользователей
func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

// Create вставляет пользователя и заполняет user.ID сгенерированным значением
// Время created_at и updated_at должно быть установлено в сервисе
func (u *UserRepo) Create(ctx context.Context, user *domain.User) error {
	const sql = `
		INSERT INTO users(email, first_name, last_name, password_hash, is_admin, is_active, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	q := db.QuerierFrom(ctx, u.pool)

	if err := q.QueryRow(ctx, sql,
		user.Email,
		user.FirstName,
		user.LastName,
		user.PasswordHash,
		user.IsAdmin,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID); err != nil {
		return apperr.NewPostgresError(err)
	}

	return nil
}

// Get возвращает пользователя по ID
// Мягко удалённые записи (deleted_at IS NOT NULL) не возвращаются
// Если запись не найдена — возвращает apperr.ErrNotFound
func (u *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const sql = `
		SELECT id, email, first_name, last_name, password_hash, is_admin, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	user := &domain.User{}
	var deletedAt *time.Time

	q := db.QuerierFrom(ctx, u.pool)

	if err := q.QueryRow(ctx, sql, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.IsAdmin,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
	); err != nil {
		return nil, apperr.NewPostgresError(err)
	}

	user.DeletedAt = deletedAt
	return user, nil
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const sql = `
		SELECT id, email, first_name, last_name, password_hash, is_admin, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	var deletedAt *time.Time

	q := db.QuerierFrom(ctx, u.pool)

	if err := q.QueryRow(ctx, sql, email).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.IsAdmin,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
	); err != nil {
		return nil, apperr.NewPostgresError(err)
	}

	user.DeletedAt = deletedAt
	return user, nil
}

// Update обновляет профиль пользователя
// Поле password_hash не трогается — для смены пароля отдельный метод будет
// Мягко удалённые записи не обновляются
func (u *UserRepo) Update(ctx context.Context, user *domain.User) error {
	const sql = `
		UPDATE users
		SET email = $1,
			first_name = $2,
			last_name = $3,
			is_admin = $4,
			is_active = $5,
			updated_at = $6
		WHERE id = $7 AND deleted_at IS NULL
	`
	q := db.QuerierFrom(ctx, u.pool)

	result, err := q.Exec(ctx, sql,
		user.Email,
		user.FirstName,
		user.LastName,
		user.IsAdmin,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return apperr.NewPostgresError(err)
	}

	if result.RowsAffected() == 0 {
		return apperr.NewNotFound(fmt.Sprintf("пользователь с ID не найден: %s", user.ID), nil)
	}

	return nil
}

// SoftDelete помечает пользователя как удалённого
// Если запись уже удалена или не существует — возвращает apperr.ErrNotFound
func (u *UserRepo) SoftDelete(ctx context.Context, userID uuid.UUID) error {
	const sql = `
		UPDATE users
		SET deleted_at = $1, updated_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	now := time.Now().UTC()

	q := db.QuerierFrom(ctx, u.pool)

	result, err := q.Exec(ctx, sql, now, userID)
	if err != nil {
		return apperr.NewPostgresError(err)
	}

	if result.RowsAffected() == 0 {
		return apperr.NewNotFound(fmt.Sprintf("пользователь с ID не найден: %s", userID), nil)
	}

	return nil
}

// SoftUnDelete восстанавливает мягко удалённого пользователя
// Если запись не удалена или не существует — возвращает apperr.ErrNotFound
func (u *UserRepo) SoftUnDelete(ctx context.Context, userID uuid.UUID) error {
	const sql = `
		UPDATE users
		SET deleted_at = NULL, updated_at = $1
		WHERE id = $2 AND deleted_at IS NOT NULL
	`
	q := db.QuerierFrom(ctx, u.pool)

	result, err := q.Exec(ctx, sql, time.Now().UTC(), userID)
	if err != nil {
		return apperr.NewPostgresError(err)
	}

	if result.RowsAffected() == 0 {
		return apperr.NewNotFound(fmt.Sprintf("пользователь с ID не найден: %s", userID), nil)
	}

	return nil
}

// ChangePassword обновляет пароль, вынес в отдельный метод
func (u *UserRepo) ChangePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	const sql = `
		UPDATE users
		SET password_hash = $1
		WHERE id = $2 AND deleted_at IS NULL
	`
	q := db.QuerierFrom(ctx, u.pool)

	result, err := q.Exec(ctx, sql, newPassword, userID)
	if err != nil {
		return apperr.NewPostgresError(err)
	}

	if result.RowsAffected() == 0 {
		return apperr.NewNotFound(fmt.Sprintf("пользователь с ID не найден: %s", userID.String()), nil)
	}

	return nil
}

// UserFilter — параметры фильтрации и пагинации для GetList
type UserFilter struct {
	Email      string
	FirstName  string
	LastName   string
	IsAdmin    filterbool.Filter
	IsActive   filterbool.Filter
	Pagination pagination.PaginationRequest
}

// GetList возвращает страницу активных пользователей и общее количество,
// Для поиска по подстроке передавайте значение с "%" на нужной позиции.
func (u *UserRepo) GetList(ctx context.Context, filter UserFilter) ([]*domain.User, int64, error) {
	const baseSelect = `
		SELECT id, email, first_name, last_name, is_admin, is_active, created_at, updated_at
		FROM users
	`
	const baseCount = `
		SELECT COUNT(id)
		FROM users
	`

	where := " WHERE deleted_at IS NULL"
	args := make([]any, 0, 5)

	if filter.Email != "" {
		args = append(args, filter.Email)
		where += fmt.Sprintf(" AND email ILIKE $%d", len(args))
	}

	if filter.FirstName != "" {
		args = append(args, filter.FirstName)
		where += fmt.Sprintf(" AND first_name ILIKE $%d", len(args))
	}

	if filter.LastName != "" {
		args = append(args, filter.LastName)
		where += fmt.Sprintf(" AND last_name ILIKE $%d", len(args))
	}

	if !filter.IsAdmin.IsAll() {
		args = append(args, *filter.IsAdmin.Bool())
		where += fmt.Sprintf(" AND is_admin = $%d", len(args))
	}

	if !filter.IsActive.IsAll() {
		args = append(args, *filter.IsActive.Bool())
		where += fmt.Sprintf(" AND is_active = $%d", len(args))
	}

	q := db.QuerierFrom(ctx, u.pool)
	// Общее количество — до применения LIMIT/OFFSET.
	var total int64
	if err := q.QueryRow(ctx, baseCount+where, args...).Scan(&total); err != nil {
		return nil, 0, apperr.NewPostgresError(err)
	}

	// Если записей нет — нет смысла делать второй запрос.
	if total == 0 {
		return []*domain.User{}, 0, nil
	}

	args = append(args, filter.Pagination.GetLimit(), filter.Pagination.GetOffset())
	sql := baseSelect + where +
		fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, apperr.NewPostgresError(err)
	}
	defer rows.Close()

	userList := make([]*domain.User, 0, filter.Pagination.GetLimit())

	for rows.Next() {
		user := &domain.User{}
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.IsAdmin,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, 0, apperr.NewPostgresError(err)
		}
		userList = append(userList, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, apperr.NewPostgresError(err)
	}

	return userList, total, nil
}
