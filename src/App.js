import React, { useState } from 'react';
import { Routes, Route, useLocation, Navigate } from 'react-router-dom';
import './App.css';

import { initialStyles } from './data/styles';
import StyleSelector from './components/StyleSelector';
import StyleAdmin from './components/StyleAdmin';

const initialCategories = [
  { id: 1, name: 'Керамогранит' },
  { id: 2, name: 'Обои' },
  { id: 3, name: 'Краска' },
  { id: 4, name: 'Ламинат' },
  { id: 5, name: 'Плинтус' },
  { id: 6, name: 'Потолочная плитка' },
];
// Мок-данные товаров
const initialProducts = [
  {
    id: 1,
    name: 'Керамогранит "Белый матовый"',
    price: 1200,
    categoryId: 1,
    description: 'Белый матовый керамогранит 60x60, износостойкость PEI IV, толщина 9 мм',
    image_url: 'https://images.unsplash.com/photo-1506744038136-46273834b3fb?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Размер': '60x60', 'Толщина': '9 мм', 'PEI': 'IV' },
  },
  {
    id: 2,
    name: 'Керамогранит "Серый глянец"',
    price: 1350,
    categoryId: 1,
    description: 'Серый глянцевый керамогранит 60x60, износостойкость PEI III, толщина 8.5 мм',
    image_url: 'https://images.unsplash.com/photo-1464983953574-0892a716854b?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Размер': '60x60', 'Толщина': '8.5 мм', 'PEI': 'III' },
  },
  {
    id: 3,
    name: 'Обои "Сканди"',
    price: 800,
    categoryId: 2,
    description: 'Обои в скандинавском стиле, флизелиновые, рулон 1.06x10 м',
    image_url: 'https://images.unsplash.com/photo-1519125323398-675f0ddb6308?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Тип': 'Флизелин', 'Размер': '1.06x10 м' },
  },
  {
    id: 4,
    name: 'Краска Dulux белая',
    price: 950,
    categoryId: 3,
    description: 'Интерьерная краска, 2.5л, матовая, моющаяся',
    image_url: 'https://images.unsplash.com/photo-1501594907352-04cda38ebc29?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Объем': '2.5 л', 'Тип': 'Матовая' },
  },
  {
    id: 5,
    name: 'Ламинат Tarkett 33 класс',
    price: 1100,
    categoryId: 4,
    description: 'Ламинат 33 класс, 8 мм, фаска 4V, дуб натуральный',
    image_url: 'https://images.unsplash.com/photo-1465101046530-73398c7f28ca?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Класс': '33', 'Толщина': '8 мм', 'Фаска': '4V' },
  },
  {
    id: 6,
    name: 'Плинтус МДФ белый',
    price: 350,
    categoryId: 5,
    description: 'Плинтус МДФ, белый, высота 80 мм, длина 2.4 м',
    image_url: 'https://images.unsplash.com/photo-1515378791036-0648a3ef77b2?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Высота': '80 мм', 'Длина': '2.4 м' },
  },
  {
    id: 7,
    name: 'Потолочная плитка "Классика"',
    price: 250,
    categoryId: 6,
    description: 'Потолочная плитка пенополистирол, 50x50 см, белая',
    image_url: 'https://images.unsplash.com/photo-1465101178521-c1a9136a3b99?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Размер': '50x50 см', 'Материал': 'Пенополистирол' },
  },
  {
    id: 8,
    name: 'Обои "Гео"',
    price: 950,
    categoryId: 2,
    description: 'Обои с геометрическим рисунком, виниловые, рулон 1.06x10 м',
    image_url: 'https://images.unsplash.com/photo-1465101046530-73398c7f28ca?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Тип': 'Винил', 'Размер': '1.06x10 м' },
  },
  {
    id: 9,
    name: 'Краска Tikkurila Euro 7',
    price: 1200,
    categoryId: 3,
    description: 'Краска для стен и потолков, 2.7л, шелковисто-матовая',
    image_url: 'https://images.unsplash.com/photo-1506744038136-46273834b3fb?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Объем': '2.7 л', 'Тип': 'Шелковисто-матовая' },
  },
  {
    id: 10,
    name: 'Ламинат Kronospan 32 класс',
    price: 900,
    categoryId: 4,
    description: 'Ламинат 32 класс, 7 мм, дуб серый',
    image_url: 'https://images.unsplash.com/photo-1519125323398-675f0ddb6308?auto=format&fit=crop&w=400&q=80',
    attributes: { 'Класс': '32', 'Толщина': '7 мм' },
  },
];

const ADMIN_KEY = 'admin123';

function useQuery() {
  return new URLSearchParams(useLocation().search);
}

