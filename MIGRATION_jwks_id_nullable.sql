-- Migration: Make subordinates.jwks_id nullable
-- Purpose: Allow creating pending subordinates without JWKS
-- 
-- This fixes the foreign key constraint error when creating subordinates
-- with status=pending that don't have JWKS yet.
--
-- Database-specific instructions:
-- =================================

-- For MariaDB/MySQL:
-- -------------------
ALTER TABLE subordinates MODIFY COLUMN jwks_id BIGINT UNSIGNED NULL;

-- Optional: Clean up any rows with jwks_id=0 (these were invalid anyway)
UPDATE subordinates SET jwks_id = NULL WHERE jwks_id = 0;

-- For PostgreSQL:
-- ---------------
-- ALTER TABLE subordinates ALTER COLUMN jwks_id DROP NOT NULL;

-- For SQLite:
-- -----------
-- No migration needed - SQLite doesn't enforce FK constraints by default
-- The Go model change is sufficient

-- After running this migration:
-- =============================
-- 1. Restart your application (GORM AutoMigrate will update the schema)
-- 2. You can now create pending subordinates without JWKS
-- 3. Active subordinates still require JWKS (enforced at application level)
