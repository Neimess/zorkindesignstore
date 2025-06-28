import React, { useState } from 'react';
import { productAPI } from '../../services/api';

/**
 * Компонент для управления товарами в админ-панели
 * Позволяет добавлять и удалять товары
 * 
 * @param {Object} props - Свойства компонента
 * @param {Array} props.categories - Список категорий
 * @param {Array} props.products - Список товаров
 * @param {Function} props.setProducts - Функция для обновления списка товаров
 * @param {Function} props.getAdminToken - Функция для получения токена администратора
 * @param {Function} props.showMessage - Функция для отображения сообщений
 * @param {Object} props.styles - Объект со стилями для элементов интерфейса
 */
function ProductManager({ categories, products, setProducts, getAdminToken, showMessage, styles }) {
  const [prod, setProd] = useState({ 
    name: '', 
    price: '', 
    categoryId: categories[0]?.id || 1, 
    description: '', 
    image_url: '', 
    attributes: '' 
  });
  
  const { inputStyle, buttonStyle, deleteButtonStyle } = styles;

  /**
   * Добавляет новый товар
   */
  const addProduct = async () => {
    if (!prod.name.trim() || !prod.price || !prod.categoryId) return;
    
    try {
      const token = await getAdminToken();
      if (!token) return;
      
      // Подготавливаем данные для API
      const productData = {
        name: prod.name,
        price: Number(prod.price),
        category_id: Number(prod.categoryId),
        description: prod.description,
        image_url: prod.image_url,
attributes: prod.attributes
  ? prod.attributes.split(';').map(attr => {
      const [id, value] = attr.split(':').map(s => s.trim());
      const attribute_id = Number(id);
      if (!attribute_id || !value) return null;
      return {
        attribute_id,
        value
      };
    }).filter(Boolean)
  : []

      };
      
      const response = await productAPI.create(productData, token);
      console.log('🚀 productData ->', productData);

      // Обновляем локальный список товаров
 const newProduct = {
   id: (response && response.id) ? response.id : Date.now(),
        name: prod.name,
        price: Number(prod.price),
        categoryId: Number(prod.categoryId),
        description: prod.description,
        image_url: prod.image_url,
        attributes: prod.attributes
          ? Object.fromEntries(prod.attributes.split(';').map((a) => a.split(':').map((s) => s.trim())))
          : {},
      };
      
      setProducts([...products, newProduct]);
      setProd({ name: '', price: '', categoryId: categories[0]?.id || 1, description: '', image_url: '', attributes: '' });
      showMessage('Товар успешно добавлен');
    } catch (error) {
      console.error('Ошибка создания товара:', error);
      showMessage('Ошибка при создании товара', true);
    }
  };
  
  /**
   * Удаляет товар по ID
   * @param {number} id - ID товара для удаления
   */
  const removeProduct = async (id) => {
    try {
      const token = await getAdminToken();
      if (!token) return;
      
      await productAPI.delete(id, token);
      
      // Обновляем локальные данные
      setProducts(products.filter((p) => p.name !== id));
      showMessage('Товар успешно удален');
    } catch (error) {
      console.error('Ошибка удаления товара:', error);
      showMessage('Ошибка при удалении товара', true);
    }
  };

  return (
    <div className="AdminSection" style={{ marginTop: 40 }}>
      <h2 style={{ 
        fontSize: '1.5rem', 
        color: '#f8fafc', 
        marginBottom: '20px', 
        position: 'relative',
        paddingBottom: '10px'
      }}>
        Товары
        <span style={{ 
          position: 'absolute', 
          bottom: 0, 
          left: 0, 
          width: '60px', 
          height: '3px', 
          background: 'linear-gradient(90deg, #3b82f6, #60a5fa)', 
          borderRadius: '2px' 
        }}></span>
      </h2>
      
      <div style={{ 
        display: 'grid', 
        gridTemplateColumns: '1fr 1fr', 
        gap: 16, 
        marginBottom: 24,
        background: 'rgba(30, 41, 59, 0.5)',
        padding: '20px',
        borderRadius: '12px',
        border: '1px solid #334155'
      }}>
        <input 
          value={prod.name} 
          onChange={e => setProd({ ...prod, name: e.target.value })} 
          placeholder="Название" 
          style={inputStyle} 
        />
        <input 
          value={prod.price} 
          onChange={e => setProd({ ...prod, price: e.target.value })} 
          placeholder="Цена" 
          type="number" 
          style={inputStyle} 
        />
        <select 
          value={prod.categoryId} 
          onChange={e => setProd({ ...prod, categoryId: e.target.value })} 
          style={{
            ...inputStyle,
            appearance: 'none',
            backgroundImage: 'url("data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%2712%27 height=%278%27 viewBox=%270 0 12 8%27%3E%3Cpath fill=%27%2360a5fa%27 d=%27M10.6.6L6 5.2 1.4.6.6 1.4 6 6.8l5.4-5.4z%27/%3E%3C/svg%3E")',
            backgroundRepeat: 'no-repeat',
            backgroundPosition: 'right 16px center',
            paddingRight: '40px'
          }}
        >
          {categories.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
        </select>
        <input 
          value={prod.image_url} 
          onChange={e => setProd({ ...prod, image_url: e.target.value })} 
          placeholder="URL картинки" 
          style={inputStyle} 
        />
        <input 
          value={prod.description} 
          onChange={e => setProd({ ...prod, description: e.target.value })} 
          placeholder="Описание" 
          style={inputStyle} 
        />
        <input 
          value={prod.attributes} 
          onChange={e => setProd({ ...prod, attributes: e.target.value })} 
          placeholder="Атрибуты (пример: 2:1.25; 3:Матовый)"
 
          style={inputStyle} 
        />
        <button 
          onClick={addProduct} 
          style={{
            ...buttonStyle,
            gridColumn: '1 / -1',
            marginTop: '10px',
            padding: '14px'
          }}
        >
          <i className="fas fa-plus" style={{ marginRight: '8px' }}></i>
          Добавить товар
        </button>
      </div>
      
      <ul style={{ 
        marginBottom: 30, 
        listStyle: 'none', 
        padding: 0,
        background: 'rgba(30, 41, 59, 0.5)',
        borderRadius: '12px',
        overflow: 'hidden',
        border: '1px solid #334155',
        maxHeight: '400px',
        overflowY: 'auto'
      }}>
        {products.map((p) => (
          // console.log('styles:', p.name),
          <li key={p.name} style={{ 
            padding: '14px 20px', 
            borderBottom: '1px solid rgba(51, 65, 85, 0.5)',
            display: 'flex',
            alignItems: 'center',
            gap: 16,
            transition: 'all 0.3s ease'
          }}>
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
                boxShadow: '0 4px 10px rgba(0,0,0,0.2)'
              }} 
            />
            <div style={{ flex: 1 }}>
              <div style={{ fontWeight: 600, fontSize: '1.1rem', marginBottom: '4px' }}>{p.name}</div>
              <div style={{ color: '#94a3b8', fontSize: '0.9rem' }}>
                {categories.find(c => c.id === p.categoryId)?.name || '—'}
              </div>
            </div>
            <div style={{ 
              fontWeight: 700, 
              fontSize: '1.2rem', 
              color: '#38bdf8',
              display: 'flex',
              alignItems: 'center',
              gap: '6px'
            }}>
              <span style={{ 
                width: '8px', 
                height: '8px', 
                background: '#38bdf8', 
                borderRadius: '50%',
                display: 'inline-block'
              }}></span>
              {p.price} ₽
            </div>
            <button 
              onClick={() => removeProduct(p.id)} 
              style={deleteButtonStyle}
            >
              <i className="fas fa-trash-alt" style={{ marginRight: '6px' }}></i>
              Удалить
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default ProductManager;