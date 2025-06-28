import React, { useState } from 'react';

export default function StyleAdmin({ products, styles = [], setStyles }) {
  const [name, setName] = useState('');
  const [selectedIds, setSelectedIds] = useState([]);

  const toggleProduct = (id) => {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    );
  };

  const addStyle = () => {
    if (!name.trim() || selectedIds.length === 0) return;
    setStyles([...styles, { id: Date.now(), name, productIds: selectedIds }]);
    setName('');
    setSelectedIds([]);
  };

  const removeStyle = (id) => {
    setStyles(styles.filter(s => s.id !== id));
  };

  const inputStyle = {
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
  };

  const buttonStyle = {
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
  };

  const deleteButtonStyle = {
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
  };

  const checkboxStyle = {
    position: 'relative',
    display: 'inline-flex',
    alignItems: 'center',
    padding: '8px 14px',
    background: 'rgba(30, 41, 59, 0.7)',
    borderRadius: '8px',
    border: '1px solid #334155',
    margin: '5px',
    cursor: 'pointer',
    transition: 'all 0.3s ease',
    fontSize: '0.95rem',
    color: '#cbd5e1'
  };

  const selectedCheckboxStyle = {
    ...checkboxStyle,
    background: 'rgba(59, 130, 246, 0.2)',
    borderColor: '#3b82f6',
    color: '#60a5fa',
    boxShadow: '0 4px 12px rgba(37, 99, 235, 0.15)'
  };

  return (
    <div className="AdminSection" style={{ marginTop: 50 }}>
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
        marginBottom: '30px'
      }}>
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Название стиля"
          style={{ ...inputStyle, marginBottom: 20 }}
        />
        
        <div style={{ marginBottom: 10, color: '#94a3b8', fontSize: '0.95rem' }}>
          <i className="fas fa-info-circle" style={{ marginRight: '8px' }}></i>
          Выберите товары, которые будут включены в этот стиль:
        </div>
        
        <div style={{ 
          display: 'flex', 
          flexWrap: 'wrap', 
          gap: 6, 
          marginBottom: 20,
          maxHeight: '200px',
          overflowY: 'auto',
          padding: '10px',
          background: 'rgba(15, 23, 42, 0.3)',
          borderRadius: '8px'
        }}>
          {products.map((p) => (
            <label 
              key={p.id} 
              style={selectedIds.includes(p.id) ? selectedCheckboxStyle : checkboxStyle}
            >
              <input
                type="checkbox"
                checked={selectedIds.includes(p.id)}
                onChange={() => toggleProduct(p.id)}
                style={{ 
                  marginRight: '8px',
                  accentColor: '#3b82f6',
                  width: '16px',
                  height: '16px'
                }}
              />
              {p.name}
            </label>
          ))}
        </div>
        
        <button 
          onClick={addStyle} 
          style={buttonStyle}
          disabled={!name.trim() || selectedIds.length === 0}
        >
          <i className="fas fa-plus" style={{ marginRight: '8px' }}></i>
          Добавить стиль
        </button>
      </div>

      <div style={{ 
        background: 'rgba(30, 41, 59, 0.5)',
        borderRadius: '12px',
        border: '1px solid #334155',
        overflow: 'hidden'
      }}>
        <div style={{ 
          padding: '15px 20px', 
          borderBottom: '1px solid rgba(51, 65, 85, 0.5)',
          color: '#94a3b8',
          fontSize: '0.95rem',
          fontWeight: 600,
          display: 'flex',
          justifyContent: 'space-between'
        }}>
          <span>Название стиля</span>
          <span>Действия</span>
        </div>
        
        {styles.length === 0 ? (
          <div style={{ 
            padding: '30px 20px', 
            textAlign: 'center', 
            color: '#94a3b8',
            fontSize: '1rem'
          }}>
            <i className="fas fa-info-circle" style={{ fontSize: '24px', marginBottom: '10px', display: 'block', opacity: 0.7 }}></i>
            Стили еще не добавлены
          </div>
        ) : (
          <ul style={{ 
            listStyle: 'none', 
            padding: 0,
            margin: 0
          }}>
            {styles.map((s) => (
              <li key={s.id} style={{ 
                padding: '14px 20px', 
                borderBottom: '1px solid rgba(51, 65, 85, 0.5)',
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                transition: 'all 0.3s ease'
              }}>
                <div>
                  <div style={{ fontWeight: 600, fontSize: '1.1rem', marginBottom: '4px' }}>{s.name}</div>
                  <div style={{ color: '#94a3b8', fontSize: '0.9rem' }}>
                    <span style={{ 
                      display: 'inline-block',
                      width: '8px',
                      height: '8px',
                      background: '#60a5fa',
                      borderRadius: '50%',
                      marginRight: '8px'
                    }}></span>
                    {s.productIds.length} товаров
                  </div>
                </div>
                <button 
                  onClick={() => removeStyle(s.id)} 
                  style={deleteButtonStyle}
                >
                  <i className="fas fa-trash-alt" style={{ marginRight: '6px' }}></i>
                  Удалить
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
