import React, { useState, createContext, useContext, useEffect } from "react";
import axios from "axios";
import * as rp from "../global/routerPaths/routerPaths.jsx";
import { useNavigate } from "react-router-dom";
import { useAuth } from "./authContext.jsx";

const ApiContext = createContext();

export function ApiProvider({ children }) {
  const { isAuth, logout, refresh, clearAuth } = useAuth();

  const navigate = useNavigate();

  const apiClient = axios.create({
    baseURL: "http://localhost:3000/api",
    withCredentials: true,
  });

  apiClient.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error.config;

      if (
        error.response &&
        error.response.status == 401 &&
        !originalRequest._retry
      ) {
        originalRequest._retry = true; // Это нужно для того, чтобы обработка ошибки не ушла в вечную рекурсию

        if (!isAuth) {
          navigate(rp.LOGIN_ROUTE);
          return Promise.reject(error);
        } else {
          if (originalRequest.url == "/refresh") {
            clearAuth();
            navigate(rp.LOGIN_ROUTE);
            return Promise.reject(error);
          }

          return refresh(apiClient)
            .then(async () => {
              return await apiClient(originalRequest); // Повторяем оригинальный запрос
            })
            .catch((refreshError) => {
              clearAuth();
              navigate(rp.LOGIN_ROUTE);
              return Promise.reject(refreshError);
            });
        }
      }
      return Promise.reject(error);
    },
  );
  return (
    <ApiContext.Provider value={{ apiClient }}>{children}</ApiContext.Provider>
  );
}

export function useApi() {
  return useContext(ApiContext);
}
