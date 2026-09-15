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
	if err := os.Chdir("../../../"); err != nil {
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
	testRepo = user.NewUserRepository(q, db)

	code := m.Run()

	// Closing koneksi DB secara eksplisit
	sqlDB.Close()
	os.Exit(code)
}

func TestCreateUser_Success(t *testing.T) {
	ctx := context.Background()

	uniqueEmail := fmt.Sprintf("dummy_%d@example.com", time.Now().UnixNano())

	dummyUser := &model.User{
		Name:      "Test Dummy User",
		Email:     uniqueEmail,
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
			Name:      fmt.Sprintf("Batch User %d", i),
			Email:     fmt.Sprintf("batch_%d_%d@example.com", now, i),
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

func TestTransactionSuccess(t *testing.T) {
	ctx := context.Background()
	uniqueEmail := fmt.Sprintf("dummy_%d@example.com", time.Now().UnixNano())
	dummyUser := &model.User{
		Name:      "Test Dummy User",
		Email:     uniqueEmail,
		CreatedAt: time.Now(),
	}

	err := testRepo.Transaction(ctx, dummyUser)
	assert.NoError(t, err)
	assert.NotZero(t, dummyUser.ID, "ID must not zero value")
	assert.False(t, dummyUser.CreatedAt.IsZero())
}
func TestTransactionFail(t *testing.T) {
	ctx := context.Background()
	uniqueEmail := fmt.Sprintf("dummy_%d@example.com", time.Now().UnixNano())
	dummyUser := &model.User{
		Name:      "",
		Email:     uniqueEmail,
		CreatedAt: time.Now(),
	}

	err := testRepo.Transaction(ctx, dummyUser)
	assert.NoError(t, err)
	assert.NotZero(t, dummyUser.ID, "ID must not zero value")
	assert.False(t, dummyUser.CreatedAt.IsZero())
}

func TestTransactionManual(t *testing.T) {
	ctx := context.Background()
	uniqueEmail := fmt.Sprintf("dummy_%d@example.com", time.Now().UnixNano())
	dummyUser := &model.User{
		Name:      "Test Dummy User",
		Email:     uniqueEmail,
		CreatedAt: time.Now(),
	}
	err := testRepo.ManualTransaction(ctx, dummyUser)
	assert.NoError(t, err)
	assert.NotZero(t, dummyUser.ID, "ID must be available")
	assert.False(t, dummyUser.CreatedAt.IsZero())
}

func TestGetByEmail(t *testing.T) {
	ctx := context.Background()
	result, err := testRepo.GetFirstEmail(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Dummy User", result.Name)
	assert.Equal(t, "dummy@example.com", result.Email)
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()
	user, err := testRepo.GetFirstID(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Dummy User", user.Name)
	assert.Equal(t, "dummy@example.com", user.Email)
}

func TestGetLastEmail(t *testing.T) {
	ctx := context.Background()
	user, err := testRepo.GetLastEmail(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Test Dummy User", user.Name)
	assert.Equal(t, "dummy_1789454181092507971@example.com", user.Email)
}

func TestGetAll(t *testing.T) {
	ctx := context.Background()
	user, err := testRepo.GetAllFind(ctx)
	assert.NoError(t, err)
	assert.Len(t, user, 38)
}
func TestGetAllPaginate(t *testing.T) {
	ctx := context.Background()
	user, err := testRepo.GetAllPaginate(ctx, 38, 0)
	assert.NoError(t, err)
	assert.Len(t, user, 38)
}

func TestSearch(t *testing.T) {
	ctx := context.Background()
	u, err := testRepo.SearchUser(ctx, "Dummy U","dummy@mail.com")
	assert.NoError(t, err)
	assert.Len(t, u, 1)
}

func TestSearchNameOrEmail(t *testing.T) {
	ctx := context.Background()
	u,err := testRepo.SearcUserOrEmail(ctx,"Dummy U","dummy@mail.com")
	assert.NoError(t,err)
	assert.Equal(t,u.Name,"Dummy U")
}

func TestSearchNotName(t *testing.T) {
	ctx := context.Background()
	u,err := testRepo.SearchNotName(ctx,"Dummy u")
	assert.NoError(t,err)
	assert.Len(t,u,38)
}
func TestFieldName(t *testing.T) {
	ctx :=context.Background()
	u,err := testRepo.SelectFieldName(ctx)
	assert.NoError(t,err)
	assert.Len(t,u,39)
}

func TestStructCondition(t *testing.T) {
	ctx := context.Background()
	u,err := testRepo.StructCondition(ctx,"Dummy U","dummy@mail.com")
	assert.NoError(t,err)
	assert.Equal(t,"Dummy U",u.Name)
	assert.Equal(t,"dummy@mail.com",u.Email)
}

func TestPaginationOrder(t *testing.T) {
	ctx := context.Background()
	page1, err := testRepo.PaginationOrderName(ctx, 10, 0)
    assert.NoError(t, err)
    assert.Len(t, page1, 10)

    page2, err := testRepo.PaginationOrderName(ctx, 10, 10)
    assert.NoError(t, err)
    assert.Len(t, page2, 10)

}