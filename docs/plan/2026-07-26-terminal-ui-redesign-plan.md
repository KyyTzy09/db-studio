# Implementation Plan - Terminal CLI UI Redesign (CMD_DESIGN.md Alignment)

**Date**: 2026-07-26  
**Scope**: CLI Terminal Output Aesthetics & User Experience  

---

## 1. Overview & Understanding Summary

Dokumen ini mendokumentasikan perencanaan teknis untuk memperbarui tampilan antarmuka terminal CLI pada **DBStudio** agar 100% selaras dengan prinsip-prinsip desain dalam **[CMD_DESIGN.md](file:///D:/Coding%20Mandiri/tools/dbstudio/CMD_DESIGN.md)**.

Output terminal akan diubah menjadi **Fast, Professional, Calm, Confident, & Minimal** dengan header box border rounded, pewarnaan ANSI semantik murni (Hijau, Kuning, Merah, Biru, Abu-abu), ikon Unicode ringkas, serta alur startup yang tenang tanpa spam log.

---

## 2. Decision Log

| Tanggal | Keputusan | Alternatif | Alasan |
| :--- | :--- | :--- | :--- |
| 26 Jul 2026 | **Zero-Dependency Native ANSI Terminal Formatter (`internal/ui`)** | Menambah dependensi `lipgloss` / `color` | Menjaga kecepatan startup di bawah 300ms dan meminimalkan jumlah dependensi biner Go. |
| 26 Jul 2026 | **Rounded Box Drawing Banner (`╭──...──╮`)** | Large ASCII Art Banner / Plain text | Memberikan kesan modern dan tenang (*calm & confident*) mirip Bun / GitHub CLI / Supabase CLI. |
| 26 Jul 2026 | **Concise Single-Line Statements** | Verbose multi-line status logs | Komunikasi singkat (`✔ Connected`, `🚀 Starting Studio...`, `✨ Ready.`) memberikan persepsi alat yang sangat cepat. |

---

## 3. Architecture & File Breakdown

### 1. Terminal UI Formatter (`internal/ui/terminal.go`)
- **Banner Formatter**: `PrintBanner(version string)`
- **Color Helpers**:
  - `Green(text)` / `Success(text)`
  - `Red(text)` / `Error(text)`
  - `Yellow(text)` / `Warning(text)`
  - `Blue(text)` / `Info(text)`
  - `Gray(text)` / `Secondary(text)`
- **Standard Messages**:
  - `Scanning()` ➔ `🔍 Scanning project...`
  - `Connected(name, driver)` ➔ `✔ Connected to PostgreSQL (dbstudio)`
  - `Starting()` ➔ `🚀 Starting Studio...`
  - `Listening(url)` ➔ `✔ Listening on http://localhost:8080`
  - `OpeningBrowser()` ➔ `🌐 Opening browser...`
  - `Ready()` ➔ `✨ Ready.`

### 2. Command Updates (`cmd/`)
- **`cmd/root.go`**: Menerapkan alur startup sekuensial lengkap dengan banner dan status semantik.
- **`cmd/connect.go`**: Memperbarui prompt wizard agar rapi.
- **`cmd/doctor.go`**: Menata kartu diagnosa health-check dengan garis pembatas minimalis dan warna latency.
- **`cmd/version.go`**: Mencetak versi dalam format clean banner.
- **`internal/http/server.go`**: Menyerap formatter pesan port listener.

---

## 4. Execution Phases

### Phase A: Terminal Formatter Module
- [ ] Buat `internal/ui/terminal.go` dengan fungsi-fungsi format ANSI & Box Banner.

### Phase B: Command Refactoring
- [ ] Refactor `cmd/version.go`.
- [ ] Refactor `cmd/doctor.go`.
- [ ] Refactor `cmd/root.go` & `internal/http/server.go`.

### Phase C: Build & CLI Verification
- [ ] Uji `go run main.go version`.
- [ ] Uji `go run main.go doctor`.
- [ ] Uji `go run main.go`.
- [ ] Jalankan `scripts/build.ps1` untuk kompilasi ulang seluruh binary rilis.
