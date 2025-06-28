import React, { useState, useEffect } from 'react';
import { Routes, Route, useLocation, Navigate } from 'react-router-dom';
import './App.css';

import { initialStyles } from './data/styles';
import StyleSelector from './components/StyleSelector';
import StyleAdmin from './components/StyleAdmin';
import { categoryAPI, productAPI, presetAPI, authAPI, tokenUtils, categoryAttributeAPI } from './services/api';


const initialCategories = [
  { id: 1, name: 'Напольные покрытия' },
  { id: 2, name: 'Настенные покрытия' },
  { id: 3, name: 'Потолочные решения' },
  { id: 4, name: 'Двери' },
  { id: 5, name: 'Окна' },
  { id: 6, name: 'Освещение' },
  { id: 7, name: 'Сантехника' },
  { id: 8, name: 'Мебель' },
  { id: 9, name: 'Декор' },
];

// Мок-данные товаров
const initialProducts = [
  {
    id: 1,
    name: 'Ламинат "Дуб Натуральный"',
    price: 1200,
    categoryId: 1,
    description: 'Ламинат 33 класса износостойкости, толщина 8 мм, фаска 4V, натуральный оттенок',
    image_url: 'https://images.unsplash.com/photo-1565538810643-b5bdb714032a?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Класс': '33', 'Толщина': '8 мм', 'Фаска': '4V' },
  },
  {
    id: 2,
    name: 'Паркетная доска "Венге"',
    price: 2500,
    categoryId: 1,
    description: 'Паркетная доска из массива дерева, темный оттенок венге, матовое покрытие',
    image_url: 'https://images.unsplash.com/photo-1581858726788-75bc0f6a952d?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Материал': 'Массив', 'Покрытие': 'Матовое', 'Толщина': '15 мм' },
  },
  {
    id: 3,
    name: 'Керамогранит "Мрамор Белый"',
    price: 1800,
    categoryId: 1,
    description: 'Керамогранит с имитацией мрамора, глянцевая поверхность, размер 60x60 см',
    image_url: 'https://images.unsplash.com/photo-1518458028785-8fbcd101ebb9?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Размер': '60x60 см', 'Поверхность': 'Глянцевая', 'PEI': 'IV' },
  },
  {
    id: 4,
    name: 'Обои "Скандинавский стиль"',
    price: 950,
    categoryId: 2,
    description: 'Флизелиновые обои в скандинавском стиле, светлый фон с геометрическим узором',
    image_url: 'https://images.unsplash.com/photo-1519710164239-da123dc03ef4?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Тип': 'Флизелиновые', 'Размер рулона': '1.06x10 м' },
  },
  {
    id: 5,
    name: 'Декоративная штукатурка "Венецианская"',
    price: 2200,
    categoryId: 2,
    description: 'Декоративная штукатурка с эффектом мрамора, для создания роскошного интерьера',
    image_url: 'https://images.unsplash.com/photo-1509644851169-2acc08aa25b5?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Расход': '1.5 кг/м²', 'Эффект': 'Мрамор', 'Цвет': 'Белый с серым' },
  },
  {
    id: 6,
    name: 'Краска интерьерная "Матовый шелк"',
    price: 1500,
    categoryId: 2,
    description: 'Интерьерная краска премиум-класса с эффектом матового шелка, моющаяся',
    image_url: 'https://images.unsplash.com/photo-1580500155991-4f3c0a0e2a05?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Объем': '5 л', 'Расход': '10 м²/л', 'Степень блеска': 'Матовая' },
  },
  {
    id: 7,
    name: 'Натяжной потолок "Сатин"',
    price: 650,
    categoryId: 3,
    description: 'Натяжной потолок сатинового типа, элегантный внешний вид, простота в уходе',
    image_url: 'https://images.unsplash.com/photo-1565538420870-da08ff96a207?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Тип': 'Сатиновый', 'Цена за': '1 м²', 'Цвет': 'Белый' },
  },
  {
    id: 8,
    name: 'Потолочная плитка "Классика"',
    price: 320,
    categoryId: 3,
    description: 'Потолочная плитка из полистирола с классическим узором, легкий монтаж',
    image_url: 'https://images.unsplash.com/photo-1513694203232-719a280e022f?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Материал': 'Полистирол', 'Размер': '50x50 см', 'Толщина': '8 мм' },
  },
  {
    id: 9,
    name: 'Дверь межкомнатная "Модерн"',
    price: 7500,
    categoryId: 4,
    description: 'Межкомнатная дверь в современном стиле, шпон натурального дерева, со стеклом',
    image_url: 'https://images.unsplash.com/photo-1549781726-3f5a6e6d3355?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Материал': 'Шпон', 'Цвет': 'Венге', 'Ширина': '80 см' },
  },
  {
    id: 10,
    name: 'Дверь входная "Стальной страж"',
    price: 25000,
    categoryId: 4,
    description: 'Входная металлическая дверь с повышенной шумоизоляцией и взломостойкостью',
    image_url: 'https://images.unsplash.com/photo-1506377295352-e3154d43ea9e?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Материал': 'Сталь', 'Толщина': '2 мм', 'Утепление': 'Минеральная вата' },
  },
  {
    id: 11,
    name: 'Окно ПВХ "Теплый дом"',
    price: 12000,
    categoryId: 5,
    description: 'Пластиковое окно с двухкамерным стеклопакетом, высокая теплоизоляция',
    image_url: 'https://images.unsplash.com/photo-1503708928676-1cb796a0891e?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Профиль': 'ПВХ', 'Стеклопакет': 'Двухкамерный', 'Фурнитура': 'Roto' },
  },
  {
    id: 12,
    name: 'Люстра "Хрустальная симфония"',
    price: 15000,
    categoryId: 6,
    description: 'Хрустальная люстра с подвесками, создает изысканную игру света',
    image_url: 'https://images.unsplash.com/photo-1543330732-9dc6c9845a62?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Материал': 'Хрусталь', 'Количество ламп': '6', 'Тип цоколя': 'E14' },
  },
  {
    id: 13,
    name: 'Точечные светильники "Орбита"',
    price: 850,
    categoryId: 6,
    description: 'Встраиваемые точечные светильники с поворотным механизмом, LED',
    image_url: 'https://images.unsplash.com/photo-1507723820574-2a979730dbe9?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Тип': 'LED', 'Мощность': '7 Вт', 'Цветовая температура': '4000K' },
  },
  {
    id: 14,
    name: 'Ванна акриловая "Комфорт"',
    price: 18000,
    categoryId: 7,
    description: 'Акриловая ванна прямоугольной формы, прочная и долговечная',
    image_url: 'https://images.unsplash.com/photo-1507652313519-d4e9174996dd?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Материал': 'Акрил', 'Размер': '170x70 см', 'Объем': '180 л' },
  },
  {
    id: 15,
    name: 'Смеситель для раковины "Водопад"',
    price: 5500,
    categoryId: 7,
    description: 'Смеситель с эффектом водопада, хромированное покрытие, керамический картридж',
    image_url: 'https://images.unsplash.com/photo-1584622650111-993a426fbf0a?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Материал': 'Латунь', 'Покрытие': 'Хром', 'Тип': 'Однорычажный' },
  },
  {
    id: 16,
    name: 'Диван угловой "Комфорт Люкс"',
    price: 45000,
    categoryId: 8,
    description: 'Угловой диван с функцией раскладывания, обивка из экокожи, ортопедический матрас',
    image_url: 'https://images.unsplash.com/photo-1555041469-a586c61ea9bc?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Материал': 'Экокожа', 'Механизм': 'Еврокнижка', 'Размер': '270x170 см' },
  },
  {
    id: 17,
    name: 'Шкаф-купе "Стильный"',
    price: 35000,
    categoryId: 8,
    description: 'Шкаф-купе с зеркальными дверями, вместительный, современный дизайн',
    image_url: 'https://images.unsplash.com/photo-1595428774223-ef52624120d2?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Материал': 'ЛДСП', 'Фасад': 'Зеркало', 'Размер': '200x60x240 см' },
  },
  {
    id: 18,
    name: 'Картина "Абстракция"',
    price: 8500,
    categoryId: 9,
    description: 'Картина в абстрактном стиле, холст, акрил, подойдет для современного интерьера',
    image_url: 'https://images.unsplash.com/photo-1549887534-1541e9326642?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Материал': 'Холст, акрил', 'Размер': '80x60 см', 'Рама': 'Деревянная' },
  },
  {
    id: 19,
    name: 'Ваза декоративная "Элегант"',
    price: 3200,
    categoryId: 9,
    description: 'Декоративная ваза из керамики, ручная роспись, станет акцентом в интерьере',
    image_url: 'https://images.unsplash.com/photo-1612196808214-b7e239e5f6b7?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Материал': 'Керамика', 'Высота': '35 см', 'Стиль': 'Современный' },
  },
  {
    id: 20,
    name: 'Зеркало настенное "Венеция"',
    price: 7800,
    categoryId: 9,
    description: 'Настенное зеркало в декоративной раме, итальянский стиль',
    image_url: 'https://images.unsplash.com/photo-1618220252344-8ec99ec624b1?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Форма': 'Овальная', 'Размер': '80x120 см', 'Рама': 'Полиуретан, золото' },
  },
];

