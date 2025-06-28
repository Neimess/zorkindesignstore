import React from 'react';

export default function StyleSelector({ styles, onSelect }) {
  return (
    <div className="Styles" style={{ 
      marginBottom: 60,
      padding: '30px 20px',
      background: 'linear-gradient(135deg, rgba(15, 23, 42, 0.8), rgba(30, 41, 59, 0.9))',
      borderRadius: 16,
      boxShadow: '0 10px 25px rgba(0,0,0,0.2)',
      border: '1px solid rgba(71, 85, 105, 0.3)'
    }}>
      <h2 style={{
        fontSize: '1.8rem',
        fontWeight: 700,
        color: '#f8fafc',
        textAlign: 'center',
        marginBottom: 30,
        position: 'relative',
        paddingBottom: 15,
        fontFamily: '"Montserrat", sans-serif'
      }}>
        Выберите стиль интерьера
        <div style={{
          position: 'absolute',
          bottom: 0,
          left: '50%',
          transform: 'translateX(-50%)',
          width: 80,
          height: 3,
          background: 'linear-gradient(90deg, #60a5fa, #3b82f6)',
          borderRadius: 2
        }}></div>
      </h2>
      <p style={{
        textAlign: 'center',
        color: '#cbd5e1',
        marginBottom: 25,
        fontSize: '1.1rem',
        maxWidth: '700px',
        margin: '0 auto 30px'
      }}>
        Выберите стиль, который лучше всего подходит для вашего помещения. Каждый стиль предлагает уникальное сочетание материалов и элементов дизайна.
      </p>
      <div style={{ 
        display: 'flex', 
        gap: 20, 
        flexWrap: 'wrap',
        justifyContent: 'center' 
      }}>
        {styles.map((s) => (
          <button
            key={s.id}
            onClick={() => onSelect(s)}
            style={{
              padding: '18px 28px',
              borderRadius: 14,
              background: 'linear-gradient(135deg, #0f172a, #1e293b)',
              border: '1px solid #334155',
              color: '#f8fafc',
              cursor: 'pointer',
              fontSize: '1.2rem',
              fontWeight: 600,
              transition: 'all 0.3s ease',
              boxShadow: '0 4px 15px rgba(0,0,0,0.15)',
              position: 'relative',
              overflow: 'hidden',
              width: '220px',
              height: '120px',
              display: 'flex',
              flexDirection: 'column',
              justifyContent: 'center',
              alignItems: 'center',
              textAlign: 'center'
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.transform = 'translateY(-5px) scale(1.03)';
              e.currentTarget.style.boxShadow = '0 12px 25px rgba(0,0,0,0.2)';
              e.currentTarget.style.borderColor = '#60a5fa';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.transform = 'translateY(0) scale(1)';
              e.currentTarget.style.boxShadow = '0 4px 15px rgba(0,0,0,0.15)';
              e.currentTarget.style.borderColor = '#334155';
            }}
          >
            <span style={{ 
              position: 'relative', 
              zIndex: 2,
              fontSize: '1.3rem',
              marginBottom: '8px',
              fontFamily: '"Montserrat", sans-serif'
            }}>{s.name}</span>
            <span style={{
              fontSize: '0.85rem',
              color: '#94a3b8',
              zIndex: 2,
              fontWeight: 400
            }}>
              {s.id === 1 && 'Светлые тона, натуральные материалы'}
              {s.id === 2 && 'Индустриальный стиль, открытые коммуникации'}
              {s.id === 3 && 'Простота форм, функциональность'}
              {s.id === 4 && 'Изысканность, симметрия, роскошь'}
              {s.id === 5 && 'Инновационные материалы, чистые линии'}
              {s.id === 6 && 'Пастельные тона, винтажные элементы'}
            </span>
            <div style={{
              position: 'absolute',
              top: 0,
              left: 0,
              width: '6px',
              height: '100%',
              background: 'linear-gradient(135deg, #3b82f6, #2563eb)',
              zIndex: 1
            }}></div>
            <div style={{
              position: 'absolute',
              bottom: 0,
              right: 0,
              width: '40%',
              height: '4px',
              background: 'linear-gradient(90deg, transparent, #3b82f6)',
              zIndex: 1,
              opacity: 0.7
            }}></div>
          </button>
        ))}
      </div>
    </div>
  );
}
