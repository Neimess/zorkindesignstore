CREATE TABLE categories (
    category_id smallserial PRIMARY KEY,
    name varchar(255) NOT NULL
);

CREATE TYPE attr_data_type AS ENUM ('string','int','float','bool','enum');

CREATE TABLE attributes (
    id            bigserial PRIMARY KEY,
    name          varchar(255) NOT NULL,
    slug          varchar(100) UNIQUE NOT NULL,
    data_type     attr_data_type NOT NULL,
    unit          varchar(50), -- Единицы измерения 
    is_filterable boolean DEFAULT false
);

CREATE TABLE products (
    id          bigserial PRIMARY KEY,
    name        varchar(255) NOT NULL,
    price       numeric(10,2),
    description text,
    category_id bigint REFERENCES categories(id) ON DELETE RESTRICT,
    image_url   text,
    created_at  timestamptz DEFAULT now()
);

CREATE TABLE category_attributes (
    id           bigserial PRIMARY KEY,
    category_id  bigint REFERENCES categories(id) ON DELETE CASCADE,
    attribute_id bigint REFERENCES attributes(id) ON DELETE CASCADE,
    is_required  boolean DEFAULT false,
    UNIQUE (category_id, attribute_id)
);

CREATE TABLE category_attribute_priority (
    id           bigserial PRIMARY KEY,
    category_id  bigint REFERENCES categories(id) ON DELETE CASCADE,
    attribute_id bigint REFERENCES attributes(id) ON DELETE CASCADE,
    priority     int NOT NULL CHECK (priority BETWEEN 1 AND 10),
    UNIQUE (category_id, attribute_id)
);

CREATE TABLE product_attributes (
    id            bigserial PRIMARY KEY,
    product_id    bigint REFERENCES products(id) ON DELETE CASCADE,
    attribute_id  bigint REFERENCES attributes(id) ON DELETE CASCADE,
    value_string  text,
    value_int     bigint,
    value_float   numeric,
    value_bool    boolean,
    value_enum    text,
    UNIQUE (product_id, attribute_id)
);