import React, { useState } from 'react';
import { presetAPI, tokenUtils } from '../services/api';

/**
 * Компонент‑админка для управления «стилями» (пресетами товаров).
 *
 * props:
 *   products  – массив товаров из API (product_id, name, price)
 *   styles    – массив существующих стилей (preset_id, name, items[] …)
 *   setStyles – setState‑функция из родителя, хранящего глобальный список
 */
function StyleAdmin({ products, styles, setStyles }) {
  const [styleName, setStyleName] = useState('');
  const [selectedProductIds, setSelectedProductIds] = useState([]);
  const [imageUrl, setImageUrl] = useState('');
  const [editingStyle, setEditingStyle] = useState(null);

  // Bearer‑токен администратора
  const token = React.useMemo(() => tokenUtils.get(), []);

  /* ------------------------------------------------------------ */
  /* helpers                                                      */
  /* ------------------------------------------------------------ */
  const normalizePreset = (p) => ({ ...p, preset_id: p.preset_id ?? p.id });

  const resetForm = () => {
    setStyleName('');
    setSelectedProductIds([]);
    setImageUrl('');
    setEditingStyle(null);
  };

  const toggleProductSelection = (id) => {
    setSelectedProductIds((prev) =>
      prev.includes(id) ? prev.filter((pid) => pid !== id) : [...prev, id],
    );
  };

  /* ------------------------------------------------------------ */
  /* save / delete                                                */
  /* ------------------------------------------------------------ */
  const saveStyle = async () => {
    if (!styleName.trim() || selectedProductIds.length === 0) return;

    const items = selectedProductIds.map((id) => ({ product_id: Number(id) }));
    const total_price = products
      .filter((p) => selectedProductIds.includes(p.product_id))
      .reduce((sum, p) => sum + p.price, 0);

    const payload = {
      name: styleName.trim(),
      description: 'Обновлено через UI',
      image_url: imageUrl.trim(),
      items,
      total_price,
    };

    try {
      const id = editingStyle?.preset_id;
      if (editingStyle && Number.isInteger(id)) {
        // --- UPDATE ---------------------------------------------------------
        const updated = await presetAPI.update(id, payload, token);
        setStyles((prev) =>
          prev.map((s) =>
            normalizePreset(s).preset_id === updated.preset_id ? updated : s,
          ),
        );
      } else {
        // --- CREATE ---------------------------------------------------------
        const created = await presetAPI.create(payload, token);
        setStyles((prev) => {
          const pid = normalizePreset(created).preset_id;
          // защита от дубликатов
          const uniq = prev.filter((s) => normalizePreset(s).preset_id !== pid);
          return [...uniq, created];
        });
      }
      resetForm();
    } catch (err) {
      console.error('Ошибка при сохранении стиля:', err?.message || err);
    }
  };

  const removeStyle = async (id) => {
    try {
      setStyles((prev) =>
        prev.filter((s) => normalizePreset(s).preset_id !== id),
      );
      await presetAPI.delete(id, token);
      if (editingStyle?.preset_id === id) resetForm();
    } catch (err) {
      console.error('Ошибка при удалении:', err);
    }
  };

  /* ------------------------------------------------------------ */
  /* UI styles                                                    */
  /* ------------------------------------------------------------ */
  const ui = {
    input: {
      padding: '12px 16px',
      borderWidth: 1,
      borderStyle: 'solid',
      borderColor: '#334155',
      borderRadius: 10,
      background: 'rgba(15,23,42,0.6)',
      color: '#f1f5f9',
      fontSize: '1rem',
      width: '100%',
      transition: 'all .3s ease',
      boxShadow: '0 4px 10px rgba(0,0,0,.1)',
      outline: 'none',
    },
    btn: {
      background: 'linear-gradient(135deg,#3b82f6,#2563eb)',
      color: '#fff',
      border: 0,
      borderRadius: 10,
      padding: '12px 20px',
      fontSize: '1rem',
      fontWeight: 600,
      cursor: 'pointer',
      transition: 'all .3s ease',
      boxShadow: '0 4px 12px rgba(37,99,235,.3)',
      textTransform: 'uppercase',
      letterSpacing: '.5px',
    },
    delBtn: {
      background: 'rgba(185,28,28,.1)',
      color: '#f87171',
      borderWidth: 1,
      borderStyle: 'solid',
      borderColor: 'rgba(185,28,28,.3)',
      borderRadius: 8,
      padding: '8px 16px',
      fontSize: '.9rem',
      fontWeight: 500,
      cursor: 'pointer',
      transition: 'all .3s ease',
      marginLeft: 10,
    },
    checkbox: {
      display: 'flex',
      alignItems: 'center',
      padding: '10px 15px',
      borderWidth: 1,
      borderStyle: 'solid',
      borderColor: '#334155',
      borderRadius: 8,
      background: 'rgba(15,23,42,0.6)',
      marginBottom: 8,
      cursor: 'pointer',
      transition: 'all .3s ease',
    },
    checkboxChecked: {
      background: 'rgba(59,130,246,.1)',
      borderColor: '#3b82f6',
    },
    checkboxInput: {
      marginRight: 10,
      cursor: 'pointer',
    },
  };

  /* ------------------------------------------------------------ */
  /* render                                                       */
  /* ------------------------------------------------------------ */
  return (
    <div className="AdminSection" style={{ marginTop: 40 }}>
      {/* ---------------- Form create / edit ---------------- */}
      <div
        style={{
          background: 'rgba(30,41,59,.5)',
          padding: 20,
          borderRadius: 12,
          border: '1px solid #334155',
          marginBottom: 24,
        }}
      >
        <input
          value={styleName}
          onChange={(e) => setStyleName(e.target.value)}
          placeholder="Название стиля"
          style={ui.input}
        />
        <input
          value={imageUrl}
          onChange={(e) => setImageUrl(e.target.value)}
          placeholder="Ссылка на изображение (https://...)"
          style={{ ...ui.input, marginTop: 15 }}
        />

        <div style={{ marginTop: 20, marginBottom: 15 }}>
          <div style={{ marginBottom: 10, color: '#94a3b8' }}>
            Выберите товары для стиля:
          </div>
          <div style={{ maxHeight: 300, overflowY: 'auto', padding: 5 }}>
            {products.map((p) => {
              const checked = selectedProductIds.includes(p.product_id);
              return (
                <div
                  key={p.product_id}
                  style={{
                    ...ui.checkbox,
                    ...(checked ? ui.checkboxChecked : {}),
                  }}
                  onClick={() => toggleProductSelection(p.product_id)}
                >
                  <input
                    type="checkbox"
                    checked={checked}
                    onChange={() => toggleProductSelection(p.product_id)}
                    style={ui.checkboxInput}
                  />
                  <div>
                    <div style={{ fontWeight: 500, color: '#f1f5f9' }}>
                      {p.name}
                    </div>
                    <div
                      style={{
                        fontSize: '.9rem',
                        color: '#94a3b8',
                        marginTop: 3,
                      }}
                    >
                      {p.price} ₽
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        <button
          onClick={saveStyle}
          style={{ ...ui.btn, marginTop: 15, width: '100%' }}
          disabled={!styleName.trim() || selectedProductIds.length === 0}
        >
          {editingStyle ? '💾 Сохранить изменения' : '➕ Добавить стиль'}
        </button>
      </div>

      {/* ---------------- List of presets ---------------- */}
      <div
        style={{
          background: 'rgba(30,41,59,.5)',
          borderRadius: 12,
          overflow: 'hidden',
          border: '1px solid #334155',
        }}
      >
        <div
          style={{
            padding: '15px 20px',
            borderBottom: '1px solid #334155',
            color: '#94a3b8',
          }}
        >
          Существующие стили
        </div>
        <ul style={{ listStyle: 'none', margin: 0, padding: 0 }}>
          {styles.map((s) => {
            const pid = normalizePreset(s).preset_id;
            const isEditing = editingStyle?.preset_id === pid;
            return (
              <li
                key={pid}
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  padding: '14px 20px',
                  borderBottom: '1px solid rgba(51,65,85,.5)',
                  borderLeft: isEditing ? '4px solid #60a5fa' : 'none',
                  backgroundColor: isEditing
                    ? 'rgba(59,130,246,.05)'
                    : 'transparent',
                }}
              >
                <div>
                  <div
                    style={{
                      fontWeight: 500,
                      color: '#f1f5f9',
                      fontSize: '1.1rem',
                    }}
                  >
                    {s.name}
                  </div>
                  <div
                    style={{
                      color: '#94a3b8',
                      fontSize: '.9rem',
                      marginTop: 5,
                    }}
                  >
                    Товаров: {s.items?.length || 0}
                  </div>
                </div>
                <div style={{ display: 'flex', gap: 8 }}>
                  <button
                    onClick={() => {
                      const norm = normalizePreset(s);
                      setEditingStyle(norm);
                      setStyleName(norm.name);
                      setImageUrl(norm.image_url || '');
                      setSelectedProductIds(
                        norm.items?.map((i) => i.product_id ?? i.product?.id) ||
                          [],
                      );
                    }}
                    style={{
                      ...ui.delBtn,
                      background: 'rgba(34,197,94,.1)',
                      color: '#4ade80',
                    }}
                  >
                    ✎ Редактировать
                  </button>

                  <button onClick={() => removeStyle(pid)} style={ui.delBtn}>
                    🗑 Удалить
                  </button>
                </div>
              </li>
            );
          })}

          {styles.length === 0 && (
            <li style={{ padding: 20, textAlign: 'center', color: '#94a3b8' }}>
              Нет добавленных стилей
            </li>
          )}
        </ul>
      </div>
    </div>
  );
}

export default StyleAdmin;
