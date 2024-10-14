import React, { createContext, useState, useContext, useCallback } from "react";
import ToastContainer from "../components/ToastContainer/ToastContainer";

// Создаем контекст
const ToastContext = createContext();

export const useToast = () => useContext(ToastContext);

// Провайдер для оборачивания всего приложения
export const ToastProvider = ({ children }) => {
  const [toasts, setToasts] = useState([]);

  const addToast = useCallback((message, type = "success", seconds = 5) => {
    const id = Math.random().toString(36).substr(2, 9); // Генерируем уникальный ID для каждого тоста
    setToasts((prevToasts) => [...prevToasts, { id, message, type }]);

    // Удаляем тост автоматически по истечении времени
    setTimeout(() => {
      removeToast(id);
    }, seconds * 1000);
  }, []);

  const removeToast = useCallback((id) => {
    setToasts((prevToasts) => prevToasts.filter((toast) => toast.id !== id));
  }, []);

  return (
    <ToastContext.Provider value={{ addToast, removeToast }}>
      <ToastContainer toasts={toasts} removeToast={removeToast} />
      {children}
    </ToastContext.Provider>
  );
};