const ADMIN_KEY = 'V2patTbDXS1wuqbqpyZGwg2vq70cem2wk3ElHO6y9l2FhfgNfN';

function useQuery() {
  return new URLSearchParams(useLocation().search);
}

function AdminPage({ categories, setCategories, products, setProducts, styles, setStyles }) {
  const query = useQuery();
  const key = query.get('key');
  const [catName, setCatName] = useState('');
  const [prod, setProd] = useState({ name: '', price: '', categoryId: categories[0]?.id || 1, description: '', image_url: '', attributes: '' });
  const [adminToken, setAdminToken] = useState(tokenUtils.get());
  const [isLoading, setIsLoading] = useState(false);
  const [message, setMessage] = useState('');

  // Функция для получения токена админа
  const getAdminToken = async () => {
    if (adminToken) return adminToken;
    
    try {
      const response = await authAPI.login(ADMIN_KEY);
      const token = response.token;
      tokenUtils.save(token);
      setAdminToken(token);
      return token;
    } catch (error) {
      console.error('Ошибка получения токена:', error);
      setMessage('Ошибка авторизации');
      return null;
    }
  };

  // Функция для показа сообщений
  const showMessage = (msg, isError = false) => {
    setMessage(msg);
    setTimeout(() => setMessage(''), 3000);
  };

  if (key !== ADMIN_KEY) {
    return (
      <div className="Configurator" style={{ maxWidth: 600, margin: '100px auto' }}>
        <div style={{ 
          padding: 40, 
          textAlign: 'center', 
          color: '#f8fafc', 
          fontSize: 24,
          background: 'rgba(185, 28, 28, 0.1)',
          borderRadius: '12px',
          border: '1px solid rgba(185, 28, 28, 0.3)',
          boxShadow: '0 10px 25px rgba(185, 28, 28, 0.15)'
        }}>
          <i className="fas fa-lock" style={{ fontSize: 48, marginBottom: 20, color: '#b91c1c' }}></i>
          <div>Доступ запрещён</div>
        </div>
      </div>
    );
  }

  const addCategory = async () => {
    if (!catName.trim()) return;
    
    setIsLoading(true);
    try {
      const token = await getAdminToken();
      if (!token) return;
      
      await categoryAPI.create({ name: catName }, token);
      
      // Обновляем локальный список категорий
      const newCategory = { id: Date.now(), name: catName };
      setCategories([...categories, newCategory]);
      setCatName('');
      showMessage('Категория успешно добавлена');
    } catch (error) {
      console.error('Ошибка создания категории:', error);
      showMessage('Ошибка при создании категории', true);
    } finally {
      setIsLoading(false);
    }
  };

  const addProduct = async () => {
    if (!prod.name.trim() || !prod.price || !prod.categoryId) return;
    
    setIsLoading(true);
    try {
      const token = await getAdminToken();
      if (!token) return;
      
      // Подготавливаем данные для API
      const productData = {
        name: prod.name,
        price: Number(prod.price),
        category_id: Number(prod.categoryId),
        description: prod.description,
        image_url: prod.image_url,
        attributes: prod.attributes
          ? prod.attributes.split(';').map(attr => {
              const [name, value] = attr.split(':').map(s => s.trim());
              return {
                attribute_id: 1, // Временно используем ID 1, в реальном приложении нужно получать из API
                value: `${name}: ${value}`
              };
            })
          : []
      };
      
      const response = await productAPI.create(productData, token);
      
      // Обновляем локальный список товаров
      const newProduct = {
        id: response.id || Date.now(),
        name: prod.name,
        price: Number(prod.price),
        categoryId: Number(prod.categoryId),
        description: prod.description,
        image_url: prod.image_url,
        attributes: prod.attributes
          ? Object.fromEntries(prod.attributes.split(';').map((a) => a.split(':').map((s) => s.trim())))
          : {},
      };
      
      setProducts([...products, newProduct]);
      setProd({ name: '', price: '', categoryId: categories[0]?.id || 1, description: '', image_url: '', attributes: '' });
      showMessage('Товар успешно добавлен');
    } catch (error) {
      console.error('Ошибка создания товара:', error);
      showMessage('Ошибка при создании товара', true);
    } finally {
      setIsLoading(false);
    }
  };

  const removeCategory = async (id) => {
    setIsLoading(true);
    try {
      const token = await getAdminToken();
      if (!token) return;
      
      await categoryAPI.delete(id, token);
      
      // Обновляем локальные данные
      setCategories(categories.filter((c) => c.id !== id));
      setProducts(products.filter((p) => p.categoryId !== id));
      showMessage('Категория успешно удалена');
    } catch (error) {
      console.error('Ошибка удаления категории:', error);
      showMessage('Ошибка при удалении категории', true);
    } finally {
      setIsLoading(false);
    }
  };
  
  const removeProduct = async (id) => {
    setIsLoading(true);
    try {
      const token = await getAdminToken();
      if (!token) return;
      
      await productAPI.delete(id, token);
      
      // Обновляем локальные данные
      setProducts(products.filter((p) => p.id !== id));
      showMessage('Товар успешно удален');
    } catch (error) {
      console.error('Ошибка удаления товара:', error);
      showMessage('Ошибка при удалении товара', true);
    } finally {
      setIsLoading(false);
    }
  };

  const inputStyle = {
    padding: '12px 16px',
    borderRadius: '10px',
    border: '1px solid #334155',
    background: 'rgba(15, 23, 42, 0.6)',
    color: '#f1f5f9',
    fontSize: '1rem',
    width: '100%',
    transition: 'all 0.3s ease',
    boxShadow: '0 4px 10px rgba(0,0,0,0.1)',
    outline: 'none'
  };

  const buttonStyle = {
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
    letterSpacing: '0.5px'
  };

  const deleteButtonStyle = {
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
  };

  return (
    <div className="Configurator" style={{ maxWidth: 1000 }}>
      <h1>АДМИН-ПАНЕЛЬ</h1>
      
      {/* Индикатор загрузки */}
      {isLoading && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          background: 'rgba(0, 0, 0, 0.5)',
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          zIndex: 1000
        }}>
          <div style={{
            background: '#1e293b',
            padding: '20px',
            borderRadius: '10px',
            color: '#f1f5f9',
            fontSize: '1.2rem'
          }}>
            Загрузка...
          </div>
        </div>
      )}
      
      {/* Сообщения */}
      {message && (
        <div style={{
          position: 'fixed',
          top: '20px',
          right: '20px',
          background: message.isError ? '#dc2626' : '#059669',
          color: 'white',
          padding: '15px 20px',
          borderRadius: '8px',
          zIndex: 1001,
          boxShadow: '0 4px 12px rgba(0, 0, 0, 0.3)'
        }}>
          {message.text}
        </div>
      )}
      
      <div className="AdminSection">
        <h2 style={{ 
          fontSize: '1.5rem', 
          color: '#f8fafc', 
          marginBottom: '20px', 
          position: 'relative',
          paddingBottom: '10px'
        }}>
          Категории
          <span style={{ 
            position: 'absolute', 
            bottom: 0, 
            left: 0, 
            width: '60px', 
            height: '3px', 
            background: 'linear-gradient(90deg, #3b82f6, #60a5fa)', 
            borderRadius: '2px' 
          }}></span>
        </h2>
        
        <div style={{ display: 'flex', gap: 12, marginBottom: 24, alignItems: 'center' }}>
          <input 
            value={catName} 
            onChange={e => setCatName(e.target.value)} 
            placeholder="Новая категория" 
            style={inputStyle} 
          />
          <button 
            onClick={addCategory} 
            style={buttonStyle}
          >
            Добавить
          </button>
        </div>
        
        <ul style={{ 
          marginBottom: 30, 
          listStyle: 'none', 
          padding: 0,
          background: 'rgba(30, 41, 59, 0.5)',
          borderRadius: '12px',
          overflow: 'hidden',
          border: '1px solid #334155'
        }}>
          {categories.map((c) => (
            <li key={c.id} style={{ 
              padding: '14px 20px', 
              borderBottom: '1px solid rgba(51, 65, 85, 0.5)',
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              transition: 'all 0.3s ease'
            }}>
              <span style={{ fontSize: '1.1rem', fontWeight: 500 }}>{c.name}</span>
              <button 
                onClick={() => removeCategory(c.id)} 
                style={deleteButtonStyle}
              >
                <i className="fas fa-trash-alt" style={{ marginRight: '6px' }}></i>
                Удалить
              </button>
            </li>
          ))}
        </ul>
      </div>
      
      <div className="AdminSection" style={{ marginTop: 40 }}>
        <h2 style={{ 
          fontSize: '1.5rem', 
          color: '#f8fafc', 
          marginBottom: '20px', 
          position: 'relative',
          paddingBottom: '10px'
        }}>
          Товары
          <span style={{ 
            position: 'absolute', 
            bottom: 0, 
            left: 0, 
            width: '60px', 
            height: '3px', 
            background: 'linear-gradient(90deg, #3b82f6, #60a5fa)', 
            borderRadius: '2px' 
          }}></span>
        </h2>
        
        <div style={{ 
          display: 'grid', 
          gridTemplateColumns: '1fr 1fr', 
          gap: 16, 
          marginBottom: 24,
          background: 'rgba(30, 41, 59, 0.5)',
          padding: '20px',
          borderRadius: '12px',
          border: '1px solid #334155'
        }}>
          <input 
            value={prod.name} 
            onChange={e => setProd({ ...prod, name: e.target.value })} 
            placeholder="Название" 
            style={inputStyle} 
          />
          <input 
            value={prod.price} 
            onChange={e => setProd({ ...prod, price: e.target.value })} 
            placeholder="Цена" 
            type="number" 
            style={inputStyle} 
          />
          <select 
            value={prod.categoryId} 
            onChange={e => setProd({ ...prod, categoryId: e.target.value })} 
            style={{
              ...inputStyle,
              appearance: 'none',
              backgroundImage: 'url("data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%2712%27 height=%278%27 viewBox=%270 0 12 8%27%3E%3Cpath fill=%27%2360a5fa%27 d=%27M10.6.6L6 5.2 1.4.6.6 1.4 6 6.8l5.4-5.4z%27/%3E%3C/svg%3E")',
              backgroundRepeat: 'no-repeat',
              backgroundPosition: 'right 16px center',
              paddingRight: '40px'
            }}
          >
            {categories.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
          </select>
          <input 
            value={prod.image_url} 
            onChange={e => setProd({ ...prod, image_url: e.target.value })} 
            placeholder="URL картинки" 
            style={inputStyle} 
          />
          <input 
            value={prod.description} 
            onChange={e => setProd({ ...prod, description: e.target.value })} 
            placeholder="Описание" 
            style={inputStyle} 
          />
          <input 
            value={prod.attributes} 
            onChange={e => setProd({ ...prod, attributes: e.target.value })} 
            placeholder="Атрибуты (пример: Размер:60x60; Толщина:9 мм)" 
            style={inputStyle} 
          />
          <button 
            onClick={addProduct} 
            style={{
              ...buttonStyle,
              gridColumn: '1 / -1',
              marginTop: '10px',
              padding: '14px'
            }}
          >
            <i className="fas fa-plus" style={{ marginRight: '8px' }}></i>
            Добавить товар
          </button>
        </div>
        
        <ul style={{ 
          marginBottom: 30, 
          listStyle: 'none', 
          padding: 0,
          background: 'rgba(30, 41, 59, 0.5)',
          borderRadius: '12px',
          overflow: 'hidden',
          border: '1px solid #334155',
          maxHeight: '400px',
          overflowY: 'auto'
        }}>
          {products.map((p) => (
            <li key={p.id} style={{ 
              padding: '14px 20px', 
              borderBottom: '1px solid rgba(51, 65, 85, 0.5)',
              display: 'flex',
              alignItems: 'center',
              gap: 16,
              transition: 'all 0.3s ease'
            }}>
              <img 
                src={p.image_url} 
                alt={p.name} 
                style={{ 
                  width: 60, 
                  height: 60, 
                  objectFit: 'cover', 
                  borderRadius: 10, 
                  background: '#1e293b',
                  border: '1px solid #334155',
                  boxShadow: '0 4px 10px rgba(0,0,0,0.2)'
                }} 
              />
              <div style={{ flex: 1 }}>
                <div style={{ fontWeight: 600, fontSize: '1.1rem', marginBottom: '4px' }}>{p.name}</div>
                <div style={{ color: '#94a3b8', fontSize: '0.9rem' }}>
                  {categories.find(c => c.id === p.categoryId)?.name || '—'}
                </div>
              </div>
              <div style={{ 
                fontWeight: 700, 
                fontSize: '1.2rem', 
                color: '#38bdf8',
                display: 'flex',
                alignItems: 'center',
                gap: '6px'
              }}>
                <span style={{ 
                  width: '8px', 
                  height: '8px', 
                  background: '#38bdf8', 
                  borderRadius: '50%',
                  display: 'inline-block'
                }}></span>
                {p.price} ₽
              </div>
              <button 
                onClick={() => removeProduct(p.id)} 
                style={deleteButtonStyle}
              >
                <i className="fas fa-trash-alt" style={{ marginRight: '6px' }}></i>
                Удалить
              </button>
            </li>
          ))}
        </ul>
      </div>
      
      <StyleAdmin products={products} styles={styles} setStyles={setStyles} />
    </div>
  );
}

