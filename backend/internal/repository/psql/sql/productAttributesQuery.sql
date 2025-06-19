SELECT
    pa.product_id,
    pa.attribute_id,
    pa.value_string,
    pa.value_int,
    pa.value_float,
    pa.value_bool,
    pa.value_enum,
    a.attribute_id            AS "attribute.attribute_id",
    a.name                    AS "attribute.name",
    a.slug                    AS "attribute.slug",
    a.data_type               AS "attribute.data_type",
    a.unit                    AS "attribute.unit",
    a.is_filterable           AS "attribute.is_filterable"
FROM product_attributes pa
JOIN attributes a ON pa.attribute_id = a.attribute_id
WHERE pa.product_id = $1
