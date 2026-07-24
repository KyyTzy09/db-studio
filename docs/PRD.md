# PRD - DBStudio

## 1. Overview

### Project Name
**DBStudio**

### Tagline
> **One command database studio for every developer.**

### Short Description
DBStudio adalah developer tool berbentuk CLI yang secara otomatis mendeteksi database pada suatu project, membuat koneksi ke database lokal, kemudian menjalankan web studio yang dapat digunakan untuk melakukan manajemen database.

Target utamanya adalah menghilangkan proses setup yang membosankan seperti membuka pgAdmin, phpMyAdmin, Adminer, atau aplikasi database lainnya hanya untuk melihat isi database. Developer cukup mengetik `dbstudio` dan seketika UI manajemen database terbuka di browser.

---

## 2. Understanding Summary
* **Apa yang dibangun:** DBStudio, sebuah tool CLI (*zero-configuration*) yang mendeteksi database di sebuah project secara otomatis dan meluncurkan antarmuka Web UI lokal. Di-distribusikan sebagai Single Binary via Go + SvelteKit.
* **Kenapa dibangun:** Untuk memotong *context switching* dan proses panjang saat membuka aplikasi klien DB yang berat hanya untuk eksplorasi tabel.
* **Untuk siapa:** Backend/Fullstack developer, student, pengguna Docker, dan open-source developer.
* **Batasan Utama (Key Constraints):** 
  * Waktu *startup* instan (< 3 detik).
  * Distribusi *single executable binary*.
  * 100% *offline* (tanpa cloud, tanpa akun).
* **Explicit Non-goals:** Tidak untuk menggantikan kompleksitas *database management system* skala enterprise (seperti manajemen permission role DB). Fokus utamanya adalah kecepatan dan eksplorasi data.

---

## 3. Assumptions
* Konfigurasi tersimpan secara global (misal: `~/.config/dbstudio/connections.json`) agar aman dan tidak bocor ke dalam git repository project.
* Pendekatan **UI-First (Lazy Connection)**: CLI langsung menyala dan membuka browser tanpa menunggu koneksi DB *establish*. Error koneksi atau pemilihan DB multi-koneksi dilakukan melalui Web UI.
* **Keamanan Destructive Query**: SQL Editor memiliki layer pengamanan global. Segala perintah mutasi (seperti `DELETE`, `UPDATE`) akan selalu meminta persetujuan *user* melalui pop-up modal, terlepas dari keakuratan klausa *WHERE*.

---

## 4. Decision Log

| Tanggal | Keputusan | Alternatif yang Dipertimbangkan | Alasan Pemilihan |
| :--- | :--- | :--- | :--- |
| 22 Jul 2026 | **Multi-Database Handling via UI "Select Database"** | Prompt CLI / Auto-pick DB pertama. | Memberikan *user experience* terbaik dan mendukung arsitektur multi-DB atau microservices secara intuitif di UI. |
| 22 Jul 2026 | **Penyimpanan Manual Config secara Global** | Disimpan di `.dbstudio.json` di dalam folder project. | Kredensial bersifat sensitif. Menyimpannya di dalam project berisiko ter-commit ke public repo, sedangkan penyimpanan global (OS-level) jauh lebih aman dan bersih. |
| 22 Jul 2026 | **Pencegahan Destructive Query via Global Keyword Matching** | Regex parsing mendalam untuk klausa `WHERE` / Mode Read-only. | Parsing SQL secara mendalam (*AST*) untuk berbagai *dialect* sangat kompleks dan rawan lolos. Keyword matching (DELETE, UPDATE, DROP) + Global Warning adalah titik tengah yang aman dan mudah diimplementasi. |
| 22 Jul 2026 | **UI-First (Lazy Connection) Flow** | CLI-First (Eager Connection) menunggu ping koneksi sebelum buka browser. | *Startup* instan adalah core value DBStudio. UI Wizard untuk setup lebih ramah ketimbang terminal CLI. |
| 22 Jul 2026 | **Aturan Clean Code (Go & SvelteKit)** | *Free-form* / Tidak diatur secara ketat. | Penting agar mudah di-maintain dan siap untuk open source. Backend Go dipisah layernya (CLI, API, Driver), Svelte dipisah (Logic API vs Presentational). |