// Компонент шапки сайта
function Header() {
  return (
    <header style={{
      background: '#0f172a',
      padding: '15px 30px',
      display: 'flex',
      justifyContent: 'space-between',
      alignItems: 'center',
      borderBottom: '1px solid #334155',
      boxShadow: '0 4px 20px rgba(0,0,0,0.1)',
      position: 'sticky',
      top: 0,
      zIndex: 100
    }}>
      <div style={{ display: 'flex', alignItems: 'center' }}>
        <div style={{ 
          fontSize: '1.8rem', 
          fontWeight: 800, 
          color: '#f1f5f9',
          display: 'flex',
          alignItems: 'center',
          gap: '10px'
        }}>
          <i className="fas fa-home" style={{ color: '#3b82f6' }}></i>
          <span>ROYAL<span style={{ color: '#3b82f6' }}>INTERIORS</span></span>
        </div>
      </div>
      <nav>
        <ul style={{ 
          display: 'flex', 
          gap: '25px', 
          listStyle: 'none',
          margin: 0,
          padding: 0
        }}>
          <li><a href="#" style={{ color: '#60a5fa', textDecoration: 'none', fontWeight: 600 }}>ДИЗАЙН ИНТЕРЬЕРА</a></li>
          <li><a href="#" style={{ color: '#f1f5f9', textDecoration: 'none' }}>РЕМОНТ</a></li>
          <li><a href="#" style={{ color: '#f1f5f9', textDecoration: 'none' }}>МАТЕРИАЛЫ</a></li>
          <li><a href="#" style={{ color: '#f1f5f9', textDecoration: 'none' }}>ПОРТФОЛИО</a></li>
          <li><a href="#" style={{ color: '#f1f5f9', textDecoration: 'none', display: 'flex', alignItems: 'center', gap: '5px' }}>
            <span>КОНФИГУРАТОР</span>
            <i className="fas fa-tools" style={{ color: '#3b82f6' }}></i>
          </a></li>
        </ul>
      </nav>
      <div style={{ display: 'flex', alignItems: 'center', gap: '20px' }}>
        <div style={{ position: 'relative' }}>
          <i className="fas fa-shopping-cart" style={{ color: '#f1f5f9', fontSize: '1.2rem' }}></i>
          <span style={{ 
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
            fontWeight: 'bold'
          }}>2</span>
        </div>
        <i className="fas fa-user" style={{ color: '#f1f5f9', fontSize: '1.2rem' }}></i>
      </div>
    </header>
  );
}

