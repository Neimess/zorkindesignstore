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

  return (
    <div style={{ marginTop: 40 }}>
      <h2>Стили ремонта</h2>
      <input
        value={name}
        onChange={(e) => setName(e.target.value)}
        placeholder="Название стиля"
        style={{ marginBottom: 12, padding: 8, width: '100%' }}
      />
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: 10, marginBottom: 16 }}>
        {products.map((p) => (
          <label key={p.id} style={{ fontSize: 14 }}>
            <input
              type="checkbox"
              checked={selectedIds.includes(p.id)}
              onChange={() => toggleProduct(p.id)}
            />{' '}
            {p.name}
          </label>
        ))}
      </div>
      <button onClick={addStyle}>Добавить стиль</button>
      <ul style={{ marginTop: 20 }}>
        {styles.map((s) => (
          <li key={s.id}>{s.name} ({s.productIds.length} товаров)</li>
        ))}
      </ul>
    </div>
  );
}
