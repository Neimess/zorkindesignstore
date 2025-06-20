-- This migration script creates the initial schema for the online merch store database.
CREATE TABLE categories (
    category_id SMALLSERIAL PRIMARY KEY,
    name  VARCHAR(255) UNIQUE NOT NULL
);
CREATE TABLE attributes (
    attribute_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    unit VARCHAR(50),
    is_filterable BOOLEAN DEFAULT FALSE
);
CREATE TABLE products (
    product_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2),
    description TEXT,
    category_id SMALLINT NOT NULL REFERENCES categories(category_id) ON DELETE RESTRICT,
    image_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE product_attributes (
    product_attribute_id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(product_id) ON DELETE CASCADE,
    attribute_id BIGINT NOT NULL REFERENCES attributes(attribute_id) ON DELETE CASCADE,
    value VARCHAR(100) NOT NULL,
    UNIQUE (product_id, attribute_id)
);
CREATE TABLE category_attributes (
    category_attribute_id BIGSERIAL PRIMARY KEY,
    category_id SMALLINT NOT NULL REFERENCES categories(category_id) ON DELETE CASCADE,
    attribute_id BIGINT NOT NULL REFERENCES attributes(attribute_id) ON DELETE CASCADE,
    is_required BOOLEAN NOT NULL DEFAULT FALSE,
    priority SMALLINT NOT NULL CHECK (
        priority BETWEEN 1 AND 10
    ),
    UNIQUE (category_id, attribute_id)
);
CREATE TABLE presets (
    preset_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    total_price NUMERIC(10, 2) NOT NULL,
    image_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE preset_items (
    preset_item_id BIGSERIAL PRIMARY KEY,
    preset_id BIGINT NOT NULL REFERENCES presets(preset_id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(product_id) ON DELETE RESTRICT,
    UNIQUE (preset_id, product_id)
);