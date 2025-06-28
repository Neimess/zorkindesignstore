package attribute

import (
	"context"
	"testing"
	"time"

	attr "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbName = "testdb"
	dbUser = "testuser"
	dbPass = "testpass"
)

type testContainer struct {
	container testcontainers.Container
	db        *sqlx.DB
}

func setupTestDB(t *testing.T) *testContainer {
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx, "postgres:15-alpine", postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Minute),
		))
	require.NoError(t, err)

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sqlx.Connect("postgres", connStr)
	require.NoError(t, err)

	err = createSchema(db)
	require.NoError(t, err)

	return &testContainer{
		container: postgresContainer,
		db:        db,
	}
}

func (tc *testContainer) teardown() {
	if tc.db != nil {
		tc.db.Close()
	}
	if tc.container != nil {
		tc.container.Terminate(context.Background())
	}
}

func createSchema(db *sqlx.DB) error {
	schema := `
    CREATE TABLE IF NOT EXISTS categories (
        category_id BIGSERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS attributes (
        attribute_id BIGSERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        unit VARCHAR(50),
        category_id BIGINT NOT NULL REFERENCES categories(category_id) ON DELETE CASCADE,
        UNIQUE(name, category_id)
    );
    `
	_, err := db.Exec(schema)
	return err
}

func createTestCategory(t *testing.T, db *sqlx.DB, name string) int64 {
	var id int64
	err := db.QueryRow("INSERT INTO categories (name) VALUES ($1) RETURNING category_id", name).Scan(&id)
	require.NoError(t, err)
	return id
}

func TestPGAttributeRepository_SaveBatch(t *testing.T) {
	tc := setupTestDB(t)
	defer tc.teardown()

	repo := NewPGAttributeRepository(tc.db)
	ctx := context.Background()

	categoryID := createTestCategory(t, tc.db, "test_category")

	t.Run("успешное сохранение батча атрибутов", func(t *testing.T) {
		unit1 := "кг"
		unit2 := "шт"
		attrs := []*attr.Attribute{
			{Name: "Вес", Unit: &unit1, CategoryID: categoryID},
			{Name: "Количество", Unit: &unit2, CategoryID: categoryID},
			{Name: "Цвет", Unit: nil, CategoryID: categoryID},
		}

		err := repo.SaveBatch(ctx, attrs)
		require.NoError(t, err)

		for i, attr := range attrs {
			assert.NotZero(t, attr.ID, "ID не присвоен для атрибута %d", i)
		}

		ids := make(map[int64]bool)
		for _, attr := range attrs {
			assert.False(t, ids[attr.ID], "Дублированный ID: %d", attr.ID)
			ids[attr.ID] = true
		}
	})

	t.Run("пустой батч", func(t *testing.T) {
		var attrs []*attr.Attribute
		err := repo.SaveBatch(ctx, attrs)
		assert.NoError(t, err)
	})

	t.Run("ошибка при дублировании имени в категории", func(t *testing.T) {
		unit := "кг"
		attrs := []*attr.Attribute{
			{Name: "Дублированное имя", Unit: &unit, CategoryID: categoryID},
			{Name: "Дублированное имя", Unit: &unit, CategoryID: categoryID},
		}

		err := repo.SaveBatch(ctx, attrs)
		assert.Error(t, err)
	})
}

func TestPGAttributeRepository_Save(t *testing.T) {
	tc := setupTestDB(t)
	defer tc.teardown()

	repo := NewPGAttributeRepository(tc.db)
	ctx := context.Background()

	categoryID := createTestCategory(t, tc.db, "test_category")

	t.Run("успешное сохранение атрибута", func(t *testing.T) {
		unit := "кг"
		attr := &attr.Attribute{
			Name:       "Вес",
			Unit:       &unit,
			CategoryID: categoryID,
		}

		err := repo.Save(ctx, attr)
		require.NoError(t, err)
		assert.NotZero(t, attr.ID)
	})

	t.Run("сохранение атрибута без единицы измерения", func(t *testing.T) {
		attr := &attr.Attribute{
			Name:       "Цвет",
			Unit:       nil,
			CategoryID: categoryID,
		}

		err := repo.Save(ctx, attr)
		require.NoError(t, err)
		assert.NotZero(t, attr.ID)
	})

	t.Run("ошибка при дублировании имени в категории", func(t *testing.T) {
		unit := "см"
		attr := &attr.Attribute{
			Name:       "Вес",
			Unit:       &unit,
			CategoryID: categoryID,
		}

		err := repo.Save(ctx, attr)
		assert.Error(t, err)
	})
}

