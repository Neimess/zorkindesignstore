package category

import (
    "context"
    "testing"
    "time"

    "github.com/Neimess/zorkin-store-project/internal/domain"
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

    postgresContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15-alpine"),
        postgres.WithDatabase(dbName),
        postgres.WithUsername(dbUser),
        postgres.WithPassword(dbPass),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(5*time.Minute),
        ),
    )
    require.NoError(t, err)

    connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
    require.NoError(t, err)

    db, err := sqlx.Connect("postgres", connStr)
    require.NoError(t, err)

    // Создание схемы БД
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

    CREATE TABLE IF NOT EXISTS products (
        product_id BIGSERIAL PRIMARY KEY,
        category_id BIGINT NOT NULL REFERENCES categories(category_id),
        name VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `
    _, err := db.Exec(schema)
    return err
}

func TestPGCategoryRepository_Create(t *testing.T) {
    tc := setupTestDB(t)
    defer tc.teardown()

    repo := NewPGCategoryRepository(tc.db)
    ctx := context.Background()

    t.Run("успешное создание категории", func(t *testing.T) {
        id, err := repo.Create(ctx, "Мебель")
        require.NoError(t, err)
        assert.NotZero(t, id)

        // Проверка что категория действительно создалась
        cat, err := repo.GetByID(ctx, id)
        require.NoError(t, err)
        assert.Equal(t, "Мебель", cat.Name)
    })

    t.Run("пустое имя", func(t *testing.T) {
        _, err := repo.Create(ctx, "")
        assert.ErrorIs(t, err, domain.ErrCategoryNameEmpty)
    })

    t.Run("дублирование имени", func(t *testing.T) {
        _, err := repo.Create(ctx, "Мебель")
        assert.Error(t, err)
    })
}

func TestPGCategoryRepository_GetByID(t *testing.T) {
    tc := setupTestDB(t)
    defer tc.teardown()

    repo := NewPGCategoryRepository(tc.db)
    ctx := context.Background()

    t.Run("получение существующей категории", func(t *testing.T) {
        id, err := repo.Create(ctx, "Электроника")
        require.NoError(t, err)

        cat, err := repo.GetByID(ctx, id)
        require.NoError(t, err)
        assert.Equal(t, id, cat.ID)
        assert.Equal(t, "Электроника", cat.Name)
    })

    t.Run("несуществующая категория", func(t *testing.T) {
        _, err := repo.GetByID(ctx, 9999)
        assert.Error(t, err)
    })
}

func TestPGCategoryRepository_Update(t *testing.T) {
    tc := setupTestDB(t)
    defer tc.teardown()

    repo := NewPGCategoryRepository(tc.db)
    ctx := context.Background()

    t.Run("успешное обновление", func(t *testing.T) {
        id, err := repo.Create(ctx, "Техника")
        require.NoError(t, err)

        err = repo.Update(ctx, id, "Бытовая техника")
        require.NoError(t, err)

        cat, err := repo.GetByID(ctx, id)
        require.NoError(t, err)
        assert.Equal(t, "Бытовая техника", cat.Name)
    })

    t.Run("обновление несуществующей категории", func(t *testing.T) {
        err := repo.Update(ctx, 9999, "Несуществующая")
        assert.ErrorIs(t, err, domain.ErrCategoryNotFound)
    })

    t.Run("обновление с дублированием имени", func(t *testing.T) {
        id1, err := repo.Create(ctx, "Категория 1")
        require.NoError(t, err)

        _, err = repo.Create(ctx, "Категория 2")
        require.NoError(t, err)

        err = repo.Update(ctx, id1, "Категория 2")
        assert.Error(t, err) // Ожидаем ошибку дублирования
    })
}

func TestPGCategoryRepository_Delete(t *testing.T) {
    tc := setupTestDB(t)
    defer tc.teardown()

    repo := NewPGCategoryRepository(tc.db)
    ctx := context.Background()

    t.Run("успешное удаление", func(t *testing.T) {
        id, err := repo.Create(ctx, "Для удаления")
        require.NoError(t, err)

        err = repo.Delete(ctx, id)
        require.NoError(t, err)

        // Проверяем что категория удалена
        _, err = repo.GetByID(ctx, id)
        assert.Error(t, err)
    })

    t.Run("удаление несуществующей категории", func(t *testing.T) {
        err := repo.Delete(ctx, 9999)
        assert.ErrorIs(t, err, domain.ErrCategoryNotFound)
    })

    t.Run("удаление категории с товарами", func(t *testing.T) {
        // Создаем категорию
        id, err := repo.Create(ctx, "С товарами")
        require.NoError(t, err)

        // Добавляем товар в эту категорию
        _, err = tc.db.ExecContext(ctx, 
            "INSERT INTO products (category_id, name) VALUES ($1, $2)",
            id, "Тестовый товар")
        require.NoError(t, err)

        // Пытаемся удалить категорию с товарами
        err = repo.Delete(ctx, id)
        assert.Error(t, err) // Должна возникнуть ошибка из-за внешнего ключа
    })
}

func TestPGCategoryRepository_List(t *testing.T) {
    tc := setupTestDB(t)
    defer tc.teardown()

    repo := NewPGCategoryRepository(tc.db)
    ctx := context.Background()

    t.Run("список пустой", func(t *testing.T) {
        // Проверяем что изначально список пуст
        list, err := repo.List(ctx)
        require.NoError(t, err)
        assert.Empty(t, list)
    })

    t.Run("список с категориями", func(t *testing.T) {
        // Создаём категории в разнобой, чтобы проверить сортировку
        _, err := repo.Create(ctx, "Зетта")
        require.NoError(t, err)
        
        _, err = repo.Create(ctx, "Альфа")
        require.NoError(t, err)
        
        _, err = repo.Create(ctx, "Омега")
        require.NoError(t, err)

        list, err := repo.List(ctx)
        require.NoError(t, err)
        assert.Len(t, list, 3)

        // Проверяем сортировку по имени
        assert.Equal(t, "Альфа", list[0].Name)
        assert.Equal(t, "Зетта", list[1].Name)
        assert.Equal(t, "Омега", list[2].Name)
    })
}

func TestPGCategoryRepository_Integration(t *testing.T) {
    tc := setupTestDB(t)
    defer tc.teardown()

    repo := NewPGCategoryRepository(tc.db)
    ctx := context.Background()

    t.Run("полный жизненный цикл", func(t *testing.T) {
        // Создаем категории
        id1, err := repo.Create(ctx, "Категория 1")
        require.NoError(t, err)
        
        id2, err := repo.Create(ctx, "Категория 2")
        require.NoError(t, err)
        
        // Получаем список
        list, err := repo.List(ctx)
        require.NoError(t, err)
        assert.Len(t, list, 2)
        
        // Обновляем категорию
        err = repo.Update(ctx, id1, "Обновленная категория")
        require.NoError(t, err)
        
        updated, err := repo.GetByID(ctx, id1)
        require.NoError(t, err)
        assert.Equal(t, "Обновленная категория", updated.Name)
        
        // Удаляем категорию
        err = repo.Delete(ctx, id2)
        require.NoError(t, err)
        
        // Проверяем итоговое состояние
        finalList, err := repo.List(ctx)
        require.NoError(t, err)
        assert.Len(t, finalList, 1)
        assert.Equal(t, "Обновленная категория", finalList[0].Name)
    })
}