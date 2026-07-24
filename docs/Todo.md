# Todo - DBStudio (Go + SvelteKit)

## Legend
- 🔴 P0 Critical path, MVP blocker
- 🟡 P1 High value, not MVP blocker
- 🟢 P2 Nice to have
- 🔵 P3 Future

---

## Phase 0: Project Initialization
### 0.1 Backend Setup (Go)
- [ ] 🔴 Inisialisasi Go module (`go mod init github.com/KyyTzy09/dbstudio`)
- [ ] 🔴 Install dependencies CLI & Config: `cobra`, `viper`
- [ ] 🔴 Install dependencies Router: `chi`
- [ ] 🔴 Install dependencies Database Drivers: `pgx/v5` (Postgres), `go-sql-driver/mysql` (MySQL), `modernc/sqlite` (SQLite)
- [ ] 🔴 Setup arsitektur folder Go (Clean Code): `cmd/` (CLI command), `internal/scanner` (deteksi config), `internal/api` (HTTP endpoints), `internal/db` (Driver Interface)
- [ ] 🔴 Konfigurasi OS global path config (misal: `~/.config/dbstudio/connections.json`)

### 0.2 Frontend Setup (SvelteKit)
- [ ] 🔴 Inisialisasi SvelteKit project (di dalam subfolder project, misal: `web/` atau `ui/`)
- [ ] 🔴 Konfigurasi `adapter-static` untuk SPA (Single Page Application) build target
- [ ] 🔴 Install dependencies UI: `tailwindcss`, `@tanstack/svelte-table` (TanStack Table), `codemirror` (CodeMirror 6)
- [ ] 🔴 Terapkan aturan *Clean Code*: Pemisahan komponen visual (Dumb Components) dengan komponen pengambil data (Smart Components)

### 0.3 Build & Embed Setup
- [ ] 🔴 Setup `go:embed` pada backend untuk membaca folder build statis dari SvelteKit (`web/build/`)
- [ ] 🔴 Buat Makefile / build script untuk otomatisasi: *Compile Frontend* -> *Embed ke Go* -> *Build Single Binary Go*

---

## Phase 1: Core CLI & Scanner Layer
### 1.1 Command Line Interface (Cobra)
- [ ] 🔴 Implementasi root command `dbstudio` (Menjalankan Scanner -> Start Chi Server -> Buka Browser)
- [ ] 🔴 Implementasi command `dbstudio connect` (Menyimpan koneksi manual ke global config)
- [ ] 🔴 Implementasi command `dbstudio doctor` (Mengecek *health* koneksi)
- [ ] 🔴 Implementasi command `dbstudio version`
- [ ] 🔴 Implementasi OS-native *browser launcher* untuk membuka `http://localhost:port` otomatis

### 1.2 Auto-Detection Scanner
- [ ] 🔴 Parsing file `.env` (Mencari `DATABASE_URL`, `DB_HOST`, `DB_USER`, dll)
- [ ] 🔴 Parsing konfigurasi dari `docker-compose.yml` atau `compose.yaml`
- [ ] 🟡 Deteksi *running* Docker container untuk fallback pencarian kredensial

---

## Phase 2: Database Driver Layer (Go)
### 2.1 Interface Definition
- [ ] 🔴 Buat base Interface standar `Database` (e.g. `Connect()`, `GetTables()`, `GetColumns()`, `ExecuteQuery()`, `Insert()`, `Update()`, `Delete()`)

### 2.2 Driver Implementation
- [ ] 🔴 Implementasi antarmuka untuk PostgreSQL
- [ ] 🔴 Implementasi antarmuka untuk MySQL
- [ ] 🔴 Implementasi antarmuka untuk SQLite
- [ ] 🔴 Implementasi logika **Lazy Connection**: Jangan *ping* DB saat CLI menyala, lakukan *ping* saat HTTP Request pertama kali masuk dari Web UI.

---

## Phase 3: REST API Layer (Chi)
### 3.1 Connection API
- [ ] 🔴 Endpoint `GET /api/connection/status` (Memicu *lazy connect* dan cek *state* database)
- [ ] 🔴 Endpoint `POST /api/connection/save` (Simpan form manual DB config ke global OS path)

