import React, { useState } from "react";
import { productAPI } from "../../services/api";

/**
 * Безопасно достаём ID товара, т.к. бэкенд иногда присылает
 *  { id }, а иногда { product_id }
 */
const getProductId = (prod) => prod?.id ?? prod?.product_id ?? null;

/**
 * Парсинг строки атрибутов вида “2:1.25; 3:Матовый” в массив
 * [{ attribute_id: 2, value: "1.25" }, { attribute_id: 3, value: "Матовый" }]
 */
const parseAttributes = (raw) =>
  raw
    .split(";")
    .map((chunk) => {
      const [idPart, valuePart] = chunk.split(":").map((s) => s.trim());
      const attribute_id = Number(idPart);
      if (!attribute_id || !valuePart) return null;
      return { attribute_id, value: valuePart };
    })
    .filter(Boolean);

/**
 * Компонент «ProductManager» – CRUD для товаров в админ‑панели
 */
function ProductManager({
  categories,
  products,
  setProducts,
  getAdminToken,
  showMessage,
  styles,
}) {
  /* ------------------------ state формы нового товара -------------------- */
  const [form, setForm] = useState({
    name: "",
    price: "",
    description: "",
    image_url: "",
    categoryId: categories[0]?.id ?? 1,
    attributes: "",
  });

  const resetForm = () =>
    setForm({
      name: "",
      price: "",
      description: "",
      image_url: "",
      categoryId: categories[0]?.id ?? 1,
      attributes: "",
    });

  const { inputStyle, buttonStyle, deleteButtonStyle } = styles;

  /* ----------------------------- CREATE ---------------------------------- */
  const addProduct = async () => {
    if (!form.name.trim() || !form.price) return;

    try {
      const token = await getAdminToken();
      if (!token) return;

      const payload = {
        name: form.name.trim(),
        price: Number(form.price),
        description: form.description.trim(),
        image_url: form.image_url.trim(),
        category_id: Number(form.categoryId),
        attributes: form.attributes.trim() ? parseAttributes(form.attributes) : [],
      };

      console.log("➡️  Создаём товар:", payload);
      const created = await productAPI.create(payload, token);
      console.log("✅  Ответ create:", created);

      if (!created || !created.product_id) {
  showMessage("Ошибка: сервер не вернул ID", true);
  return;
}


      // Бэкенд теперь возвращает полную модель → просто кладём её в store
      const newList = [...products, { ...created, categoryId: created.category_id }];
      setProducts(newList);
      resetForm();
      showMessage("Товар успешно добавлен");
    } catch (err) {
      const detail = err?.response?.message ?? err.message;
      showMessage(`Ошибка сервера: ${detail}`, true);
      console.error("❌ addProduct:", err);
      showMessage(err.message || "Ошибка при создании", true);
    }
  };

  /* ----------------------------- DELETE ---------------------------------- */
  const removeProduct = async (prod) => {
    const id = getProductId(prod);
    if (!id) {
      showMessage("Некорректный ID", true);
      return;
    }

    try {
      const token = await getAdminToken();
      if (!token) return;

      console.log("🗑️  Удаляем товар", id);
      await productAPI.delete(id, token);
      setProducts((prev) => prev.filter((p) => getProductId(p) !== id));
      showMessage("Товар удалён");
    } catch (err) {
      console.error("❌ removeProduct:", err);
      showMessage(err.message || "Ошибка при удалении", true);
    }
  };

  /* ------------------------------ UI ------------------------------------- */
  return (
    <div className="AdminSection" style={{ marginTop: 40 }}>
      {/* ====== header ====== */}
      <h2 style={{ fontSize: "1.5rem", color: "#f8fafc", marginBottom: 20 }}>
        Товары
        <span
          style={{
            display: "block",
            width: 60,
            height: 3,
            marginTop: 4,
            background: "linear-gradient(90deg,#3b82f6,#60a5fa)",
            borderRadius: 2,
          }}
        />
      </h2>

      {/* ====== form ====== */}
      <div
        style={{
          display: "grid",
          gridTemplateColumns: "1fr 1fr",
          gap: 16,
          marginBottom: 24,
          background: "rgba(30,41,59,0.5)",
          padding: 20,
          borderRadius: 12,
          border: "1px solid #334155",
        }}
      >
        <input
          value={form.name}
          placeholder="Название"
          onChange={(e) => setForm({ ...form, name: e.target.value })}
          style={inputStyle}
        />
        <input
          type="number"
          value={form.price}
          placeholder="Цена"
          onChange={(e) => setForm({ ...form, price: e.target.value })}
          style={inputStyle}
        />
        <select
          value={form.categoryId}
          onChange={(e) => setForm({ ...form, categoryId: e.target.value })}
          style={{ ...inputStyle, appearance: "none", paddingRight: 40 }}
        >
          {categories.map((c) => (
            <option key={c.id} value={c.id}>
              {c.name}
            </option>
          ))}
        </select>
        <input
          value={form.image_url}
          placeholder="URL картинки"
          onChange={(e) => setForm({ ...form, image_url: e.target.value })}
          style={inputStyle}
        />
        <input
          value={form.description}
          placeholder="Описание"
          onChange={(e) => setForm({ ...form, description: e.target.value })}
          style={inputStyle}
        />
        <input
          value={form.attributes}
          placeholder="Атрибуты: 2:1.25; 3:Матовый"
          onChange={(e) => setForm({ ...form, attributes: e.target.value })}
          style={inputStyle}
        />
        <button
          onClick={addProduct}
          style={{ ...buttonStyle, gridColumn: "1/-1", marginTop: 10, padding: 14 }}
        >
          <i className="fas fa-plus" style={{ marginRight: 8 }} /> Добавить
        </button>
      </div>

      {/* ====== list ====== */}
      <ul
        style={{
          listStyle: "none",
          padding: 0,
          margin: 0,
          background: "rgba(30,41,59,0.5)",
          borderRadius: 12,
          border: "1px solid #334155",
          maxHeight: 400,
          overflowY: "auto",
        }}
      >
        {products.length === 0 && (
          <li style={{ padding: 20, textAlign: "center", color: "#94a3b8" }}>
            Список товаров пуст
          </li>
        )}

        {products.map((p) => {
          const id = getProductId(p);
          return (
            <li
              key={id ?? `${p.name}_${p.price}`}
              style={{
                display: "flex",
                alignItems: "center",
                gap: 16,
                padding: "14px 20px",
                borderBottom: "1px solid rgba(51,65,85,0.5)",
              }}
            >
              <img
                src={p.image_url}
                alt={p.name}
                style={{
                  width: 60,
                  height: 60,
                  objectFit: "cover",
                  borderRadius: 10,
                  background: "#1e293b",
                  border: "1px solid #334155",
                }}
              />
              <div style={{ flex: 1 }}>
                <div style={{ fontWeight: 600 }}>{p.name}</div>
                <div style={{ color: "#94a3b8", fontSize: 14 }}>
                  {categories.find((c) => c.id === p.categoryId)?.name || "—"}
                </div>
              </div>
              <div style={{ fontWeight: 700, color: "#38bdf8" }}>{p.price} ₽</div>
              <button onClick={() => removeProduct(p)} style={deleteButtonStyle}>
                <i className="fas fa-trash-alt" style={{ marginRight: 6 }} /> Удалить
              </button>
            </li>
          );
        })}
      </ul>
    </div>
  );
}

export default ProductManager;
