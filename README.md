# DBStudio ⚡

> **Zero-configuration Database Web Studio CLI for PostgreSQL, MySQL, and SQLite.**

DBStudio adalah alat CLI lokal berbasis **Go + SvelteKit** yang secara otomatis mendeteksi konfigurasi database dalam proyek kamu (`.env`, `docker-compose.yml`, `connections.json`), kemudian meluncurkan Web Studio lokal yang ringan, indah (Type UI Sleek), dan responsif.

---

## 🚀 Fitur Utama

- 🔍 **Auto-detection Scanner**: Otomatis mendeteksi kredensial database dari `.env` atau `docker-compose.yml` di folder kerja kamu.
- 🧙 **CLI Interactive Wizard**: Fallback wizard interaktif di terminal jika tidak ada file konfigurasi yang ditemukan.
- 🔒 **Global Config Manager**: Menyimpan kredensial secara aman di lokasi OS pengguna (`~/.config/dbstudio/connections.json`) terikat ke path proyek.
- ⚡ **Lazy Connection & Fast Startup**: Startup cepat (< 3 detik). Koneksi fisik DB baru dipicu saat Web UI dibuka.
- 📊 **Dynamic Data Grid**: Browsing data dengan pagination, client-side filtering, sorting, & CSV/JSON export.
- ✏️ **Full CRUD Operations**: Modal tambah, edit, dan hapus baris dengan Shadcn Dialog.
- 📥 **Bulk Import CSV & JSON**: Unggah berkas CSV/JSON dengan zona *drag and drop* serta mode *Insert Only* atau *Upsert*.
- 🛡️ **Raw SQL Editor & Safety Guard (HTTP 428)**: CodeMirror 6 SQL editor yang mendeteksi kata kunci destruktif (`DROP`, `DELETE`, `UPDATE`, `TRUNCATE`) dan menampilkan pop-up peringatan keselamatan sebelum dieksekusi.
- 🎨 **Type UI Sleek Design**: Tampilan antarmuka modern dengan dukungan konsisten untuk Soft Charcoal Dark Mode (`#0F1115`) dan Soft Light Mode (`#FCFCFC`).

---

## 💻 Cara Penggunaan

### 1. Menjalankan via `npx` (Tanpa Instalasi)
```bash
npx dbstudio
```

### 2. Instalasi Global via `npm`
```bash
npm install -g dbstudio
dbstudio
```

### 3. Perintah CLI yang Tersedia

```bash
# Menjalankan auto-detection & membuka Web Studio
dbstudio

# Menjalankan CLI Interactive Wizard untuk menambah/mengedit koneksi
dbstudio connect

# Mengecek diagnosa health-check koneksi database
dbstudio doctor

# Menampilkan informasi versi rilis DBStudio
dbstudio version
```

---

## 🛠️ Panduan Build Otomatis (Cross-Platform)

Untuk mengkompilasi SvelteKit frontend dan menghasilkan single binary executable untuk seluruh OS:

### Windows (PowerShell):
```powershell
npm run build
# Atau langsung jalankan script PowerShell:
powershell -ExecutionPolicy Bypass -File .\scripts\build.ps1
```

### Linux & macOS (Bash):
```bash
npm run build:sh
# Atau langsung jalankan script Bash:
bash ./scripts/build.sh
```

Aset kompilasi akan otomatis ditempatkan di folder `bin/`:
- `bin/dbstudio.exe` (Windows local binary)
- `bin/dbstudio-win-x64.exe` (Windows 64-bit)
- `bin/dbstudio-linux-x64` (Linux 64-bit)
- `bin/dbstudio-darwin-arm64` (macOS Apple Silicon)
- `bin/dbstudio-darwin-x64` (macOS Intel)
- `bin/cli.js` (Node.js NPM Runner Launcher)

---

## 📄 Lisensi
[MIT License](LICENSE) © DBStudio Team
