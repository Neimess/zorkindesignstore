import React, { useState, useEffect } from 'react';
import Header from '../components/layout/Header';
import Footer from '../components/layout/Footer';
import StyleSelector from '../components/StyleSelector';

/**
 * Компонент главной страницы с конфигуратором
 * 
 * @param {Object} props - Свойства компонента
 * @param {Array} props.categories - Список категорий
 * @param {Array} props.products - Список товаров
 * @param {Array} props.styles - Список стилей
 */
function MainPage({ categories, products, styles }) {
  // Состояние для выбранных категорий и товаров
  const [selectedCategories, setSelectedCategories] = useState([]);
  const [selectedProducts, setSelectedProducts] = useState([]);
  const [selectedStyle, setSelectedStyle] = useState(null);
  const [showStyleModal, setShowStyleModal] = useState(false);
  
    useEffect(() => {
    document.body.style.overflow = showStyleModal ? 'hidden' : '';
    return () => { document.body.style.overflow = ''; };
  }, [showStyleModal]);

  /**
   * Обработчик выбора категории
   * @param {number} category_id - ID категории
   */
 const handleCategorySelect = (category_id) => {
    if (selectedCategories.includes(category_id)) {
      setSelectedCategories(selectedCategories.filter(id => id !== category_id));
    } else {
      setSelectedCategories([...selectedCategories, category_id]);
    }
  };
  
  /**
   * Обработчик выбора товара
   * @param {Object} product - Объект товара
   */
  const handleProductSelect = (product) => {
    if (!selectedProducts.find(p => p.id === product.product_id)) {
      setSelectedProducts([...selectedProducts, product]);
    }
  };
  
  /**
   * Обработчик удаления товара из выбранных
   * @param {number} productId - ID товара
   */
  const handleProductDeselect = (productId) => {
    setSelectedProducts(prev => {
      const index = prev.findIndex(p => p.id === productId);
      if (index !== -1) {
        const newProducts = [...prev];
        newProducts.splice(index, 1);
        return newProducts;
      }
      return prev;
    });
  };

  const handleStyleClick = (style) => {
    setSelectedStyle(style);
    setShowStyleModal(true);
  };
  const confirmStyle = () => {
    const newProducts = (selectedStyle.items || [])
      .map(item => item.product)
      .filter(product =>
        product && !selectedProducts.some(p => p.product_id === product.id)
      )
      .map(product => ({
        ...product,
        product_id: product.id
      }));

    setSelectedProducts(prev => [...prev, ...newProducts]);
    setShowStyleModal(false);
    setSelectedStyle(null);
  };
  /**
   * Обработчик выбора стиля
   * @param {Object} style - Объект стиля
   */
  const handleStyleSelect = (style) => {
    const newProducts = (style.items || [])
      .map(item => item.product)
      .filter(product =>
        product && !selectedProducts.some(p => p.product_id === product.id)
      )
      .map(product => ({
        ...product,
        product_id: product.id
      }));

    setSelectedProducts(prev => [...prev, ...newProducts]);
  };


  
  // Вычисляем общую стоимость выбранных товаров
  const totalPrice = selectedProducts.reduce((sum, product) => sum + product.price, 0);

  // Стили для элементов интерфейса
  const styles_ui = {
    categoryButton: (isSelected) => ({
      background: isSelected ? 'linear-gradient(135deg, #3b82f6, #2563eb)' : 'rgba(15, 23, 42, 0.6)',
      color: '#f1f5f9',
      border: isSelected ? 'none' : '1px solid #334155',
      borderRadius: '10px',
      padding: '12px 20px',
      margin: '0 10px 10px 0',
      fontSize: '1rem',
      fontWeight: isSelected ? 600 : 400,
      cursor: 'pointer',
      transition: 'all 0.3s ease',
      boxShadow: isSelected ? '0 4px 12px rgba(37, 99, 235, 0.3)' : 'none'
    }),
    productCard: {
      background: 'rgba(15, 23, 42, 0.6)',
      borderRadius: '12px',
      padding: '20px',
      margin: '15px 0',
      boxShadow: '0 4px 10px rgba(0, 0, 0, 0.1)',
      border: '1px solid #334155',
      transition: 'all 0.3s ease',
      cursor: 'pointer',
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'space-between'
    },
    productCardHover: {
      transform: 'translateY(-5px)',
      boxShadow: '0 10px 20px rgba(0, 0, 0, 0.2)',
      borderColor: '#3b82f6'
    },
    addButton: {
      background: 'linear-gradient(135deg, #3b82f6, #2563eb)',
      color: '#fff',
      border: 'none',
      borderRadius: '8px',
      padding: '10px 16px',
      fontSize: '0.9rem',
      fontWeight: 600,
      cursor: 'pointer',
      transition: 'all 0.3s ease',
      marginTop: '15px',
      boxShadow: '0 4px 12px rgba(37, 99, 235, 0.3)',
      textTransform: 'uppercase',
      letterSpacing: '0.5px'
    },
    removeButton: {
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
    selectedProductsPanel: {
      background: 'rgba(15, 23, 42, 0.8)',
      borderRadius: '16px',
      padding: '25px',
      marginTop: '30px',
      boxShadow: '0 10px 25px rgba(0, 0, 0, 0.2)',
      border: '1px solid #334155'
    },
    totalPrice: {
      fontSize: '1.5rem',
      fontWeight: 700,
      color: '#f1f5f9',
      marginTop: '20px',
      padding: '15px 0',
      borderTop: '1px solid #334155',
      display: 'flex',
      justifyContent: 'space-between'
    },
       popularStylesSection: {
      background: 'rgba(15, 23, 42, 0.6)',
      borderRadius: '16px',
      padding: '25px',
      marginTop: '40px',
      boxShadow: '0 10px 25px rgba(0, 0, 0, 0.2)',
      border: '1px solid #334155'
    }
  };

  return (
    <div>
      {/* Шапка сайта */}
      <Header />
      
      <div className="Configurator">
        <h1>Конфигуратор ремонта</h1>
        
        {/* Секция выбора категорий */}
        <div style={{ marginBottom: '30px' }}>
          <h2>Выберите категории</h2>
          <div style={{ display: 'flex', flexWrap: 'wrap', marginTop: '15px' }}>
            {categories.map(category => (
              <button
                key={category.id}
                onClick={() => handleCategorySelect(category.id)}
                style={styles_ui.categoryButton(selectedCategories.includes(category.id))}
              >
                {category.name}
              </button>
            ))}
          </div>
        </div>
        
        {/* Секция выбора товаров */}
        {selectedCategories.length > 0 && (
          <div style={{ marginBottom: '40px' }}>
            <h2>Выберите товары</h2>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '20px', marginTop: '20px' }}>
              {products
                .filter(product => selectedCategories.includes(product.category_id))
                .map(product => {const isSelected = selectedProducts.some(p => p.id === product.product_id);
                  return (
                    <div
  key={product.product_id}
  style={{
    ...styles_ui.productCard,
    ...(isSelected ? styles_ui.productCardHover : {}),
    opacity: isSelected ? 0.7 : 1
  }}
>
  <div>
    {/* 👇 добавляем фото */}
    {product.image_url && (
      <img
        src={product.image_url}
        alt={product.name}
        style={{
          width: '100%',
          height: '180px',
          objectFit: 'cover',
          borderRadius: '10px',
          marginBottom: '10px'
        }}
      />
    )}
    
    <h3 style={{ fontSize: '1.2rem', marginBottom: '10px', color: '#f1f5f9' }}>
      {product.name}
    </h3>
    <p style={{ color: '#94a3b8', marginBottom: '10px' }}>{product.description}</p>
    <div style={{ color: '#f1f5f9', fontWeight: 600, fontSize: '1.1rem' }}>
      {(product?.price ?? 0).toLocaleString()} ₽
    </div>
  </div>
  
  {!isSelected && (
    <button
      onClick={() => handleProductSelect(product)}
      style={styles_ui.addButton}
    >
      Добавить
    </button>
  )}
</div>
                  );
                })}
            </div>
          </div>
        )}
        
        {/* Панель выбранных товаров */}
        {selectedProducts.length > 0 && (
          <div style={styles_ui.selectedProductsPanel}>
            <h2 style={{ marginBottom: '20px', color: '#f1f5f9' }}>Добавленные товары</h2>
            {selectedProducts.map((product, index) => (
  <div key={`${product.product_id}-${index}`} style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                padding: '15px 0',
                borderBottom: '1px solid #334155'
              }}>
                <div>
                  <div style={{ fontSize: '1.1rem', fontWeight: 500, color: '#f1f5f9' }}>{product.name}</div>
                  <div style={{ color: '#94a3b8', fontSize: '0.9rem', marginTop: '5px' }}>
                    {categories.find(c => c.id === product.category_id)?.name}
                  </div>
                </div>
                <div style={{ display: 'flex', alignItems: 'center' }}>
                  <div style={{ fontWeight: 600, color: '#f1f5f9', marginRight: '15px' }}>
                    {(product?.price ?? 0).toLocaleString()} ₽

                  </div>
                  <button
                    onClick={() => handleProductDeselect(product.id)}
                    style={styles_ui.removeButton}
                  >
                    Удалить
                  </button>
                </div>
              </div>
            ))}
            <div style={styles_ui.totalPrice}>
              <span>Итого:</span>
              <span>{totalPrice.toLocaleString()} ₽</span>
            </div>
            {showStyleModal && selectedStyle && (
  <div style={{
    position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh',
    background: 'rgba(0,0,0,0.6)', display: 'flex',
    alignItems: 'center', justifyContent: 'center', zIndex: 1000
  }}>
    <div style={{
      background: '#0f172a', padding: '30px', borderRadius: '16px',
      maxWidth: '600px', width: '90%', boxShadow: '0 10px 30px rgba(0,0,0,0.4)'
    }}>
      <h2 style={{ color: '#f1f5f9', marginBottom: 20 }}>{selectedStyle.name}</h2>
      <img src={selectedStyle.image_url} alt={selectedStyle.name} style={{
        width: '100%', height: 200, objectFit: 'cover', borderRadius: 8, marginBottom: 20
      }} />
      <p style={{ color: '#cbd5e1', marginBottom: 10 }}>{selectedStyle.description}</p>
      <ul style={{ listStyle: 'none', padding: 0, marginBottom: 20 }}>
        {selectedStyle.items?.map((item, i) => (
          <li key={i} style={{
            color: '#f8fafc',
            borderBottom: '1px solid #334155',
            padding: '8px 0'
          }}>
            {item.product?.name} — {(item.product?.price ?? 0).toLocaleString()} ₽
          </li>
        ))}
      </ul>
      <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 10 }}>
        <button onClick={() => setShowStyleModal(false)} style={{
          padding: '10px 16px', background: '#475569', color: '#fff',
          border: 'none', borderRadius: 8, cursor: 'pointer'
        }}>
          Отмена
        </button>
        <button onClick={confirmStyle} style={{
          padding: '10px 16px', background: '#3b82f6', color: '#fff',
          border: 'none', borderRadius: 8, cursor: 'pointer'
        }}>
          Добавить в корзину
        </button>
      </div>
    </div>
  </div>
)}

          </div>
        )}
        
        {/* Секция популярных стилей */}
        <div style={styles_ui.popularStylesSection}>
          <h2 style={{ marginBottom: '15px', color: '#f1f5f9' }}>Популярные стили интерьера</h2>
          <p style={{ color: '#94a3b8', marginBottom: '25px' }}>
            Выберите готовый стиль интерьера, и мы автоматически добавим все необходимые товары для его реализации.
          </p>

          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(260px, 1fr))', gap: '20px' }}>
            {styles.map(style => (
              <div key={style.preset_id} style={{ background: '#1e293b', padding: 20, borderRadius: 12, border: '1px solid #334155', cursor: 'pointer' }} onClick={() => handleStyleClick(style)}>
                <img src={style.image_url} alt={style.name} style={{ width: '100%', height: 160, objectFit: 'cover', borderRadius: 8, marginBottom: 12 }} />
                <h3 style={{ color: '#f1f5f9', fontSize: '1.2rem' }}>{style.name}</h3>
                <p style={{ color: '#94a3b8', fontSize: '0.9rem', marginTop: 8 }}>{style.description}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
      <Footer />

      {showStyleModal && selectedStyle && (
        <div style={{ position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh', backgroundColor: 'rgba(0,0,0,0.6)', backdropFilter: 'blur(3px)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 9999 }}>
          <div style={{ background: '#0f172a', padding: '30px', borderRadius: '16px', maxWidth: '600px', width: '90%', boxShadow: '0 10px 30px rgba(0,0,0,0.4)' }}>
            <h2 style={{ color: '#f1f5f9', marginBottom: 20 }}>{selectedStyle.name}</h2>
            <img src={selectedStyle.image_url} alt={selectedStyle.name} style={{ width: '100%', height: 200, objectFit: 'cover', borderRadius: 8, marginBottom: 20 }} />
            <p style={{ color: '#cbd5e1', marginBottom: 10 }}>{selectedStyle.description}</p>
            <ul style={{ listStyle: 'none', padding: 0, marginBottom: 20 }}>
              {selectedStyle.items?.map((item, i) => (
                <li key={i} style={{ color: '#f8fafc', borderBottom: '1px solid #334155', padding: '8px 0' }}>
                  {item.product?.name} — {(item.product?.price ?? 0).toLocaleString()} ₽
                </li>
              ))}
            </ul>
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 10 }}>
              <button onClick={() => setShowStyleModal(false)} style={{ padding: '10px 16px', background: '#475569', color: '#fff', border: 'none', borderRadius: 8, cursor: 'pointer' }}>Отмена</button>
              <button onClick={confirmStyle} style={{ padding: '10px 16px', background: '#3b82f6', color: '#fff', border: 'none', borderRadius: 8, cursor: 'pointer' }}>Добавить в корзину</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default MainPage;