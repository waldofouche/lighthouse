-- Migration: Fix MySQL Error 1170 - BLOB/TEXT column used in key specification without a key length
-- Purpose: Add size limits to TEXT columns with UNIQUE indexes to allow proper indexing in MySQL
-- Issue: #77
--
-- GORM tries to create UNIQUE indexes on TEXT columns, but MySQL requires a prefix length.
-- This migration converts TEXT columns to VARCHAR(255) where appropriate.
--
-- Database-specific instructions:
-- =================================

-- For MariaDB/MySQL:
-- -------------------
-- Convert TEXT columns to VARCHAR(255) for fields with unique indexes

-- subordinates table
ALTER TABLE subordinates MODIFY COLUMN entity_id VARCHAR(255);

-- trust_mark_types table
ALTER TABLE trust_mark_types MODIFY COLUMN trust_mark_type VARCHAR(255);

-- trust_mark_owners table
ALTER TABLE trust_mark_owners MODIFY COLUMN entity_id VARCHAR(255);

-- trust_mark_issuers table
ALTER TABLE trust_mark_issuers MODIFY COLUMN issuer VARCHAR(255);

-- trust_mark_specs table
ALTER TABLE trust_mark_specs MODIFY COLUMN trust_mark_type VARCHAR(255);

-- trust_mark_subjects table
ALTER TABLE trust_mark_subjects MODIFY COLUMN entity_id VARCHAR(255);

-- issued_trust_mark_instances table
ALTER TABLE issued_trust_mark_instances MODIFY COLUMN trust_mark_type VARCHAR(255);
ALTER TABLE issued_trust_mark_instances MODIFY COLUMN subject VARCHAR(255);

-- subordinate_additional_claims table
ALTER TABLE subordinate_additional_claims MODIFY COLUMN claim VARCHAR(255);

-- entity_configuration_additional_claims table
ALTER TABLE entity_configuration_additional_claims MODIFY COLUMN claim VARCHAR(255);

-- authority_hints table
ALTER TABLE authority_hints MODIFY COLUMN entity_id VARCHAR(255);

-- users table
ALTER TABLE users MODIFY COLUMN username VARCHAR(255);

-- For PostgreSQL:
-- ---------------
-- No migration needed - PostgreSQL doesn't have this limitation
-- The Go model changes are sufficient

-- For SQLite:
-- -----------
-- No migration needed - SQLite doesn't enforce index key length requirements
-- The Go model changes are sufficient

-- After running this migration:
-- =============================
-- 1. Restart your application (GORM AutoMigrate will use the new column sizes)
-- 2. MySQL databases will now properly create indexes on these columns
-- 3. All existing data is preserved (assuming values fit within 255 characters)
