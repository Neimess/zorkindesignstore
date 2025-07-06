ALTER TABLE categories
ADD COLUMN parent_id SMALLINT REFERENCES categories(category_id) ON DELETE CASCADE;
