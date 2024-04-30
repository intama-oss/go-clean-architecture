# Go Clean Architecture

Go Clean Architecture adalah sebuah template project yang dibangun menggunakan bahasa pemrograman Go. Project
ini dibangun menggunakan konsep Clean Architecture yang diperkenalkan oleh Robert C. Martin. Konsep ini
memisahkan antara business logic, delivery mechanism, dan data storage.

## Konsep

Project ini dibangun menggunakan konsep Clean Architecture yang terdiri dari beberapa layer, yaitu:

### Domain / Business Logic

Layer ini berisi business logic dari aplikasi. Layer ini tidak boleh bergantung pada layer lainnya. Layer ini
berisi model, repository interface, dan service interface. Contoh dari layer ini adalah `domain/author.go` dan
`domain/article.go`.

### Repository

Layer ini berisi implementasi dari repository interface yang didefinisikan di layer domain. Layer ini berisi
implementasi untuk mengakses data dari storage. Contoh dari layer ini adalah `<domain>/mysql_repository.go`. <domain>
adalah nama domain yang didefinisikan di layer domain. Dalam contoh ini, <domain> adalah `author` atau `article`.

Penamaan file di layer ini adalah `<storage>_repository.go`. <storage> adalah jenis storage yang digunakan, seperti
`mysql`, `mongodb`, `redis`, dan lain-lain. Contoh dari layer ini adalah `mysql_repository.go`.

### Service / Usecase

Layer ini berisi implementasi dari service interface yang didefinisikan di layer domain. Layer ini berisi implementasi
untuk business logic dari aplikasi. Contoh dari layer ini adalah `<domain>/service.go`. Sama seperti
repository, <domain>
adalah nama domain yang didefinisikan di layer domain. Dalam contoh ini, <domain> adalah `author` atau `article`.

Penamaan file di layer ini cukup `service.go` dan diletakkan di dalam folder domain, contoh: `article/service.go`.

### Delivery / Presenter / Handler

Layer ini berisi implementasi untuk mengirimkan data ke client. Layer ini berisi implementasi untuk mengakses data dari
client, seperti HTTP, gRPC, dan lain-lain. Penamaan file di layer ini adalah `http_handler.go` atau `grpc_handler.go`.

Penamaan file di layer ini adalah `<delivery>_handler.go`. <delivery> adalah jenis delivery yang digunakan,
seperti `http`,
`grpc`, dan lain-lain. Kemudian untuk penempatan file, file ini diletakkan di dalam folder domain,
contoh: `article/http_handler.go`.

### Infrastructure

Layer ini berisi semua implementasi yang berhubungan dengan infrastruktur, seperti database, cache, router, dan
lain-lain.
Contoh dari layer ini adalah `gorm.go` dan `fiber.go`. Penamaan file di layer ini adalah `<infra>.go`. <infra> adalah
jenis
infrastruktur yang digunakan, seperti `gorm`, `fiber`, dan lain-lain.

### Config

Layer ini berisi code yang berhubungan dengan konfigurasi aplikasi. Layer ini berisi implementasi untuk mengakses
konfigurasi
aplikasi. Penamaan file di layer ini adalah `config.go`.

### Middleware

Layer ini berisi code yang berhubungan dengan middleware aplikasi. Layer ini berisi implementasi untuk middleware
aplikasi.

### Utilities

Layer ini berisi code yang berhubungan dengan utilities aplikasi. Layer ini berisi implementasi untuk utilities
aplikasi.

## Library

Project ini menggunakan beberapa library, yaitu:

