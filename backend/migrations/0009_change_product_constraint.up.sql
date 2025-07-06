ALTER TABLE products
DROP CONSTRAINT products_category_id_fkey;

ALTER TABLE products
ADD CONSTRAINT products_category_id_fkey
    FOREIGN KEY (category_id)
    REFERENCES categories(category_id)
    ON DELETE CASCADE;
