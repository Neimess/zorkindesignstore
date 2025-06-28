// Базовый URL API
// Используем относительный путь, так как запросы будут проксироваться через локальный сервер
const API_BASE_URL = 'https://dev.api.inspireforge.ru/api';
const apiRequest = async (endpoint, options = {}) => {
  const url = `${API_BASE_URL}${endpoint}`;
  const config = {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  };

  console.log(`[API] ➤ GET ${url}`);
  console.log('↳ Request config:', config);

  try {
    const response = await fetch(url, config);

    console.log(`↳ Response status: ${response.status}`);

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: 'Unknown error' }));
      console.error(`[API ERROR] ${url} ➤ ${response.status}: ${errorData.message}`);
      throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
    }

    // Для DELETE запросов может не быть тела ответа
    if (response.status === 204) {
      console.log(`✓ Success: No Content`);
      return null;
    }

    const data = await response.json();
    console.log('✓ Success:', data);
    return data;
  } catch (error) {
    console.error('✗ Ошибка выполнения запроса:', error.message);
    throw error;
  }
};

// Функции для работы с категориями
export const categoryAPI = {
  // Получить все категории
  getAll: () => apiRequest('/category'),
  
  // Получить категорию по ID
  getById: (id) => apiRequest(`/category/${id}`),
  
  // Создать категорию (требует авторизации)
  create: (categoryData, token) => apiRequest('/admin/category', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify(categoryData),
  }),
  
  // Обновить категорию (требует авторизации)
  update: (id, categoryData, token) => apiRequest(`/admin/category/${id}`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify(categoryData),
  }),
  
  // Удалить категорию (требует авторизации)
  delete: (id, token) => apiRequest(`/admin/category/${id}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  }),
};
// Функции для работы с атрибутами категорий
export const categoryAttributeAPI = {
  // Получить все атрибуты категории
  getAll: (categoryID) => apiRequest(`/category/${categoryID}/attribute`),

  // Получить конкретный атрибут
  getById: (categoryID, attributeID) =>
    apiRequest(`/category/${categoryID}/attribute/${attributeID}`),

  // Создать один атрибут (требует авторизации)
  create: (categoryID, attributeData, token) =>
    apiRequest(`/admin/category/${categoryID}/attribute`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify(attributeData),
    }),

  // Создать несколько атрибутов сразу (batch)
  createBatch: (categoryID, attributesArray, token) =>
    apiRequest(`/admin/category/${categoryID}/attribute/batch`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify(attributesArray),
    }),

  // Обновить атрибут (требует авторизации)
  update: (categoryID, attributeID, attributeData, token) =>
    apiRequest(`/admin/category/${categoryID}/attribute/${attributeID}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify(attributeData),
    }),

  // Удалить атрибут (требует авторизации)
  delete: (categoryID, attributeID, token) =>
    apiRequest(`/admin/category/${categoryID}/attribute/${attributeID}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    }),
};




// Функции для работы с товарами
export const productAPI = {
  // Получить товар по ID
  getById: (id) => apiRequest(`/product/${id}`),
  
  // Получить товары по категории
  getByCategory: (categoryId) => apiRequest(`/product/category/${categoryId}`),
  
  // Создать товар (требует авторизации)
  create: (productData, token) => apiRequest('/admin/product', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify(productData),
  }),
  
  // Обновить товар (требует авторизации)
  update: (id, productData, token) => apiRequest(`/admin/product/${id}`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify(productData),
  }),
  
  // Удалить товар (требует авторизации)
  delete: (id, token) => apiRequest(`/admin/product/${id}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  }),
};

// Функции для работы с пресетами (стилями)
export const presetAPI = {
  // Получить все пресеты (краткая информация)
  getAll: () => apiRequest('/presets'),
  
  // Получить все пресеты с подробной информацией
  getAllDetailed: () => apiRequest('/presets/detailed'),
  
  // Получить пресет по ID
  getById: (id) => apiRequest(`/presets/${id}`),
  
  // Создать пресет (требует авторизации)
  create: (presetData, token) => apiRequest('/admin/presets', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify(presetData),
  }),
  
  // Удалить пресет (требует авторизации)
  delete: (id, token) => apiRequest(`/admin/presets/${id}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  }),
};

// Функции для авторизации
export const authAPI = {
  // Получить токен для админа
  login: (secretKey) => apiRequest(`/admin/auth/${secretKey}`, {
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