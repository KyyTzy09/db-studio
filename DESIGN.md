# 🎨 DBStudio Design System (Type UI: Sleek)

`DBStudio` is a lightweight, zero-config local web studio for inspecting and managing databases. The user interface follows the **Type UI "Sleek"** design paradigm: refined, high-density, crisp contrast, soft dark mode backgrounds (avoiding harsh `#000000` pitch black), and zero hardcoded colors.

---

## 🔮 Design Principles (Sleek Aesthetic)

1. **High Visual Density & Precision**:
   - Compact table cells, crisp `0.5rem` (`var(--radius)`) rounded corners, subtle border strokes (`1px`).
   - Monospace typography for database object names, SQL queries, data types, and primary key badges.

2. **Ergonomic Soft Dark Mode**:
   - Avoid pure `#000000` pitch black background which causes severe eye strain in long debugging sessions.
   - Deep slate-tinted background (`#0F1115`) paired with elevated surface panels (`#171A21`) and card containers (`#1D2129`).

3. **Semantic Status Hierarchy**:
   - **Primary / Accent**: `#2563EB` (Light) / `#4F8CFF` (Dark) — Main actions, active navigation, tab highlights.
   - **Success**: `#16A34A` (Light) / `#22C55E` (Dark) — Database connected status, successful query executions, positive status badges.
   - **Warning**: `#D97706` (Light) / `#F59E0B` (Dark) — Destructive query alerts, nullable columns, unsaved changes.
   - **Danger / Destructive**: `#DC2626` (Light) / `#EF4444` (Dark) — Record deletion, drop table actions, connection failure errors.

4. **100% Tokenized CSS Variables**:
   - All components consume Tailwind CSS theme tokens (`bg-background`, `bg-card`, `bg-muted`, `border-border`, `text-foreground`, `text-muted-foreground`) defined in `app.css`.

---

## 🎨 Color Palette & Design Tokens

### Light Theme
| Token | HEX Value | Description |
| :--- | :--- | :--- |
| `--background` | `#FCFCFC` | Global page background |
| `--card` / Surface | `#FFFFFF` | Elevating card panels & modals |
| `--muted` / Surface Alt | `#F6F6F7` | Sidebar active background, table header hover |
| `--border` | `#E5E7EB` | Subtle dividing strokes |
| `--foreground` | `#111827` | High contrast primary text |
| `--muted-foreground` | `#6B7280` | Secondary labels, descriptions, metadata |
| `--primary` / Accent | `#2563EB` | Interactive buttons, active tabs, focused inputs |
| `--success` | `#16A34A` | Connection health, positive notifications |
| `--warning` | `#D97706` | Precondition warning badges |
| `--destructive` | `#DC2626` | Delete modals, destructive actions |

### Dark Theme (Soft Dark - No #000000)
| Token | HEX Value | Description |
| :--- | :--- | :--- |
| `--background` | `#0F1115` | Soft charcoal global canvas background |
| `--popover` / Surface | `#171A21` | Sidebar, dropdown popovers, headers |
| `--card` / Surface Alt | `#1D2129` | Modals, table containers, card panels |
| `--border` | `#2A2F3A` | Dark border strokes |
| `--foreground` | `#F5F5F5` | Off-white primary text |
| `--muted-foreground` | `#A1A1AA` | Slate secondary text |
| `--primary` / Accent | `#4F8CFF` | Electric blue highlights, active indicators |
| `--success` | `#22C55E` | Green connection badge |
| `--warning` | `#F59E0B` | Amber warning dialogs |
| `--destructive` | `#EF4444` | Red delete button & error banners |

---

## 🧩 Component Architecture (Shadcn-Svelte Integration)

- **Shadcn Primitives**: `Button`, `Input`, `Label`, `Select`, `Card`, `Badge`, `Popover`, `AlertDialog`, `Separator`, `Skeleton`.
- **Layout System**:
  - `Sidebar`: Database table list, search filter, doctor ping tool, dark/light mode toggle.
  - `Header`: Active table breadcrumb, Data Grid / Schema Info / SQL Workspace tab switcher.
  - `TableGrid`: Data grid with pagination, search filter, insert row modal trigger, bulk import modal trigger, export CSV/JSON utility.
  - `SchemaExplorer`: Detailed column metadata breakdown with PK badges & data type tags.
  - `SQLEditor`: CodeMirror 6 query editor with real-time execution timer and HTTP 428 Danger Zone Modal.
  - `Modals`: `InsertRowModal`, `EditRowModal`, `ImportModal`.

---

## ✒️ Typography & Spacing Rules

- **Font Family**: `Geist Variable`, sans-serif for UI; `Fira Code` / `JetBrains Mono` / `ui-monospace` for SQL, table names, and column data.
- **Base Radius**: `0.5rem` (`var(--radius)`).
- **Transitions**: `transition-all duration-150 ease-in-out` on interactive controls.
