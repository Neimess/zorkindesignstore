import React, { useState, useEffect } from 'react';
import Header from '../components/layout/Header';
import Footer from '../components/layout/Footer';
import TelegramIcon from '../assets/telegram-icon.png';
import { categoryAPI } from '../services/api'; // уже должен быть
import { serviceAPI } from '../services/api';
import { productAPI } from '../services/api';
import { coefficientAPI } from '../services/api';



/**
 * Компонент главной страницы с конфигуратором
 *
 * @param {Object} props - Свойства компонента
 * @param {Array} props.categories - Список категорий
 * @param {Array} props.products - Список товаров
 * @param {Array} props.styles - Список стилей
 */
function MainPage({ products, styles }) {
  const [categories, setCategories] = useState([]);
  const [selectedProducts, setSelectedProducts] = useState([]);
  const [selectedStyle, setSelectedStyle] = useState(null);
  const [showStyleModal, setShowStyleModal] = useState(false);
  const [selectedRoomType, setSelectedRoomType] = useState(null);
  const [selectedElement, setSelectedElement] = useState(null);
  const [selectedSubElement, setSelectedSubElement] = useState(null);
  const [modalProduct, setModalProduct] = useState(null);
const [allServices, setAllServices] = useState([]);
const [coefficients, setCoefficients] = useState([]);
const [activeCoefficient, setActiveCoefficient] = useState(1);

const handleOpenProductModal = async (product) => {
  const id    = product.product_id || product.id;
  const token = localStorage.getItem('admin_token');

  console.log('▶️ Открываю карточку товара (id):', id);
  console.log('▶️ Отправляю запрос к API /admin/product/{id}');

  try {
    const detailed = await productAPI.getAdminById(id, token);
    console.log('✅ Получен детальный товар:', detailed);

    setModalProduct({
      ...detailed,
      services: detailed.services ?? [],
    });

    console.log('📦 modalProduct установлен:', {
      ...detailed,
      services: detailed.services ?? [],
    });

  } catch (err) {
    console.error('❌ Ошибка получения детального товара:', err);

    setModalProduct({ ...product, services: [] });
  }
};

const handleSendToTelegram = () => {
  const messageLines = [];

  messageLines.push('📋 Данные проекта:');
  messageLines.push(`• Площадь: ${formData.area} м²`);
  messageLines.push(`• Тип недвижимости: ${formData.propertyType || 'не указано'}`);
  messageLines.push(`• Комнат: ${formData.rooms}`);
  messageLines.push(`• Санузлов: ${formData.bathrooms}`);
  messageLines.push('');

  messageLines.push('🛒 Товары:');
  selectedProducts.forEach(p => {
    messageLines.push(`• ${p.name} x ${p.quantity} = ${p.price * p.quantity} ₽`);
  });

  messageLines.push('');
  messageLines.push('🔧 Услуги:');
  serviceCart.forEach(s => {
    const qtyText = `${s.quantity} ${s.unit || ''}`.trim();
    messageLines.push(`• ${s.name} (${qtyText}) = ${s.price * s.quantity} ₽`);
  });

  const finalMessage = messageLines.join('\n');

  fetch('https://api.telegram.org/bot<ТОКЕН>/sendMessage', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      chat_id: '<ID_ЧАТА>',
      text: finalMessage,
    }),
  }).then(res => {
    if (res.ok) {
      alert('Отправлено в Telegram!');
    } else {
      alert('Ошибка отправки в Telegram');
    }
  });
};


const [formData, setFormData] = useState({
  area: '',
  propertyType: '',
  rooms: '',
  bathrooms: ''
});

// "Комнаты" (нет parent_id)
const rooms = categories.filter(c => !('parent_id' in c) || c.parent_id === null || c.parent_id === 0);

const [serviceCart, setServiceCart] = useState([]);

const handleAddService = (service) => {
  if (!serviceCart.find(s => s.service_id === service.id)) {
    setServiceCart(prev => [...prev, { 
      ...service, 
      service_id: service.id, 
      quantity: 1, 
      unit: '' 
    }]);
  }
};


