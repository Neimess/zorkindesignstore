import React, { useState, useEffect } from 'react';
import { coefficientAPI } from '../../services/api';

function CoefficientManager({ getAdminToken, showMessage, styles }) {
  const [coefficients, setCoefficients] = useState([]);
  const [loading, setLoading] = useState(false);
  const [form, setForm] = useState({
    name: '',
    value: '',
    description: '',
  });
  const [editingId, setEditingId] = useState(null);

  const { inputStyle, buttonStyle, deleteButtonStyle } = styles;
  const editBtnStyle = {
    ...deleteButtonStyle,
    background: 'rgba(34,197,94,.1)',
    color: '#4ade80',
  };

  // Загрузка коэффициентов при монтировании компонента
  useEffect(() => {
    fetchCoefficients();
  }, []);

  // Получение списка коэффициентов
  const fetchCoefficients = async () => {
    try {
      setLoading(true);
      const data = await coefficientAPI.getAll();
      setCoefficients(data);
    } catch (err) {
      console.error('Ошибка при загрузке коэффициентов:', err);
      showMessage('Ошибка при загрузке коэффициентов: ' + (err.message || 'Неизвестная ошибка'), true);
    } finally {
      setLoading(false);
    }
  };

  // Сброс формы
  const resetForm = () => {
    setForm({
      name: '',
      value: '',
      description: '',
    });
    setEditingId(null);
  };

  // Редактирование коэффициента
  const editCoefficient = (coef) => {
    setForm({
      name: coef.name,
      value: coef.value.toString(),
      description: coef.description || '',
    });
    setEditingId(coef.id);
  };

  // Сохранение коэффициента (создание или обновление)
  const saveCoefficient = async () => {
    if (!form.name.trim() || !form.value.trim()) {
      return showMessage('Заполните название и значение коэффициента', true);
    }

    try {
      const token = await getAdminToken();
      if (!token) return;

      const payload = {
        name: form.name.trim(),
        value: parseFloat(form.value),
        description: form.description.trim(),
      };

      if (editingId) {
        // Обновление существующего коэффициента
        const updated = await coefficientAPI.update(editingId, payload, token);
        setCoefficients((prev) =>
          prev.map((c) => (c.id === editingId ? updated : c))
        );
        showMessage('Коэффициент обновлен');
      } else {
        // Создание нового коэффициента
        const created = await coefficientAPI.create(payload, token);
        setCoefficients((prev) => [...prev, created]);
        showMessage('Коэффициент добавлен');
      }

      resetForm();
    } catch (err) {
      console.error('Ошибка при сохранении коэффициента:', err);
      showMessage('Ошибка при сохранении: ' + (err.message || 'Неизвестная ошибка'), true);
    }
  };

  // Удаление коэффициента
  const removeCoefficient = async (id) => {
    try {
      const token = await getAdminToken();
      if (!token) return;

      await coefficientAPI.delete(id, token);
      setCoefficients((prev) => prev.filter((c) => c.id !== id));
      if (editingId === id) resetForm();
      showMessage('Коэффициент удален');
    } catch (err) {
      console.error('Ошибка при удалении коэффициента:', err);
      showMessage('Ошибка при удалении: ' + (err.message || 'Неизвестная ошибка'), true);
    }
  };

  return (
    <div className="AdminSection" style={{ marginTop: 40 }}>
      <h2 style={{ fontSize: '1.5rem', color: '#f8fafc', marginBottom: 20 }}>
        Коэффициенты
        <span
          style={{
            display: 'block',
            width: 60,
            height: 3,
            marginTop: 4,
            background: 'linear-gradient(90deg,#3b82f6,#60a5fa)',
            borderRadius: 2,
          }}
        />
      </h2>

      {/* Форма добавления/редактирования */}
      <div
        style={{
          display: 'grid',
          gridTemplateColumns: '1fr 1fr 1fr auto',
          gap: 16,
          marginBottom: 24,
          background: 'rgba(30,41,59,0.5)',
          padding: 20,
          borderRadius: 12,
          border: '1px solid #334155',
        }}
      >
        <input
          value={form.name}
          onChange={(e) => setForm({ ...form, name: e.target.value })}
          placeholder="Название коэффициента"
          style={inputStyle}
        />
        <input
          value={form.value}
          onChange={(e) => setForm({ ...form, value: e.target.value })}
          placeholder="Значение"
          type="number"
          step="0.01"
          style={inputStyle}
        />
        <input
          value={form.description}
          onChange={(e) => setForm({ ...form, description: e.target.value })}
          placeholder="Описание (опционально)"
          style={inputStyle}
        />
        <div style={{ display: 'flex', gap: 8 }}>
          <button onClick={saveCoefficient} style={buttonStyle}>
            {editingId ? 'Обновить' : 'Добавить'}
          </button>
          {editingId && (
            <button onClick={resetForm} style={deleteButtonStyle}>
              Отмена
            </button>
          )}
        </div>
      </div>

      {/* Список коэффициентов */}
      {loading ? (
        <div style={{ textAlign: 'center', padding: 20 }}>Загрузка...</div>
      ) : (
        <ul
          style={{
            listStyle: 'none',
            padding: 0,
            background: 'rgba(30,41,59,0.5)',
            borderRadius: 12,
            overflow: 'hidden',
            border: '1px solid #334155',
          }}
        >
          {coefficients.length === 0 ? (
            <li style={{ padding: 20, textAlign: 'center' }}>
              Нет добавленных коэффициентов
            </li>
          ) : (
            coefficients.map((coef) => (
              <li
                key={coef.id}
                style={{
                  padding: '14px 20px',
                  borderBottom: '1px solid rgba(51,65,85,0.5)',
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                }}
              >
                <div>
                  <div style={{ fontSize: '1.1rem', fontWeight: 500 }}>
                    {coef.name}
                  </div>
                  <div style={{ fontSize: '0.9rem', color: '#94a3b8' }}>
                    Значение: {coef.value}
                    {coef.description && ` • ${coef.description}`}
                  </div>
                </div>
                <div style={{ display: 'flex', gap: 8 }}>
                  <button
                    onClick={() => editCoefficient(coef)}
                    style={editBtnStyle}
                  >
                    <i className="fas fa-edit" style={{ marginRight: 6 }}></i>
                    Изменить
                  </button>
                  <button
                    onClick={() => removeCoefficient(coef.id)}
                    style={deleteButtonStyle}
                  >
                    <i
                      className="fas fa-trash-alt"
                      style={{ marginRight: 6 }}
                    ></i>
                    Удалить
                  </button>
                </div>
              </li>
            ))
          )}
        </ul>
      )}
    </div>
  );
}

export default CoefficientManager;