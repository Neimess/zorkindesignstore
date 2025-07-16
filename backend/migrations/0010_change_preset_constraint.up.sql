ALTER TABLE preset_items
DROP CONSTRAINT IF EXISTS preset_items_product_id_fkey;

ALTER TABLE preset_items
ADD CONSTRAINT preset_items_product_id_fkey
FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE;