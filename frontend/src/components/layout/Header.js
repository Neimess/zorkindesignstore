import React from 'react';

/**
 * Компонент шапки сайта
 * Отображает логотип, навигационное меню и иконки корзины/пользователя
 */
function Header() {
  return (
    <header
      style={{
        background: '#0f172a',
        padding: '15px 30px',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        borderBottom: '1px solid #334155',
        boxShadow: '0 4px 20px rgba(0,0,0,0.1)',
        position: 'sticky',
        top: 0,
        zIndex: 100,
      }}
    >
      <div style={{ display: 'flex', alignItems: 'center' }}>
        <div
          style={{
            fontSize: '1.8rem',
            fontWeight: 800,
            color: '#f1f5f9',
            display: 'flex',
            alignItems: 'center',
            gap: '10px',
          }}
        >
          <i className="fas fa-home" style={{ color: '#3b82f6' }}></i>
          <span>
            ROYAL<span style={{ color: '#3b82f6' }}>INTERIORS</span>
          </span>
        </div>
      </div>
      <nav>
        <ul
          style={{
            display: 'flex',
            gap: '25px',
            listStyle: 'none',
            margin: 0,
            padding: 0,
          }}
        >
          <li>
            <a
              href="#"
              style={{
                color: '#60a5fa',
                textDecoration: 'none',
                fontWeight: 600,
              }}
            >
              ДИЗАЙН ИНТЕРЬЕРА
            </a>
          </li>
          <li>
            <a href="#" style={{ color: '#f1f5f9', textDecoration: 'none' }}>
              РЕМОНТ
            </a>
          </li>
          <li>
            <a href="#" style={{ color: '#f1f5f9', textDecoration: 'none' }}>
              МАТЕРИАЛЫ
            </a>
          </li>
          <li>
            <a href="#" style={{ color: '#f1f5f9', textDecoration: 'none' }}>
              ПОРТФОЛИО
            </a>
          </li>
          <li>
            <a
              href="#"
              style={{
                color: '#f1f5f9',
                textDecoration: 'none',
                display: 'flex',
                alignItems: 'center',
                gap: '5px',
              }}
            >
              <span>КОНФИГУРАТОР</span>
              <i className="fas fa-tools" style={{ color: '#3b82f6' }}></i>
            </a>
          </li>
        </ul>
      </nav>
      <div style={{ display: 'flex', alignItems: 'center', gap: '20px' }}>
        <div style={{ position: 'relative' }}>
          <i
            className="fas fa-shopping-cart"
            style={{ color: '#f1f5f9', fontSize: '1.2rem' }}
          ></i>
          <span
            style={{
              position: 'absolute',
              top: '-8px',
              right: '-8px',
              background: '#3b82f6',
              color: 'white',
              borderRadius: '50%',
              width: '18px',
              height: '18px',
              fontSize: '0.7rem',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontWeight: 'bold',
            }}
          >
            2
          </span>
        </div>
        <i
          className="fas fa-user"
          style={{ color: '#f1f5f9', fontSize: '1.2rem' }}
        ></i>
      </div>
    </header>
  );
}

export default Header;
