CREATE TABLE IF NOT EXISTS product_services (
    product_service_id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(product_id) ON DELETE CASCADE,
    service_id BIGINT NOT NULL REFERENCES services(service_id) ON DELETE CASCADE,
    UNIQUE (product_id, service_id)
);