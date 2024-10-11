import React, { useState, createContext, useContext, useEffect } from "react";
import { jwtDecode } from "jwt-decode";
import Cookies from "js-cookie";
import { useNavigate } from "react-router-dom";
import * as rp from "../global/routerPaths/routerPaths.jsx";

const AuthContext = createContext();

export function AuthProvider({ children }) {
  const [userID, setUserID] = useState(sessionStorage.getItem("userID"));
  const [role, setRole] = useState(sessionStorage.getItem("role"));
  const [isAuth, setIsAuth] = useState(sessionStorage.getItem("isAuth"));

  const navigate = useNavigate();

  console.log("Mount authContext", userID, role, isAuth);
  useEffect(() => {
    const accessToken = Cookies.get("access-token");
    console.log("effect", userID, role, isAuth);
    if (accessToken) {
      decodeToken(accessToken);
      return;
    }

    //if (isAuth) refresh();
  }, []);

  function decodeToken(token) {
    try {
      const decoded = jwtDecode(token);
      setUserID(decoded.id);
      setRole(decoded.role);
      setIsAuth(true);
      sessionStorage.setItem("isAuth", "true");
      sessionStorage.setItem("userID", `${decoded.id}`);
      sessionStorage.setItem("role", decoded.role);
    } catch (error) {
      console.error("Ошибка декодирования токена:", error);
    }
  }

  function login(apiClient, credentials) {
    return apiClient
      .post("/login", credentials)
      .then((response) => {
        const accessToken = Cookies.get("access-token");
        decodeToken(accessToken);
        console.log("login успешен");
        return response;
      })
      .catch((error) => {
        if (error.response) {
          if (error.response.status == 403) {
            console.error(
              "Данные входа корректны. Однако вы заблокированы и не можете войти в систему",
            );
          }
        }
        console.error("Не удалось выполнить login");
        return Promise.reject(error);
      });
  }

  function clearAuth() {
    setUserID(null);
    setRole(null);
    Cookies.remove("access-token");
    Cookies.remove("refresh-token");
    setIsAuth(false);
    sessionStorage.removeItem("isAuth");
    sessionStorage.removeItem("userID");
    sessionStorage.removeItem("role");
    console.log("Очистка данных авторизации успешна");
  }

  function logout(apiClient) {
    return apiClient
      .get("/logout")
      .then((response) => {
        clearAuth();
        console.log("Logout успешен", userID, role);
        return response;
      })
      .catch((error) => {
        console.error("Не удалось выполнить logout");
        return Promise.reject(error);
      });
  }

  function refresh(apiClient) {
    return apiClient
      .get("/refresh")
      .then((response) => {
        const accessToken = Cookies.get("access-token");
        decodeToken(accessToken);
        console.log("refresh токенов успешен");

        return response;
      })
      .catch((error) => {
        console.error("Не удалось сделать refresh токенов");
        return Promise.reject(error);
      });
  }

  return (
    <AuthContext.Provider
      value={{ userID, role, login, logout, refresh, isAuth, clearAuth }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
