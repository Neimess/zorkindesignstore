import React, { useState, useEffect } from 'react';
import { productAPI, categoryAttributeAPI, serviceAPI } from '../../services/api';

/* ------------------------------------------------------------
 * Helpers
 * ----------------------------------------------------------*/

const getProductId = (p) => p?.product_id ?? p?.id ?? null;

const parseAttributes = (input) => {
  if (!input.trim()) return [];

  return input.split(';').map((item) => {
    const [namePart, valuePart] = item.split(':');

    return {
      name: namePart ? namePart.trim() : '',
      unit: '', // Или добавляй единицы измерения, если нужно
      value: valuePart ? valuePart.trim() : '',
    };
  });
};
    

const stringifyAttributes = (arr = []) =>
  arr
    .map((a) => {
      // если с бэка пришло имя, показываем его, иначе fallback на id
      const label = a.name || a.attribute_name || a.attribute_id;
      return `${label}:${a.value}`;
    })
    .join('; ');


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
  const [roomId, setRoomId] = useState('');
  const [elementId, setElementId] = useState('');
  const [subId, setSubId] = useState('');

  const [availableServices, setAvailableServices] = useState([]);

useEffect(() => {
  const fetchServices = async () => {
    try {
      const token = await getAdminToken();
      console.log('Получаем токен для сервисов:', token);
      if (!token) return;
      const data = await serviceAPI.getAll(token);
      setAvailableServices(data || []);
    } catch (err) {
      console.error('Ошибка загрузки сервисов:', err);
    }
  };
  fetchServices();
}, []);

const [allowedAttrs, setAllowedAttrs] = useState([]);

useEffect(() => {
  if (!subId) return;
  (async () => {
    try {
      const arr = await categoryAttributeAPI.getAll(subId);
      setAllowedAttrs(arr);      // [{id, name, unit}, ...]
    } catch (e) {
      console.error('Не смог получить атрибуты подкатегории', e);
    }
  })();
}, [subId]);

  
  // Добавляем отладочную информацию
  console.log('ProductManager - Все категории:', categories);
  
  // Проверяем, что categories - это массив перед фильтрацией
  const rooms = Array.isArray(categories) ? categories.filter(c => c.parent_id == null) : [];
  console.log('ProductManager - Комнаты:', rooms);

  const elements = Array.isArray(categories) ? categories.filter(
    c => roomId && c.parent_id === Number(roomId)
  ) : [];
  console.log('ProductManager - Элементы для комнаты', roomId, ':', elements);

  const subs = Array.isArray(categories) ? categories.filter(
    c => elementId && c.parent_id === Number(elementId)
  ) : [];
  console.log('ProductManager - Подкатегории для элемента', elementId, ':', subs);

/* --------------------- local state ---------------------- */
  const [form, setForm] = useState({
  name: '',
  price: '',
  description: '',
  image_url: '',
  categoryId: categories[0]?.id ?? 1,
  attributes: '',      // пока строкой
  services: [],        // <-- массив
});

  const [editingId, setEditingId] = useState(null);

  const resetForm = () => {
    setForm({
      name: '',
      price: '',
      description: '',
      image_url: '',
      categoryId: categories[0]?.id ?? 1,
      attributes: '',
      services: [],
    });
    setEditingId(null);
    setRoomId('');
    setElementId('');
    setSubId('');
  };

  const { inputStyle, buttonStyle, deleteButtonStyle } = styles;
  const editBtnStyle = {
    ...deleteButtonStyle,
    background: 'rgba(34,197,94,.1)',
    color: '#4ade80',
  };
 

