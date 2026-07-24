# Implementation Plan - Full CRUD & Bulk Drag-and-Drop Import/Export System

**Date**: 2026-07-25  
**Scope**: UI Full Data Manipulation (Insert, Edit, Delete, Bulk CSV/JSON Import/Export)  

---

## 1. Overview & Understanding Summary

Rencana ini mencakup implementasi antarmuka pengolahan data tingkat lanjut di **DBStudio**:
1. Form modal **Insert Row** & **Edit Row** dinamis berdasarkan tipe data `schema.columns`.
2. **Bulk Import Modal** dengan area visual **Drag & Drop Dropzone** (parsing CSV/JSON di client-side, live data preview, dan toggle mode **"Insert Only" vs "Upsert"**).
3. **Data Export System** (sekali klik unduh CSV & JSON dari Data Grid maupun SQL Workspace).

---

## 2. Decision Log

| Tanggal | Keputusan | Alternatif | Alasan |
| :--- | :--- | :--- | :--- |
| 25 Jul 2026 | **Client-Side CSV/JSON Parsing + Batch API** | Upload file multipart ke disk server Go | Client-side parsing 100% di browser jauh lebih cepat, tidak menambah I/O disk temporary, dan memungkinkan *live preview* sebelum data dikirim ke DB. |
| 25 Jul 2026 | **Import Mode Toggle ("Insert Only" vs "Upsert")** | Hanya mode Insert biasa / Hanya mode Replace | Memberikan fleksibilitas penuh bagi pengguna untuk memilih apakah baris dengan Primary Key duplikat harus digagalkan atau di-update. |
| 25 Jul 2026 | **Drag & Drop Zone + File Browser Fallback** | Hanya HTML `<input type="file">` standar | Area *drag and drop* memberikan UX modern dan intuitif yang mempermudah pengunggahan file data secara langsung. |

---

## 3. Architectural & File Breakdown

### Frontend Components (`web/src/lib/`)
* **`exportUtils.ts`**: Helper utility client-side untuk memformat dan men-generate download file `.csv` & `.json`.
* **`InsertRowModal.svelte`**: Modal form input 1 baris baru yang menyesuaikan field dengan kolom skema tabel.
* **`EditRowModal.svelte`**: Modal form update data baris berdasar Primary Key.
* **`ImportModal.svelte`**:
  * Visual **Drag and Drop Zone** (efek hover visual saat file diseret).
  * CSV/JSON Client-side Parser.
  * Preview 5 baris data awal.
  * Toggle Mode Import (`"insert"` / `"upsert"`).
  * Batch Request Sender dengan progress bar.

### Backend Endpoints (`internal/http/` & `internal/db/`)
* `POST /api/tables/{name}/batch`: Endpoint baru untuk menerima array data `[]map[string]interface{}` dan param `mode="insert"|"upsert"`.
* Metode `BatchInsertOrUpdate` pada Interface `Database` di Go (`postgres.go`, `mysql.go`, `sqlite.go`).

---

## 4. Execution Phases

### Phase A: Batch API Backend (`Go`)
- [ ] Tambahkan metode `BatchInsertOrUpdate(ctx, tableName, rows, mode)` di Interface `Database` & driver (`postgres`, `mysql`, `sqlite`).
- [ ] Buat Endpoint `POST /api/tables/{name}/batch` di `internal/http/handlers.go`.

### Phase B: Client-Side Export Utilities (`web/`)
- [ ] Buat `web/src/lib/utils/exportUtils.ts` (Fungsi `exportToCSV` & `exportToJSON`).
- [ ] Tambahkan tombol **Export CSV** & **Export JSON** di toolbar `TableGrid.svelte` & `SQLEditor.svelte`.

### Phase C: Form Modal Insert & Edit Row (`web/`)
- [ ] Buat `InsertRowModal.svelte` dan hubungkan dengan `POST /api/tables/{name}`.
- [ ] Buat `EditRowModal.svelte` dan hubungkan dengan `PATCH /api/tables/{name}`.

### Phase D: Drag & Drop Bulk Import Modal (`web/`)
- [ ] Buat `ImportModal.svelte` dengan area **Drag & Drop Dropzone** (`on:dragover`, `on:drop`, `on:change`).
- [ ] Integrasikan CSV Parser & JSON Parser di browser.
- [ ] Tambahkan Live Data Preview (5 baris pertama) & Toggle Mode ("Insert Only" vs "Upsert").
- [ ] Hubungkan ke API `POST /api/tables/{name}/batch`.

### Phase E: Verification & Compilation
- [ ] Test Insert Row, Edit Row, Delete Row via UI.
- [ ] Test Drag & Drop Import file CSV & JSON.
- [ ] Test Export Data CSV & JSON.
- [ ] Rebuild single binary `bin/dbstudio.exe`.
