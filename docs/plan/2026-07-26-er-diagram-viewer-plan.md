# Implementation Plan - Visual ER Diagram Viewer

## 1. Overview & Understanding Summary
* **Feature Name:** Visual ER Diagram Viewer (`[ 🕸️ ER Diagram ]`)
* **Purpose:** Provides an interactive 2D relational schema diagram showing tables, columns, primary keys, and foreign key relationships with drag-and-drop node positioning, zoom/pan controls, and click-to-navigate integration.
* **Target Audience:** Developers needing instant visual clarity of complex relational database schemas without writing manual `JOIN` queries or inspecting schemas individually.
* **Key Constraints:**
  - Fast single-request schema extraction (`GET /api/schema/graph`).
  - Interactive Svelte 5 SVG engine (no heavy external canvas canvas dependencies).
  - Node position persistence in `localStorage`.

---

## 2. Architecture & Data Flow

```
[ Web UI: ERDiagramView.svelte ]
           │
           │ GET /api/schema/graph
           ▼
[ Chi REST API: table_handler.go ]
           │
           │ Driver.GetSchemaGraph(ctx)
           ▼
[ Driver Layer: postgres.go / mysql.go / sqlite.go ]
           │
           │ Queries information_schema / PRAGMA foreign_key_list
           ▼
[ Database (Postgres / MySQL / SQLite) ]
```

### JSON Payload Contract (`GET /api/schema/graph`):
```json
{
  "nodes": [
    {
      "name": "users",
      "columns": [
        { "name": "id", "data_type": "INTEGER", "is_primary_key": true, "is_nullable": false },
        { "name": "email", "data_type": "VARCHAR(255)", "is_primary_key": false, "is_nullable": false }
      ]
    },
    {
      "name": "orders",
      "columns": [
        { "name": "id", "data_type": "INTEGER", "is_primary_key": true, "is_nullable": false },
        { "name": "user_id", "data_type": "INTEGER", "is_primary_key": false, "is_nullable": false }
      ]
    }
  ],
  "edges": [
    {
      "id": "orders_user_id-users_id",
      "source_table": "orders",
      "source_column": "user_id",
      "target_table": "users",
      "target_column": "id"
    }
  ]
}
```

---

## 3. Decision Log

| Aspect | Decision | Rationale |
| :--- | :--- | :--- |
| **Rendering Engine** | **Custom Svelte SVG Node Graph** | Native Svelte 5 reactive rendering, lightweight, fully styled with Tailwind dark/light theme tokens, enables instant click navigation to Data Grid. |
| **API Architecture** | **`GET /api/schema/graph` Single Endpoint** | Extracts all nodes and foreign key constraints across the database in one single query pass, eliminating N+1 HTTP request overhead. |
| **Node Placement** | **Smart Auto-Grid + Drag Persistence (`localStorage`)** | Places nodes neatly on initial load and saves user's custom card coordinates in `localStorage` under `dbstudio_node_pos_<dbname>`. |

---

## 4. Execution Phases

### Phase 1: Backend DDL & Relationship Extraction (`internal/db`)
1. Extend `Database` interface in `internal/db/driver.go`:
   - Add `ForeignKeyRelation` struct (`SourceTable`, `SourceColumn`, `TargetTable`, `TargetColumn`).
   - Add `SchemaGraph` struct (`Nodes []TableSchema`, `Edges []ForeignKeyRelation`).
   - Add `GetSchemaGraph(ctx context.Context) (*SchemaGraph, error)` method.
2. Implement `GetSchemaGraph` in Postgres driver (`internal/db/postgres.go`):
   - Query `information_schema.key_column_usage` & `information_schema.constraint_column_usage`.
3. Implement `GetSchemaGraph` in MySQL driver (`internal/db/mysql.go`):
   - Query `information_schema.KEY_COLUMN_USAGE`.
4. Implement `GetSchemaGraph` in SQLite driver (`internal/db/sqlite.go`):
   - Query `PRAGMA foreign_key_list(table_name)`.

### Phase 2: REST API Endpoint (`internal/http`)
1. Add `HandleGetSchemaGraph` in `internal/http/handlers/table_handler.go`.
2. Register route `r.Get("/api/schema/graph", h.HandleGetSchemaGraph)` in `internal/http/routes/routes.go`.

### Phase 3: Frontend ER Diagram Components (`web/src/lib/components/diagram`)
1. Create state hook `web/src/lib/hooks/useERDiagram.svelte.ts`:
   - Manage node positions `{ [tableName]: { x, y } }`.
   - Manage canvas transform `{ zoom, panX, panY }`.
   - Handle drag events for table cards & canvas background.
   - Persistence helper for `localStorage`.
2. Create component `web/src/lib/components/diagram/TableNode.svelte`:
   - Render styled card with table title, column list, PK key icons, and FK badges.
3. Create component `web/src/lib/components/diagram/RelationshipLine.svelte`:
   - Calculate SVG Bezier curve path connecting FK column port to target PK column port.
4. Create view component `web/src/lib/components/diagram/ERDiagramView.svelte`:
   - Zoom (+ / - / Reset) & Pan control bar.
   - SVG Overlay layer for connection lines.
   - HTML Node layer for draggable table cards.
   - Click-to-navigate action (Jump to Data Grid or Schema Info).

### Phase 4: Integration & Top Navigation
1. Add `[ 🕸️ ER Diagram ]` tab button in `web/src/routes/+page.svelte` header.
2. Render `<ERDiagramView />` when `$activeTab === 'diagram'`.
3. Update `web/src/lib/stores/dbStore.ts` active tab type to include `'diagram'`.

---

## 5. Verification & Acceptance Criteria
- [ ] `GET /api/schema/graph` returns complete list of tables and foreign key relations for Postgres, MySQL, and SQLite.
- [ ] ER Diagram tab renders all tables cleanly with SVG Bezier curves connecting FK to PK columns.
- [ ] Nodes are draggable and their positions persist across page refreshes via `localStorage`.
- [ ] Zoom in/out and panning work smoothly without lag.
- [ ] Clicking a table node allows jumping directly to the Data Grid for that table.
