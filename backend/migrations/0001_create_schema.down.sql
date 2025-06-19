-- 1. Таблица, ссылающаяся на все остальные
DROP TABLE IF EXISTS product_attributes                 CASCADE;

-- 2. Таблицы с внешними ключами на categories / attributes
DROP TABLE IF EXISTS category_attribute_priority        CASCADE;
DROP TABLE IF EXISTS category_attributes                CASCADE;

-- 3. Товары (ссылается на categories)
DROP TABLE IF EXISTS products                           CASCADE;

-- 4. Справочник атрибутов
DROP TABLE IF EXISTS attributes                         CASCADE;

-- 5. Категории (самореференс parent_id)
DROP TABLE IF EXISTS categories                         CASCADE;

-- 6. Пользовательский ENUM, больше не нужен
DROP TYPE  IF EXISTS attr_data_type                     CASCADE;
