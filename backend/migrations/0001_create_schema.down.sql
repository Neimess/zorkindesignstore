-- 1. Таблица, ссылающаяся на все остальные
DROP TABLE IF EXISTS product_attributes CASCADE;
-- 2. Таблицы с внешними ключами на categories / attributes
DROP TABLE IF EXISTS category_attributes CASCADE;
-- 3. Товары (ссылается на categories)
DROP TABLE IF EXISTS products CASCADE;
-- 4. Справочник атрибутов
DROP TABLE IF EXISTS attributes CASCADE;
-- 5. Категории (самореференс parent_id)
DROP TABLE IF EXISTS categories CASCADE;
-- 6. Пользовательский ENUM, больше не нужен
DROP TYPE IF EXISTS attr_data_type CASCADE;
-- 7. Заготовленные пресеты
DROP TABLE IF EXISTS presets CASCADE;
-- 8. Товары пресета
DROP TABLE IF EXISTS preset_items CASCADE;
-- 9. Таблица presets
DROP TABLE IF EXISTS presets CASCADE;
-- 10. Таблица preset_items
DROP TABLE IF EXISTS preset_items CASCADE;