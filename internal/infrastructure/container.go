package infrastructure

import (
	"github.com/caarlos0/env/v10"
	"go-clean-architecture/internal/article"
	"go-clean-architecture/internal/author"
	"go-clean-architecture/internal/config"
	"go-clean-architecture/internal/domain"
	"go-clean-architecture/pkg/xlogger"
)

var (
	cfg config.Config

	authorRepository  domain.AuthorRepository
	articleRepository domain.ArticleRepository

	articleService domain.ArticleService
)

func init() {
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	xlogger.Setup(cfg)
	dbSetup()

	authorRepository = author.NewMysqlAuthorRepository(db)
	articleRepository = article.NewMysqlArticleRepository(db)

	articleService = article.NewArticleService(articleRepository, authorRepository)
}
