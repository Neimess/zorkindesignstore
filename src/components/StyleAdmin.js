import React, { useState } from 'react';
import { presetAPI } from '../services/api';

/**
 * Компонент для управления стилями в админ-панели
 * 
 * @param {Object} props - Свойства компонента
 * @param {Array} props.products - Список товаров
 * @param {Array} props.styles - Список стилей
 * @param {Function} props.setStyles - Функция для обновления списка стилей
 */
function StyleAdmin({ products, styles, setStyles }) {
  // Состояние для нового стиля
  const [styleName, setStyleName] = useState('');
  const [selectedProductIds, setSelectedProductIds] = useState([]);
  
  /**
   * Обработчик выбора/отмены выбора товара
   * @param {number} productId - ID товара
   */
  const toggleProductSelection = (productId) => {
    setSelectedProductIds(prev => 
      prev.includes(productId)
        ? prev.filter(id => id !== productId)
        : [...prev, productId]
    );
  };
  
  /**
   * Добавление нового стиля
   */
  const addStyle = async () => {
    if (!styleName.trim() || selectedProductIds.length === 0) return;
    
    try {
      // Создаем новый стиль
      const newStyle = {
        preset_id: Date.now(),  // Временный ID, в реальном приложении будет присвоен сервером
        name: styleName,
        productIds: selectedProductIds
      };
      
      // Обновляем список стилей
      setStyles([...styles, newStyle]);
      
      // Сбрасываем форму
      setStyleName('');
      setSelectedProductIds([]);
      
      // Здесь можно добавить вызов API для сохранения стиля на сервере
      // await presetAPI.create(newStyle);
    } catch (error) {
      console.error('Ошибка при создании стиля:', error);
    }
  };
  
  /**
   * Удаление стиля
   * @param {number} styleId - ID стиля
   */
  const removeStyle = async (styleId) => {
    try {
      // Удаляем стиль из списка
      setStyles(styles.filter(style => style.preset_id !== styleId));
      
      // Здесь можно добавить вызов API для удаления стиля на сервере
      // await presetAPI.delete(styleId);
    } catch (error) {
      console.error('Ошибка при удалении стиля:', error);
    }
  };
  
  // Стили для элементов интерфейса
  const uiStyles = {
    inputStyle: {
      padding: '12px 16px',
      borderRadius: '10px',
      border: '1px solid #334155',
      background: 'rgba(15, 23, 42, 0.6)',
      color: '#f1f5f9',
      fontSize: '1rem',
      width: '100%',
      transition: 'all 0.3s ease',
      boxShadow: '0 4px 10px rgba(0,0,0,0.1)',
      outline: 'none'
    },
    buttonStyle: {
      background: 'linear-gradient(135deg, #3b82f6, #2563eb)',
      color: '#fff',
      border: 'none',
      borderRadius: '10px',
      padding: '12px 20px',
      fontSize: '1rem',
      fontWeight: 600,
      cursor: 'pointer',
      transition: 'all 0.3s ease',
      boxShadow: '0 4px 12px rgba(37, 99, 235, 0.3)',
      textTransform: 'uppercase',
      letterSpacing: '0.5px'
    },
    deleteButtonStyle: {
      background: 'rgba(185, 28, 28, 0.1)',
      color: '#f87171',
      border: '1px solid rgba(185, 28, 28, 0.3)',
      borderRadius: '8px',
      padding: '8px 16px',
      fontSize: '0.9rem',
      fontWeight: 500,
      cursor: 'pointer',
      transition: 'all 0.3s ease',
      marginLeft: '10px'
    },
    checkboxStyle: {
      display: 'flex',
      alignItems: 'center',
      padding: '10px 15px',
      borderRadius: '8px',
      border: '1px solid #334155',
      background: 'rgba(15, 23, 42, 0.6)',
      marginBottom: '8px',
      cursor: 'pointer',
      transition: 'all 0.3s ease'
    },
    checkboxCheckedStyle: {
      background: 'rgba(59, 130, 246, 0.1)',
      borderColor: '#3b82f6'
    },
    checkboxInput: {
      marginRight: '10px',
      cursor: 'pointer'
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
        Стили ремонта
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
        background: 'rgba(30, 41, 59, 0.5)',
        padding: '20px',
        borderRadius: '12px',
        border: '1px solid #334155',
        marginBottom: '24px'
      }}>
        <input 
          value={styleName} 
          onChange={e => setStyleName(e.target.value)} 
          placeholder="Название стиля" 
          style={uiStyles.inputStyle} 
        />
        
        <div style={{ marginTop: '20px', marginBottom: '15px' }}>
          <div style={{ marginBottom: '10px', color: '#94a3b8' }}>Выберите товары для стиля:</div>
          <div style={{ maxHeight: '300px', overflowY: 'auto', padding: '5px' }}>
            {products.map(product => (
              // console.log('products:', products), // массив
              <div 
                key={product.product_id ?? product.id}
                style={{
                  ...uiStyles.checkboxStyle,
                  ...(selectedProductIds.includes(product.product_id) ? uiStyles.checkboxCheckedStyle : {})
                }}
                onClick={() => toggleProductSelection(product.product_id)}
              >
                <input 
                  type="checkbox" 
                  checked={selectedProductIds.includes(product.product_id)}
                      onChange={() => toggleProductSelection(product.product_id)} 
                  style={uiStyles.checkboxInput} 
                />
                <div>
                  <div style={{ fontWeight: 500, color: '#f1f5f9' }}>{product.name}</div>
                  <div style={{ fontSize: '0.9rem', color: '#94a3b8', marginTop: '3px' }}>
                    {product.price} ₽
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
        
        <button 
          onClick={addStyle} 
          style={{
            ...uiStyles.buttonStyle,
            marginTop: '15px',
            width: '100%'
          }}
          disabled={!styleName.trim() || selectedProductIds.length === 0}
        >
          <i className="fas fa-plus" style={{ marginRight: '8px' }}></i>
          Добавить стиль
        </button>
      </div>
      
      <div style={{ 
        background: 'rgba(30, 41, 59, 0.5)',
        borderRadius: '12px',
        overflow: 'hidden',
        border: '1px solid #334155'
      }}>
        <div style={{ padding: '15px 20px', borderBottom: '1px solid #334155', color: '#94a3b8' }}>
          Существующие стили
        </div>
        <ul style={{ 
          listStyle: 'none', 
          padding: 0,
          margin: 0
        }}>
          {styles.map(style => (
            <li key={style.preset_id} style={{ 
              padding: '14px 20px', 
              borderBottom: '1px solid rgba(51, 65, 85, 0.5)',
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center'
            }}>
              <div>
                <div style={{ fontWeight: 500, color: '#f1f5f9', fontSize: '1.1rem' }}>{style.name}</div>
                <div style={{ color: '#94a3b8', fontSize: '0.9rem', marginTop: '5px' }}>
                  Товаров: {style.productIds?.length || 0}
                </div>
              </div>
              <button 
                onClick={() => removeStyle(style.preset_id)} 
                style={uiStyles.deleteButtonStyle}
              >
                <i className="fas fa-trash-alt" style={{ marginRight: '6px' }}></i>
                Удалить
              </button>
            </li>
          ))}
          {styles.length === 0 && (
            <li style={{ 
              padding: '20px', 
              textAlign: 'center',
              color: '#94a3b8'
            }}>
              Нет добавленных стилей
            </li>
          )}
        </ul>
      </div>
    </div>
  );
}

export default StyleAdmin;
