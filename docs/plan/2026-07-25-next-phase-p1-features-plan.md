# Implementation Plan - Next Phase P1 Expansion (CLI Commands, Port Fallback & Docker Scanner)

**Date**: 2026-07-25  
**Scope**: P1 Capability Expansion  

---

## 1. Overview & Understanding Summary

Dokumen ini mendokumentasikan perencanaan teknis untuk pengembangan **Next Phase (P1 Expansion)** pada **DBStudio**, mencakup pembuatan Subcommand CLI lengkap, Automatic Port Fallback, dan Docker Compose Auto-Detection Scanner.

---

## 2. Fitur yang Direncanakan & Alasan Kebutuhan

| Fitur | Untuk Apa (Fungsi) | Mengapa Dibutuhkan (Alasan) |
| :--- | :--- | :--- |
| **`dbstudio connect`** | Membuka CLI Interactive Wizard langsung dari terminal. | Memudahkan pengguna menambah/mengedit konfigurasi koneksi DB secara manual tanpa bergantung pada auto-detection file `.env`. |
| **`dbstudio doctor`** | Diagnosa *health check* koneksi database yang tersimpan. | Membantu developer mengecek latency ping, validasi kredensial, dan mendiagnosa secara cepat jika database lokal down/terisolasi. |
| **`dbstudio version`** | Menampilkan informasi versi build, versi Go runtime, dan OS/Arch. | Memudahkan tracking rilis, pelaporan bug, dan verifikasi versi CLI yang terpasang. |
| **Automatic Port Fallback** | Pengecekan otomatis port HTTP dari 8080 hingga 8090. | Mencegah aplikasi *crash* (`bind: address already in use`) jika port 8080 sedang terpakai oleh dev server/Docker container lain. |
| **`DockerComposeScanner`** | Parsing file `docker-compose.yml` / `compose.yaml` di direktori proyek. | Mendukung proyek berarsitektur Docker/Microservices secara *out-of-the-box* tanpa perlu menulis `.env` manual. |

---

## 3. Decision Log

| Tanggal | Keputusan | Alternatif | Alasan |
| :--- | :--- | :--- | :--- |
| 25 Jul 2026 | **Dynamic Port Listener Loop (8080-8090)** | Membiarkan error `port in use` / Prompt ganti port | Pendekatan otomatisasi murni (UI-First) memberikan UX terbaik tanpa menghentikan eksekusi CLI pengguna. |
| 25 Jul 2026 | **Cobra Subcommand Structure** | Flag tunggal di `rootCmd` (misal `--doctor`) | Mengikuti standar konvensi CLI modern (`git`, `docker`, `gh`) dengan antarmuka yang bersih dan berarsitektur modular. |
| 25 Jul 2026 | **YAML/Regex Docker Compose Parser** | Membutuhkan Docker Daemon CLI running | Parsing file konfigurasi Compose jauh lebih ringan, instan, dan tidak tergantung apakah Docker service sedang berjalan atau tidak saat scan. |

---

## 4. Architectural & File Breakdown

### Backend Changes (`Go`)
* `cmd/connect.go`:
  * Mengintegrasikan `wizard.RunCLIWizard(cwd)` dan menyimpan koneksi baru ke `internal/config`.
* `cmd/doctor.go`:
  * Membaca `connections.json` untuk proyek saat ini, mencoba `driver.Ping(ctx)`, mengukur *latency ms*, dan mencetak tabel diagnosa ringkas di terminal.
* `cmd/version.go`:
  * Mencetak string versi (e.g. `DBStudio v0.1.0-mvp (windows/amd64 go1.25.5)`).
* `internal/http/server.go`:
  * Menguji binding socket TCP `net.Listen("tcp", addr)`. Jika port terpakai, otomatis *increment* port hingga menemukan port yang terbuka.
* `internal/scanner/docker.go`:
  * Scanner baru yang membaca `docker-compose.yml` atau `compose.yaml` untuk mengekstrak service database (`postgres`, `mysql`, `sqlite`, `environment:` variables, dan `ports:` mapping).

---

## 5. Execution Phases

### Phase A: CLI Subcommands Implementation
- [ ] Buat `cmd/version.go` dengan output formatted text.
- [ ] Buat `cmd/connect.go` memicu `wizard.RunCLIWizard`.
- [ ] Buat `cmd/doctor.go` dengan format output diagnosa health check.

### Phase B: Automatic Port Fallback & Network Resilience
- [ ] Refactor `ListenAndServe()` di `internal/http/server.go` menggunakan loop `net.Listen` dinamis.
- [ ] Update pemberitahuan URL di terminal jika port ter-fallback ke port baru (e.g., `http://localhost:8081`).

### Phase C: Docker Compose Auto-Detection Scanner
- [ ] Buat `internal/scanner/docker.go` dan daftarkan ke `CompositeScanner` di `cmd/root.go`.

### Phase D: Verification & Build Test
- [ ] Test perintah `dbstudio connect`, `dbstudio doctor`, `dbstudio version`.
- [ ] Test skenario bentrok port 8080.
- [ ] Kompilasi ulang single binary `bin/dbstudio.exe`.
