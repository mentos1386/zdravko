# `CONTRIBUTING.md`

Hello, and thank you for choosing to contribute to this project.

## Changing database schema.

### Create new migration file.

You can read about the syntax in the [rubenv/sql-migrate](https://github.com/rubenv/sql-migrate?tab=readme-ov-file#writing-migrations) repo. The files are located in `database/sqlite/migrations/` folder.

To create a new migration file, run `just migration-new creating-something-new`.
