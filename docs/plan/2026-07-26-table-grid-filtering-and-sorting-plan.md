# Implementation Plan - TableGrid Multi-Condition Filtering & Column Sorting

**Date**: 2026-07-26  
**Scope**: Frontend Data Management Enhancement  

---

## 1. Overview & Understanding Summary

Dokumen ini mendokumentasikan perencanaan teknis untuk penambahan fitur **Single & Multi-Condition Filtering** serta **Column Header Sorting** pada komponen `TableGrid.svelte` di **DBStudio**.

Fitur ini memberikan kemampuan pada pengguna untuk:
1. Melakukan penyaringan cepat (*Single Filter*) maupun penyaringan tingkat lanjut dengan banyak aturan (*Multi-Condition Filter* dengan logika `AND`).
2. Mengurutkan data (*Sorting*) berdasarkan kolom manapun (Ascending, Descending, Reset) dengan mengklik header tabel `<th>`.

---

## 2. Decision Log

| Tanggal | Keputusan | Alternatif | Alasan |
| :--- | :--- | :--- | :--- |
| 26 Jul 2026 | **Client-Side Derived Filtering & Sorting** | Server-side SQL WHERE & ORDER BY query | Mengingat data baris tabel sudah berada di memori browser, penyaringan & pengurutan client-side memberikan respon instan (*zero latency*) tanpa pemicu reload server. |
| 26 Jul 2026 | **3-State Cycle Sorting (`ASC` ➔ `DESC` ➔ `Reset`)** | Toggle 2-State (`ASC` ↔ `DESC` tanpa reset) | Memungkinkan pengguna mengembalikan urutan ke kondisi asli tabel dengan 3x klik. |
| 26 Jul 2026 | **Collapsible Filter Panel di Toolbar** | Floating Pop-over per Kolom Header | Filter Panel terpusat di toolbar lebih bersih, mudah digunakan untuk mengombinasikan banyak aturan filter, dan tidak mengganggu visual header tabel. |

---

## 3. Architecture & File Breakdown

### 1. Hook Extension (`web/src/lib/hooks/useTableData.svelte.ts`)
- **Filter Rule Types**:
  - `FilterOperator`: `'contains' | 'equals' | 'gt' | 'lt' | 'starts' | 'is_null'`
  - `FilterRule`: `{ id, column, operator, value }`
  - `SortState`: `{ column: string | null, direction: 'asc' | 'desc' | null }`
- **Actions**:
  - `addFilterRule(rule)`
  - `removeFilterRule(id)`
  - `clearAllFilters()`
  - `toggleSort(column)`

### 2. Table Grid Component (`web/src/lib/components/table/TableGrid.svelte`)
- **Toolbar UI**:
  - Tombol **Filter** dengan indikator badge (jumlah filter aktif).
  - Panel Collapsible Filter:
    - Dropdown Nama Kolom
    - Dropdown Operator (`contains`, `equals`, `greater than`, `less than`, `starts with`, `is null`)
    - Input Nilai Filter
    - Tombol `+ Add Rule` & `Clear All`
  - Active Filter Chips (Label filter aktif dengan tombol hapus `x`).
- **Table Header `<th>`**:
  - Interactive clickable header dengan icon indikator urutan (`ArrowUp`, `ArrowDown`, `ArrowUpDown`).

---

## 4. Execution Phases

### Phase A: Hook Logic (`useTableData.svelte.ts`)
- [ ] Tambahkan tipe `FilterRule` & `SortState`.
- [ ] Tambahkan handler `addFilterRule`, `removeFilterRule`, `clearAllFilters`, `toggleSort`.
- [ ] Terapkan `$derived` pipeline untuk penyaringan & pengurutan data.

### Phase B: TableGrid Component UI (`TableGrid.svelte`)
- [ ] Tambahkan Filter Bar Toggle & Collapsible Panel.
- [ ] Render Active Filter Chips.
- [ ] Tambahkan handler click sort & icon panah di `<th>`.

### Phase C: Build & Verification
- [ ] Jalankan `bun run build` di `web/`.
- [ ] Jalankan `go build -o bin/dbstudio.exe main.go`.
- [ ] Uji fungsionalitas filter & sort di browser.