const saveProduct = async () => {
  if (!form.name.trim() || !form.price) return;

  try {
    const token = await getAdminToken();
    if (!token) return;

    if (!subId) {
      showMessage('Выберите подкатегорию', true);
      return;
    }

  
    // Парсинг атрибутов
    // 🔻 внутри saveProduct (или где формируешь payload)
const preparedAttributes = Array.isArray(form.attributes)
  ? form.attributes
  : parseAttributes(form.attributes);   // ← твой парсер строки "Цвет: белый"

  const isUpdate = Boolean(editingId);
const attributesForApi = preparedAttributes.map((a) => {
  if (isUpdate) {
    // В update unit всегда есть
    return {
      name: a.name?.trim() || '',
      unit: (a.unit ?? '').trim(),
      value: a.value?.toString().trim() || '',
    };
  } else {
    // В create unit добавляется только если не пустой
    const attr = {
      name: a.name?.trim() || '',
      value: a.value?.toString().trim() || '',
    };

    if (a.unit && a.unit.trim().length > 0) {
      attr.unit = a.unit.trim();
    }

    return attr;
  }
});







    // Парсинг услуг
    const preparedServices = Array.isArray(form.services)
  ? form.services.map((id) => ({ service_id: Number(id) }))
  : [];


console.log('ATTR', preparedAttributes);
  console.log('SERV', preparedServices);






   const payload = {
  name        : form.name.trim(),
  price       : Number(form.price),
  description : form.description.trim(),
  image_url   : form.image_url.trim(),
  category_id : Number(subId),
  attributes  : attributesForApi,           // ← уже без пустого unit
  services    : preparedServices            // как раньше
};


 console.log("ЧЕНКНИ запрос", payload)

  /* ─── вызов API ──────────────────────────────────────────── */
  const hasExtras = payload.attributes.length || payload.services.length;
  // console.log('payload:', payload);

    if (editingId) {
  const updated = await productAPI.update(editingId, payload, token);

  setProducts(prev =>
    prev.map(p =>
      getProductId(p) === editingId
        ? { ...updated, categoryId: updated.category_id }
        : p
    )
  );

  showMessage('Товар обновлён');
}
else {
      // -------- CREATE ----------
     const created = await productAPI.create(payload, token);

      if (!created?.product_id && !created?.id) {
      showMessage('Сервер не вернул ID', true);
      return;
    }

      const productId = created.product_id || created.id;
      setProducts((prev) => [
      ...prev,
      { ...created, product_id: productId, categoryId: created.category_id },
    ]);

    showMessage('Товар добавлен');
   resetForm();
    setRoomId('');
    setElementId('');
    setSubId('');  
  }


  } catch (err) {
    console.error('saveProduct', err);
    showMessage(err.message || 'Ошибка сервера', true);
  }
};

  /* --------------------- DELETE --------------------------- */
  const removeProduct = async (prod) => {
    const id = getProductId(prod);
    if (!id) return showMessage('Некорректный ID', true);

    try {
      const token = await getAdminToken();
      if (!token) return;

      await productAPI.delete(id, token);
      setProducts((prev) => prev.filter((p) => getProductId(p) !== id));
      if (editingId === id) {
        resetForm();
        setRoomId('');
        setElementId('');
        setSubId('');
      }
      showMessage('Товар удалён');
    } catch (err) {
      console.error('removeProduct', err);
      showMessage(err.message || 'Ошибка при удалении', true);
    }
  };

  /* ---------------------- render -------------------------- */
  return (
    <div className="AdminSection" style={{ marginTop: 40 }}>
      <h2 style={{ fontSize: '1.5rem', color: '#f8fafc', marginBottom: 20 }}>
        Товары
        <span
          style={{
            display: 'block',
            width: 60,
            height: 3,
            marginTop: 4,
            background: 'linear-gradient(90deg,#3b82f6,#60a5fa)',
            borderRadius: 2,
          }}
        />
      </h2>

      {/* ------------------ form ------------------ */}
      <div
        style={{
          display: 'grid',
          gridTemplateColumns: '1fr 1fr',
          gap: 16,
          marginBottom: 24,
          background: 'rgba(30,41,59,0.5)',
          padding: 20,
          borderRadius: 12,
          border: '1px solid #334155',
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
        {/* -- Комната -- */}
        <div style={{ position: 'relative', width: '100%' }}>
          <select
            value={roomId}
            onChange={(e) => {
              setRoomId(e.target.value);
              setElementId('');
              setSubId('');
            }}
            style={{
              ...inputStyle,
              borderColor: !rooms || rooms.length === 0 ? '#ef4444' : inputStyle.borderColor
            }}
          >
            <option value="">— Комната —</option>
            {rooms && rooms.length > 0 ? (
              rooms.map((r) => (
                <option key={r.id} value={r.id}>
                  {r.name}
                </option>
              ))
            ) : (
              <option value="" disabled>Нет доступных комнат</option>
            )}
          </select>
          {(!rooms || rooms.length === 0) && (
            <div style={{ 
              color: '#ef4444', 
              fontSize: '12px', 
              marginTop: '4px',
              position: 'absolute',
              bottom: '-20px',
              left: '0'
            }}>
              Необходимо создать комнаты в разделе "Категории"
            </div>
          )}
        </div>

        {/* — Элемент — */}
        <div style={{ position: 'relative', width: '100%' }}>
          <select
            value={elementId}
            onChange={e => {
              setElementId(e.target.value);
              setSubId('');
            }}
            disabled={!roomId}
            style={{
              ...inputStyle,
              opacity: !roomId ? 0.7 : 1,
              borderColor: roomId && (!elements || elements.length === 0) ? '#ef4444' : inputStyle.borderColor
            }}
          >
            <option value="">— Элемент —</option>
            {elements && elements.length > 0 ? (
              elements.map(el => (
                <option key={el.id} value={el.id}>
                  {el.name}
                </option>
              ))
            ) : (
              <option value="" disabled>{roomId ? 'Нет элементов для этой комнаты' : 'Сначала выберите комнату'}</option>
            )}
          </select>
          {roomId && (!elements || elements.length === 0) && (
            <div style={{ 
              color: '#ef4444', 
              fontSize: '12px', 
              marginTop: '4px',
              position: 'absolute',
              bottom: '-20px',
              left: '0'
            }}>
              Необходимо создать элементы для этой комнаты
            </div>
          )}
        </div>

        {/* -- Подкатегория -- */}
        <div style={{ position: 'relative', width: '100%' }}>
          <select
            value={subId}
            onChange={(e) => setSubId(e.target.value)}
            disabled={!elementId}
            style={{
              ...inputStyle,
              opacity: !elementId ? 0.7 : 1,
              borderColor: elementId && (!subs || subs.length === 0) ? '#ef4444' : inputStyle.borderColor
            }}
          >
            <option value="">— Подкатегория —</option>
            {subs && subs.length > 0 ? (
              subs.map((s) => (
                <option key={s.id} value={s.id}>
                  {s.name}
                </option>
              ))
            ) : (
              <option value="" disabled>{elementId ? 'Нет подкатегорий для этого элемента' : 'Сначала выберите элемент'}</option>
            )}
          </select>
          {elementId && (!subs || subs.length === 0) && (
            <div style={{ 
              color: '#ef4444', 
              fontSize: '12px', 
              marginTop: '4px',
              position: 'absolute',
              bottom: '-20px',
              left: '0'
            }}>
              Необходимо создать подкатегории для этого элемента
            </div>
          )}
        </div>

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
        <select
  multiple
  value={form.services}
  onChange={(e) => {
    const selected = Array.from(e.target.selectedOptions).map((opt) => opt.value);
    setForm({ ...form, services: selected });
  }}
  style={inputStyle}
>
  {availableServices.map((s) => (
    <option key={s.id} value={s.id}>
      {s.name} ({s.price} ₽)
    </option>
  ))}
</select>


        <button
          onClick={saveProduct}
          style={{
            ...buttonStyle,
            gridColumn: '1/-1',
            marginTop: 10,
            padding: 14,
          }}
        >
          {editingId ? '💾 Сохранить' : '➕ Добавить'}
        </button>
      </div>

      {/* ------------------ list ------------------ */}
      <ul
        style={{
          listStyle: 'none',
          padding: 0,
          margin: 0,
          background: 'rgba(30,41,59,0.5)',
          borderRadius: 12,
          border: '1px solid #334155',
          maxHeight: 400,
          overflowY: 'auto',
        }}
      >
        {products.length === 0 && (
          <li style={{ padding: 20, textAlign: 'center', color: '#94a3b8' }}>
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
                display: 'flex',
                alignItems: 'center',
                gap: 16,
                padding: '14px 20px',
                borderBottom: '1px solid rgba(51,65,85,0.5)',
                borderLeft: isEditing ? '4px solid #60a5fa' : 'none',
                backgroundColor: isEditing
                  ? 'rgba(59,130,246,.05)'
                  : 'transparent',
              }}
            >
              <img
                src={p.image_url}
                alt={p.name}
                style={{
                  width: 60,
                  height: 60,
                  objectFit: 'cover',
                  borderRadius: 10,
                  background: '#1e293b',
                  border: '1px solid #334155',
                }}
              />
              <div style={{ flex: 1 }}>
                <div style={{ fontWeight: 600 }}>{p.name}</div>
                <div style={{ color: '#94a3b8', fontSize: 14 }}>
                  {categories.find((c) => c.id === p.categoryId)?.name || '—'}
                </div>
                {p.attributes && p.attributes.length > 0 && (
                  <div style={{ color: '#94a3b8', fontSize: 12, marginTop: 4 }}>
                    <span style={{ color: '#60a5fa' }}>Атрибуты:</span>{' '}
                    {p.attributes.map((attr, idx) => (
                      <span key={attr.attribute_id}>
                        {attr.name}: {attr.value}{attr.unit ? ` ${attr.unit}` : ''}
                        {idx < p.attributes.length - 1 ? ', ' : ''}
                      </span>
                    ))}
                  </div>
                )}
                {p.services && p.services.length > 0 && (
                  <div style={{ color: '#94a3b8', fontSize: 12, marginTop: 2 }}>
                    <span style={{ color: '#60a5fa' }}>Сервисы:</span>{' '}
                    {p.services.map((service, idx) => (
                      <span key={service.id}>
                        {service.name} ({service.price} ₽)
                        {idx < p.services.length - 1 ? ', ' : ''}
                      </span>
                    ))}
                  </div>
                )}
              </div>
              <div style={{ fontWeight: 700, color: '#38bdf8' }}>
                {p.price} ₽
              </div>
              <button
                onClick={() => {
                  const categoryId = p.categoryId ?? p.category_id;
                  const category = categories.find(c => c.id === categoryId);
                  
                  // Находим родительские категории
                  if (category) {
                    const element = categories.find(c => c.id === category.parent_id);
                    if (element) {
                      const room = categories.find(c => c.id === element.parent_id);
                      
                      // Устанавливаем значения в правильном порядке
                      if (room) setRoomId(room.id.toString());
                      setElementId(element.id.toString());
                      setSubId(category.id.toString());
                    }
                  }
                  
                  setForm({
  name: p.name,
  price: p.price,
  description: p.description ?? '',
  image_url: p.image_url ?? '',
  categoryId: categoryId,
  attributes: stringifyAttributes(p.attributes),   // пока оставим строкой
  services: (p.services || []).map((s) => String(s.id)), // <-- массив строк
});

                  setEditingId(id);
                }}
                style={editBtnStyle}
              >
                ✎ Редактировать
              </button>
              <button
                onClick={() => removeProduct(p)}
                style={deleteButtonStyle}
              >
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
