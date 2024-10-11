import * as rp from "../../global/routerPaths/routerPaths";
import { Link, useLocation } from "react-router-dom";
import "./Navbar.scss";
import { useAuth } from "../../context/authContext";

function Navbar() {
  let curLoc = useLocation();
  const { isAuth, role } = useAuth();

  return (
    <>
      <nav className="nav">
        <ul className="nav__list">
          <li className="nav__item">
            <Link
              className={curLoc.pathname == rp.MAIN_ROUTE ? "active" : ""}
              to={rp.MAIN_ROUTE}
            >
              Главная{" "}
            </Link>
          </li>
          {isAuth && (
            <>
              {role == "admin" && (
                <li className="nav__item">
                  <Link
                    className={
                      curLoc.pathname == rp.ADMIN_HOME_ROUTE ? "active" : ""
                    }
                    to={rp.ADMIN_HOME_ROUTE}
                  >
                    Страница админа
                  </Link>
                </li>
              )}

              <li className="nav__item">
                <Link
                  className={
                    curLoc.pathname == rp.USER_HOME_ROUTE ? "active" : ""
                  }
                  to={rp.USER_HOME_ROUTE}
                >
                  Страница пользователя
                </Link>
              </li>
            </>
          )}
        </ul>
      </nav>
    </>
  );
}

export default Navbar;
