import React from "react";
import { Route, Routes } from "react-router-dom";
import MainPage from "../../pages/MainPage/MainPage";
import LoginPage from "../../pages/LoginPage/LoginPage";
import AdminHomePage from "../../pages/AdminHomePage/AdminHomePage";
import UserHomePage from "../../pages/UserHomePage/UserHomePage";
import * as rp from "../../global/routerPaths/routerPaths";
import { ProtectedRoute } from "./ProtectedRoute";
import { BaseErrorPage } from "../../pages/ErrorPages/BaseErorrPage/BaseErrorPage";

function Router() {
  return (
    <>
      <Routes>
        <Route path={rp.MAIN_ROUTE} element={<MainPage />} />
        <Route path={rp.LOGIN_ROUTE} element={<LoginPage />} />
        <Route
          path={rp.ADMIN_HOME_ROUTE}
          element={
            <ProtectedRoute requiredRole={"admin"}>
              <AdminHomePage />
            </ProtectedRoute>
          }
        />
        <Route
          path={rp.USER_HOME_ROUTE}
          element={
            <ProtectedRoute requiredRole={"user"}>
              <UserHomePage />
            </ProtectedRoute>
          }
        />
        <Route
          path={rp.FORBIDDEN_ERROR_ROUTE}
          element={
            <BaseErrorPage
              errorCode={403}
              message="Вы не имеете доступа к данному ресурсу"
            />
          }
        />
        <Route
          path={rp.INTERNAL_ERROR_ROUTE}
          element={
            <BaseErrorPage
              errorCode={500}
              message="Произошла ошибка на сервере"
            />
          }
        />
        <Route
          path={rp.NOT_FOUND_ERROR_ROUTE}
          element={
            <BaseErrorPage
              errorCode={404}
              message="Запрашиваемый ресурс не найден"
            />
          }
        />
      </Routes>
    </>
  );
}

export default Router;