function AdminPage({ categories, setCategories, products, setProducts, styles, setStyles }) {
  const query = useQuery();
  const key = query.get('key');
  const [catName, setCatName] = useState('');
  const [prod, setProd] = useState({ name: '', price: '', categoryId: categories[0]?.id || 1, description: '', image_url: '', attributes: '' });

  if (key !== ADMIN_KEY) {
    return <div style={{ padding: 40, textAlign: 'center', color: '#b91c1c', fontSize: 22 }}>Доступ запрещён</div>;
  }

  const addCategory = () => {
    if (!catName.trim()) return;
    setCategories([...categories, { id: Date.now(), name: catName }]);
    setCatName('');
  };

  const addProduct = () => {
    if (!prod.name.trim() || !prod.price || !prod.categoryId) return;
    setProducts([
      ...products,
      {
        id: Date.now(),
        name: prod.name,
        price: Number(prod.price),
        categoryId: Number(prod.categoryId),
        description: prod.description,
        image_url: prod.image_url,
        attributes: prod.attributes
          ? Object.fromEntries(prod.attributes.split(';').map((a) => a.split(':').map((s) => s.trim())))
          : {},
      },
    ]);
    setProd({ name: '', price: '', categoryId: categories[0]?.id || 1, description: '', image_url: '', attributes: '' });
  };

  const removeCategory = (id) => {
    setCategories(categories.filter((c) => c.id !== id));
    setProducts(products.filter((p) => p.categoryId !== id));
  };
  const removeProduct = (id) => setProducts(products.filter((p) => p.id !== id));

  return (
    <div className="Configurator" style={{ maxWidth: 900 }}>
      <h1>Админ-панель</h1>
      <h2>Категории</h2>
      <div style={{ display: 'flex', gap: 8, marginBottom: 16 }}>
        <input value={catName} onChange={e => setCatName(e.target.value)} placeholder="Новая категория" />
        <button onClick={addCategory}>Добавить</button>
      </div>
      <ul style={{ marginBottom: 24 }}>
        {categories.map((c) => (
          <li key={c.id} style={{ marginBottom: 4 }}>
            {c.name} <button onClick={() => removeCategory(c.id)} style={{ color: '#b91c1c' }}>Удалить</button>
          </li>
        ))}
      </ul>
      <h2>Товары</h2>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 8, marginBottom: 16 }}>
        <input value={prod.name} onChange={e => setProd({ ...prod, name: e.target.value })} placeholder="Название" />
        <input value={prod.price} onChange={e => setProd({ ...prod, price: e.target.value })} placeholder="Цена" type="number" />
        <select value={prod.categoryId} onChange={e => setProd({ ...prod, categoryId: e.target.value })}>
          {categories.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
        </select>
        <input value={prod.image_url} onChange={e => setProd({ ...prod, image_url: e.target.value })} placeholder="URL картинки" />
        <input value={prod.description} onChange={e => setProd({ ...prod, description: e.target.value })} placeholder="Описание" />
        <input value={prod.attributes} onChange={e => setProd({ ...prod, attributes: e.target.value })} placeholder="Атрибуты (пример: Размер:60x60; Толщина:9 мм)" />
        <button onClick={addProduct}>Добавить товар</button>
      </div>
      <ul>
        {products.map((p) => (
          <li key={p.id} style={{ marginBottom: 8, display: 'flex', alignItems: 'center', gap: 8 }}>
            <img src={p.image_url} alt={p.name} style={{ width: 40, height: 40, objectFit: 'cover', borderRadius: 6, background: '#eee' }} />
            <span>{p.name} ({categories.find(c => c.id === p.categoryId)?.name || '—'}) — {p.price} ₽</span>
            <button onClick={() => removeProduct(p.id)} style={{ color: '#b91c1c' }}>Удалить</button>
          </li>
        ))}
      </ul>
      <StyleAdmin products={products} styles={styles} setStyles={setStyles} />
    </div>
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
    <div className="Configurator">
      <h1>Конфигуратор помещений</h1>

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
          <h2>Добавленные товары</h2>
          {addedProducts.length === 0 ? (
            <div className="Empty">Нет выбранных товаров</div>
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
        </div>
      </div>
    </div>
  );
}

function App() {
  const [categories, setCategories] = useState(initialCategories);
  const [products, setProducts] = useState(initialProducts);
  const [styles, setStyles] = useState(initialStyles);

  return (
    <Routes>
      <Route path="/" element={<MainPage categories={categories} products={products} styles={styles} />} />
      <Route path="/admin" element={
        <AdminPage
          categories={categories}
          setCategories={setCategories}
          products={products}
          setProducts={setProducts}
          styles={styles}
          setStyles={setStyles}
        />} />
      <Route path="*" element={<Navigate to="/" />} />
    </Routes>
  );
}

export default App;