import React, { useState, useEffect } from 'react';
import { serviceAPI } from '../../services/api';

const ServicesManager = ({ getAdminToken, showMessage, styles }) => {
  const [services, setServices] = useState([]);
  const [form, setForm] = useState({ name: '', description: '', price: '' });
  const [editingId, setEditingId] = useState(null);

  const fetchServices = async () => {
    const token = await getAdminToken();
    if (!token) return;
    try {
      const data = await serviceAPI.getAll(token);
      setServices(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error('Ошибка загрузки услуг:', err);
      showMessage('Ошибка загрузки услуг', true);
    }
  };

  useEffect(() => { fetchServices(); }, []);

  const saveService = async () => {
    const token = await getAdminToken();
    if (!token) return;
    const payload = {
      name: form.name.trim(),
      description: form.description.trim(),
      price: Number(form.price),
    };
    try {
      if (editingId) {
        await serviceAPI.update(editingId, payload, token);
        showMessage('Услуга обновлена');
      } else {
        await serviceAPI.create(payload, token);
        showMessage('Услуга добавлена');
      }
      setForm({ name: '', description: '', price: '' });
      setEditingId(null);
      fetchServices();
    } catch (err) {
      console.error('Ошибка сохранения услуги:', err);
      showMessage('Ошибка сохранения услуги', true);
    }
  };

  const deleteService = async (id) => {
    if (!window.confirm('Удалить услугу?')) return;
    const token = await getAdminToken();
    try {
      await serviceAPI.remove(id, token);
      showMessage('Услуга удалена');
      fetchServices();
    } catch (err) {
      console.error('Ошибка удаления услуги:', err);
      showMessage('Ошибка удаления', true);
    }
  };

  const startEdit = (service) => {
    setForm({
      name: service.name,
      description: service.description,
      price: service.price
    });
    setEditingId(service.id);
  };

  return (
    <div style={{ marginTop: 40 }}>
      <h2>Управление услугами</h2>
      <div style={{ display: 'flex', gap: 10, marginBottom: 20 }}>
        <input
          type="text"
          placeholder="Название"
          value={form.name}
          onChange={(e) => setForm({ ...form, name: e.target.value })}
          style={styles.inputStyle}
        />
        <input
          type="text"
          placeholder="Описание"
          value={form.description}
          onChange={(e) => setForm({ ...form, description: e.target.value })}
          style={styles.inputStyle}
        />
        <input
          type="number"
          placeholder="Цена"
          value={form.price}
          onChange={(e) => setForm({ ...form, price: e.target.value })}
          style={styles.inputStyle}
        />
        <button onClick={saveService} style={styles.buttonStyle}>
          {editingId ? 'Сохранить' : 'Добавить'}
        </button>
      </div>

      <table style={{ width: '100%', color: '#f1f5f9' }}>
        <thead>
          <tr>
            <th>ID</th>
            <th>Название</th>
            <th>Описание</th>
            <th>Цена</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          {services.map((s) => (
            <tr key={s.id}>
              <td>{s.id}</td>
              <td>{s.name}</td>
              <td>{s.description}</td>
              <td>{s.price} ₽</td>
              <td>
                <button onClick={() => startEdit(s)} style={styles.buttonStyle}>✏️</button>
                <button onClick={() => deleteService(s.id)} style={styles.deleteButtonStyle}>Удалить</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default ServicesManager;
