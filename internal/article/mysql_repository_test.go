package article

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

func TestMysqlArticleRepository_Fetch(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "SELECT * FROM `articles` WHERE title LIKE ? ORDER BY created_at DESC LIMIT ?"

	expectedTitle := "title"
	expectedContent := "content"
	expectedAuthorID := 1
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
		AddRow(1, expectedTitle, expectedContent, expectedAuthorID, expectedCreatedAt, expectedUpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("%"+expectedTitle+"%", 10).
		WillReturnRows(rows)

	repo := NewMysqlArticleRepository(db)

	articles, _, err := repo.Fetch(1, 10, &domain.Article{Title: expectedTitle})
	assert.NoError(t, err)
	assert.NotNil(t, articles)
}

func TestMysqlArticleRepository_Fetch_Error(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "SELECT * FROM `articles` WHERE title LIKE ? ORDER BY created_at DESC LIMIT 10"

	expectedTitle := "title"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("%" + expectedTitle + "%").
		WillReturnError(assert.AnError)

	repo := NewMysqlArticleRepository(db)

	articles, _, err := repo.Fetch(1, 10, &domain.Article{Title: expectedTitle})
	assert.Error(t, err)
	assert.Nil(t, articles)
}

func TestMysqlArticleRepository_GetByID(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	queryArticle := "SELECT * FROM `articles` WHERE `articles`.`id` = ? ORDER BY `articles`.`id` LIMIT ?"

	expectedArticleID := 1
	expectedTitle := "title"
	expectedContent := "content"
	expectedAuthorID := 1
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()
	rowsArticle := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
		AddRow(expectedArticleID, expectedTitle, expectedContent, expectedAuthorID, expectedCreatedAt, expectedUpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(queryArticle)).
		WithArgs(expectedArticleID, 1).
		WillReturnRows(rowsArticle)

	queryAuthor := "SELECT * FROM `authors` WHERE `authors`.`id` = ?"
	expectedAuthorName := "author"
	rowsAuthor := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(expectedAuthorID, expectedAuthorName, expectedCreatedAt, expectedUpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(queryAuthor)).
		WithArgs(expectedAuthorID).
		WillReturnRows(rowsAuthor)

	repo := NewMysqlArticleRepository(db)

	article, err := repo.GetByID(uint(expectedArticleID))
	assert.NoError(t, err)
	assert.NotNil(t, article)
}

func TestMysqlArticleRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	queryArticle := "SELECT * FROM `articles` WHERE `articles`.`id` = ? ORDER BY `articles`.`id` LIMIT 1"

	expectedArticleID := 1
	mock.ExpectQuery(regexp.QuoteMeta(queryArticle)).
		WithArgs(expectedArticleID).
		WillReturnError(gorm.ErrRecordNotFound)

	repo := NewMysqlArticleRepository(db)

	article, err := repo.GetByID(uint(expectedArticleID))
	assert.Error(t, err)
	assert.Nil(t, article)
}

func TestMysqlArticleRepository_Count(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "SELECT count(*) FROM `articles` WHERE title LIKE ?"

	expectedTitle := "title"
	expectedCount := int64(1)
	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(expectedCount)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("%" + expectedTitle + "%").
		WillReturnRows(rows)

	repo := NewMysqlArticleRepository(db)

	count, err := repo.Count(&domain.Article{Title: expectedTitle})
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}

func TestMysqlArticleRepository_Count_WithFilter(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "SELECT count(*) FROM `articles` WHERE title LIKE ?"

	expectedTitle := "title"
	expectedCount := int64(1)
	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(expectedCount)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("%" + expectedTitle + "%").
		WillReturnRows(rows)

	repo := NewMysqlArticleRepository(db)

	count, err := repo.Count(&domain.Article{Title: expectedTitle})
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}

func TestMysqlArticleRepository_Count_Error(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "SELECT count(*) FROM `articles` WHERE title LIKE ?"

	expectedTitle := "title"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("%" + expectedTitle + "%").
		WillReturnError(assert.AnError)

	repo := NewMysqlArticleRepository(db)

	count, err := repo.Count(&domain.Article{Title: expectedTitle})
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
}

