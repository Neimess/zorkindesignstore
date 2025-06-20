import React from 'react';

export default function StyleSelector({ styles, onSelect }) {
  return (
    <div className="Styles" style={{ marginBottom: 24 }}>
      <h2>Стили ремонта</h2>
      <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap' }}>
        {styles.map((s) => (
          <button
            key={s.id}
            onClick={() => onSelect(s)}
            style={{
              padding: '10px 16px',
              borderRadius: 8,
              background: '#fef9c3',
              border: '1px solid #facc15',
              color: '#854d0e',
              cursor: 'pointer',
            }}
          >
            {s.name}
          </button>
        ))}
      </div>
    </div>
  );
}
