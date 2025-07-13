import React, { useState } from 'react';
import { categoryAPI } from '../../services/api';

function CategoryManager({
  categories,
  setCategories,
  getAdminToken,
  showMessage,
  styles,
  onViewCategoryProducts,
}) {
  const [catName, setCatName] = useState('');
  const [catType, setCatType] = useState('room');

  const [roomType, setRoomType] = useState(''); // «Спальня», «Ванная» …
  const [parentRoom, setParentRoom] = useState('');
  const [parentElement, setParentElement] = useState('');
  const { inputStyle, buttonStyle, deleteButtonStyle } = styles;
const [isSubmitting, setIsSubmitting] = useState(false);
  

const addCategory = async () => {
  console.log('🔔 addCategory вызван');

  if (!catName.trim()) return;

  const token = await getAdminToken();
  if (!token) return;

  const parent_id =
    catType === 'room'
      ? null
      : catType === 'element'
        ? Number(parentRoom)
        : Number(parentElement);

  const payload = {
    name: catName.trim(),
    parent_id,
    description: catType === 'room' ? roomType.trim() : undefined,
  };

  try {
    const created = await categoryAPI.create(payload, token);
    setCategories((prev) => [...prev, created]);
    showMessage('Категория успешно добавлена');

    // сбрасываем форму
    setCatName('');
    setRoomType('');
    setCatType('room');
    setParentRoom('');
    setParentElement('');
  } catch (e) {
    console.error('Ошибка создания категории:', e);
    showMessage(e.message || 'Ошибка при создании категории', true);
  }
};


  const removeCategory = async (id) => {
    console.log('Удаляем категорию с id:', id);
    try {
      const token = await getAdminToken();
      if (!token) return;
      await categoryAPI.delete(id, token);
      setCategories(categories.filter((c) => c.id !== id));
      showMessage('Категория успешно удалена');
    } catch (error) {
      console.error('Ошибка удаления категории:', error);
      showMessage('Ошибка при удалении категории', true);
    }
  };

  return (
    <div className="AdminSection">
      <h2
        style={{
          fontSize: '1.5rem',
          color: '#f8fafc',
          marginBottom: '20px',
          position: 'relative',
          paddingBottom: '10px',
        }}
      >
        Категории
        <span
          style={{
            position: 'absolute',
            bottom: 0,
            left: 0,
            width: '60px',
            height: '3px',
            background: 'linear-gradient(90deg, #3b82f6, #60a5fa)',
            borderRadius: '2px',
          }}
        ></span>
      </h2>

      <div
        style={{
          display: 'flex',
          gap: 12,
          marginBottom: 24,
          alignItems: 'center',
        }}
      >
        <select
          value={catType}
          onChange={(e) => {
            setCatType(e.target.value);
            setParentRoom('');
            setParentElement('');
          }}
          style={{ ...inputStyle }}
        >
          <option value="room">🏠 Комната</option>
          <option value="element">📦 Элемент (внутри комнаты)</option>
          <option value="sub">🔹 Подкатегория (внутри элемента)</option>
        </select>

        {catType !== 'room' && (
          <select
            value={parentRoom}
            onChange={(e) => setParentRoom(e.target.value)}
            style={inputStyle}
          >
            <option value="">— Выбери комнату —</option>
            {categories
              .filter((c) => c.parent_id === null)
              .map((room) => (
                <option key={room.id} value={room.id}>
                  {room.name}
                </option>
              ))}
          </select>
        )}

        {catType === 'sub' && parentRoom && (
          <select
            value={parentElement}
            onChange={(e) => setParentElement(e.target.value)}
            style={inputStyle}
          >
            <option value="">— Выбери элемент —</option>
            {categories
              .filter((c) => c.parent_id === Number(parentRoom))
              .map((elem) => (
                <option key={elem.id} value={elem.id}>
                  {elem.name}
                </option>
              ))}
          </select>
        )}

        {catType === 'room' && (
          <input
            value={roomType}
            onChange={(e) => setRoomType(e.target.value)}
            placeholder="Тип комнаты (например: Гостиная)"
            style={inputStyle}
          />
        )}

        <input
          value={catName}
          onChange={(e) => setCatName(e.target.value)}
          placeholder="Название категории"
          style={inputStyle}
        />

        <button onClick={addCategory} style={buttonStyle}>
          Добавить
        </button>
      </div>

      <ul
        style={{
          marginBottom: 30,
          listStyle: 'none',
          padding: 0,
          background: 'rgba(30, 41, 59, 0.5)',
          borderRadius: '12px',
          overflow: 'hidden',
          border: '1px solid #334155',
        }}
      >
        {categories.map((c) => (
          <li
            key={c.id}
            style={{
              padding: '14px 20px',
              borderBottom: '1px solid rgba(51, 65, 85, 0.5)',
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              transition: 'all 0.3s ease',
            }}
          >
            <span style={{ fontSize: '1.1rem', fontWeight: 500 }}>
              {c.name}
            </span>
            <div style={{ display: 'flex', gap: 8 }}>
              <button
                onClick={() => onViewCategoryProducts(c.id)}
                style={{
                  ...buttonStyle,
                  backgroundColor: '#0ea5e9',
                  borderColor: '#0284c7',
                }}
              >
                Товары
              </button>
              <button
                onClick={() => removeCategory(c.id)}
                style={deleteButtonStyle}
              >
                <i
                  className="fas fa-trash-alt"
                  style={{ marginRight: '6px' }}
                ></i>
                Удалить
              </button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default CategoryManager;
