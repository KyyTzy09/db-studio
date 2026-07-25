# Implementation Plan - Schema Alteration UI (Table Creator & Column Structure Editor)

**Date**: 2026-07-26  
**Scope**: Fullstack DDL Schema Alteration (Backend DDL Engine & Frontend Visual Editor)  

---

## 1. Overview & Understanding Summary

Dokumen ini merepresentasikan perancangan teknis untuk fitur **Schema Alteration UI** pada DBStudio. Fitur ini memungkinkan pengguna membuat tabel baru secara visual melalui modal interaktif (`CreateTableModal.svelte`), serta melihat dan mengedit struktur skema tabel (menambah kolom baru, menghapus kolom) melalui tampilan tab **Structure & Columns** pada `TableSchemaView.svelte`.

Seluruh DDL diproses secara otomatis oleh driver backend Go (`internal/db/`) sesuai dengan dialek database aktif (PostgreSQL, MySQL, SQLite).

---

## 2. Decision Log

| Tanggal | Keputusan | Alternatif | Alasan |
| :--- | :--- | :--- | :--- |
| 26 Jul 2026 | **Abstraksi DDL SQL Builder di `internal/db`** | Mengirim Raw SQL dari Frontend | Menjaga keamanan backend dan mencegah SQL Injection saat membuat/mengedit skema tabel. |
| 26 Jul 2026 | **Live DDL Preview pada `CreateTableModal`** | Mengoperasikan tanpa preview SQL | Memberikan transparansi penuh bagi pengembang untuk melihat SQL `CREATE TABLE` yang akan dieksekusi sebelum commit. |
| 26 Jul 2026 | **Tab Navigation di Table View (`Data Grid` vs `Structure & Columns`)** | Halaman/Routing terpisah | User Experience yang seamless dan instan tanpa reload halaman saat berganti dari melihat baris data ke struktur tabel. |

---

## 3. Architecture & File Breakdown

### 1. Backend DDL Engine (`internal/db/` & `internal/http/`)
- **Driver Interface Extension (`internal/db/database.go`)**:
  - `CreateTable(ctx context.Context, req CreateTableRequest) error`
  - `AddColumn(ctx context.Context, table string, col ColumnSpec) error`
  - `DropColumn(ctx context.Context, table string, colName string) error`
- **Dialect DDL Builders**:
  - `internal/db/postgres.go`: Generasi SQL DDL PostgreSQL (`SERIAL`, `VARCHAR`, `TIMESTAMP WITH TIME ZONE`).
  - `internal/db/mysql.go`: Generasi SQL DDL MySQL (`AUTO_INCREMENT`, `VARCHAR(255)`, `DATETIME`).
  - `internal/db/sqlite.go`: Generasi SQL DDL SQLite (`INTEGER PRIMARY KEY AUTOINCREMENT`).
- **REST API Endpoints**:
  - `POST /api/tables` ➔ `CreateTableHandler`
  - `POST /api/tables/{name}/columns` ➔ `AddColumnHandler`
  - `DELETE /api/tables/{name}/columns/{column}` ➔ `DropColumnHandler`

### 2. Frontend Hooks & Visual Components (`web/src/lib/`)
- **`useSchemaEditor.svelte.ts`**:
  - Hook controller untuk mengelola state draft tabel baru (nama tabel, baris kolom, tipe data, PK, nullable, default).
  - Fungsi generasi DDL SQL Client-Side untuk live preview.
- **`CreateTableModal.svelte`**:
  - Pop-up modal pembuatan tabel dengan tabel pengisian kolom dinamis, switch toggle (`PK`, `Null`, `Auto Inc`), dan tab **Live SQL Preview**.
- **`AddColumnModal.svelte`**:
  - Pop-up modal penambahan kolom baru ke tabel yang sudah ada.
- **`TableSchemaView.svelte`**:
  - Tampilan tabel metadata skema (Nama Kolom, Tipe Data, Nullable, Default, Key) dengan aksi **+ Add Column** dan **Drop Column**.
- **`TableGrid.svelte`**:
  - Menambahkan Header Tab Switcher: `[ 📊 Data Grid ]` | `[ 🛠️ Structure & Columns ]`.
- **`Sidebar.svelte`**:
  - Menambahkan tombol **`+ New Table`** di bagian atas daftar tabel.

---

## 4. Execution Phases

### Phase A: Backend DDL Layer
- [ ] Tambahkan `ColumnSpec` dan `CreateTableRequest` DTO.
- [ ] Implementasikan method DDL di interface `Database` & driver (`postgres.go`, `mysql.go`, `sqlite.go`).
- [ ] Tambahkan REST API endpoints & router handler di `internal/http/`.

### Phase B: Frontend Schema Controller & Components
- [ ] Buat `useSchemaEditor.svelte.ts`.
- [ ] Buat `CreateTableModal.svelte` dan `AddColumnModal.svelte`.
- [ ] Buat `TableSchemaView.svelte`.
- [ ] Integrasikan Tab Switcher di `TableGrid.svelte` & tombol `+ New Table` di Sidebar.

### Phase C: Build & End-to-End Verification
- [ ] Uji pembuatan tabel baru via `CreateTableModal`.
- [ ] Uji penambahan & penghapusan kolom via `TableSchemaView`.
- [ ] Build SvelteKit static assets & Go single binaries (`scripts/build.ps1`).
