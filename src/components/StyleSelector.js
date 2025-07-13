import React, { useState } from 'react';

/**
 * Компонент для выбора стиля интерьера на главной странице
 *
 * @param {Object} props - Свойства компонента
 * @param {Array} props.styles - Список доступных стилей
 * @param {Function} props.onSelect - Функция обработки выбора стиля
 */
function StyleSelector({ styles, onSelect }) {
  // Описания стилей для отображения
  const styleDescriptions = {
    1: 'Светлые тона, натуральные материалы, функциональность и минимализм',
    2: 'Индустриальный характер, открытые коммуникации, кирпич и металл',
    3: 'Роскошь, изысканные материалы, классические элементы декора',
    4: 'Яркие акценты, геометрические узоры, смелые цветовые решения',
    5: 'Природные материалы, спокойные тона, экологичность и уют',
    6: 'Плавные линии, яркие цвета, необычные формы и текстуры',
  };
  const [expandedStyleId, setExpandedStyleId] = useState(null);
  // Стили для элементов интерфейса
  const uiStyles = {
    container: {
      marginTop: 60,
      marginBottom: 80,
      padding: '0 20px',
    },
    title: {
      fontSize: '2rem',
      fontWeight: 700,
      color: '#f8fafc',
      marginBottom: '15px',
      textAlign: 'center',
    },
    description: {
      fontSize: '1.1rem',
      color: '#94a3b8',
      maxWidth: '700px',
      margin: '0 auto 40px',
      textAlign: 'center',
      lineHeight: 1.6,
    },
    stylesGrid: {
      display: 'grid',
      gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
      gap: '25px',
      maxWidth: '1200px',
      margin: '0 auto',
    },
    styleButton: {
      position: 'relative',
      padding: '30px 25px',
      borderRadius: '16px',
      border: '1px solid #334155',
      background: 'rgba(30, 41, 59, 0.6)',
      cursor: 'pointer',
      transition: 'all 0.3s ease',
      overflow: 'hidden',
      height: '100%',
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'space-between',
    },
    styleButtonHover: {
      transform: 'translateY(-5px)',
      boxShadow: '0 15px 30px rgba(0, 0, 0, 0.2)',
      borderColor: '#60a5fa',
    },
    styleName: {
      fontSize: '1.4rem',
      fontWeight: 600,
      color: '#f1f5f9',
      marginBottom: '15px',
      position: 'relative',
      zIndex: 1,
    },
    styleDesc: {
      fontSize: '1rem',
      color: '#cbd5e1',
      lineHeight: 1.6,
      position: 'relative',
      zIndex: 1,
    },
    gradientOverlay: {
      position: 'absolute',
      top: 0,
      left: 0,
      right: 0,
      bottom: 0,
      opacity: 0.15,
      transition: 'opacity 0.3s ease',
      zIndex: 0,
    },
    decorBorder: {
      position: 'absolute',
      width: '60px',
      height: '60px',
      borderRadius: '30px',
      transition: 'all 0.3s ease',
      opacity: 0.6,
    },
  };
  // console.log('styles:', styles);
  return (
    <section style={uiStyles.container}>
      <h2 style={uiStyles.title}>Популярные стили интерьера</h2>
      <p style={uiStyles.description}>
        Выберите один из готовых стилей, чтобы автоматически добавить
        рекомендуемые товары для ремонта. Каждый стиль содержит тщательно
        подобранные материалы и элементы, которые хорошо сочетаются между собой.
      </p>

      <div style={uiStyles.stylesGrid}>
        {styles.map((style) => (
          <div
            key={style.preset_id}
            onClick={() => {
              onSelect(style);
              setExpandedStyleId((prev) =>
                prev === style.preset_id ? null : style.preset_id,
              );
            }}
            style={uiStyles.styleButton}
            onMouseEnter={(e) => {
              // Добавляем эффекты при наведении
              Object.assign(e.currentTarget.style, uiStyles.styleButtonHover);
              // Увеличиваем яркость градиента
              e.currentTarget.querySelector('.gradient-overlay').style.opacity =
                '0.3';
            }}
            onMouseLeave={(e) => {
              // Удаляем эффекты при уходе курсора
              e.currentTarget.style.transform = '';
              e.currentTarget.style.boxShadow = '';
              e.currentTarget.style.borderColor = '#334155';
              // Возвращаем исходную яркость градиента
              e.currentTarget.querySelector('.gradient-overlay').style.opacity =
                '0.15';
            }}
          >
            {/* Декоративные элементы */}
            <div
              className="gradient-overlay"
              style={{
                ...uiStyles.gradientOverlay,
                background: `linear-gradient(135deg, #3b82f6, #8b5cf6)`,
              }}
            ></div>
            <div
              style={{
                ...uiStyles.decorBorder,
                top: '-20px',
                right: '-20px',
                border: '2px solid rgba(59, 130, 246, 0.3)',
              }}
            ></div>
            <div
              style={{
                ...uiStyles.decorBorder,
                bottom: '-30px',
                left: '-30px',
                border: '2px solid rgba(139, 92, 246, 0.3)',
                width: '90px',
                height: '90px',
                borderRadius: '45px',
              }}
            ></div>

            {/* Содержимое */}
            <div>
              <h3 style={uiStyles.styleName}>{style.name}</h3>
              <p style={uiStyles.styleDesc}>
                {styleDescriptions[style.preset_id] ||
                  'Уникальный стиль интерьера с подобранными товарами'}
              </p>
            </div>

            <div
              style={{
                marginTop: '20px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                position: 'relative',
                zIndex: 1,
              }}
            >
              {expandedStyleId === style.preset_id && (
                <div
                  style={{
                    marginTop: '15px',
                    paddingTop: '10px',
                    borderTop: '1px solid #334155',
                  }}
                >
                  {(style.items || []).map((item, index) => (
                    <div
                      key={index}
                      style={{ marginBottom: '10px', color: '#cbd5e1' }}
                    >
                      • {item.product?.name} —{' '}
                      {(item.product?.price ?? 0).toLocaleString()} ₽
                    </div>
                  ))}
                </div>
              )}
              <span
                style={{
                  fontSize: '0.9rem',
                  color: '#94a3b8',
                }}
              >
                {style.items?.length || 0} товаров
              </span>
              <span
                style={{
                  fontSize: '0.9rem',
                  color: '#60a5fa',
                  fontWeight: 500,
                }}
              >
                Выбрать
              </span>
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}

export default StyleSelector;
