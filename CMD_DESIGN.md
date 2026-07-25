# Terminal Design Guidelines

## Philosophy

The terminal is the first interaction users have with DBStudio.

It should feel:

- Fast
- Professional
- Calm
- Confident
- Minimal

Avoid flashy hacker aesthetics.

The goal is to communicate progress clearly while keeping output clean.

---

# Inspiration

- Bun
- Docker
- Git
- Wrangler
- Fly.io
- GitHub CLI
- Railway CLI
- Supabase CLI

---

# Visual Principles

## Minimal Colors

Use only semantic colors.

Green
- Success
- Connected
- Completed

Yellow
- Warning
- Missing optional configuration

Red
- Errors
- Connection failures

Blue
- Information
- Startup messages

Gray
- Secondary text

Avoid rainbow output.

---

## Icons

Prefer simple Unicode icons.

✔ Success

✖ Error

⚠ Warning

ℹ Information

🔍 Scanning

🚀 Starting

🌐 Server

📦 Package

🗄 Database

🧩 Plugin

Do not overuse emojis.

Maximum one icon per line.

---

## Typography

Keep spacing consistent.

Example:

✔ Connected PostgreSQL

Instead of:

[ SUCCESS ] Connected PostgreSQL Successfully!!!

---

## Borders

Use box drawing characters sparingly.

Good:

╭──────────────────────────────╮
│        DBStudio v0.1         │
╰──────────────────────────────╯

Avoid ASCII art logos.

Avoid giant banners.

---

## Progress

Show meaningful progress.

Good:

🔍 Scanning project...

✔ docker-compose.yml detected

✔ PostgreSQL detected

✔ DATABASE_URL loaded

🚀 Starting Studio...

✔ Listening on http://localhost:3333

🌐 Opening browser...

Bad:

Loading...

Loading...

Loading...

Done.

---

## Errors

Errors should explain:

- What happened
- Why
- How to fix

Good:

✖ DATABASE_URL not found

Create a .env file or specify:

dbstudio --database-url <url>

Bad:

Fatal Error

Exit code: 1

---

## Startup Flow

Recommended sequence:

DBStudio Banner

↓

Version

↓

Project Scan

↓

Environment Detection

↓

Database Connection

↓

HTTP Server

↓

Browser Launch

↓

Ready

---

## Example Output

╭────────────────────────────────────╮
│           DBStudio v0.1            │
╰────────────────────────────────────╯

🔍 Scanning project...

✔ .env found

✔ docker-compose.yml found

✔ PostgreSQL detected

✔ Connected

🚀 Starting Studio...

✔ Listening on http://localhost:3333

🌐 Opening browser...

✨ Ready.

---

## Logging

Do not spam logs.

Users should not see internal implementation details.

Hide:

Initializing HTTP router...

Loading middleware...

Registering handlers...

Injecting dependencies...

Unless running with:

--verbose

---

## Verbose Mode

Verbose mode may include:

Database driver

Execution time

HTTP routes

Environment resolution

Detected project type

Driver selection

---

## Performance

Startup should appear responsive.

If an operation exceeds 300ms,
show a spinner.

If it exceeds 1 second,
display progress information.

---

## Recommended Libraries

CLI
- Cobra

Styling
- Lip Gloss

Spinner
- Bubbles

Interactive Forms
- Huh

Terminal Detection
- termenv

Color Support
- Lip Gloss Adaptive Colors

Browser Opening
- pkg/browser

---

## Tone

Messages should be concise.

Good

✔ Connected

✔ Opening browser

✔ Ready

Bad

Successfully connected to PostgreSQL database server.

Launching browser now...

Everything has completed successfully.

---

## Animation

Use subtle animations only.

Spinner

Progress Bar

No blinking text.

No rainbow effects.

No matrix animations.

---

## Accessibility

Output should remain readable on:

Dark terminals

Light terminals

Windows Terminal

iTerm2

Ghostty

Warp

VSCode Terminal

---

## Design Goal

Users should immediately feel:

"This tool is fast."

Not:

"This terminal looks cool."

Speed and clarity always come before visual effects.