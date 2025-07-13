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
  const { inputStyle, buttonStyle, deleteButtonStyle } = styles;

  const addCategory = async () => {
    if (!catName.trim()) return;
    try {
      const token = await getAdminToken();
      if (!token) return;
      const response = await categoryAPI.create({ name: catName }, token);
      const newCategory = {
        id: response.id,
        name: response.name,
      };
      setCategories([...categories, newCategory]);
      setCatName('');
      showMessage('Категория успешно добавлена');
    } catch (error) {
      console.error('Ошибка создания категории:', error);
      showMessage('Ошибка при создании категории', true);
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
        <input
          value={catName}
          onChange={(e) => setCatName(e.target.value)}
          placeholder="Новая категория"
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
