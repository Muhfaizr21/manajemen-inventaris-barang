# manajemen-inventaris-barang

![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?style=for-the-badge&logo=go)
![Gin Framework](https://img.shields.io/badge/Gin-Gonic-008ECF?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15%2B-316192?style=for-the-badge&logo=postgresql)
![GORM](https://img.shields.io/badge/ORM-GORM-00ADD8?style=for-the-badge)
![Architecture](https://img.shields.io/badge/Architecture-Layered%20MVC%20%26%20SOLID-success?style=for-the-badge)

RESTful API **Sistem Manajemen Inventaris Barang** yang dibangun menggunakan **Golang (Gin Framework)**, **GORM**, dan **PostgreSQL**. Proyek ini dirancang dengan standar **Production-Ready Clean Code**, prinsip **SOLID**, **Constructor Dependency Injection**, serta pengamanan transaksi konkurensi (**Pessimistic Row-Level Locking**).

---

## Daftar Isi

1. [Fitur Utama](#-1-fitur-utama)
2. [Arsitektur & Prinsip Desain (Clean Code & SOLID)](#️-2-arsitektur--prinsip-desain-clean-code--solid)
3. [Spesifikasi Endpoint API](#-3-spesifikasi-endpoint-api)
4. [Dokumentasi API (Postman & OpenAPI/Swagger)](#-4-dokumentasi-api-postman--openapiswagger)
5. [Panduan Instalasi & Menjalankan Secara Lokal](#-5-panduan-instalasi--menjalankan-secara-lokal)
6. [Jawaban Pertanyaan Pemahaman Teknis (WAJIB)](#-6-jawaban-pertanyaan-pemahaman-teknis-wajib)
   - [6.1 Alur Perjalanan Request & Alasan Pemisahan Layer](#61-alur-perjalanan-request--alasan-pemisahan-layer)
   - [6.2 Keamanan & Penyimpanan Token JWT di Frontend](#62-keamanan--penyimpanan-token-jwt-di-frontend)
   - [6.3 Penanganan Konkurensi (Race Condition Stok Barang)](#63-penanganan-konkurensi-race-condition-stok-barang)

---

## 🚀 1. Fitur Utama

- **Authentication & Authorization (JWT)**:
  - Registrasi pengguna baru dengan enkripsi password menggunakan algoritma `bcrypt`.
  - Login pengguna yang menghasilkan **JSON Web Token (JWT)** (`HS256`) berbatas waktu.
  - Middleware proteksi endpoint berbasis `Bearer Token`.
- **Manajemen Inventaris Barang (CRUD)**:
  - **Create (`POST /items`)**: Menambahkan barang baru (`name`, `description`, `stock`, `price`) disertai pencatatan pembuat (`created_by`).
  - **Read (`GET /items` & `GET /items/:id`)**: Melihat daftar barang dilengkapi fitur **Search** (*case-insensitive* pada nama & deskripsi menggunakan PostgreSQL `ILIKE`) dan **Pagination** (`page`, `limit`, `total_items`, `total_pages`).
  - **Update (`PUT /items/:id`)**: Memperbarui informasi dan stok barang secara aman di dalam transaksi database dengan penguncian baris (*row-level lock*).
  - **Delete (`DELETE /items/:id`)**: Menghapus barang (*soft-delete*) dari daftar aktif.
- **Validasi Data Berlapis (*Defense in Depth*)**:
  - **Lapis 1 (DTO Validator)**: Validasi tipe data, format email, dan aturan `stock >= 0` & `price >= 0` di level HTTP Request.
  - **Lapis 2 (Service Guard)**: Validasi aturan bisnis di layer Service sebelum menyentuh Repository.
  - **Lapis 3 (Database Constraint)**: Constraint `CHECK (stock >= 0)` dan `CHECK (price >= 0)` langsung pada skema tabel PostgreSQL.

---

## 🏗️ 2. Arsitektur & Prinsip Desain (Clean Code & SOLID)

Kode sumber berada di dalam direktori [`backend/`](./backend) dengan struktur **Layered MVC Architecture**:

```text
manajemen-inventaris-barang/
├── backend/
│   ├── config/
│   │   └── config.go               # Konfigurasi environment & inisialisasi koneksi PostgreSQL
│   ├── controllers/
│   │   ├── auth_controller.go      # Controller: HTTP handler untuk Register & Login
│   │   └── item_controller.go      # Controller: HTTP handler untuk CRUD Inventaris Barang
│   ├── docs/
│   │   ├── openapi.yaml            # Dokumentasi spesifikasi OpenAPI 3.0 (Swagger)
│   │   └── postman_collection.json # Postman Collection siap pakai (auto-save JWT token)
│   ├── middlewares/
│   │   └── auth_middleware.go      # Middleware autentikasi Bearer JWT
│   ├── models/
│   │   ├── item.go                 # Model & DTO validasi entitas Barang (Item)
│   │   └── user.go                 # Model & DTO validasi entitas Pengguna (User)
│   ├── repositories/
│   │   ├── item_repository.go      # Repository: Query database & locking transaksi Item
│   │   └── user_repository.go      # Repository: Query database User
│   ├── routes/
│   │   └── routes.go               # Definisi routing & pengelompokan middleware
│   ├── services/
│   │   ├── auth_service.go         # Service: Logika bisnis autentikasi, bcrypt, & JWT
│   │   └── item_service.go         # Service: Logika bisnis inventaris & kalkulasi paginasi
│   ├── views/
│   │   └── response_view.go        # View: Standarisasi struktur JSON response & formatter validasi
│   ├── .env.example                # Contoh konfigurasi environment variable
│   ├── go.mod
│   ├── go.sum
│   └── main.go                     # Entry point & Dependency Injection Wiring
└── README.md
```

### Penerapan Prinsip SOLID:
1. **Single Responsibility Principle (SRP)**:
   - `Controller` hanya menangani protokol HTTP (parsing parameter/body dan memanggil `View`).
   - `Service` hanya menangani logika bisnis murni (hashing, pembuatan token, kalkulasi paginasi, aturan domain).
   - `Repository` hanya bertanggung jawab atas komunikasi ke PostgreSQL melalui GORM.
2. **Open/Closed Principle (OCP)**:
   - Setiap layer berkomunikasi melalui kontrak `interface`. Penambahan perilaku baru (seperti caching layer atau audit logging) dapat ditambahkan melalui dekorator implementasi tanpa memodifikasi kode pemanggil.
3. **Liskov Substitution Principle (LSP)**:
   - Implementasi konkret (`userRepository`, `itemRepository`, `authService`, `itemService`) sepenuhnya memenuhi kontrak interface-nya dan dapat digantikan dengan *mock object* saat pengujian.
4. **Interface Segregation Principle (ISP)**:
   - Interface dirancang spesifik dan ramping sesuai domain masing-masing (`UserRepository`, `ItemRepository`, `AuthService`, `ItemService`).
5. **Dependency Inversion Principle (DIP)**:
   - Modul level tinggi (`Controller`) tidak bergantung pada modul level rendah secara langsung, melainkan pada abstraksi (`Service` & `Repository` interface) yang dirakit di [`backend/main.go`](./backend/main.go) melalui **Constructor Dependency Injection** (tanpa variabel global DB).

---

## 📌 3. Spesifikasi Endpoint API

| Method | Endpoint | Proteksi | Deskripsi |
| :--- | :--- | :---: | :--- |
| `POST` | `/register` | Publik | Registrasi akun pengguna baru |
| `POST` | `/login` | Publik | Login pengguna & mendapatkan token JWT |
| `POST` | `/items` | **Bearer JWT** | Menambahkan barang baru ke inventaris |
| `GET` | `/items` | **Bearer JWT** | Mengambil daftar barang (`?search=...&page=1&limit=10`) |
| `GET` | `/items/:id` | **Bearer JWT** | Mengambil detail satu barang berdasarkan ID |
| `PUT` | `/items/:id` | **Bearer JWT** | Memperbarui data & stok barang berdasarkan ID |
| `DELETE` | `/items/:id` | **Bearer JWT** | Menghapus barang berdasarkan ID |

> **Catatan**: Seluruh endpoint di atas juga tersedia dengan prefix `/api` (contoh: `/api/register`, `/api/login`, `/api/items`) untuk kemudahan standarisasi API Gateway.

### Contoh Request & Response

#### 1. Login (`POST /login`)
**Request Body:**
```json
{
  "email": "faiz@example.com",
  "password": "rahasia123"
}
```
**Response (`200 OK`):**
```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "Muhammad Faiz",
      "email": "faiz@example.com",
      "created_at": "2026-10-30T13:18:04+07:00",
      "updated_at": "2026-10-30T13:18:04+07:00"
    }
  }
}
```

#### 2. Tambah Barang (`POST /items`)
**Headers:** `Authorization: Bearer <token>`
**Request Body:**
```json
{
  "name": "MacBook Pro M3 14 Inch",
  "description": "Laptop operasional tim engineering RAM 16GB SSD 512GB",
  "stock": 15,
  "price": 28500000
}
```
**Response (`201 Created`):**
```json
{
  "success": true,
  "message": "Barang berhasil ditambahkan",
  "data": {
    "id": 1,
    "name": "MacBook Pro M3 14 Inch",
    "description": "Laptop operasional tim engineering RAM 16GB SSD 512GB",
    "stock": 15,
    "price": 28500000,
    "created_by": 1,
    "created_at": "2026-10-30T13:18:04+07:00",
    "updated_at": "2026-10-30T13:18:04+07:00"
  }
}
```

#### 3. Contoh Penolakan Stok Negatif (`POST /items` dengan `stock: -5`)
**Response (`400 Bad Request`):**
```json
{
  "success": false,
  "message": "Validasi data barang gagal",
  "errors": {
    "stock": "Field 'stock' tidak boleh kurang dari 0 (tidak boleh negatif)"
  }
}
```

#### 4. Lihat Daftar Barang dengan Search & Pagination (`GET /items?search=MacBook&page=1&limit=10`)
**Response (`200 OK`):**
```json
{
  "success": true,
  "message": "Berhasil mengambil daftar barang",
  "data": [
    {
      "id": 1,
      "name": "MacBook Pro M3 14 Inch",
      "description": "Laptop operasional tim engineering RAM 16GB SSD 512GB",
      "stock": 15,
      "price": 28500000,
      "created_by": 1,
      "created_at": "2026-10-30T13:18:04+07:00",
      "updated_at": "2026-10-30T13:18:04+07:00"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

---

## 📄 4. Dokumentasi API (Postman & OpenAPI/Swagger)

Seluruh file dokumentasi API tersedia di folder [`backend/docs/`](./backend/docs):
1. **Postman Collection**: [`backend/docs/postman_collection.json`](./backend/docs/postman_collection.json)
   - Import file tersebut ke Postman.
   - Endpoint **Login User** sudah dilengkapi *Post-Response Test Script* otomatis yang langsung menyimpan token JWT ke Collection Variable `{{token}}`, sehingga seluruh endpoint CRUD Barang dapat langsung diuji dengan sekali klik.
2. **OpenAPI 3.0 (Swagger)**: [`backend/docs/openapi.yaml`](./backend/docs/openapi.yaml)
   - Spesifikasi standar OpenAPI 3.0.3 yang dapat divisualisasikan di Swagger UI / Swagger Editor.

---

## 💻 5. Panduan Instalasi & Menjalankan Secara Lokal

Tersedia dua cara menjalankan: **Docker Compose** (direkomendasikan, tanpa install Go/PostgreSQL manual) atau **Manual** (Go + PostgreSQL lokal).

### Langkah 1: Clone Repository
```bash
git clone https://github.com/Muhfaizr21/manajemen-inventaris-barang.git
cd manajemen-inventaris-barang
```

### Opsi A: Jalankan dengan Docker Compose

**Prasyarat:** [Docker Desktop](https://www.docker.com/products/docker-desktop/) / Docker Engine + Compose.

```bash
docker compose up -d --build
```

Perintah ini otomatis:
1. Membangun image backend dari [`backend/Dockerfile`](./backend/Dockerfile) (multi-stage build Go → Alpine).
2. Menjalankan container **PostgreSQL 16** dengan healthcheck dan persistent volume `pgdata`.
3. Menjalankan container **API** yang menunggu database *healthy* sebelum start.
4. GORM melakukan **Auto-Migrate** tabel `users` dan `items` (termasuk constraint `CHECK (stock >= 0)`).

Server siap diakses pada `http://localhost:8080`.

**Perintah berguna:**
```bash
docker compose logs -f api      # Lihat log backend
docker compose down             # Stop & hapus container (data DB tetap di volume)
docker compose down -v --remove-orphans   # Stop + HAPUS volume database
```

**Override konfigurasi (opsional):**
```bash
DB_USER=postgres DB_PASS=rahasia JWT_SECRET=custom-secret-prod docker compose up -d --build
```

### Opsi B: Jalankan Manual (Tanpa Docker)

**Prasyarat**
- **Golang** `1.21` atau versi lebih baru
- **PostgreSQL** `13+` berjalan di mesin lokal

**Langkah 2: Buat Database PostgreSQL**
Pastikan PostgreSQL sudah aktif, lalu buat database bernama `tes_backend`:
```sql
CREATE DATABASE tes_backend;
```

**Langkah 3: Konfigurasi Environment Variable**
Aplikasi memiliki nilai konfigurasi bawaan (*default fallback*) yang siap jalan sebagaimana tercantum pada [`backend/.env.example`](./backend/.env.example):

| Variabel | Nilai Default | Deskripsi |
| :--- | :--- | :--- |
| `PORT` | `8080` | Port HTTP Server |
| `DB_HOST` | `localhost` | Host PostgreSQL |
| `DB_PORT` | `5432` | Port PostgreSQL |
| `DB_USER` | `muhfaiizr` | Username PostgreSQL |
| `DB_PASS` | `admin` | Password PostgreSQL |
| `DB_NAME` | `tes_backend` | Nama Database PostgreSQL |
| `JWT_SECRET` | `super-secret-inventory-jwt-key-2026` | Secret key untuk signing JWT |
| `JWT_EXPIRY_HOURS` | `24` | Masa berlaku token JWT (jam) |

Apabila kredensial PostgreSQL lokal Anda berbeda, Anda dapat meng-export variabel environment sebelum menjalankan aplikasi:
```bash
export DB_USER=postgres
export DB_PASS=password_anda
export DB_NAME=tes_backend
```

**Langkah 4: Jalankan Aplikasi**
```bash
cd backend
go mod tidy
go run main.go
```
GORM akan secara otomatis melakukan **Auto-Migrate** tabel `users` dan `items` beserta constraint `CHECK (stock >= 0)`. Server siap diakses pada `http://localhost:8080`.

---

## 🧠 6. Jawaban Pertanyaan Pemahaman Teknis (WAJIB)

### 6.1 Alur Perjalanan Request & Alasan Pemisahan Layer

#### A. Alur Perjalanan Request (Dari Client hingga Tersimpan di Database)
Berikut adalah siklus hidup sebuah HTTP request (misalnya `POST /items`):

```text
[Client / Frontend]
       │  HTTP POST /items (Header: Authorization: Bearer <JWT>, Body: JSON)
       ▼
1. [Router Layer] (backend/routes/routes.go)
       │  Mencocokkan method & URI path, lalu meneruskan request ke rantai Middleware grup `/items`.
       ▼
2. [Middleware Layer] (backend/middlewares/auth_middleware.go)
       │  Mengekstrak Bearer token dari header Authorization dan memvalidasi signature serta klaim
       │  kedaluwarsa (`exp`) melalui `AuthService.ValidateToken`.
       │  • Jika tidak valid/kosong: Menghentikan siklus (`c.Abort()`) dan merespons `401 Unauthorized`.
       │  • Jika valid: Menyimpan `user_id` ke dalam `gin.Context` dan memanggil `c.Next()`.
       ▼
3. [Controller Layer] (backend/controllers/item_controller.go)
       │  Melakukan deserialisasi JSON body ke struct DTO `CreateItemRequest` sekaligus menjalankan
       │  validasi skema input (`required`, `gte=0`).
       │  • Jika validasi gagal: Memanggil `views.JSONError` untuk merespons `400 Bad Request`.
       │  • Jika lolos: Mengambil `user_id` dari context dan memanggil `ItemService.CreateItem(ctx, userID, req)`.
       ▼
4. [Service Layer] (backend/services/item_service.go)
       │  Mengeksekusi aturan bisnis domain (normalisasi/trim string, memastikan stok & harga tidak negatif,
       │  mengasosiasikan `CreatedBy` dengan `userID`), merakit entitas `models.Item`, lalu memanggil
       │  `ItemRepository.Create(ctx, item)`.
       ▼
5. [Repository Layer] (backend/repositories/item_repository.go)
       │  Meneruskan `context.Context` ke GORM (`r.db.WithContext(ctx).Create(item)`) untuk mengeksekusi
       │  parameterized SQL query `INSERT INTO "items" ... RETURNING "id"` ke PostgreSQL.
       ▼
6. [PostgreSQL Database]
       │  Memverifikasi constraint tabel (`CHECK (stock >= 0)`), menyimpan baris ke disk/WAL, dan
       │  mengembalikan ID baru ke Repository -> Service -> Controller.
       ▼
7. [View Layer] (backend/views/response_view.go)
          Controller memanggil `views.JSONSuccess` untuk membungkus data ke format standar `APIResponse`
          dengan status `201 Created` dan mengirimkannya kembali ke Client.
```

#### B. Mengapa Memisahkan Logic ke Dalam Layer-Layer Tersebut?
1. **Separation of Concerns (SoC) & Maintainability**: Setiap layer memiliki satu tanggung jawab spesifik (*Single Responsibility Principle*). Perubahan pada format JSON atau framework HTTP hanya berdampak pada layer `Controller`/`View`, sedangkan perubahan pada struktur query atau migrasi database hanya berdampak pada `Repository` tanpa merusak logika bisnis di `Service`.
2. **High Testability (Mudah Di-Unit Test)**: Dengan menghubungkan antar-layer menggunakan kontrak `interface` (*Dependency Inversion Principle*), kita dapat menguji logika bisnis di `ItemService` dan `AuthService` menggunakan *mock repository* secara cepat di memori tanpa perlu menyalakan database PostgreSQL.
3. **Reusability & Transport Agnostic**: Karena `Service` tidak mengimpor atau bergantung pada `*gin.Context`, logika bisnis yang sama dapat dipanggil kembali oleh antarmuka lain di masa depan (seperti CLI command, cron job, message queue consumer RabbitMQ/Kafka, atau gRPC server) tanpa menduplikasi kode.

---

### 6.2 Keamanan & Penyimpanan Token JWT di Frontend

#### Rekomendasi: **HttpOnly, Secure, SameSite Cookie** (Dikombinasikan dengan In-Memory Short-Lived Access Token)

Untuk aplikasi web modern berstandar produksi, menyimpan token autentikasi di **`HttpOnly Cookie`** jauh lebih aman dibandingkan menyimpannya di **`Local Storage` (`localStorage` / `sessionStorage`)**.

#### Analisis Perbandingan & Risiko Keamanan:

1. **Mengapa `Local Storage` Berisiko Tinggi terhadap XSS (Cross-Site Scripting)?**
   - `localStorage` dapat diakses secara bebas oleh kode JavaScript apa pun yang berjalan pada domain yang sama melalui `window.localStorage.getItem('token')`.
   - Apabila aplikasi frontend memiliki satu saja celah **XSS** (misalnya dari sanitasi input yang kurang sempurna, atau serangan *supply-chain* pada paket pihak ketiga `npm` yang disusupi script berbahaya), penyerang dapat langsung mengekstrak seluruh token JWT milik user dan mengirimkannya ke server eksternal (*account takeover*).

2. **Mengapa `HttpOnly Cookie` Lebih Aman & Bagaimana Mitigasi Risikonya?**
   - **Imun terhadap Pencurian Token via XSS**: Cookie yang disetel oleh server dengan atribut `HttpOnly` **tidak dapat dibaca oleh JavaScript** (`document.cookie` tidak akan menampilkan cookie tersebut). Sekalipun terjadi celah XSS, script penyerang tidak bisa mencuri string token secara langsung.
   - **Enkripsi Transmisi (`Secure` Flag)**: Dengan menambahkan flag `Secure`, browser menjamin cookie hanya dikirim melalui koneksi terenkripsi HTTPS.
   - **Risiko CSRF (Cross-Site Request Forgery) & Mitigasinya**: Kelemahan tradisional cookie adalah browser otomatis mengirimkannya pada setiap request ke domain tujuan, yang membuka celah CSRF. Namun, pada browser modern risiko ini dapat **dimitigasi secara tuntas** dengan:
     1. Menambahkan atribut **`SameSite=Strict`** atau **`SameSite=Lax`** pada cookie.
     2. Menerapkan pola **Refresh Token di `HttpOnly Cookie`** + **Short-Lived Access Token (5–15 menit) di dalam memori variabel JavaScript (React/Vue State)** yang dikirim via header `Authorization: Bearer <token>`. Pendekatan *hybrid* ini menggabungkan ketahanan terhadap XSS sekaligus ketahanan penuh terhadap CSRF.

---

### 6.3 Penanganan Konkurensi (Race Condition Stok Barang)

#### Skenario Masalah (*Race Condition* / *Lost Update*):
Barang **X** memiliki sisa `stock = 1`. Pada milidetik yang **tepat bersamaan**, **User A** dan **User B** mengirim request untuk mengurangi stok Barang X sebanyak `1`.
Jika aplikasi hanya membaca stok di memori Go (`stock == 1`), lalu mengurangi `1 - 1` dan melakukan `UPDATE`, kedua request sama-sama melihat `stock = 1` sebelum salah satu sempat melakukan `COMMIT`. Akibatnya, kedua request lolos dan stok akhir bisa menjadi `-1` (atau terjadi *overselling*).

#### Solusi 3 Lapis yang Diterapkan pada Proyek Ini:

##### 1. Pessimistic Row-Level Locking (`SELECT ... FOR UPDATE`) di dalam Database Transaction
Di dalam [`backend/repositories/item_repository.go`](./backend/repositories/item_repository.go) (pada method `AdjustStockAtomically` dan `Update`), pembacaan dan pembaruan stok dibungkus dalam satu transaksi ACID menggunakan klausa eksklusif `FOR UPDATE`:

```go
err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // 1. Mengunci baris spesifik di PostgreSQL: SELECT * FROM items WHERE id = $1 FOR UPDATE
    if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&updatedItem, id).Error; err != nil {
        return err
    }

    // 2. Mengecek ketersediaan stok setelah lock didapat
    if updatedItem.Stock + delta < 0 {
        return ErrInsufficientStock
    }

    // 3. Melakukan pembaruan stok secara atomik
    return tx.Model(&updatedItem).
        Where("id = ? AND stock + ? >= 0", id, delta).
        Update("stock", gorm.Expr("stock + ?", delta)).Error
})
```
**Cara Kerjanya saat 2 Request Bersamaan:**
- Ketika **Request A** dan **Request B** masuk bersamaan, PostgreSQL hanya mengizinkan **satu transaksi** (misal Request A) memperoleh *Row-Level Exclusive Lock* pada baris Barang X.
- **Request B** akan **di-block (antre menunggu)** di level database sampai transaksi Request A selesai (`COMMIT` atau `ROLLBACK`).
- Setelah Request A mengurangi stok dari `1` menjadi `0` dan melakukan `COMMIT`, kunci dilepaskan.
- Request B kemudian melanjutkan eksekusi `SELECT ... FOR UPDATE` dan membaca nilai terbaru yang sudah ter-commit yaitu **`stock = 0`**.
- Pengecekan `if updatedItem.Stock - 1 < 0` pada Request B langsung bernilai `true`, sehingga transaksi Request B otomatis di-`ROLLBACK` dan mengembalikan error `stok barang tidak mencukupi` tanpa pernah menyentuh angka `-1`.

##### 2. Atomic Conditional Update (`UPDATE ... WHERE stock >= jumlah`)
Selain penguncian baris, query perubahan stok menggunakan ekspresi SQL atomik disertai klausa guard:
```sql
UPDATE items
SET stock = stock - 1, updated_at = NOW()
WHERE id = 1 AND stock - 1 >= 0;
```
Karena operasi `UPDATE` tunggal bersifat atomik di mesin PostgreSQL, apabila stok bernilai `0`, kondisi `WHERE stock - 1 >= 0` tidak terpenuhi (`RowsAffected == 0`), sehingga stok mustahil berkurang menjadi `-1`.

##### 3. Database-Level `CHECK` Constraint (`CHECK (stock >= 0)`)
Pada definisi model [`backend/models/item.go`](./backend/models/item.go), kolom `stock` memiliki tag `gorm:"not null;default:0;check:stock >= 0"`. Ini membuat aturan permanen di level engine PostgreSQL. Sekalipun ada query eksternal yang mencoba memaksa nilai negatif, PostgreSQL akan langsung menolak transaksi tersebut dengan *constraint violation error*.
