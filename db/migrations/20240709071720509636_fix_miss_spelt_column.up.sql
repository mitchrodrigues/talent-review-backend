-- Up Migration 20240709071720509636 fix_miss_spelt_column

ALTER TABLE feedback_details
  RENAME COLUMN strenghts TO strengths;

