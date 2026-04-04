---
name: Migrations strategy
description: Edit the initial migration file instead of creating new ones — no prod DB yet
type: feedback
---

Edit the original migration file (00001_init_schema.sql) instead of creating a new migration file.

**Why:** The project has no production database yet, so there's no need for incremental migrations. Editing the initial schema keeps history clean.

**How to apply:** When adding new columns or tables, modify the existing migration files directly rather than creating new numbered migration files.
