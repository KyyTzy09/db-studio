# Todo - DBStudio (Go + SvelteKit)

## Legend
- 🔴 P0 Critical path, MVP blocker
- 🟡 P1 High value, not MVP blocker
- 🟢 P2 Nice to have
- 🔵 P3 Future

---

## Phase 0: Project Initialization & Monorepo Setup
### 0.1 Backend Setup (Go)
- [x] 🔴 Inisialisasi Go module (`go mod init db-studio-go`)
- [x] 🔴 Install dependencies CLI & Config: `github.com/spf13/cobra`, `github.com/spf13/viper`
- [x] 🔴 Install dependencies Router: `github.com/go-chi/chi/v5`
- [x] 🔴 Install dependencies Database Drivers: `github.com/jackc/pgx/v5` (Postgres), `github.com/go-sql-driver/mysql` (MySQL), `modernc.org/sqlite` (SQLite)
- [x] 🔴 Setup arsitektur folder Go (Modular Pipeline):
  - `cmd/` (CLI Commands: `root.go`, `connect.go`, `doctor.go`, `version.go`)
  - `internal/config/` (OS Global Path Config Manager)
  - `internal/scanner/` (Auto-detection Scanners)
  - `internal/wizard/` (CLI Interactive Wizard Fallback)
  - `internal/db/` (Database Interface & Drivers)
  - `internal/http/` (Chi REST API Router & Middleware)
- [x] 🔴 Konfigurasi OS global path config (`~/.config/dbstudio/connections.json`)

### 0.2 Frontend Setup (SvelteKit)
- [x] 🔴 Inisialisasi SvelteKit project di subfolder `web/` (`bun create svelte@latest web`)
- [x] 🔴 Konfigurasi `@sveltejs/adapter-static` untuk SPA (Single Page Application) build target di `web/vite.config.ts`
- [x] 🔴 Install dependencies UI: `@tailwindcss/vite` / `tailwindcss`, `@tanstack/svelte-table`, `@codemirror/state`, `@codemirror/view`, `@codemirror/lang-sql`
- [ ] 🔴 Setup konversi komponen visual (Dumb Presentational Components) dan komponen pengambil data (Smart Components / Stores)

### 0.3 Build & Embed Integration
- [x] 🔴 Setup `go:embed` pada `main.go` dan handler Chi untuk membaca folder build statis dari SvelteKit (`web/build/`)
- [ ] 🔴 Buat Makefile / build script otomatisasi: `npm --prefix web run build` ➔ `go build -o bin/dbstudio main.go`

---

## Phase 1: Modular Config, Scanner & CLI Wizard Pipeline
### 1.1 Global Config Manager (`internal/config`)
- [x] 🔴 Implementasi pembacaan/penulisan `connections.json` di global OS directory (`os.UserConfigDir()`)
- [x] 🔴 Terapkan pemetaan koneksi berdasarkan hash/path absolut dari direktori proyek saat ini (`cwd`)

### 1.2 Auto-Detection Scanner Pipeline (`internal/scanner`)
- [x] 🔴 `SavedConfigScanner`: Cek apakah proyek ini sudah tersimpan di global config. Jika ADA ➔ langsung gunakan.
- [x] 🔴 `EnvScanner`: Parsing file `.env` lokal (deteksi `DATABASE_URL`, `DB_HOST`, `DB_USER`, `DB_PASS`, `DB_NAME`, `DB_PORT`)
- [ ] 🟡 `DockerComposeScanner`: Parsing `docker-compose.yml` / `compose.yaml` untuk mengekstrak kredensial database container
- [x] 🔴 Logic Handler Hasil Scanning:
  - Jika **1 Config Ditemukan**: Otomatis pilih & jalankan server.
  - Jika **>1 Config Ditemukan**: Tampilkan pilihan prompt konfirmasi di CLI.
  - Jika **0 Config Ditemukan**: Pemicu fallback ke CLI Wizard.

### 1.3 CLI Interactive Wizard (`internal/wizard`)
- [x] 🔴 Implementasi prompt interaktif di terminal (Tipe DB: Postgres/MySQL/SQLite, Host, Port, User, Pass, DB Name) saat auto-detect gagal
- [x] 🔴 Simpan hasil input wizard ke global config untuk penggunaan berikutnya

### 1.4 CLI Commands & Browser Launcher (`cmd/`)
- [x] 🔴 Root Command `dbstudio`: Eksekusi Pipeline (Saved Config ➔ Scanner ➔ Wizard Fallback ➔ Start Chi Server ➔ Buka Browser)
- [ ] 🔴 Command `dbstudio connect`: Langsung memicu CLI Wizard untuk menambah/mengedit koneksi
- [ ] 🔴 Command `dbstudio doctor`: Mengecek *health* koneksi database yang tersimpan
- [ ] 🔴 Command `dbstudio version`: Menampilkan informasi versi build
- [x] 🔴 Utility *Browser Launcher* cross-platform (Windows `rundll32`, macOS `open`, Linux `xdg-open`)

---

