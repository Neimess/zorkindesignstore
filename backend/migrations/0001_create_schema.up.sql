CREATE TABLE categories (
    category_id smallserial PRIMARY KEY,
    name varchar(255) NOT NULL
);
CREATE TYPE attr_data_type AS ENUM ('string', 'int', 'float', 'bool', 'enum');
CREATE TABLE attributes (
    attribute_id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL,
    slug varchar(100) UNIQUE NOT NULL,
    data_type attr_data_type NOT NULL,
    unit varchar(50),
    is_filterable boolean DEFAULT false
);
CREATE TABLE products (
    product_id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL,
    price numeric(10, 2),
    description text,
    category_id bigint REFERENCES categories(category_id) ON DELETE RESTRICT,
    image_url text,
    created_at timestamptz DEFAULT now()
);
CREATE TABLE category_attributes (
    category_attribute_priority_id bigserial PRIMARY KEY,
    category_id bigint REFERENCES categories(category_id) ON DELETE CASCADE,
    attribute_id bigint REFERENCES attributes(attribute_id) ON DELETE CASCADE,
    is_required boolean DEFAULT false,
    UNIQUE (category_id, attribute_id)
);
CREATE TABLE category_attribute_priority (
    category_attribute_priority_id bigserial PRIMARY KEY,
    category_id bigint REFERENCES categories(category_id) ON DELETE CASCADE,
    attribute_id bigint REFERENCES attributes(attribute_id) ON DELETE CASCADE,
    priority int NOT NULL CHECK (
        priority BETWEEN 1 AND 10
    ),
    UNIQUE (category_id, attribute_id)
);
CREATE TABLE product_attributes (
    product_attribute_id bigserial PRIMARY KEY,
    product_id bigint REFERENCES products(product_id) ON DELETE CASCADE,
    attribute_id bigint REFERENCES attributes(attribute_id) ON DELETE CASCADE,
    value_string text,
    value_int bigint,
    value_float numeric,
    value_bool boolean,
    value_enum text,
    UNIQUE (product_id, attribute_id)
);
CREATE TABLE presets (
    preset_id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL,
    description text,
    total_price numeric(10, 2) NOT NULL,
    image_url text,
    created_at timestamptz DEFAULT now()
);
CREATE TABLE preset_items (
    preset_item_id bigserial PRIMARY KEY,
    preset_id bigint REFERENCES presets(preset_id) ON DELETE CASCADE,
    product_id bigint REFERENCES products(product_id) ON DELETE RESTRICT,
    quantity numeric(10, 2) NOT NULL,
    -- Учитываем, что может быть 1.5 м² плитки
    unit varchar(20),
    -- м², шт., рулон, л, кг
    note text,
    -- Дополнительно: "для пола", "для стены"
    UNIQUE (preset_id, product_id)
);