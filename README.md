# DBStudio ⚡

[![npm version](https://img.shields.io/npm/v/%40kyytzy09%2Fdbstudio.svg?color=22c55e)](https://www.npmjs.com/package/@kyytzy09/dbstudio)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?logo=go)](https://go.dev)
[![SvelteKit](https://img.shields.io/badge/SvelteKit-5-FF3E00?logo=svelte)](https://kit.svelte.dev)

> **Zero-configuration Database Web Studio CLI for PostgreSQL, MySQL, and SQLite.**

**DBStudio** is a lightweight, ultra-fast local database studio powered by a **Go** backend and a **SvelteKit** single-page web interface. It automatically scans your local project workspace for existing database configurations (`.env`, `docker-compose.yml`, `connections.json`) and instantly launches a sleek local web interface for browsing, querying, and altering your database.

---

## ✨ Features

- 🔍 **Zero-Config Auto-Scanner**: Detects database connection strings from `.env` or `docker-compose.yml` automatically in your current working directory.
- 🧙 **Charm TUI Terminal Engine**: Built with [Charm](https://charm.sh) libraries (`lipgloss`, `huh`, `log`). Enjoy a clean, beautiful terminal banner, calm output, and intuitive interactive forms.
- 🕸️ **Interactive 2D ER Diagram Viewer**: Visualize table relationships with drag-and-drop table nodes, smooth SVG Bezier curve connectors, pan/zoom canvas controls, and automatic layout persistence.
- 🛠️ **Visual Schema Alteration & Table Creator**: Design new tables and add columns visually with Primary Key badges, Auto-Increment toggles, and relational Foreign Key dropdown pickers.
- ⚡ **Lazy Connection & Instant Startup**: Launches in under 3 seconds! Physical database connections are lazily established only when the Web UI triggers HTTP requests.
- 📊 **Dynamic Data Grid**: Browse records with fast client-side sorting, column filtering, pagination, and instant CSV/JSON exports.
- ✏️ **Full Data CRUD**: Add, edit, and delete database rows effortlessly with explicit PK/FK indicators powered by Shadcn UI.
- 📥 **Drag-and-Drop Bulk Data Import**: Upload CSV or JSON files via drag-and-drop zones with selectable **Insert Only** or **Upsert** modes.
- 🛡️ **Raw SQL Editor & Safety Guard (HTTP 428)**: CodeMirror 6 SQL editor with syntax highlighting, autocomplete, and a built-in safety guard interceptor that blocks destructive queries (`DROP`, `DELETE`, `TRUNCATE`, `ALTER`) until explicitly confirmed.
- 🎨 **Type UI Sleek Dark & Light Mode**: Modern, ultra-premium UI styled in Soft Charcoal Dark Mode (`#0F1115`) and Soft Light Mode (`#FCFCFC`).

---

## 💻 Quickstart

### 1. Instant Run via `npx` (No Installation Required)
Navigate to any project directory containing a `.env` or `docker-compose.yml` file and run:
```bash
npx @kyytzy09/dbstudio
```

### 2. Global Installation via `npm`
```bash
npm install -g @kyytzy09/dbstudio
dbstudio
```

---

## 🕹️ CLI Commands

```bash
# Auto-detect database credentials in cwd & launch Web Studio
dbstudio

# Launch the interactive Charm CLI wizard to manually add or select database connections
dbstudio connect

# Run diagnostic health-checks on configured databases
dbstudio doctor

# Display DBStudio build version information
dbstudio version

# Enable debug logging output
dbstudio --verbose
```

---

## 🤝 Open Source & Contributing

DBStudio is **100% Open Source**! Whether you are a database enthusiast, Go developer, or frontend hacker who loves clean tools, **contributions are warmly welcomed**! 💙

### How You Can Help:
- 🐛 **Report Issues & Bugs**: Found a bug or edge case? Open an issue on GitHub.
- 💡 **Suggest Features**: Have an idea for Redis support, AI query assistance, or new UI themes? Start a discussion!
- 🔀 **Submit Pull Requests**: Feel free to pick up open tasks, improve docs, or add test coverage.

#### Local Development Setup:
```bash
# Clone the repository
git clone https://github.com/KyyTzy09/db-studio.git
cd db-studio

# Install web dependencies
cd web && bun install && cd ..

# Build Web UI & compile Go binary locally
npm run build

# Run the compiled binary
./bin/dbstudio
```

---

## 🚀 Automated CI/CD Deployment

DBStudio utilizes GitHub Actions to automatically build cross-platform binaries (Windows, Linux, macOS ARM64/x64) and publish to the [NPM Registry](https://www.npmjs.com/package/@kyytzy09/dbstudio) whenever a release tag is pushed:

```bash
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

> **Note for Maintainers**: Ensure `NPM_TOKEN` is set in your GitHub repository secrets (`Settings -> Secrets and variables -> Actions`).

---

## 📄 License

Distributed under the **MIT License**. Created with ❤️ by [KyyTzy09](https://github.com/KyyTzy09).

See [LICENSE](LICENSE) for more information.