func TestMysqlArticleRepository_Store(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "INSERT INTO `articles` (`title`,`content`,`author_id`,`created_at`,`updated_at`) VALUES (?,?,?,?,?)"

	article := &domain.Article{
		Title:    "title",
		Content:  "content",
		AuthorID: 1,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(article.Title, article.Content, article.AuthorID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewMysqlArticleRepository(db)

	err = repo.Store(article)
	assert.NoError(t, err)
}

func TestMysqlArticleRepository_Store_Error(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "INSERT INTO `articles` (`title`,`content`,`author_id`,`created_at`,`updated_at`) VALUES (?,?,?,?,?)"

	expectedTitle := "title"
	expectedContent := "content"
	expectedAuthorID := uint(1)
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(expectedTitle, expectedContent, expectedAuthorID, expectedCreatedAt, expectedUpdatedAt).
		WillReturnError(assert.AnError)

	repo := NewMysqlArticleRepository(db)

	err = repo.Store(&domain.Article{
		Title:     expectedTitle,
		Content:   expectedContent,
		AuthorID:  expectedAuthorID,
		CreatedAt: expectedCreatedAt,
		UpdatedAt: expectedUpdatedAt,
	})
	assert.Error(t, err)
}

func TestMysqlArticleRepository_Update(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "UPDATE `articles` SET `title`=?,`content`=?,`author_id`=?,`updated_at`=? WHERE `id` = ?"

	article := &domain.Article{
		ID:       1,
		Title:    "title",
		Content:  "content",
		AuthorID: 1,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(article.Title, article.Content, article.AuthorID, sqlmock.AnyArg(), article.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewMysqlArticleRepository(db)

	err = repo.Update(article)
	assert.NoError(t, err)
}

func TestMysqlArticleRepository_Update_Error(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "UPDATE `articles` SET `title`=?,`content`=?,`author_id`=?,`created_at`=?,`updated_at`=? WHERE `id` = ?"

	expectedID := uint(1)
	expectedTitle := "title"
	expectedContent := "content"
	expectedAuthorID := uint(1)
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(expectedTitle, expectedContent, expectedAuthorID, expectedCreatedAt, expectedUpdatedAt, expectedID).
		WillReturnError(assert.AnError)

	repo := NewMysqlArticleRepository(db)

	err = repo.Update(&domain.Article{
		ID:        expectedID,
		Title:     expectedTitle,
		Content:   expectedContent,
		AuthorID:  expectedAuthorID,
		CreatedAt: expectedCreatedAt,
		UpdatedAt: expectedUpdatedAt,
	})
	assert.Error(t, err)
}

func TestMysqlArticleRepository_Delete(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "DELETE FROM `articles` WHERE `articles`.`id` = ?"

	expectedID := uint(1)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(expectedID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewMysqlArticleRepository(db)

	err = repo.Delete(expectedID)
	assert.NoError(t, err)
}

func TestMysqlArticleRepository_Delete_Error(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "DELETE FROM `articles` WHERE `id` = ?"

	expectedID := uint(1)

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(expectedID).
		WillReturnError(assert.AnError)

	repo := NewMysqlArticleRepository(db)

	err = repo.Delete(expectedID)
	assert.Error(t, err)
}

func TestMysqlArticleRepository_GetByAuthorID(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "SELECT * FROM `articles` WHERE author_id = ?"

	expectedAuthorID := uint(1)
	expectedTitle := "title"
	expectedContent := "content"
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
		AddRow(1, expectedTitle, expectedContent, expectedAuthorID, expectedCreatedAt, expectedUpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(expectedAuthorID).
		WillReturnRows(rows)

	repo := NewMysqlArticleRepository(db)

	articles, err := repo.GetByAuthorID(expectedAuthorID)
	assert.NoError(t, err)
	assert.NotNil(t, articles)
}

func TestMysqlArticleRepository_GetByAuthorID_Error(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "SELECT * FROM `articles` WHERE author_id = ?"

	expectedAuthorID := uint(1)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(expectedAuthorID).
		WillReturnError(assert.AnError)

	repo := NewMysqlArticleRepository(db)

	articles, err := repo.GetByAuthorID(expectedAuthorID)
	assert.Error(t, err)
	assert.Nil(t, articles)
}

func TestMysqlArticleRepository_GetByTitle(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "SELECT * FROM `articles` WHERE title LIKE ?"

	expectedTitle := "title"
	expectedContent := "content"
	expectedAuthorID := 1
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
		AddRow(1, expectedTitle, expectedContent, expectedAuthorID, expectedCreatedAt, expectedUpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("%" + expectedTitle + "%").
		WillReturnRows(rows)

	repo := NewMysqlArticleRepository(db)

	articles, err := repo.GetByTitle(expectedTitle)
	assert.NoError(t, err)
	assert.NotNil(t, articles)
}

func TestMysqlArticleRepository_GetByTitle_Error(t *testing.T) {
	db, mock, err := mockDBConnection()
	assert.NoError(t, err)

	query := "SELECT * FROM `articles` WHERE title LIKE ?"

	expectedTitle := "title"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("%" + expectedTitle + "%").
		WillReturnError(assert.AnError)

	repo := NewMysqlArticleRepository(db)

	articles, err := repo.GetByTitle(expectedTitle)
	assert.Error(t, err)
	assert.Nil(t, articles)
}