function MainPage({ categories, products, styles }) {
  const [selectedCategory, setSelectedCategory] = useState(categories[0].id);
  const [selectedProducts, setSelectedProducts] = useState([]);

  const handleCategoryChange = (id) => setSelectedCategory(id);

  const handleProductToggle = (product) => {
    setSelectedProducts((prev) =>
      prev.includes(product.id)
        ? prev.filter((pid) => pid !== product.id)
        : [...prev, product.id]
    );
  };

  const handleStyleSelect = (style) => {
    setSelectedProducts((prev) => Array.from(new Set([...prev, ...style.productIds])));
  };

  const total = selectedProducts.reduce(
    (sum, pid) => sum + (products.find((p) => p.id === pid)?.price || 0),
    0
  );

  const addedProducts = products.filter((p) => selectedProducts.includes(p.id));

  return (
    <>
      <Header />
      <div className="Configurator">
        <h1>КОНФИГУРАТОР ПОМЕЩЕНИЙ</h1>
        
        <div style={{ 
          marginBottom: 40, 
          display: 'flex', 
          flexDirection: 'column', 
          gap: '10px',
          color: '#94a3b8',
          fontSize: '0.95rem',
          maxWidth: '800px',
          margin: '0 auto 40px auto',
          textAlign: 'center',
          lineHeight: '1.6'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px', justifyContent: 'center' }}>
            <i className="fas fa-check-circle" style={{ color: '#3b82f6' }}></i>
            <span>Онлайн-конфигуратор помещений с подбором совместимых материалов</span>
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px', justifyContent: 'center' }}>
            <i className="fas fa-check-circle" style={{ color: '#3b82f6' }}></i>
            <span>Создайте дизайн своей мечты с профессиональными материалами</span>
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px', justifyContent: 'center' }}>
            <i className="fas fa-check-circle" style={{ color: '#3b82f6' }}></i>
            <span>Выберите стиль интерьера и подберите подходящие материалы</span>
          </div>
        </div>

        <StyleSelector styles={styles} onSelect={handleStyleSelect} />

      <div className="Categories">
        {categories.map((cat) => (
          <button
            key={cat.id}
            className={cat.id === selectedCategory ? 'active' : ''}
            onClick={() => handleCategoryChange(cat.id)}
          >
            {cat.name}
          </button>
        ))}
      </div>
      <div className="MainContent">
        <div className="Products">
          {products
            .filter((p) => p.categoryId === selectedCategory)
            .map((p) => (
              <div key={p.id} className="ProductCard">
                <img src={p.image_url} alt={p.name} className="ProductImg" />
                <div className="ProductInfo">
                  <div className="ProductName">{p.name}</div>
                  <div className="ProductDesc">{p.description}</div>
                  <div className="ProductAttrs">
                    {Object.entries(p.attributes).map(([k, v]) => (
                      <span key={k} className="Attr">{k}: {v}</span>
                    ))}
                  </div>
                  <div className="ProductPrice">{p.price} ₽</div>
                </div>
                <button
                  className={selectedProducts.includes(p.id) ? 'selected' : ''}
                  onClick={() => handleProductToggle(p)}
                >
                  {selectedProducts.includes(p.id) ? 'Убрать' : 'Добавить'}
                </button>
              </div>
            ))}
        </div>
        <div className="AddedPanel">
          <h2>Выбранные материалы</h2>
          {addedProducts.length === 0 ? (
            <div className="Empty">Нет выбранных материалов</div>
          ) : (
            <ul>
              {addedProducts.map((p) => (
                <li key={p.id} className="AddedItem">
                  <img src={p.image_url} alt={p.name} className="AddedImg" />
                  <div>
                    <div className="AddedName">{p.name}</div>
                    <div className="AddedPrice">{p.price} ₽</div>
                  </div>
                  <button onClick={() => handleProductToggle(p)}>Убрать</button>
                </li>
              ))}
            </ul>
          )}
          <div className="TotalPanel">
            <span>Итого: </span>
            <b>{total} ₽</b>
          </div>
          <button style={{
            background: 'linear-gradient(135deg, #3b82f6, #2563eb)',
            color: '#fff',
            border: 'none',
            borderRadius: '10px',
            padding: '16px 30px',
            fontSize: '1.1rem',
            fontWeight: '600',
            cursor: 'pointer',
            transition: 'all 0.3s ease',
            marginTop: '20px',
            width: '100%',
            boxShadow: '0 4px 15px rgba(37, 99, 235, 0.3)',
            textTransform: 'uppercase',
            letterSpacing: '1px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            gap: '10px'
          }}>
            <i className="fas fa-shopping-cart"></i>
            Добавить в корзину
          </button>
        </div>
      </div>
      
      {/* Информационный блок */}
      <div style={{
        background: '#1e293b',
        borderRadius: '16px',
        padding: '30px',
        marginTop: '50px',
        border: '1px solid #334155',
        boxShadow: '0 8px 24px rgba(0,0,0,0.15)',
      }}>
        <h2 style={{
          fontSize: '1.5rem',
          fontWeight: '700',
          color: '#f1f5f9',
          marginBottom: '20px',
          textAlign: 'center',
          position: 'relative',
          paddingBottom: '15px'
        }}>
          Популярные стили интерьера
          <div style={{
            position: 'absolute',
            bottom: '0',
            left: '50%',
            transform: 'translateX(-50%)',
            width: '60px',
            height: '3px',
            background: '#60a5fa',
            borderRadius: '2px'
          }}></div>
        </h2>
        
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))',
          gap: '20px',
          marginTop: '30px'
        }}>
          {['Скандинавский', 'Лофт', 'Минимализм', 'Классический', 'Современный', 'Прованс', 'Хай-тек', 'Эко-стиль', 'Индустриальный', 'Японский'].map(style => (
            <div key={style} style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              padding: '12px 16px',
              background: '#0f172a',
              borderRadius: '10px',
              border: '1px solid #334155'
            }}>
              <span style={{ color: '#f1f5f9', fontWeight: '500' }}>{style}</span>
              <span style={{ color: '#38bdf8', fontWeight: '600' }}>Подробнее</span>
            </div>
          ))}
        </div>
      </div>
    </div>
    
    {/* Подвал сайта */}
    <footer style={{
      background: '#0f172a',
      padding: '40px 30px',
      marginTop: '60px',
      borderTop: '1px solid #334155',
      color: '#94a3b8'
    }}>
      <div style={{
        maxWidth: '1200px',
        margin: '0 auto',
        display: 'flex',
        justifyContent: 'space-between',
        flexWrap: 'wrap',
        gap: '30px'
      }}>
        <div style={{ flex: '1', minWidth: '250px' }}>
          <div style={{ 
            fontSize: '1.5rem', 
            fontWeight: '800', 
            color: '#f1f5f9',
            display: 'flex',
            alignItems: 'center',
            gap: '10px',
            marginBottom: '15px'
          }}>
            <i className="fas fa-home" style={{ color: '#3b82f6' }}></i>
            <span>ROYAL<span style={{ color: '#3b82f6' }}>INTERIORS</span></span>
          </div>
          <p style={{ lineHeight: '1.6', marginBottom: '20px' }}>Профессиональный дизайн и ремонт помещений. Индивидуальный подход к каждому клиенту. Гарантия качества на все работы.</p>
          <div style={{ display: 'flex', gap: '15px' }}>
            <i className="fab fa-vk" style={{ color: '#3b82f6', fontSize: '1.3rem' }}></i>
            <i className="fab fa-telegram" style={{ color: '#3b82f6', fontSize: '1.3rem' }}></i>
            <i className="fab fa-youtube" style={{ color: '#3b82f6', fontSize: '1.3rem' }}></i>
          </div>
        </div>
        
        <div style={{ flex: '1', minWidth: '250px' }}>
          <h3 style={{ color: '#f1f5f9', marginBottom: '15px', fontSize: '1.1rem' }}>Каталог</h3>
          <ul style={{ listStyle: 'none', padding: '0', display: 'flex', flexDirection: 'column', gap: '10px' }}>
            <li><a href="#" style={{ color: '#94a3b8', textDecoration: 'none', transition: 'color 0.3s' }}>Дизайн интерьера</a></li>
            <li><a href="#" style={{ color: '#94a3b8', textDecoration: 'none', transition: 'color 0.3s' }}>Ремонт квартир</a></li>
            <li><a href="#" style={{ color: '#94a3b8', textDecoration: 'none', transition: 'color 0.3s' }}>Отделочные материалы</a></li>
            <li><a href="#" style={{ color: '#94a3b8', textDecoration: 'none', transition: 'color 0.3s' }}>Мебель на заказ</a></li>
            <li><a href="#" style={{ color: '#94a3b8', textDecoration: 'none', transition: 'color 0.3s' }}>Конфигуратор помещений</a></li>
          </ul>
        </div>
        
        <div style={{ flex: '1', minWidth: '250px' }}>
          <h3 style={{ color: '#f1f5f9', marginBottom: '15px', fontSize: '1.1rem' }}>Контакты</h3>
          <ul style={{ listStyle: 'none', padding: '0', display: 'flex', flexDirection: 'column', gap: '10px' }}>
            <li style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              <i className="fas fa-map-marker-alt" style={{ color: '#3b82f6', width: '20px' }}></i>
              <span>г. Москва, ул. Дизайнерская, 42</span>
            </li>
            <li style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              <i className="fas fa-phone" style={{ color: '#3b82f6', width: '20px' }}></i>
              <span>+7 (999) 123-45-67</span>
            </li>
            <li style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              <i className="fas fa-envelope" style={{ color: '#3b82f6', width: '20px' }}></i>
              <span>info@royalinteriors.ru</span>
            </li>
            <li style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              <i className="fas fa-clock" style={{ color: '#3b82f6', width: '20px' }}></i>
              <span>Пн-Пт: 9:00-20:00, Сб-Вс: 10:00-18:00</span>
            </li>
          </ul>
        </div>
      </div>
      <div style={{ borderTop: '1px solid #334155', marginTop: '30px', paddingTop: '20px', textAlign: 'center' }}>
        © 2023 ROYAL INTERIORS. Все права защищены.
      </div>
    </footer>
    </>
  );
}

