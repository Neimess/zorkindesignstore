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
  console.log('Текущие категории:', categories);  // ← вставь сюда
  const [catName, setCatName] = useState('');
  const [catType, setCatType] = useState('room');

  const [roomType, setRoomType] = useState(''); // «Спальня», «Ванная» …
  const [parentRoom, setParentRoom] = useState('');
  const [parentElement, setParentElement] = useState('');
  const { inputStyle, buttonStyle, deleteButtonStyle } = styles;
const [isSubmitting, setIsSubmitting] = useState(false);
  

const addCategory = async () => {
  console.log('🔔 addCategory вызван');
  console.log('Текущие категории:', categories);
  console.log('Тип категории:', catType);
  console.log('Родительская комната:', parentRoom);
  console.log('Родительский элемент:', parentElement);

  if (!catName.trim()) {
    showMessage('Введите название категории', true);
    return;
  }

  const token = await getAdminToken();
  if (!token) return;

  const parent_id =
    catType === 'room'
      ? null
      : catType === 'element'
        ? Number(parentRoom)
        : Number(parentElement);

  if ((catType === 'element' && !parentRoom) || (catType === 'sub' && !parentElement)) {
    showMessage('Выберите родительскую категорию', true);
    return;
  }

  const payload = {
    name: catName.trim(),
    parent_id,
    description: catType === 'room' ? roomType.trim() : undefined,
  };

  console.log('Отправляемые данные:', payload);

  try {
      setIsSubmitting(true);
    const created = await categoryAPI.create(payload, token);
    console.log('Созданная категория:', created);
    setCategories((prev) => {
      const newCategories = [...prev, created];
      console.log('Обновленные категории:', newCategories);
      return newCategories;
    });
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
  } finally {
    setIsSubmitting(false);
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
      showMessage(error.message || 'Ошибка при удалении категории', true);
    }
  };

  // Функция для создания тестовых категорий
  const createTestCategories = async () => {
    const token = await getAdminToken();
    if (!token) return;

    try {
      setIsSubmitting(true);
      showMessage('Создание тестовых категорий...');

      // Создаем комнаты
      const room1 = await categoryAPI.create({ name: 'Гостиная', parent_id: null, description: 'Гостиная' }, token);
      const room2 = await categoryAPI.create({ name: 'Спальня', parent_id: null, description: 'Спальня' }, token);
      
      // Создаем элементы для комнат
      const element1 = await categoryAPI.create({ name: 'Диван', parent_id: room1.id }, token);
      const element2 = await categoryAPI.create({ name: 'Стол', parent_id: room1.id }, token);
      const element3 = await categoryAPI.create({ name: 'Кровать', parent_id: room2.id }, token);
      
      // Создаем подкатегории для элементов
      await categoryAPI.create({ name: 'Угловой диван', parent_id: element1.id }, token);
      await categoryAPI.create({ name: 'Прямой диван', parent_id: element1.id }, token);
      await categoryAPI.create({ name: 'Обеденный стол', parent_id: element2.id }, token);
      await categoryAPI.create({ name: 'Двуспальная кровать', parent_id: element3.id }, token);

      // Обновляем список категорий
      const updatedCategories = await categoryAPI.getAll();
      setCategories(updatedCategories);

      showMessage('Тестовые категории успешно созданы!');
    } catch (e) {
      console.error('Ошибка создания тестовых категорий:', e);
      showMessage('Ошибка при создании тестовых категорий', true);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="AdminSection">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
        <h2
          style={{
            fontSize: '1.5rem',
            color: '#f8fafc',
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
        
        <button
          onClick={createTestCategories}
          disabled={isSubmitting}
          style={{
            background: 'linear-gradient(135deg, #10b981, #059669)',
            color: '#fff',
            border: 'none',
            borderRadius: '8px',
            padding: '8px 16px',
            fontSize: '0.9rem',
            fontWeight: 600,
            cursor: isSubmitting ? 'not-allowed' : 'pointer',
            transition: 'all 0.3s ease',
            boxShadow: '0 4px 12px rgba(16, 185, 129, 0.3)',
            opacity: isSubmitting ? 0.7 : 1,
          }}
        >
          {isSubmitting ? 'Создание...' : 'Создать тестовые категории'}
        </button>
      </div>

      <div
        style={{
          display: 'flex',
          gap: 12,
          marginBottom: 24,
          alignItems: 'center',
          flexWrap: 'wrap'
        }}
      >
        <select
          value={catType}
          onChange={(e) => {
            setCatType(e.target.value);
            // Сбрасываем значения при изменении типа
            setParentRoom('');
            setParentElement('');
          }}
          style={{ ...inputStyle, minWidth: '200px' }}
        >
          <option value="room">🏠 Комната</option>
          <option value="element">📦 Элемент (внутри комнаты)</option>
          <option value="sub">🔹 Подкатегория (внутри элемента)</option>
        </select>

        {/* Выпадающий список для выбора комнаты при создании элемента или подкатегории */}
        {(catType === 'element' || catType === 'sub') && (
          <div style={{ position: 'relative', width: '100%', marginTop: '10px' }}>
            <select
              value={parentRoom}
              onChange={(e) => setParentRoom(e.target.value)}
              style={{
                ...inputStyle,
                borderColor: !categories.filter(c => c.parent_id === null).length ? '#ef4444' : inputStyle.borderColor
              }}
            >
              <option value="">— Выбери комнату —</option>
              {categories
  .filter((c) => !('parent_id' in c) || c.parent_id === null || c.parent_id === 0)
  .map((room) => (
    <option key={room.id} value={room.id}>
      {room.name}
    </option>
  ))}

              {!categories.filter(c => c.parent_id === null).length && (
                <option value="" disabled>Нет доступных комнат</option>
              )}
            </select>
            {!categories.filter(c => c.parent_id === null).length && (
              <div style={{ 
                color: '#ef4444', 
                fontSize: '12px', 
                marginTop: '4px',
                position: 'absolute',
                bottom: '-20px',
                left: '0'
              }}>
                Сначала создайте комнаты
              </div>
            )}
          </div>
        )}

        {catType === 'sub' && parentRoom && (
          <div style={{ position: 'relative', width: '100%', marginTop: '10px' }}>
            <select
              value={parentElement}
              onChange={(e) => setParentElement(e.target.value)}
              style={{
                ...inputStyle,
                borderColor: !categories.filter(c => c.parent_id === Number(parentRoom)).length ? '#ef4444' : inputStyle.borderColor
              }}
            >
              <option value="">— Выбери элемент —</option>
              {categories
                .filter((c) => c.parent_id === Number(parentRoom))
                .map((elem) => (
                  <option key={elem.id} value={elem.id}>
                    {elem.name}
                  </option>
                ))}
              {!categories.filter(c => c.parent_id === Number(parentRoom)).length && (
                <option value="" disabled>Нет доступных элементов</option>
              )}
            </select>
            {!categories.filter(c => c.parent_id === Number(parentRoom)).length && (
              <div style={{ 
                color: '#ef4444', 
                fontSize: '12px', 
                marginTop: '4px',
                position: 'absolute',
                bottom: '-20px',
                left: '0'
              }}>
                Сначала создайте элементы в этой комнате
              </div>
            )}
          </div>
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