---

## 5. Final Design

### 5.1 Architecture & Data Flow
DBStudio dibangun sebagai **Monolithic Single Binary**:
* **Backend:** Go (menggunakan `Cobra` untuk CLI command, dan `Chi` untuk HTTP Router yang sangat ringan).
* **Frontend:** SvelteKit yang di-build menjadi *static assets* (SPA menggunakan `adapter-static`).
* **Packaging:** *Static assets* dari SvelteKit akan disisipkan (di-embed) langsung ke dalam binary Go menggunakan package `go:embed`.

**Alur Data Utama:**
1. **Init:** Developer mengetik `dbstudio` di terminal.
2. **Scan:** Modul Scanner pada Go membaca file project lokal (`.env`, `docker-compose.yml`, dsb) dan mengekstrak kredensial (seperti `DATABASE_URL`).
3. **Serve & Open:** Go menjalankan server `Chi` di localhost (misal *port* 8080), menyajikan aset UI SvelteKit yang di-embed, dan memicu OS untuk membuka browser.
4. **Connect (Lazy):** Web UI mengirim *request* API `GET /api/connection/status` ke backend Go. 
5. **Driver Layer:** Backend menggunakan *driver* (pgx/mysql/sqlite) untuk mencoba membangun koneksi. Jika berhasil, state disimpan, jika gagal, UI menampilkan Wizard konfigurasi.

**Manajemen Komponen Driver:**
*Driver Layer* pada Go menggunakan pola *Interface* untuk standarisasi (e.g. `Connect()`, `GetTables()`, `ExecuteQuery()`), sehingga mudah menambah dukungan database baru (MongoDB, dll) di masa depan.

### 5.2 Error Handling & Edge Cases
* **Tidak Ada Kredensial (Miss):** CLI tetap menjalankan Web Server dan membuka UI. User disambut halaman "Connect to Database" (Wizard).
* **Kredensial Ditemukan Tapi Koneksi Gagal:** API Go mengembalikan status `401/403/503`. UI menampilkan Wizard *pre-filled* dan menampilkan pesan *error* *native* dari database.
* **Keamanan Kueri Berbahaya:** Jika *request* `POST /api/query` mengandung perintah seperti `DELETE` atau `UPDATE`, API merespons dengan kode *warning* (`428 Precondition Required`). UI SvelteKit memunculkan *modal* "Danger Zone". Jika dikonfirmasi, *request* dikirim dengan `force=true`.

### 5.3 MVP Scope
#### Fitur CLI
* `dbstudio` (Scan, Start, Open Browser)
* `dbstudio connect` (Menyimpan koneksi secara manual)
* `dbstudio doctor` (Mengecek koneksi)
* `dbstudio version`

#### Fitur Web Studio
* **Sidebar:** Tables, Views, Functions.
* **Table Explorer:** Columns, Rows, Indexes, Foreign Keys.
* **Data Viewer:** Pagination, Sorting, Filtering, Search.
* **CRUD:** Insert, Update, Delete via UI Tabel.
* **SQL Editor:** Mendukung sintaks *highlighting* (CodeMirror) dengan fitur eksekusi kueri.

#### Database Support (MVP)
✅ PostgreSQL  
✅ MySQL  
✅ SQLite  

### 5.4 Distribution Strategy
* `go install github.com/fiky/dbstudio@latest`
* Homebrew (Mac/Linux)
* Scoop (Windows)
* NPM (`npm install -g dbstudio`) -> *Wrapper* untuk mengunduh Go Binary.

### 5.5 Tech Stack Wajib
* **Go** (Cobra, Viper, Chi, `pgx/v5`, `go-sql-driver/mysql`, `modernc/sqlite`)
* **SvelteKit** (Static SPA)
* **TailwindCSS** (Styling)
* **TanStack Table** (Data Grid)
* **CodeMirror 6** (SQL Editor)
* Wajib mengimplementasikan *Clean Code* (pemisahan layer di Backend, pemisahan *dumb/smart components* di Frontend).
