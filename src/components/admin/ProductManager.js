import React, { useState } from "react";
import { productAPI } from "../../services/api";

/* ------------------------------------------------------------
 * Helpers
 * ----------------------------------------------------------*/
const getProductId = (p) => p?.product_id ?? p?.id ?? null;

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

const stringifyAttributes = (arr = []) =>
  arr.map((a) => `${a.attribute_id}:${a.value}`).join("; ");

/* ============================================================
 * ProductManager
 * ==========================================================*/
function ProductManager({
  categories,
  products,
  setProducts,
  getAdminToken,
  showMessage,
  styles,
}) {
  /* --------------------- local state ---------------------- */
  const [form, setForm] = useState({
    name: "",
    price: "",
    description: "",
    image_url: "",
    categoryId: categories[0]?.id ?? 1,
    attributes: "",
  });
  const [editingId, setEditingId] = useState(null);

  const resetForm = () => {
    setForm({
      name: "",
      price: "",
      description: "",
      image_url: "",
      categoryId: categories[0]?.id ?? 1,
      attributes: "",
    });
    setEditingId(null);
  };

  const { inputStyle, buttonStyle, deleteButtonStyle } = styles;
  const editBtnStyle = {
    ...deleteButtonStyle,
    background: "rgba(34,197,94,.1)",
    color: "#4ade80",
  };

  /* ---------------- CREATE / UPDATE ----------------------- */
  const saveProduct = async () => {
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
        attributes: form.attributes.trim()
          ? parseAttributes(form.attributes)
          : [],
      };

      if (editingId) {
        /* -------- UPDATE ---------- */
        const updated = await productAPI.update(editingId, payload, token);
        setProducts((prev) =>
          prev.map((p) =>
            getProductId(p) === editingId ? { ...updated, categoryId: updated.category_id } : p
          )
        );
        showMessage("Товар обновлён");
      } else {
        /* -------- CREATE ---------- */
        const created = await productAPI.create(payload, token);
        if (!created?.product_id) {
          showMessage("Сервер не вернул ID", true);
          return;
        }
        setProducts((prev) => [...prev, { ...created, categoryId: created.category_id }]);
        showMessage("Товар добавлен");
      }
      resetForm();
    } catch (err) {
      console.error("saveProduct", err);
      showMessage(err.message || "Ошибка сервера", true);
    }
  };

  /* --------------------- DELETE --------------------------- */
  const removeProduct = async (prod) => {
    const id = getProductId(prod);
    if (!id) return showMessage("Некорректный ID", true);

    try {
      const token = await getAdminToken();
      if (!token) return;

      await productAPI.delete(id, token);
      setProducts((prev) => prev.filter((p) => getProductId(p) !== id));
      if (editingId === id) resetForm();
      showMessage("Товар удалён");
    } catch (err) {
      showMessage(err.message || "Ошибка при удалении", true);
    }
  };

  /* ---------------------- render -------------------------- */
  return (
    <div className="AdminSection" style={{ marginTop: 40 }}>
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

      {/* ------------------ form ------------------ */}
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
          onClick={saveProduct}
          style={{ ...buttonStyle, gridColumn: "1/-1", marginTop: 10, padding: 14 }}
        >
          {editingId ? "💾 Сохранить" : "➕ Добавить"}
        </button>
      </div>

      {/* ------------------ list ------------------ */}
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
          const isEditing = editingId === id;
          return (
            <li
              key={id ?? `${p.name}_${p.price}`}
              style={{
                display: "flex",
                alignItems: "center",
                gap: 16,
                padding: "14px 20px",
                borderBottom: "1px solid rgba(51,65,85,0.5)",
                borderLeft: isEditing ? "4px solid #60a5fa" : "none",
                backgroundColor: isEditing ? "rgba(59,130,246,.05)" : "transparent",
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
              <button
                onClick={() => {
                  setForm({
                    name: p.name,
                    price: p.price,
                    description: p.description ?? "",
                    image_url: p.image_url ?? "",
                    categoryId: p.categoryId ?? p.category_id,
                    attributes: stringifyAttributes(p.attributes),
                  });
                  setEditingId(id);
                }}
                style={editBtnStyle}
              >
                ✎ Редактировать
              </button>
              <button onClick={() => removeProduct(p)} style={deleteButtonStyle}>
                🗑 Удалить
              </button>
            </li>
          );
        })}
      </ul>
    </div>
  );
}

export default ProductManager;