## Phase 2: Database Driver & Lazy Connection Layer (`internal/db`)
### 2.1 Interface Standard
- [ ] 🔴 Buat base interface `Database`:
  ```go
  type Database interface {
      Connect(ctx context.Context) error
      Disconnect() error
      Ping(ctx context.Context) error
      GetTables(ctx context.Context) ([]TableInfo, error)
      GetSchema(ctx context.Context, tableName string) (*TableSchema, error)
      ExecuteQuery(ctx context.Context, query string, force bool) (*QueryResult, error)
      InsertRow(ctx context.Context, table string, data map[string]interface{}) error
      UpdateRow(ctx context.Context, table string, primaryKey map[string]interface{}, data map[string]interface{}) error
      DeleteRow(ctx context.Context, table string, primaryKey map[string]interface{}) error
  }
  ```

### 2.2 Driver Implementations
- [x] 🔴 Implementasi driver PostgreSQL (`pgx/v5`)
- [x] 🔴 Implementasi driver MySQL (`go-sql-driver/mysql`)
- [x] 🔴 Implementasi driver SQLite (`modernc.org/sqlite`)
- [x] 🔴 Terapkan **Lazy Connection**: Inisialisasi struct driver tanpa melakukan blocking `Ping()` saat CLI menyala. Ping & pembuat koneksi fisik hanya dipicu saat request HTTP pertama masuk dari Web UI.

---

## Phase 3: REST API Layer (`internal/http`)
### 3.1 Connection & Info Endpoints
- [x] 🔴 `GET /api/connection/status`: Trigger lazy connection & kembalikan status DB
- [x] 🔴 `GET /api/databases`: Mengembalikan daftar database yang terdeteksi/tersedia
- [x] 🔴 `GET /api/tables`: Mengembalikan daftar tabel, views, dan fungsi
- [x] 🔴 `GET /api/tables/{name}/schema`: Mengembalikan metadata kolom, tipe data, primary key, index, dan foreign key

### 3.2 Data CRUD Endpoints
- [x] 🔴 `GET /api/tables/{name}/data`: Mengembalikan baris data (Pagination, Sorting, Filtering, Global Search)
- [x] 🔴 `POST /api/tables/{name}`: Insert baris data baru
- [x] 🔴 `PATCH /api/tables/{name}`: Update baris data berdasarkan Primary Key
- [x] 🔴 `DELETE /api/tables/{name}`: Hapus baris data berdasarkan Primary Key

### 3.3 Raw Query Endpoint & Safety Guard Middleware
- [x] 🔴 `POST /api/query`: Memproses Raw SQL dari CodeMirror
- [x] 🔴 Safety Guard Interceptor: Match regex keyword destruktif `(?i)\b(DROP|DELETE|UPDATE|TRUNCATE|ALTER)\b`
- [x] 🔴 Kembalikan HTTP `428 Precondition Required` jika query mengandung kata kunci destruktif dan payload tidak memiliki `force: true`.

---

## Phase 4: Web Studio Frontend (`web/`)
### 4.1 Layout & State Management
- [x] 🔴 Setup Sidebar Navigasi (Search Table, List Tables, Views, Functions)
- [x] 🔴 Store Svelte untuk mengelola state koneksi aktif dan daftar tabel

### 4.2 Data Grid & CRUD Components
- [x] 🔴 Implementasi `TableGrid.svelte` (Data Grid dengan pagination, sorting header, & quick filter)
- [x] 🔴 Modal Form Insert & Edit Data
- [x] 🔴 Modal Konfirmasi Delete Data

### 4.3 Schema Explorer Components
- [x] 🔴 Tab Viewer Columns (Nama kolom, tipe data, nullable, default value)
- [x] 🟡 Tab Viewer Indexes & Foreign Keys

### 4.4 SQL Editor Component
- [x] 🔴 Integrasi CodeMirror 6 dengan SQL syntax highlighting
- [x] 🔴 Grid Result Viewer untuk hasil kueri `SELECT`
- [x] 🔴 **Danger Zone Modal**: Menangkap respon HTTP `428` dari server dan memunculkan pop-up peringatan berbahaya. Jika user setuju "Run Destructive Query", kirim ulang request dengan `{ force: true }`.

---

## Phase 5: Testing, Optimization & Distribution
### 5.1 Performance & Offline Testing
- [x] 🔴 Verifikasi startup time target (< 3 detik dari command CLI hingga browser terbuka)
- [x] 🔴 Pengujian alur 100% offline (tanpa koneksi internet)
- [ ] 🟡 Automatic Port Fallback: Jika port default `8080` terpakai, pilih random available port secara otomatis

### 5.2 Binary Distribution
- [x] 🔴 Compilation test untuk Windows (.exe), macOS, dan Linux
- [ ] 🟡 Buat NPM wrapper package (`npm install -g dbstudio`)
- [ ] 🟡 Rilis dokumen README.md & panduan instalasi Homebrew/Scoop

---

## Phase 6: Post-MVP Roadmap 🔵
- [ ] 🔵 Support Redis, MongoDB, MariaDB
- [ ] 🔵 Visual ER Diagram Viewer
- [ ] 🔵 Data Export / Import (CSV, JSON)
- [ ] 🔵 AI Query Assistant (Natural Language to SQL)
- [ ] 🔵 Dark Mode / Light Mode Toggle
