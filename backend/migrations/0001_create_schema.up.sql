-- This migration script creates the initial schema for the online merch store database.
CREATE TABLE IF NOT EXISTS categories (
    category_id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS attributes (
    attribute_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    unit VARCHAR(50),
    category_id SMALLINT NOT NULL REFERENCES categories(category_id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS products (
    product_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2),
    description TEXT,
    category_id SMALLINT NOT NULL REFERENCES categories(category_id) ON DELETE RESTRICT,
    image_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS product_attributes (
    product_attribute_id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(product_id) ON DELETE CASCADE,
    attribute_id BIGINT NOT NULL REFERENCES attributes(attribute_id) ON DELETE CASCADE,
    value VARCHAR(100) NOT NULL,
    UNIQUE (product_id, attribute_id)
);
CREATE TABLE IF NOT EXISTS presets (
    preset_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    total_price NUMERIC(10, 2) NOT NULL,
    image_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS preset_items (
    preset_item_id BIGSERIAL PRIMARY KEY,
    preset_id BIGINT NOT NULL REFERENCES presets(preset_id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(product_id) ON DELETE RESTRICT,
    UNIQUE (preset_id, product_id)
);