import React from 'react';
import { BrowserRouter, Route, Routes, Link } from 'react-router-dom';
import MainPage from '../../pages/MainPage';

const MAIN_ROUTE = "/"
const LOGIN_ROUTE = "/login"
const ADMIN_HOME_ROUTE = "/admin-profile"
const USER_HOME_ROUTE = "/user-profile"

function Router() {
  return (
    <>
      <BrowserRouter>
        <nav>
          <ul>
            <li><Link to={MAIN_ROUTE}>Главная</Link></li>
            <li><Link to={LOGIN_ROUTE}>Войти</Link></li>
            <li><Link to={ADMIN_HOME_ROUTE}>Страница админа</Link></li>
            <li><Link to={USER_HOME_ROUTE}>Страница пользователя</Link></li>
          </ul>
        </nav>
        <Routes>
          <Route path={MAIN_ROUTE} element={<MainPage></MainPage>} />
        </Routes>
      </BrowserRouter>
    </>
  )

}

export default Router