### 3.2 Database Info API
- [ ] 🔴 Endpoint `GET /api/databases` (List multiple DB jika auto-detect menemukan lebih dari satu)
- [ ] 🔴 Endpoint `GET /api/tables` (List tabel, views, functions)
- [ ] 🔴 Endpoint `GET /api/tables/:name/schema` (Mendapatkan meta data: columns, data type, indexes, foreign keys)

### 3.3 Data & CRUD API
- [ ] 🔴 Endpoint `GET /api/tables/:name/data` (Pagination, Sorting, Filtering, Search)
- [ ] 🔴 Endpoint `POST /api/tables/:name` (Insert baris baru)
- [ ] 🔴 Endpoint `PATCH /api/tables/:name` (Update data)
- [ ] 🔴 Endpoint `DELETE /api/tables/:name` (Delete baris)

### 3.4 Raw Query API & Security Guard
- [ ] 🔴 Endpoint `POST /api/query` (Eksekusi Raw SQL dari CodeMirror)
- [ ] 🔴 Implementasi *Middleware/Interceptor*: Regex matching keywords destruktif (`(?i)\b(DROP|DELETE|UPDATE|TRUNCATE|ALTER)\b`)
- [ ] 🔴 Return HTTP `428 Precondition Required` jika kueri berbahaya ditemukan dan *payload* API tidak menyertakan `force=true`.

---

## Phase 4: Web Studio (Frontend - SvelteKit)
### 4.1 Layout & Initial Screens
- [ ] 🔴 Halaman Wizard "Connect to Database" (Form fallback jika deteksi CLI gagal)
- [ ] 🔴 Halaman "Select Database" (Jika multi-DB config terdeteksi)
- [ ] 🔴 Setup Layouting Utama & Sidebar (Daftar Tables, Views, Functions)

### 4.2 Data Viewer & CRUD (TanStack Table)
- [ ] 🔴 Implementasi Data Grid utama (Pagination, Sort, Filter, Search params)
- [ ] 🔴 UI form/modal untuk Insert Data
- [ ] 🔴 UI *inline edit* atau modal untuk Update Data
- [ ] 🔴 UI modal konfirmasi untuk Delete Data dari tabel

### 4.3 Schema Explorer
- [ ] 🔴 UI untuk Tab "Columns" (Menampilkan meta data tipe kolom tabel)
- [ ] 🔴 UI untuk Tab "Indexes" & "Foreign Keys"

### 4.4 SQL Editor (CodeMirror)
- [ ] 🔴 Integrasi CodeMirror 6 dengan tema dan *SQL syntax highlighting*
- [ ] 🔴 Area output *results* berupa Data Grid
- [ ] 🔴 Implementasi "Danger Zone Modal": Menangkap HTTP Status `428` dari *query request*, lalu memunculkan modal peringatan. Jika "Yes, Execute", kirim ulang request dengan `force=true`.

---

## Phase 5: Polish & Finalisasi MVP
### 5.1 Optimization & Testing
- [ ] 🔴 Verifikasi *startup time* target (< 3 detik) untuk *flow* mengetik command hingga browser terbuka
- [ ] 🔴 Testing *offline-only flow*
- [ ] 🟡 Penanganan konflik *port*: Jika port default (e.g. 8080) terpakai, otomatis *fallback* mencari random open port.

### 5.2 Packaging & Distribution
- [ ] 🔴 Optimasi ukuran Go binary build
- [ ] 🟡 Publikasi NPM wrapper package (`npm install -g dbstudio`)
- [ ] 🟡 Buat repository, README.md instruksi rilis untuk `go install`, `brew`, dan `scoop`.

---

## Phase 6: Future Capabilities (Post-MVP) 🔵
### 6.1 Roadmap Mendatang
- [ ] 🔵 Dukungan Database Tambahan (MongoDB, MariaDB, SQL Server, Redis)
- [ ] 🔵 ER Diagram Viewer
- [ ] 🔵 Export & Import Data (CSV/JSON)
- [ ] 🔵 AI Query Assistant
- [ ] 🔵 Dark Mode
- [ ] 🔵 Database Diff / Schema Compare
- [ ] 🔵 Migration History Viewer
