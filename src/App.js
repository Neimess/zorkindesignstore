import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css';

// Импорт страниц
import MainPage from './pages/MainPage';
import AdminPage from './pages/AdminPage';

// Импорт API сервисов
import { categoryAPI, productAPI, presetAPI } from './services/api';

/**
 * Главный компонент приложения
 * Управляет глобальным состоянием и маршрутизацией
 */
function App() {
  // Состояние для хранения данных
  const [categories, setCategories] = useState([]);
  const [products, setProducts] = useState([]);
  const [styles, setStyles] = useState([]);

  // Состояние для отображения загрузки и ошибок
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  /**
   * Загрузка данных при монтировании компонента
   */

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);

        // 1. Загрузка категорий
        console.log('Начинаем загрузку категорий...');
        const categoriesData = await categoryAPI.getAll();
        console.log('Получены категории от API:', categoriesData);
        
        if (!Array.isArray(categoriesData)) {
          console.error('Ошибка: categoriesData не является массивом:', categoriesData);
          setCategories([]);
        } else {
          setCategories(categoriesData);
          console.log('Категории установлены в state:', categoriesData);
        }

        // 2. Загрузка товаров по категориям
        console.log('Загруженные категории:', categoriesData);

        const productLists = [];
        if (Array.isArray(categoriesData)) {
          for (const cat of categoriesData) {
            try {
              console.log(`Загрузка товаров для категории ID ${cat.id}...`);
              const products = await productAPI.getByCategory(cat.id);
              console.log(`Получены товары для категории ID ${cat.id}:`, products);
              productLists.push(products);
            } catch (err) {
              console.warn(
                `Ошибка при загрузке товаров категории ID ${cat.id}:`,
                err.message,
              );
              productLists.push([]); // чтобы не ломать структуру, даже если запрос не удался
            }
          }
        } else {
          console.error('Невозможно загрузить товары: categoriesData не является массивом');
        }

        setProducts(productLists.flat()); // объединяем все массивы
        console.log('Товары установлены в state:', productLists.flat());

        // 3. Загрузка стилей
        console.log('Начинаем загрузку стилей...');
        const stylesData = await presetAPI.getAllDetailed();
        console.log('Получены стили от API:', stylesData);
        setStyles(stylesData);

        setLoading(false);
      } catch (err) {
        console.error('Ошибка при загрузке данных:', err);
        setError(
          'Произошла ошибка при загрузке данных. Пожалуйста, попробуйте позже.',
        );
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  // Отображение индикатора загрузки
  if (loading) {
    return (
      <div className="loading-container">
        <div className="loading-spinner"></div>
        <div className="loading-text">Загрузка данных...</div>
      </div>
    );
  }

  // Отображение ошибки
  if (error) {
    return (
      <div className="error-container">
        <div className="error-icon">⚠️</div>
        <div className="error-message">{error}</div>
        <button
          className="error-retry-button"
          onClick={() => window.location.reload()}
        >
          Попробовать снова
        </button>
      </div>
    );
  }

  return (
    <Routes>
      <Route
        path="/"
        element={
          <MainPage
            categories={categories}
            products={products}
            styles={styles}
          />
        }
      />

      <Route
        path="/admin"
        element={
          <AdminPage
            categories={categories}
            setCategories={setCategories}
            products={products}
            setProducts={setProducts}
            styles={styles}
            setStyles={setStyles}
          />
        }
      />
    </Routes>
  );
}

export default App;
