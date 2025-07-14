// Базовый URL API

// Используем относительный путь, так как запросы будут проксироваться через локальный сервер
const API_BASE_URL = 'https://dev.api.inspireforge.ru/api';

const apiRequest = async (endpoint, options = {}) => {
  const url = `${API_BASE_URL}${endpoint}`;
  const method = (options.method || 'GET').toUpperCase();
  const hasBody = ['POST', 'PUT', 'PATCH'].includes(method);
  const isForm = options.body instanceof FormData;
  const config = {
    ...options,
    headers: {
      Accept: 'application/json',
      ...(hasBody && !isForm ? { 'Content-Type': 'application/json' } : {}),
      ...(options.headers || {}),
    },
  };

  console.log(`[API] ➤ ${options.method || 'GET'} ${url}`);
  console.log('↳ Request config:', config);
  console.log(`[API] ➤ ${method} ${url}`, config);
  try {
    const response = await fetch(url, config);
    const rawText = await response.text();

    if (!response.ok) {
      let err;
      try {
        err = JSON.parse(rawText);
      } catch {
        err = { message: rawText || 'Unknown error' };
      }
      console.error(`[API ERROR] ${url}`, response.status, err);
      throw new Error(err.message || `HTTP ${response.status}`, { cause: err });
    }

    if (!rawText) return null; // 204 или пустой body
    try {
      return JSON.parse(rawText);
    } catch {
      // обычный JSON
      return rawText;
    } // plain text / html
  } catch (e) {
    console.error('✗ Запрос упал:', e.message);
    throw e;
  }
};

// Функции для работы с категориями
export const categoryAPI = {
  // Получить все категории
  getAll: () => apiRequest('/category'),

  // Получить категорию по ID
  getById: (id) => apiRequest(`/category/${id}`),

  // Создать категорию
  create: (data, token) =>
  apiRequest('/admin/category', {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  }),

  // Обновить категорию
  update: (id, data, token) =>
  apiRequest(`/admin/category/${id}`, {
    method: 'PUT',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  }),

  // Удалить категорию
  delete: (id, token) =>
  apiRequest(`/admin/category/${id}`, {
    method: 'DELETE',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  }),

  // 🔥 Получить все категории с parent_id
  getAllWithParents: async () => {
    const flat = await apiRequest('/category');
    console.log('🔵 ШАГ 1: flat', flat);

    const detailed = await Promise.all(
      flat.map(async (cat) => {
        try {
          const full = await apiRequest(`/category/${cat.id}`);
          console.log('   ↳ 🔹 деталь', full); // <‑‑ вот это покажет, есть ли parent_id
          return full;
        } catch (err) {
          console.error(`Не получил категорию ${cat.id}`, err);
          return cat;
        }
      }),
    );

    return detailed;
  },
};

// Функции для работы с атрибутами категорий
export const categoryAttributeAPI = {
  // Получить все атрибуты категории
  getAll: (categoryID) => apiRequest(`/category/${categoryID}/attribute`),

  // Получить конкретный атрибут
  getById: (categoryID, attributeID) =>
    apiRequest(`/category/${categoryID}/attribute/${attributeID}`),

  // Создать один атрибут
  create: (categoryID, attributeData, token) =>
    apiRequest(`/admin/category/${categoryID}/attribute`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(attributeData),
    }),

  // Создать несколько атрибутов сразу (batch)
  createBatch: (categoryID, attributesArray, token) =>
    apiRequest(`/admin/category/${categoryID}/attribute/batch`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(attributesArray),
    }),

  // Обновить атрибут (требует авторизации)
  update: (categoryID, attributeID, attributeData, token) =>
    apiRequest(`/admin/category/${categoryID}/attribute/${attributeID}`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(attributeData),
    }),

  // Удалить атрибут (требует авторизации)
  delete: (categoryID, attributeID, token) =>
    apiRequest(`/admin/category/${categoryID}/attribute/${attributeID}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }),
};

// Функции для работы с товарами
export const productAPI = {
  // Получить товар по ID
  getById: (id) => apiRequest(`/product/${id}`),
  getAll: () => apiRequest('/product'),
  // Получить товары по категории
  getByCategory: (categoryId) => apiRequest(`/product/category/${categoryId}`),

  // Создать товар (требует авторизации)
  create: (productData, token) =>
    apiRequest('/admin/product', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(productData),
    }),

  // Создать товар с атрибутами
  createDetailed: (productData, token) =>
    apiRequest('/admin/product/detailed', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(productData),
    }),

  // Обновить товар (требует авторизации)
  update: (id, productData, token) =>
    apiRequest(`/admin/product/${id}`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(productData),
    }),

  // Удалить товар (требует авторизации)
  delete: (id, token) =>
    apiRequest(`/admin/product/${id}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }),
};

// Функции для работы с коэффициентами
export const coefficientAPI = {
  // Получить все коэффициенты
  getAll: (token) =>
    apiRequest('/admin/coefficients', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }),
  
  // Получить коэффициент по ID
  getById: (id) => apiRequest(`/admin/coefficients/${id}`),
  
  // Создать коэффициент
  create: (data, token) =>
    apiRequest('/admin/coefficients', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    }),
  
  // Обновить коэффициент
  update: (id, data, token) =>
    apiRequest(`/admin/coefficients/${id}`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    }),
  
  // Удалить коэффициент
  delete: (id, token) =>
    apiRequest(`/admin/coefficients/${id}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }),
};

// Функции для работы с пресетами (стилями)
export const presetAPI = {
  create: async (data, token) =>
    await apiRequest('/admin/presets', {
      method: 'POST',
      headers: { 
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    }),
  delete: async (id, token) =>
    await apiRequest(`/admin/presets/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` },
    }),
  update: async (id, data, token) =>
    await apiRequest(`/admin/presets/${id}`, {
      method: 'PUT',
      headers: { 
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    }),

  list: async () => await apiRequest('/presets'),
  getAllDetailed: async () => await apiRequest('/presets/detailed'),
  getById: async (id) => await apiRequest(`/presets/${id}`),
};

// Функции для авторизации
export const authAPI = {
  // Получить токен для админа
  login: (secretKey) =>
    apiRequest(`/admin/auth/${secretKey}`, {
      method: 'GET',
    }),
};

// Утилиты для работы с токеном
export const tokenUtils = {
  // Сохранить токен в localStorage
  save: (token) => {
    localStorage.setItem('admin_token', token);
  },

  // Получить токен из localStorage
  get: () => {
    return localStorage.getItem('admin_token');
  },

  // Удалить токен из localStorage
  remove: () => {
    localStorage.removeItem('admin_token');
  },

  // Проверить, есть ли токен
  exists: () => {
    return !!localStorage.getItem('admin_token');
  },
};