function App() {
  const [categories, setCategories] = useState([]);
  const [products, setProducts] = useState([]);
  const [styles, setStyles] = useState(initialStyles);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Загрузка данных при монтировании компонента
  useEffect(() => {
    const loadData = async () => {
      try {
        setLoading(true);
        setError(null);
        
        // Загружаем категории
        const categoriesData = await categoryAPI.getAll();
        setCategories(categoriesData || []);
        
        // Загружаем все товары по категориям
        const allProducts = [];
        for (const category of categoriesData || []) {
          try {
            const categoryProducts = await productAPI.getByCategory(category.id);
            allProducts.push(...(categoryProducts || []));
          } catch (err) {
            console.warn(`Не удалось загрузить товары для категории ${category.id}:`, err);
          }
        }
        setProducts(allProducts);
        
        // Загружаем пресеты (стили)
        try {
          const presetsData = await presetAPI.getAllDetailed();
          if (presetsData && presetsData.length > 0) {
            // Преобразуем пресеты в формат стилей
            const stylesFromPresets = presetsData.map(preset => ({
              id: preset.preset_id,
              name: preset.name,
              description: preset.description,
              products: preset.items?.map(item => item.product?.id).filter(Boolean) || []
            }));
            setStyles([...initialStyles, ...stylesFromPresets]);
          }
        } catch (err) {
          console.warn('Не удалось загрузить пресеты:', err);
        }
        
      } catch (err) {
        console.error('Ошибка загрузки данных:', err);
        setError('Не удалось загрузить данные. Используются локальные данные.');
        // Fallback к локальным данным
        setCategories(initialCategories);
        setProducts(initialProducts);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

  if (loading) {
    return (
      <div className="App" style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        background: '#0f172a',
        color: '#f1f5f9'
      }}>
        <div style={{ textAlign: 'center' }}>
          <div style={{ 
            width: '50px', 
            height: '50px', 
            border: '3px solid #334155',
            borderTop: '3px solid #3b82f6',
            borderRadius: '50%',
            animation: 'spin 1s linear infinite',
            margin: '0 auto 20px'
          }}></div>
          <div>Загрузка данных...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="App">
      {error && (
        <div style={{
          background: 'rgba(239, 68, 68, 0.1)',
          border: '1px solid rgba(239, 68, 68, 0.3)',
          color: '#fca5a5',
          padding: '10px 20px',
          textAlign: 'center',
          fontSize: '14px'
        }}>
          {error}
        </div>
      )}
      <Routes>
        <Route path="/" element={<MainPage categories={categories} products={products} styles={styles} />} />
        <Route path="/admin" element={<AdminPage categories={categories} setCategories={setCategories} products={products} setProducts={setProducts} styles={styles} setStyles={setStyles} />} />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </div>
  );
}

export default App;