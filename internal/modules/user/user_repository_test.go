package user_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"go-gorm/config"
	"go-gorm/internal/database/model"
	"go-gorm/internal/database/query"
	"go-gorm/internal/modules/user" // sesuaikan import ini

	"github.com/stretchr/testify/assert"
)

var testRepo user.RepositoryUser

func TestMain(m *testing.M) {
	if err := os.Chdir("../../../"); 
	err != nil {
		fmt.Printf("Gagal berpindah ke root directory: %v\n", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Gagal load config: %v\n", err)
		os.Exit(1)
	}

	sqlDB := config.InitDB(cfg)
	db, err := config.InitGORM(sqlDB)
	if err != nil {
		fmt.Printf("Gagal koneksi ke test DB: %v\n", err)
		sqlDB.Close()
		os.Exit(1)
	}

	q := query.Use(db)
	testRepo = user.NewUserRepository(q)

	code := m.Run()

	// Closing koneksi DB secara eksplisit
	sqlDB.Close()
	os.Exit(code)
}

func TestCreateUser_Success(t *testing.T) {
	ctx := context.Background()

	uniqueEmail := fmt.Sprintf("dummy_%d@example.com", time.Now().UnixNano())

	dummyUser := &model.User{
		Name:  "Test Dummy User",
		Email: uniqueEmail,
		CreatedAt: time.Now(),
	}

	// Act
	err := testRepo.Create(ctx, dummyUser)

	// Assert
	assert.NoError(t, err)
	assert.NotZero(t, dummyUser.ID, "ID auto increment harus terisi pada struct")
	assert.False(t, dummyUser.CreatedAt.IsZero())
}
func TestCreateBatchUser_Success(t *testing.T) {
	ctx := context.Background()

	// 1. Buat data dummy sebanyak 5 user
	var dummyUsers []*model.User
	now := time.Now().UnixNano()

	for i := 1; i <= 5; i++ {
		dummyUsers = append(dummyUsers, &model.User{
			Name:     fmt.Sprintf("Batch User %d", i),
			Email:    fmt.Sprintf("batch_%d_%d@example.com", now, i),
			CreatedAt: time.Now(),
		})
	}

	// 2. Act: Eksekusi CreateBatch dengan batchSize 2
	// 5 data dengan batchSize 2 akan menghasilkan 3 query SQL INSERT (2 + 2 + 1)
	err := testRepo.CreateBatch(ctx, dummyUsers, 2)

	// 3. Assert
	assert.NoError(t, err)

	// Pastikan semua user di slice dummyUsers mendapatkan ID dari MySQL
	for _, user := range dummyUsers {
		assert.NotZero(t, user.ID, "ID auto increment setiap user harus terisi")
	}
}
