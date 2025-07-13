import React from 'react';

/**
 * Компонент подвала сайта
 * Отображает информацию о компании, ссылки на разделы и контактную информацию
 */
function Footer() {
  return (
    <footer
      style={{
        background: '#0f172a',
        padding: '40px 30px',
        marginTop: '60px',
        borderTop: '1px solid #334155',
        color: '#94a3b8',
      }}
    >
      <div
        style={{
          maxWidth: '1200px',
          margin: '0 auto',
          display: 'flex',
          justifyContent: 'space-between',
          flexWrap: 'wrap',
          gap: '30px',
        }}
      >
        <div style={{ flex: '1', minWidth: '250px' }}>
          <div
            style={{
              fontSize: '1.5rem',
              fontWeight: '800',
              color: '#f1f5f9',
              display: 'flex',
              alignItems: 'center',
              gap: '10px',
              marginBottom: '15px',
            }}
          >
            <i className="fas fa-home" style={{ color: '#3b82f6' }}></i>
            <span>
              ROYAL<span style={{ color: '#3b82f6' }}>INTERIORS</span>
            </span>
          </div>
          <p style={{ lineHeight: '1.6', marginBottom: '20px' }}>
            Профессиональный дизайн и ремонт помещений. Индивидуальный подход к
            каждому клиенту. Гарантия качества на все работы.
          </p>
          <div style={{ display: 'flex', gap: '15px' }}>
            <i
              className="fab fa-vk"
              style={{ color: '#3b82f6', fontSize: '1.3rem' }}
            ></i>
            <i
              className="fab fa-telegram"
              style={{ color: '#3b82f6', fontSize: '1.3rem' }}
            ></i>
            <i
              className="fab fa-youtube"
              style={{ color: '#3b82f6', fontSize: '1.3rem' }}
            ></i>
          </div>
        </div>

        <div style={{ flex: '1', minWidth: '250px' }}>
          <h3
            style={{
              color: '#f1f5f9',
              marginBottom: '15px',
              fontSize: '1.1rem',
            }}
          >
            Каталог
          </h3>
          <ul
            style={{
              listStyle: 'none',
              padding: '0',
              display: 'flex',
              flexDirection: 'column',
              gap: '10px',
            }}
          >
            <li>
              <a
                href="#"
                style={{
                  color: '#94a3b8',
                  textDecoration: 'none',
                  transition: 'color 0.3s',
                }}
              >
                Дизайн интерьера
              </a>
            </li>
            <li>
              <a
                href="#"
                style={{
                  color: '#94a3b8',
                  textDecoration: 'none',
                  transition: 'color 0.3s',
                }}
              >
                Ремонт квартир
              </a>
            </li>
            <li>
              <a
                href="#"
                style={{
                  color: '#94a3b8',
                  textDecoration: 'none',
                  transition: 'color 0.3s',
                }}
              >
                Отделочные материалы
              </a>
            </li>
            <li>
              <a
                href="#"
                style={{
                  color: '#94a3b8',
                  textDecoration: 'none',
                  transition: 'color 0.3s',
                }}
              >
                Мебель на заказ
              </a>
            </li>
            <li>
              <a
                href="#"
                style={{
                  color: '#94a3b8',
                  textDecoration: 'none',
                  transition: 'color 0.3s',
                }}
              >
                Конфигуратор помещений
              </a>
            </li>
          </ul>
        </div>

        <div style={{ flex: '1', minWidth: '250px' }}>
          <h3
            style={{
              color: '#f1f5f9',
              marginBottom: '15px',
              fontSize: '1.1rem',
            }}
          >
            Контакты
          </h3>
          <ul
            style={{
              listStyle: 'none',
              padding: '0',
              display: 'flex',
              flexDirection: 'column',
              gap: '10px',
            }}
          >
            <li style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              <i
                className="fas fa-map-marker-alt"
                style={{ color: '#3b82f6', width: '20px' }}
              ></i>
              <span>г. Москва, ул. Дизайнерская, 42</span>
            </li>
            <li style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              <i
                className="fas fa-phone"
                style={{ color: '#3b82f6', width: '20px' }}
              ></i>
              <span>+7 (999) 123-45-67</span>
            </li>
            <li style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              <i
                className="fas fa-envelope"
                style={{ color: '#3b82f6', width: '20px' }}
              ></i>
              <span>info@royalinteriors.ru</span>
            </li>
            <li style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              <i
                className="fas fa-clock"
                style={{ color: '#3b82f6', width: '20px' }}
              ></i>
              <span>Пн-Пт: 9:00-20:00, Сб-Вс: 10:00-18:00</span>
            </li>
          </ul>
        </div>
      </div>
      <div
        style={{
          borderTop: '1px solid #334155',
          marginTop: '30px',
          paddingTop: '20px',
          textAlign: 'center',
        }}
      >
        © 2023 ROYAL INTERIORS. Все права защищены.
      </div>
    </footer>
  );
}

export default Footer;
