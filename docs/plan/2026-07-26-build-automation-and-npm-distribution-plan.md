# Implementation Plan - DBStudio Build Automation & NPM Global Distribution

**Date**: 2026-07-26  
**Scope**: Build Script Automation & NPM Package Launcher Integration  

---

## 1. Overview & Understanding Summary

Dokumen ini mendokumentasikan perencanaan teknis untuk **Automated Cross-Platform Build Pipeline** dan **NPM Package Launcher** untuk **DBStudio**.

Dengan sistem ini, proses kompilasi frontend SvelteKit dan *cross-compilation* single binary Go (Windows, Linux, macOS) dapat dijalankan dalam 1 perintah otomatis (`scripts/build.ps1` atau `scripts/build.sh`). Selain itu, aplikasi siap dipublikasikan ke registri NPM sehingga pengguna dapat langsung menjalankan `npx dbstudio` atau `npm install -g dbstudio`.

---

## 2. Decision Log

| Tanggal | Keputusan | Alternatif | Alasan |
| :--- | :--- | :--- | :--- |
| 26 Jul 2026 | **Cross-Platform Shell Scripts (`build.ps1` & `build.sh`)** | Mengharuskan `make` / `nmake` terpasang di OS | PowerShell (`.ps1`) didukung native di Windows, sedangkan Shell (`.sh`) native di macOS/Linux. Keduanya tidak memerlukan dependensi build tambahan. |
| 26 Jul 2026 | **Pure Node.js Binary Launcher (`bin/cli.js`)** | Meminta user menginstal Go runtime di komputer mereka | `cli.js` menggunakan `child_process.spawn` untuk mengeksekusi single binary native sesuai OS/Arch pengguna tanpa syarat dependensi Go. |
| 26 Jul 2026 | **Root `package.json` Manifest (`bin: {"dbstudio": "./bin/cli.js"}`)** | Terpisah di subfolder terisolasi | Memasang `package.json` di root repositori memudahkan rilis publik ke NPM (`npm publish`) langsung dari root folder. |

---

## 3. Architecture & File Breakdown

### File Outputs & Structure
```
dbstudio/
├── scripts/
│   ├── build.ps1                 # Script build otomatis untuk Windows PowerShell
│   └── build.sh                  # Script build otomatis untuk Bash / Linux / macOS
├── bin/
│   ├── cli.js                    # NPM Runner script (Node.js launcher)
│   ├── dbstudio.exe              # Local Windows binary
│   ├── dbstudio-win-x64.exe      # Release Windows x64
│   ├── dbstudio-linux-x64        # Release Linux x64
│   ├── dbstudio-darwin-arm64     # Release macOS Apple Silicon (M1/M2/M3)
│   └── dbstudio-darwin-x64       # Release macOS Intel
├── package.json                  # Global NPM manifest
└── README.md                     # Panduan penggunaan & rilis
```

### 1. Script Automasi (`scripts/build.ps1` & `scripts/build.sh`)
- Otomatis kompilasi SvelteKit: `bun --cwd web run build`
- Menyiapkan folder output `bin/`
- Compiling Go Binaries (dengan `CGO_ENABLED=0`):
  - Windows x64: `GOOS=windows GOARCH=amd64 go build -o bin/dbstudio-win-x64.exe main.go` (dan `bin/dbstudio.exe`)
  - Linux x64: `GOOS=linux GOARCH=amd64 go build -o bin/dbstudio-linux-x64 main.go`
  - macOS ARM64: `GOOS=darwin GOARCH=arm64 go build -o bin/dbstudio-darwin-arm64 main.go`
  - macOS x64: `GOOS=darwin GOARCH=amd64 go build -o bin/dbstudio-darwin-x64 main.go`

### 2. NPM Node.js Binary Launcher (`bin/cli.js`)
- Membaca `process.platform` (`win32`, `linux`, `darwin`) & `process.arch` (`x64`, `arm64`).
- Menentukan path ke executable binary di folder `bin/`.
- Memanggil `spawn(binaryPath, process.argv.slice(2), { stdio: 'inherit' })`.

### 3. Root Package Manifest (`package.json`)
- `"name": "dbstudio"`
- `"bin": { "dbstudio": "./bin/cli.js" }`
- `"scripts": { "build": "powershell -ExecutionPolicy Bypass -File ./scripts/build.ps1" }`

---

## 4. Execution Phases

### Phase A: Shell & PowerShell Build Scripts
- [ ] Buat `scripts/build.ps1` (PowerShell).
- [ ] Buat `scripts/build.sh` (Bash).

### Phase B: NPM Binary Launcher & Package Manifest
- [ ] Buat `bin/cli.js` dengan pengenalan platform & spawn runner.
- [ ] Buat / Update `package.json` di root repo.

### Phase C: Execution & Multi-Platform Verification
- [ ] Eksekusi `scripts/build.ps1` untuk menghasilkan seluruh binary target.
- [ ] Uji fungsi `node bin/cli.js` dan `npx .` di terminal lokal.
- [ ] Verifikasi `dbstudio doctor` dan `dbstudio connect` melalui runner script.