func TestPGAttributeRepository_GetByID(t *testing.T) {
	tc := setupTestDB(t)
	defer tc.teardown()

	repo := NewPGAttributeRepository(tc.db)
	ctx := context.Background()

	categoryID := createTestCategory(t, tc.db, "test_category")

	t.Run("успешное получение атрибута", func(t *testing.T) {
		unit := "кг"
		original := &attr.Attribute{
			Name:       "Тестовый атрибут",
			Unit:       &unit,
			CategoryID: categoryID,
		}

		err := repo.Save(ctx, original)
		require.NoError(t, err)

		found, err := repo.GetByID(ctx, original.ID)
		require.NoError(t, err)
		assert.Equal(t, original.ID, found.ID)
		assert.Equal(t, original.Name, found.Name)
		assert.Equal(t, *original.Unit, *found.Unit)
		assert.Equal(t, original.CategoryID, found.CategoryID)
	})

	t.Run("атрибут не найден", func(t *testing.T) {
		_, err := repo.GetByID(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestPGAttributeRepository_FindByCategory(t *testing.T) {
	tc := setupTestDB(t)
	defer tc.teardown()

	repo := NewPGAttributeRepository(tc.db)
	ctx := context.Background()

	categoryID := createTestCategory(t, tc.db, "test_category")

	t.Run("получение атрибутов категории", func(t *testing.T) {
		unit1 := "кг"
		unit2 := "шт"
		attrs := []*attr.Attribute{
			{Name: "Атрибут А", Unit: &unit1, CategoryID: categoryID},
			{Name: "Атрибут Б", Unit: &unit2, CategoryID: categoryID},
			{Name: "Атрибут В", Unit: nil, CategoryID: categoryID},
		}

		err := repo.SaveBatch(ctx, attrs)
		require.NoError(t, err)

		found, err := repo.FindByCategory(ctx, categoryID)
		require.NoError(t, err)
		assert.Len(t, found, 3)

		// Проверяем сортировку по имени
		assert.Equal(t, "Атрибут А", found[0].Name)
		assert.Equal(t, "Атрибут Б", found[1].Name)
		assert.Equal(t, "Атрибут В", found[2].Name)
	})

	t.Run("пустая категория", func(t *testing.T) {
		emptyCategoryID := createTestCategory(t, tc.db, "empty_category")

		found, err := repo.FindByCategory(ctx, emptyCategoryID)
		require.NoError(t, err)
		assert.Empty(t, found)
	})
}

func TestPGAttributeRepository_Update(t *testing.T) {
	tc := setupTestDB(t)
	defer tc.teardown()

	repo := NewPGAttributeRepository(tc.db)
	ctx := context.Background()

	categoryID := createTestCategory(t, tc.db, "test_category")

	t.Run("успешное обновление атрибута", func(t *testing.T) {
		unit := "кг"
		attr := &attr.Attribute{
			Name:       "Исходное имя",
			Unit:       &unit,
			CategoryID: categoryID,
		}

		err := repo.Save(ctx, attr)
		require.NoError(t, err)

		newUnit := "г"
		attr.Name = "Новое имя"
		attr.Unit = &newUnit

		err = repo.Update(ctx, attr)
		require.NoError(t, err)

		updated, err := repo.GetByID(ctx, attr.ID)
		require.NoError(t, err)
		assert.Equal(t, "Новое имя", updated.Name)
		assert.Equal(t, "г", *updated.Unit)
	})

	t.Run("обновление несуществующего атрибута", func(t *testing.T) {
		attr := &attr.Attribute{
			ID:         99999,
			Name:       "Новое имя",
			CategoryID: categoryID,
		}

		err := repo.Update(ctx, attr)
		assert.Error(t, err)
	})
}

func TestPGAttributeRepository_Delete(t *testing.T) {
	tc := setupTestDB(t)
	defer tc.teardown()

	repo := NewPGAttributeRepository(tc.db)
	ctx := context.Background()

	categoryID := createTestCategory(t, tc.db, "test_category")

	t.Run("успешное удаление атрибута", func(t *testing.T) {
		unit := "кг"
		attr := &attr.Attribute{
			Name:       "Удаляемый атрибут",
			Unit:       &unit,
			CategoryID: categoryID,
		}

		err := repo.Save(ctx, attr)
		require.NoError(t, err)

		err = repo.Delete(ctx, attr.ID)
		require.NoError(t, err)

		// Проверяем, что атрибут действительно удален
		_, err = repo.GetByID(ctx, attr.ID)
		assert.Error(t, err)
	})

	t.Run("удаление несуществующего атрибута", func(t *testing.T) {
		err := repo.Delete(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestPGAttributeRepository_Integration(t *testing.T) {
	tc := setupTestDB(t)
	defer tc.teardown()

	repo := NewPGAttributeRepository(tc.db)
	ctx := context.Background()

	categoryID := createTestCategory(t, tc.db, "integration_category")

	t.Run("полный жизненный цикл атрибутов", func(t *testing.T) {
		unit1 := "кг"
		unit2 := "шт"
		attrs := []*attr.Attribute{
			{Name: "Вес", Unit: &unit1, CategoryID: categoryID},
			{Name: "Количество", Unit: &unit2, CategoryID: categoryID},
			{Name: "Цвет", Unit: nil, CategoryID: categoryID},
		}

		err := repo.SaveBatch(ctx, attrs)
		require.NoError(t, err)

		found, err := repo.FindByCategory(ctx, categoryID)
		require.NoError(t, err)
		assert.Len(t, found, 3)

		newUnit := "г"
		attrs[0].Unit = &newUnit
		err = repo.Update(ctx, attrs[0])
		require.NoError(t, err)

		updated, err := repo.GetByID(ctx, attrs[0].ID)
		require.NoError(t, err)
		assert.Equal(t, "г", *updated.Unit)

		err = repo.Delete(ctx, attrs[1].ID)
		require.NoError(t, err)

		final, err := repo.FindByCategory(ctx, categoryID)
		require.NoError(t, err)
		assert.Len(t, final, 2)
	})
}
