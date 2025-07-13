import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import CategoryManager from '../components/admin/CategoryManager';
import ProductManager from '../components/admin/ProductManager';
import StyleAdmin from '../components/StyleAdmin';
import { authAPI, tokenUtils, productAPI } from '../services/api';

function useQuery() {
  return new URLSearchParams(useLocation().search);
}

function AdminPage({
  categories,
  setCategories,
  products,
  setProducts,
  styles,
  setStyles,
}) {
  const query = useQuery();
  const key = query.get('key');
  const [adminToken, setAdminToken] = useState(tokenUtils.get());
  const [isLoading, setIsLoading] = useState(false);
  const [message, setMessage] = useState('');
  const [modalProducts, setModalProducts] = useState([]);
  const [modalVisible, setModalVisible] = useState(false);
  const [categoryForm, setCategoryForm] = useState({
    name: '',
    type: 'room', // room | element | sub
    parentRoom: '',
    parentElement: '',
  });

  const ADMIN_KEY = 'V2patTbDXS1wuqbqpyZGwg2vq70cem2wk3ElHO6y9l2FhfgNfN';

  const getAdminToken = async () => {
    if (adminToken) return adminToken;

    try {
      setIsLoading(true);
      const response = await authAPI.login(ADMIN_KEY);
      const token = response.token;
      tokenUtils.save(token);
      setAdminToken(token);
      return token;
    } catch (error) {
      console.error('Ошибка получения токена:', error);
      showMessage('Ошибка авторизации', true);
      return null;
    } finally {
      setIsLoading(false);
    }
  };

  const showMessage = (msg, isError = false) => {
    setMessage({ text: msg, isError });
    setTimeout(() => setMessage(''), 3000);
  };

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
      outline: 'none',
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
      letterSpacing: '0.5px',
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
      marginLeft: '10px',
    },
  };

  if (key !== ADMIN_KEY) {
    return (
      <div
        className="Configurator"
        style={{ maxWidth: 600, margin: '100px auto' }}
      >
        <div
          style={{
            padding: 40,
            textAlign: 'center',
            color: '#f8fafc',
            fontSize: 24,
            background: 'rgba(185, 28, 28, 0.1)',
            borderRadius: '12px',
            border: '1px solid rgba(185, 28, 28, 0.3)',
            boxShadow: '0 10px 25px rgba(185, 28, 28, 0.15)',
          }}
        >
          <i
            className="fas fa-lock"
            style={{ fontSize: 48, marginBottom: 20, color: '#b91c1c' }}
          ></i>
          <div>Доступ запрещён</div>
        </div>
      </div>
    );
  }

  const handleShowCategoryProducts = async (categoryId) => {
    try {
      const products = await productAPI.getByCategory(categoryId);
      setModalProducts(products);
      setModalVisible(true);
    } catch (err) {
      showMessage('Ошибка загрузки товаров категории', true);
    }
  };

  return (
    <div className="Configurator" style={{ maxWidth: 1000 }}>
      <h1>АДМИН-ПАНЕЛЬ</h1>

      {isLoading && (
        <div
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            background: 'rgba(0, 0, 0, 0.5)',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            zIndex: 1000,
          }}
        >
          <div
            style={{
              background: '#1e293b',
              padding: '20px',
              borderRadius: '10px',
              color: '#f1f5f9',
              fontSize: '1.2rem',
            }}
          >
            Загрузка...
          </div>
        </div>
      )}

      {message && (
        <div
          style={{
            position: 'fixed',
            top: '20px',
            right: '20px',
            background: message.isError ? '#dc2626' : '#059669',
            color: 'white',
            padding: '15px 20px',
            borderRadius: '8px',
            zIndex: 1001,
            boxShadow: '0 4px 12px rgba(0, 0, 0, 0.3)',
          }}
        >
          {message.text}
        </div>
      )}

      <div
        style={{
          background: 'rgba(30,41,59,0.5)',
          padding: 20,
          marginTop: 40,
          border: '1px solid #334155',
          borderRadius: 12,
        }}
      ></div>

      <CategoryManager
        categories={categories}
        setCategories={setCategories}
        getAdminToken={getAdminToken}
        showMessage={showMessage}
        styles={uiStyles}
        onViewCategoryProducts={handleShowCategoryProducts}
      />

      <ProductManager
        categories={categories}
        products={products}
        setProducts={setProducts}
        getAdminToken={getAdminToken}
        showMessage={showMessage}
        styles={uiStyles}
      />

      <StyleAdmin products={products} styles={styles} setStyles={setStyles} />

      {modalVisible && (
        <div
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            background: 'rgba(0, 0, 0, 0.6)',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            zIndex: 1002,
          }}
        >
          <div
            style={{
              background: '#1e293b',
              padding: 30,
              borderRadius: 12,
              maxWidth: 600,
              width: '90%',
              color: '#f1f5f9',
            }}
          >
            <h2 style={{ marginBottom: 20 }}>📦 Товары категории</h2>
            <ul style={{ listStyle: 'none', padding: 0 }}>
              {modalProducts.map((p) => (
                <li
                  key={p.product_id}
                  style={{
                    marginBottom: 10,
                    borderBottom: '1px solid #334155',
                    paddingBottom: 6,
                  }}
                >
                  <strong>{p.name}</strong> — {p.price} ₽
                </li>
              ))}
              {modalProducts.length === 0 && <li>Нет товаров в категории</li>}
            </ul>
            <button
              onClick={() => setModalVisible(false)}
              style={{ marginTop: 20, ...uiStyles.buttonStyle }}
            >
              Закрыть
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

export default AdminPage;