const handleRemoveService = (serviceId) => {
  setServiceCart(prev => prev.filter(s => s.service_id !== serviceId));
};

useEffect(() => {
  const fetchServices = async () => {
    try {
      const token = localStorage.getItem('admin_token'); // Или метод получения токена
      const services = await serviceAPI.getAll(token);
      setAllServices(services || []);
    } catch (err) {
      console.error('Ошибка загрузки услуг:', err);
    }
  };

  fetchServices();
}, []);

useEffect(() => {
  const fetchCoefficients = async () => {
    try {
      const token = localStorage.getItem('admin_token');
      const data = await coefficientAPI.getAll(token);
      setCoefficients(data || []);
    } catch (err) {
      console.error('Ошибка загрузки коэффициентов:', err);
    }
  };

  fetchCoefficients();
}, []);




// Элементы внутри выбранной комнаты
const elements = categories.filter(c => String(c.parent_id) === String(selectedRoomType?.id));

// Подкатегории внутри выбранного элемента
const subElements = categories.filter(c => String(c.parent_id) === String(selectedElement?.id));

  useEffect(() => {
    document.body.style.overflow = showStyleModal ? 'hidden' : '';
    return () => {
      document.body.style.overflow = '';
    };
  }, [showStyleModal]);

// Получение плоского списка категорий
useEffect(() => {
  const fetchCategories = async () => {
    try {
      const flat = await categoryAPI.getAll();
      console.log('Загруженные категории:', flat);
      setCategories(flat); // Храним плоский массив, не дерево
    } catch (e) {
      console.error('Ошибка при загрузке категорий:', e);
    }
  };

  fetchCategories();
}, []);

  /**
   * Обработчик выбора категории
   * @param {number} category_id - ID категории
   */

  /**
   * Обработчик выбора товара
   * @param {Object} product - Объект товара
   */
  const handleProductSelect = (product) => {
    if (!selectedProducts.find((p) => p.id === product.product_id)) {
      setSelectedProducts([...selectedProducts, { ...product, quantity: 1 }]);
    }
  };

  /**
   * Обработчик удаления товара из выбранных
   * @param {number} productId - ID товара
   */
  const handleProductDeselect = (productId) => {
    setSelectedProducts((prev) => {
      const index = prev.findIndex((p) => p.id === productId);
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
  // const confirmStyle = () => {
  //   const newProducts = (selectedStyle.items || [])
  //     .map((item) => item.product)
  //     .filter(
  //       (product) =>
  //         product && !selectedProducts.some((p) => p.product_id === product.id),
  //     )
  //     .map((product) => ({
  //       ...product,
  //       product_id: product.id,
  //     }));

  //   setSelectedProducts((prev) => [...prev, ...newProducts]);
  //   setShowStyleModal(false);
  //   setSelectedStyle(null);
  // };
  const confirmStyle = () => {
  const newProducts = (selectedStyle.items || [])
    .map(item => item.product)
    .filter(product => product)
    .filter(product => !selectedProducts.some(p => p.product_id === product.id))
    .map(product => ({
      ...product,
      product_id: product.id,           // унифицированный идентификатор
      quantity: 1,                       // начальное количество
      price: product.price ?? 0          // цена по умолчанию
    }));

  setSelectedProducts(prev => [...prev, ...newProducts]);
  setShowStyleModal(false);
  setSelectedStyle(null);
};

  /**
   * Обработчик выбора стиля
   * @param {Object} style - Объект стиля
   */


  // Вычисляем общую стоимость выбранных товаров
  // const totalPrice = selectedProducts.reduce(
  //   (sum, product) => sum + product.price * product.quantity,
  //   0,
  // );
const totalPrice = selectedProducts.reduce(
  (sum, product) =>
    sum + (Number(product.price) || 0) * (Number(product.quantity) || 1),
  0
);

  // Стили для элементов интерфейса
  const styles_ui = {
    categoryButton: (isSelected) => ({
      background: isSelected
        ? 'linear-gradient(135deg, #3b82f6, #2563eb)'
        : 'rgba(15, 23, 42, 0.6)',
      color: '#f1f5f9',
      border: isSelected ? 'none' : '1px solid #334155',
      borderRadius: '10px',
      padding: '12px 20px',
      margin: '0 10px 10px 0',
      fontSize: '1rem',
      fontWeight: isSelected ? 600 : 400,
      cursor: 'pointer',
      transition: 'all 0.3s ease',
      boxShadow: isSelected ? '0 4px 12px rgba(37, 99, 235, 0.3)' : 'none',
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
      justifyContent: 'space-between',
    },
    productCardHover: {
      transform: 'translateY(-5px)',
      boxShadow: '0 10px 20px rgba(0, 0, 0, 0.2)',
      borderColor: '#3b82f6',
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
      letterSpacing: '0.5px',
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
      marginLeft: '10px',
    },
    selectedProductsPanel: {
      background: 'rgba(15, 23, 42, 0.8)',
      borderRadius: '16px',
      padding: '25px',
      marginTop: '30px',
      boxShadow: '0 10px 25px rgba(0, 0, 0, 0.2)',
      border: '1px solid #334155',
    },
    totalPrice: {
      fontSize: '1.5rem',
      fontWeight: 700,
      color: '#f1f5f9',
      marginTop: '20px',
      padding: '15px 0',
      borderTop: '1px solid #334155',
      display: 'flex',
      justifyContent: 'space-between',
    },
    popularStylesSection: {
      background: 'rgba(15, 23, 42, 0.6)',
      borderRadius: '16px',
      padding: '25px',
      marginTop: '40px',
      boxShadow: '0 10px 25px rgba(0, 0, 0, 0.2)',
      border: '1px solid #334155',
    },
  };
// console.log('modalProduct.services:', modalProduct.services);
// console.log('allServices:', allServices);
console.log('🛑 modalProduct в модалке:', modalProduct);
console.log('📋 modalProduct.services:', modalProduct?.services);

  return (
    <div>
      {/* Шапка сайта */}
      <Header />

      <div className="Configurator">
  <h1>Конфигуратор ремонта</h1>
  <div
    style={{
      background: 'rgba(15, 23, 42, 0.6)',
      borderRadius: '16px',
      padding: '25px',
      margin: '30px 0',
      boxShadow: '0 10px 25px rgba(0, 0, 0, 0.2)',
      border: '1px solid #334155',
      color: '#f1f5f9',
    }}
  >
    <h2 style={{ marginBottom: '20px' }}>
      Перед началом ответьте на несколько вопросов
    </h2>
    <div style={{ display: 'grid', gap: '15px' }}>
      
      {/* Площадь помещения */}
      <label>
        Площадь помещения (м²):
        <input
          type="number"
          name="area"
          value={formData.area}
          onChange={(e) => setFormData({ ...formData, area: e.target.value })}
          style={{
            width: '100%',
            padding: '8px',
            borderRadius: '8px',
            marginTop: '5px',
            background: '#1e293b',
            border: '1px solid #334155',
            color: '#f1f5f9',
          }}
        />
      </label>

      {/* Тип недвижимости */}
      <label>
        Тип недвижимости:
        <select
          name="propertyType"
  value={formData.propertyType}
  onChange={(e) => {
    const type = e.target.value;

    setFormData({ ...formData, propertyType: type });

    if (type === 'primary') {
      const coeff = coefficients.find(c => c.name.toLowerCase().includes('первич'))?.value || 1;
      setActiveCoefficient(coeff);
    } else if (type === 'secondary') {
      const coeff = coefficients.find(c => c.name.toLowerCase().includes('вторич'))?.value || 1;
      setActiveCoefficient(coeff);
    } else {
      setActiveCoefficient(1);
    }
  }}
          style={{
            width: '100%',
            padding: '8px',
            borderRadius: '8px',
            marginTop: '5px',
            background: '#1e293b',
            border: '1px solid #334155',
            color: '#f1f5f9',
          }}


        >
          <option value="">-- Выберите --</option>
          <option value="primary">Первичный рынок</option>
          <option value="secondary">Вторичный рынок</option>
        </select>
      </label>

      {/* Количество комнат */}
      <label>
        Количество комнат:
        <input
          type="number"
          name="rooms"
          min="1"
          value={formData.rooms}
          onChange={(e) => setFormData({ ...formData, rooms: e.target.value })}
          style={{
            width: '100%',
            padding: '8px',
            borderRadius: '8px',
            marginTop: '5px',
            background: '#1e293b',
            border: '1px solid #334155',
            color: '#f1f5f9',
          }}
        />
      </label>

      {/* Количество санузлов */}
      <label>
        Количество санузлов:
        <input
          type="number"
          name="bathrooms"
          min="1"
          value={formData.bathrooms}
          onChange={(e) => setFormData({ ...formData, bathrooms: e.target.value })}
          style={{
            width: '100%',
            padding: '8px',
            borderRadius: '8px',
            marginTop: '5px',
            background: '#1e293b',
            border: '1px solid #334155',
            color: '#f1f5f9',
          }}
        />
      </label>

      <button
        style={{
          marginTop: '20px',
          padding: '12px 20px',
          borderRadius: '8px',
          border: 'none',
          background: '#3b82f6',
          color: '#fff',
          cursor: 'pointer',
          fontWeight: 600,
        }}
        onClick={() => alert('Форма сохранена. Продолжайте конфигурацию.')}
      >
        Начать конфигурацию
      </button>
    </div>
  </div>



        {/* Секция выбора категорий */}
        <div style={{ marginBottom: '30px' }}>
          <h2>Выберите комнату, элемент и подкатегорию</h2>

             {/* Комнаты */}
{rooms.map(room => (
  <button
    key={room.id}
    onClick={() => {
      setSelectedRoomType(room);
      setSelectedElement(null);
      setSelectedSubElement(null);
    }}
    style={styles_ui.categoryButton(selectedRoomType?.id === room.id)}
  >
    {room.name}
  </button>
))}

{/* Элементы */}
{selectedRoomType && (
  <>
    <h3>Выберите элемент</h3>
    {elements.map(elem => (
      <button
        key={elem.id}
        onClick={() => {
          setSelectedElement(elem);
          setSelectedSubElement(null);
        }}
        style={styles_ui.categoryButton(selectedElement?.id === elem.id)}
      >
        {elem.name}
      </button>
    ))}
  </>
)}

{/* Подкатегории */}
{selectedElement && (
  <>
    <h3>Выберите подкатегорию</h3>
    {subElements.map(sub => (
      <button
        key={sub.id}
        onClick={() => setSelectedSubElement(sub)}
        style={styles_ui.categoryButton(selectedSubElement?.id === sub.id)}
      >
        {sub.name}
      </button>
    ))}
  </>
)}
 
          
        </div>

        {/* Сюда?  */}

        {/* Секция выбора товаров */}
        {selectedSubElement && (
  <div style={{ marginTop: '30px', marginBottom: '40px' }}>
    <h2>Выберите товары</h2>
    <div
      style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
        gap: '20px',
        marginTop: '20px',
      }}
    >
      {/* {products
        .filter(product => product.category_id === selectedSubElement?.id)
        .map(product => {
          const isSelected = selectedProducts.some(
            (p) => p.id === product.product_id
          );

          return (
            <div
              key={product.product_id}
              style={{
                ...styles_ui.productCard,
                ...(isSelected ? styles_ui.productCardHover : {}),
                opacity: isSelected ? 0.7 : 1,
              }}
            >
              <div>
                {product.image_url && (
                  <img
                    src={product.image_url}
                    alt={product.name}
                    style={{
                      width: '100%',
                      height: '180px',
                      objectFit: 'cover',
                      borderRadius: '10px',
                      marginBottom: '10px',
                    }}
                  />
                )}
                <h3 style={{ fontSize: '1.2rem', marginBottom: '10px', color: '#f1f5f9' }}>
                  {product.name}
                </h3>
                <p style={{ color: '#94a3b8', marginBottom: '10px' }}>
                  {product.description}
                </p>
                <div style={{ color: '#f1f5f9', fontWeight: 600, fontSize: '1.1rem' }}>
                  {(product?.price ?? 0).toLocaleString()} ₽
                </div>
              </div>
              {!isSelected && (
                <button
                  onClick={() => handleProductSelect({ ...product, quantity: 1 })}
                  style={styles_ui.addButton}
                >
                  Добавить
                </button>
              )}
            </div>
          );
        })} */}
        {products
  .filter(product => product.category_id === selectedSubElement?.id)
  .map(product => {
    const isSelected = selectedProducts.some(
      (p) => p.product_id === product.product_id
    );

    return (
      <div
        key={product.product_id}
        onClick={() => handleOpenProductModal(product)}



        style={{
          ...styles_ui.productCard,
          ...(isSelected ? styles_ui.productCardHover : {}),
          opacity: isSelected ? 0.7 : 1,
          cursor: 'pointer',
        }}
      >
        <div>
          {product.image_url && (
            <img
              src={product.image_url}
              alt={product.name}
              style={{
                width: '100%',
                height: '180px',
                objectFit: 'cover',
                borderRadius: '10px',
                marginBottom: '10px',
              }}
            />
          )}
          <h3>{product.name}</h3>
          <p>{product.description}</p>
          <div style={{ fontWeight: 'bold' }}>
            {(product.price ?? 0).toLocaleString()} ₽
          </div>
        </div>

        {!isSelected && (
          <button
            onClick={(e) => {
              e.stopPropagation(); // важно: чтобы клик по кнопке не открыл модалку
              handleProductSelect({ ...product, quantity: 1 });
            }}
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
            <h2 style={{ marginBottom: '20px', color: '#f1f5f9' }}>
              Добавленные товары
            </h2>
            {selectedProducts.map((product, index) => (
              <div
                key={`${product.product_id}-${index}`}
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  padding: '15px 0',
                  borderBottom: '1px solid #334155',
                }}
              >
                <div>
                  <div
                    style={{
                      fontSize: '1.1rem',
                      fontWeight: 500,
                      color: '#f1f5f9',
                    }}
                  >
                    {product.name}
                  </div>
                  <div
                    style={{
                      color: '#94a3b8',
                      fontSize: '0.9rem',
                      marginTop: '5px',
                    }}
                  >
                    {categories.find((c) => c.id === product.category_id)?.name}
                  </div>
                </div>
                <div
                  style={{ display: 'flex', alignItems: 'center', gap: '10px' }}
                >
                  <input
                    type="number"
                    min="1"
                    value={product.quantity}
                    // onChange={(e) => {
                    //   const newQty = parseInt(e.target.value, 10) || 1;
                    //   setSelectedProducts((prev) =>
                    //     prev.map((p) =>
                    //       p.id === product.id ? { ...p, quantity: newQty } : p,
                    //     ),
                    //   );
                    // }}
                    onChange={(e) => {
  const newQty = parseInt(e.target.value, 10) || 1;
  setSelectedProducts(prev =>
    prev.map(p =>
      p.product_id === product.product_id
        ? { ...p, quantity: newQty }
        : p
    )
  );
}}

                    style={{
                      width: '60px',
                      padding: '5px',
                      borderRadius: '6px',
                      background: '#1e293b',
                      border: '1px solid #334155',
                      color: '#f1f5f9',
                    }}
                  />
                  <div style={{ fontWeight: 600, color: '#f1f5f9' }}>
                    {(product.price * product.quantity).toLocaleString()} ₽
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
            <div
              style={{
                display: 'flex',
                justifyContent: 'flex-end',
                marginTop: '20px',
              }}
            >
              <button
                onClick={handleSendToTelegram}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: 8,
                  padding: '10px 16px',
                  background: '#229ED9',
                  color: '#fff',
                  border: 'none',
                  borderRadius: 8,
                  cursor: 'pointer',
                }}
              >
                <img
                  src={TelegramIcon}
                  alt="Telegram"
                  style={{ width: 20, height: 20 }}
                />
                Отправить
              </button>
            </div>

            {showStyleModal && selectedStyle && (
              <div
                style={{
                  position: 'fixed',
                  top: 0,
                  left: 0,
                  width: '100vw',
                  height: '100vh',
                  background: 'rgba(0,0,0,0.6)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  zIndex: 1000,
                }}
              >
                <div
                  style={{
                    background: '#0f172a',
                    padding: '30px',
                    borderRadius: '16px',
                    maxWidth: '600px',
                    width: '90%',
                    boxShadow: '0 10px 30px rgba(0,0,0,0.4)',
                  }}
                >
                  <h2 style={{ color: '#f1f5f9', marginBottom: 20 }}>
                    {selectedStyle.name}
                  </h2>
                  <img
                    src={selectedStyle.image_url}
                    alt={selectedStyle.name}
                    style={{
                      width: '100%',
                      height: 200,
                      objectFit: 'cover',
                      borderRadius: 8,
                      marginBottom: 20,
                    }}
                  />
                  <p style={{ color: '#cbd5e1', marginBottom: 10 }}>
                    {selectedStyle.description}
                  </p>
                  <ul
                    style={{ listStyle: 'none', padding: 0, marginBottom: 20 }}
                  >
                    {selectedStyle.items?.map((item, i) => (
                      <li
                        key={i}
                        style={{
                          color: '#f8fafc',
                          borderBottom: '1px solid #334155',
                          padding: '8px 0',
                        }}
                      >
                        {item.product?.name} —{' '}
                        {(item.product?.price ?? 0).toLocaleString()} ₽
                      </li>
                    ))}
                  </ul>
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'flex-end',
                      gap: 10,
                    }}
                  >
                    <button
                      onClick={() => setShowStyleModal(false)}
                      style={{
                        padding: '10px 16px',
                        background: '#475569',
                        color: '#fff',
                        border: 'none',
                        borderRadius: 8,
                        cursor: 'pointer',
                      }}
                    >
                      Отмена
                    </button>
                    <button
                      onClick={confirmStyle}
                      style={{
                        padding: '10px 16px',
                        background: '#3b82f6',
                        color: '#fff',
                        border: 'none',
                        borderRadius: 8,
                        cursor: 'pointer',
                      }}
                    >
                      Добавить в корзину
                    </button>
                  </div>
                </div>
              </div>
            )}
          </div>
        )}

{serviceCart.length > 0 && (
  <div style={styles_ui.selectedProductsPanel}>
    <h2 style={{ marginBottom: '20px', color: '#f1f5f9' }}>
      Добавленные услуги
    </h2>

    {serviceCart.map((service, index) => (
      <div
        key={`${service.service_id}-${index}`}
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          padding: '15px 0',
          borderBottom: '1px solid #334155'
        }}
      >
        <div style={{ flex: 1 }}>
          <div
            style={{
              fontSize: '1.1rem',
              fontWeight: 500,
              color: '#f1f5f9'
            }}
          >
            {service.name}
          </div>

          <div style={{ display: 'flex', gap: '10px', marginTop: '8px' }}>
            <input
              type="number"
              min="1"
              value={service.quantity}
              onChange={(e) => {
                const qty = parseInt(e.target.value) || 1;
                setServiceCart(prev => prev.map(s =>
                  s.service_id === service.service_id
                    ? { ...s, quantity: qty }
                    : s
                ));
              }}
              style={{
                width: '60px',
                padding: '5px',
                borderRadius: '6px',
                background: '#1e293b',
                border: '1px solid #334155',
                color: '#f1f5f9'
              }}
            />

            <input
              type="text"
              placeholder="м² / этажи / шт"
              value={service.unit || ''}
              onChange={(e) => {
                const unit = e.target.value;
                setServiceCart(prev => prev.map(s =>
                  s.service_id === service.service_id
                    ? { ...s, unit }
                    : s
                ));
              }}
              style={{
                flex: 1,
                padding: '5px',
                borderRadius: '6px',
                background: '#1e293b',
                border: '1px solid #334155',
                color: '#f1f5f9'
              }}
            />
          </div>
        </div>

        <div style={{ fontWeight: 600, color: '#f1f5f9', minWidth: '90px', textAlign: 'right' }}>
          {(service.price * (service.quantity || 1) * activeCoefficient).toLocaleString()
} ₽
        </div>

        <button
          onClick={() => handleRemoveService(service.service_id)}
          style={styles_ui.removeButton}
        >
          Удалить
        </button>
      </div>
    ))}
  </div>
)}



        {/* Секция популярных стилей */}
        <div style={styles_ui.popularStylesSection}>
          <h2 style={{ marginBottom: '15px', color: '#f1f5f9' }}>
            Популярные стили интерьера
          </h2>
          <p style={{ color: '#94a3b8', marginBottom: '25px' }}>
            Выберите готовый стиль интерьера, и мы автоматически добавим все
            необходимые товары для его реализации.
          </p>

          <div
            style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fill, minmax(260px, 1fr))',
              gap: '20px',
            }}
          >
            {styles.map((style) => (
              <div
                key={style.preset_id}
                style={{
                  background: '#1e293b',
                  padding: 20,
                  borderRadius: 12,
                  border: '1px solid #334155',
                  cursor: 'pointer',
                }}
                onClick={() => handleStyleClick(style)}
              >
                <img
                  src={style.image_url}
                  alt={style.name}
                  style={{
                    width: '100%',
                    height: 160,
                    objectFit: 'cover',
                    borderRadius: 8,
                    marginBottom: 12,
                  }}
                />
                <h3 style={{ color: '#f1f5f9', fontSize: '1.2rem' }}>
                  {style.name}
                </h3>
                <p
                  style={{ color: '#94a3b8', fontSize: '0.9rem', marginTop: 8 }}
                >
                  {style.description}
                </p>
              </div>
            ))}
          </div>
        </div>
      </div>
      <Footer />

      {showStyleModal && selectedStyle && (
        <div
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            width: '100vw',
            height: '100vh',
            backgroundColor: 'rgba(0,0,0,0.6)',
            backdropFilter: 'blur(3px)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            zIndex: 9999,
          }}
        >
          <div
            style={{
              background: '#0f172a',
              padding: '30px',
              borderRadius: '16px',
              maxWidth: '600px',
              width: '90%',
              boxShadow: '0 10px 30px rgba(0,0,0,0.4)',
            }}
          >
            <h2 style={{ color: '#f1f5f9', marginBottom: 20 }}>
              {selectedStyle.name}
            </h2>
            <img
              src={selectedStyle.image_url}
              alt={selectedStyle.name}
              style={{
                width: '100%',
                height: 200,
                objectFit: 'cover',
                borderRadius: 8,
                marginBottom: 20,
              }}
            />
            <p style={{ color: '#cbd5e1', marginBottom: 10 }}>
              {selectedStyle.description}
            </p>
            <ul style={{ listStyle: 'none', padding: 0, marginBottom: 20 }}>
              {selectedStyle.items?.map((item, i) => (
                <li
                  key={i}
                  style={{
                    color: '#f8fafc',
                    borderBottom: '1px solid #334155',
                    padding: '8px 0',
                  }}
                >
                  {item.product?.name} —{' '}
                  {(item.product?.price ?? 0).toLocaleString()} ₽
                </li>
              ))}
            </ul>
            <div
              style={{ display: 'flex', justifyContent: 'flex-end', gap: 10 }}
            >
              <button
                onClick={() => setShowStyleModal(false)}
                style={{
                  padding: '10px 16px',
                  background: '#475569',
                  color: '#fff',
                  border: 'none',
                  borderRadius: 8,
                  cursor: 'pointer',
                }}
              >
                Отмена
              </button>
              <button
                onClick={confirmStyle}
                style={{
                  padding: '10px 16px',
                  background: '#3b82f6',
                  color: '#fff',
                  border: 'none',
                  borderRadius: 8,
                  cursor: 'pointer',
                }}
              >
                Добавить в корзину
              </button>
            </div>
          </div>
        </div>
        
      )}
      {modalProduct && (
  <div
    style={{
      position: 'fixed',
      top: 0,
      left: 0,
      width: '100vw',
      height: '100vh',
      backgroundColor: 'rgba(0,0,0,0.6)',
      backdropFilter: 'blur(2px)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      zIndex: 9999,
    }}
    onClick={() => setModalProduct(null)} // клик вне закрывает модалку
  >
    <div
      style={{
        background: '#1e293b',
        padding: '30px',
        borderRadius: '16px',
        maxWidth: '600px',
        width: '90%',
        color: '#f1f5f9',
        boxShadow: '0 10px 30px rgba(0,0,0,0.4)',
      }}
      onClick={(e) => e.stopPropagation()} // клик внутри — не закрывает
    >
      <h2 style={{ fontSize: '1.4rem', marginBottom: '10px' }}>{modalProduct.name}</h2>

      {modalProduct.image_url && (
        <img
          src={modalProduct.image_url}
          alt={modalProduct.name}
          style={{
            width: '100%',
            height: '200px',
            objectFit: 'cover',
            borderRadius: '8px',
            marginBottom: '15px',
          }}
        />
      )}

      <p style={{ marginBottom: '15px' }}>
        {modalProduct.description || 'Описание отсутствует.'}
      </p>

      <div style={{ fontWeight: 'bold', fontSize: '1.2rem', marginBottom: '20px' }}>
        {(modalProduct.price ?? 0).toLocaleString()} ₽
      </div>

      {Array.isArray(modalProduct.attributes) && modalProduct.attributes.length > 0 && (
  <div style={{ marginBottom: '20px' }}>
    <h3 style={{ fontSize: '1rem', marginBottom: '10px' }}>Характеристики:</h3>
    <ul>
      {modalProduct.attributes.map((attr, index) => (
        <li key={index}>
          {attr.name}: {attr.value}
          {attr.unit ? ` ${attr.unit}` : ''}
        </li>
      ))}
    </ul>
  </div>
)}


{Array.isArray(modalProduct.services) && modalProduct.services.length > 0 && (
  <div style={{ marginBottom: '20px' }}>
    <h3>Сопутствующие услуги:</h3>
    <ul>
      {modalProduct.services.map(rawService => {
        const serviceId = rawService.service_id || rawService.id;
        const fullService = allServices.find(s => s.id === serviceId);

        if (!fullService) {
          return (
            <li key={serviceId}>
              Услуга ID {serviceId}
            </li>
          );
        }

        return (
          <li key={fullService.id} style={{ marginBottom: '8px' }}>
            {fullService.name} — {fullService.price} ₽
            <button
              onClick={() => handleAddService(fullService)}
              style={{
                marginLeft: '10px',
                padding: '4px 8px',
                background: '#3b82f6',
                color: '#fff',
                border: 'none',
                borderRadius: '6px',
                fontSize: '0.8rem',
                cursor: 'pointer',
              }}
            >
              Добавить в корзину услуг
            </button>
          </li>
        );
      })}
    </ul>
  </div>
)}



      <button
        onClick={() => setModalProduct(null)}
        style={{
          marginTop: '10px',
          padding: '10px 16px',
          background: '#475569',
          color: '#fff',
          border: 'none',
          borderRadius: '8px',
          cursor: 'pointer',
          width: '100%',
        }}
      >
        Закрыть
      </button>
    </div>
  </div>
)}


    </div>
    
  );
}

export default MainPage;
