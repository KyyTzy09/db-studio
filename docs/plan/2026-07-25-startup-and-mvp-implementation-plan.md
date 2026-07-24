# Implementation Plan - DBStudio Startup Flow & Modular MVP

**Date**: 2026-07-25  
**Architecture Mode**: Option 1 - Modular Pipeline Architecture  

---

## 1. Overview & Understanding Summary

DBStudio adalah tool CLI (*zero-configuration*) berbasis Go dan SvelteKit yang secara otomatis mendeteksi database dalam sebuah proyek, kemudian meluncurkan Web Studio lokal yang aman dan ringan.

### Alur Startup Utama
```
dbstudio command
       │
       ▼
[Cek Global Saved Config] ──(Ditemukan)──► Run HTTP Server & Open Web UI
       │
   (Tidak Ada)
       ▼
[Auto-Detection Scanner]
  ├── Ditemukan 1 Config ───────────────► Run HTTP Server & Open Web UI
  ├── Ditemukan >1 Configs ──────────────► User Choice / Prompt di CLI/UI
  └── Ditemukan 0 Config (Gagal) ───────► Fallback: CLI Interactive Wizard
                                                │
                                                ▼
                                    Simpan ke Global Config & Run Web UI
```

---

## 2. Decision Log

| Tanggal | Keputusan | Alternatif | Alasan |
| :--- | :--- | :--- | :--- |
| 25 Jul 2026 | **Modular Pipeline Architecture (Option 1)** | Monolithic Handler di `cmd/root.go` | Memisahkan tanggung jawab antara `config`, `scanner`, `wizard`, `db driver`, dan `http handler` untuk memudahan pengujian (*unit testing*) dan maintainability. |
| 25 Jul 2026 | **CLI Interactive Wizard Fallback** | Langsung buka UI Wizard di Browser | Jika auto-detect 0 config, prompt CLI interaktif di terminal langsung menyelesaikan input kredensial sebelum browser dibuka, memastikan koneksi siap saat UI menyala. |
| 25 Jul 2026 | **Global OS Path for Credentials** | Menyimpan di folder repo proyek `.env` | Menghindari risiko kredensial ter-commit ke repositori publik Git. |
| 25 Jul 2026 | **Destructive Query Protection (HTTP 428)** | AST SQL Parsing / Mode Read-only murni | Keyword matching + HTTP 428 Precondition Required adalah titik tengah yang aman, responsif, dan mudah dipahami pengguna. |

---

## 3. Modular Architecture Breakdown

### Backend Layer (`Go`)
* `cmd/`: Entry point CLI Cobra (`root.go`, `connect.go`, `doctor.go`, `version.go`).
* `internal/config/`: Membaca/menulis `~/.config/dbstudio/connections.json`. Menyimpan koneksi dengan *key* berupa absolut path proyek.
* `internal/scanner/`: Logika ekstraksi variabel database dari `.env` dan `docker-compose.yml`.
* `internal/wizard/`: Prompt interaktif CLI jika *scanner* gagal menemukan koneksi.
* `internal/db/`: Interface standar `Database` (Lazy Connection) dan driver (`postgres`, `mysql`, `sqlite`).
* `internal/http/`: HTTP Router Chi, REST API endpoints, middleware safety guard (`HTTP 428`), dan handler `go:embed` untuk SvelteKit static build.

### Frontend Layer (`SvelteKit`)
* `web/`: SvelteKit dengan `adapter-static` (SPA).
* `web/src/lib/components/`:
  * `TableGrid.svelte`: Data Grid TanStack Table.
  * `SQLEditor.svelte`: CodeMirror 6 dengan Danger Zone modal.
  * `SchemaExplorer.svelte`: Viewer kolom, tipe data, index, dan foreign key.
  * `ConnectWizard.svelte`: Fallback UI untuk manajemen koneksi.

---

## 4. Execution Phases

### Phase 0: Project Setup & Monorepo Structure
- Setup Go modules & dependencies (`cobra`, `chi`, `pgx/v5`, `go-sql-driver/mysql`, `modernc/sqlite`).
- Setup SvelteKit di `web/` dengan `adapter-static`, Tailwind CSS v4, TanStack Table, CodeMirror 6.
- Setup `go:embed` handler dan Makefile build pipeline.

### Phase 1: Config, Scanner & CLI Wizard
- Implementasi `internal/config` (OS Global Path Manager).
- Implementasi `internal/scanner` (`SavedConfigScanner`, `EnvScanner`, `DockerComposeScanner`).
- Implementasi `internal/wizard` (CLI interactive prompt fallback).
- Implementasi `cmd/` Cobra commands & browser launcher.

### Phase 2: Database Driver & Lazy Connection
- Interface `Database` di `internal/db/driver.go`.
- Implementasi driver PostgreSQL, MySQL, SQLite dengan pola *Lazy Connection*.

### Phase 3: REST API & Safety Guard
- Endpoint status, database list, tables, schema, data CRUD.
- Middleware safety guard HTTP `428` untuk query destruktif (`DELETE`, `UPDATE`, `DROP`, dsb.).

### Phase 4: SvelteKit Web Studio Development
- Implementasi Layouting, Data Grid, Schema Explorer, SQL Editor + Danger Zone Modal.

### Phase 5: Verification, Build & Distribution
- Verifikasi startup time (< 3 detik).
- Multi-platform binary compilation.
