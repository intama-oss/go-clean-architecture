package author

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go-clean-architecture/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"regexp"
	"testing"
	"time"
)

func mockDBConnection() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})

	if err != nil {
		return nil, nil, err
	}

	return gdb, mock, nil
}

func TestMysqlAuthorRepository_GetByID(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "SELECT * FROM `authors` WHERE `authors`.`id` = ? ORDER BY `authors`.`id` LIMIT ?"

	userID := uint(1)
	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(1, "Author 1", time.Now(), time.Now())
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userID, 1).WillReturnRows(rows)

	repo := NewMysqlAuthorRepository(db)
	author, err := repo.GetByID(userID)
	assert.NoError(t, err)
	assert.NotNil(t, author)
}

func TestMysqlAuthorRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "SELECT * FROM `authors` WHERE `authors`.`id` = ? ORDER BY `authors`.`id` LIMIT 1"

	userID := uint(1)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userID).WillReturnError(gorm.ErrRecordNotFound)

	repo := NewMysqlAuthorRepository(db)
	author, err := repo.GetByID(userID)
	assert.Error(t, err)
	assert.Nil(t, author)
}

func TestMysqlAuthorRepository_Count(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "SELECT count(*) FROM `authors`"

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	repo := NewMysqlAuthorRepository(db)
	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(10), count)
}

func TestMysqlAuthorRepository_Count_Error(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "SELECT count(*) FROM `authors`"

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(gorm.ErrInvalidDB)

	repo := NewMysqlAuthorRepository(db)
	count, err := repo.Count()
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
}

func TestMysqlAuthorRepository_Store(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "INSERT INTO `authors`"

	author := &domain.Author{
		ID:   0,
		Name: "Author 1",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewMysqlAuthorRepository(db)

	err = repo.Store(author)
	assert.NoError(t, err)
}