- [Fiber](https://gofiber.io) sebagai web framework untuk membuat aplikasi web.
- [GORM](https://gorm.io) sebagai ORM untuk mengakses database.
- [Gorm SQLite](https://github.com/glebarez/sqlite) sebagai library untuk menggunakan SQLite dengan GORM.
- [Gorm MySQL](https://gorm.io/driver/mysql) sebagai library untuk menggunakan MySQL dengan GORM.
- [SQL Mock](https://github.com/DATA-DOG/go-sqlmock) sebagai library untuk mocking SQL.
- [carlos0/env](https://github.com/caarlos0/env) sebagai library untuk mengakses environment.
- [rs/zerolog](https://github.com/rs/zerolog) sebagai library untuk logging zero allocation.
- [stretchr/testify](https://github.com/stretchr/testify) sebagai library untuk testing.
- [go-faker/faker](https://github.com/go-faker/faker) sebagai library untuk membuat data palsu pada testing.
- [go-playground/validator](https://github.com/go-playground/validator) sebagai library untuk validasi data.
- [swaggo/swag](https://github.com/swaggo/swag) sebagai library untuk generate swagger.

## Struktur Folder

Struktur folder dari project ini adalah sebagai berikut:

```plaintext
.
├── cmd
│   └── app
│       └── main.go
├── docs
│   ├── docs.go
│   └── swagger.json
├── internal
│   ├── domain
│   │   ├── article.go
│   │   ├── author.go
│   │   └── config.go
│   ├── middleware
│   │   └── <middleware-name>
│   │       └── <middleware-name>.go
│   ├── author
│   │   └── mysql_repository.go
│   ├── article
│   │   ├── http_handler.go
│   │   ├── mysql_repository.go
│   │   └── service.go
│   ├── <domain>
│   │   ├── http_handler.go
│   │   ├── middleware.go
│   │   ├── <storage>_repository.go
│   │   └── service.go
│   ├── config
│   │   └── config.go
│   ├── infrastructure
│   │   ├── gorm.go
│   │   └── fiber.go
│   └── utilities
│       └── <utility-name>.go
├── mocks
│   ├── <domain>_repository.go
│   └── <domain>_service.go
├── go.mod
└── go.sum
```

## Cara Penggunaan

Project ini menggunakan Go Modules, sehingga tidak perlu melakukan `go get` untuk mengunduh library yang digunakan.

Untuk menjalankan project ini, silahkan jalankan perintah berikut:

```bash
go run cmd/app/main.go
```

Untuk melakukan build project ini, silahkan jalankan perintah berikut:

```bash
go build -o bin/<app-name> cmd/app/main.go
```

> Pastikan untuk mengganti `<app-name>` dengan nama aplikasi yang diinginkan. Jika anda menggunakan Windows, maka
> perintah di atas akan menghasilkan file `bin/<app-name>.exe`.

## Environment

Daftar environment yang digunakan pada project ini.

| Key               | Description                          | Example                                                                                              | Default                                  |
|-------------------|--------------------------------------|------------------------------------------------------------------------------------------------------|------------------------------------------|
| `Host`            | Alamat untuk binding service         | `localhost`                                                                                          |                                          |
| `Port`            | Port untuk binding service           | `3000`                                                                                               | `3000`                                   |
| `IS_DEVELOPMENT`  | Mode development                     | `true`                                                                                               | `false`                                  |
| `PROXY_HEADER`    | Header untuk mendapatkan IP asli     | `X-Real-IP` atau `X-Forwarded-For`                                                                   |                                          |
| `LOG_FIELDS`      | Field yang akan ditampilkan pada log | `method,path,ip` lihat [disini](https://github.com/gofiber/contrib/blob/main/fiberzerolog/config.go) | `latency,status,method,url,error`        |
| `DATABASE_DRIVER` | Driver database                      | `mysql` atau `sqlite`                                                                                | `sqlite` (in memory)                     |
| `DATABASE_DSN`    | Data source name database            | `user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local`                  | `file::memory:?cache=shared` (in memory) |

## Testing

Project ini menggunakan library [stretchr/testify](https://github.com/stretchr/testify) sebagai library untuk testing. Untuk
melakukan testing, silahkan jalankan perintah berikut:

```bash
go test ./...
```

## Swagger Generation

Project ini menggunakan library [swaggo/swag](https://github.com/swaggo/swag) sebagai library untuk generate swagger. Untuk
mengenerate swagger, silahkan jalankan perintah berikut:

```bash
go generate ./...
```
