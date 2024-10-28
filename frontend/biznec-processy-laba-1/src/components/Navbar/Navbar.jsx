import * as rp from "../../global/routerPaths/routerPaths";
import { Link, useLocation } from "react-router-dom";
import "./Navbar.scss";
import { useAuth } from "../../context/authContext";

const links = [
  {
    path: rp.USER_HOME_ROUTE,
    needAuth: true,
    needRole: "admin",
    title: "Пользователь",
  },
  {
    path: rp.ADMIN_HOME_ROUTE,
    needAuth: true,
    needRole: "admin",
    title: "Админ",
  },
  {
    path: rp.SCHEDULE_ROUTE,
    needAuth: true,
    needRole: "admin",
    title: "Расписания",
  },
  {
    path: rp.FLIGHT_BOOKING_ROUTE,
    needAuth: true,
    needRole: "admin",
    title: "Бронирования",
  },
  {
    path: rp.SURVEY_PAGE,
    needAuth: true,
    needRole: "admin",
    title: "Опросы",
  },
  {
    path: rp.AMINITIES_ROUTE,
    needAuth: true,
    needRole: "admin",
    title: "Услуги",
  },
  {
    path: rp.AMINITIES_REPORT_ROUTE,
    needAuth: true,
    needRole: "admin",
    title: "Отчет по услугам",
  },
  {
    path: rp.SUMMARY_ROUTE,
    needAuth: true,
    needRole: "admin",
    title: "Статистика",
  },
];

function Navbar() {
  let curLoc = useLocation();
  const { isAuth, role } = useAuth();

  return (
    <>
      <nav className="nav">
        <ul className="nav__list">
          {links.map((link, index) => {
            if (link.needAuth && !isAuth) return "";
            if (link.needRole && role != link.needRole) return "";

            return (
              <li className="nav__item" key={index}>
                <Link
                  className={curLoc.pathname == link.path ? "active" : ""}
                  to={link.path}
                >
                  {link.title}
                </Link>
              </li>
            );
          })}
        </ul>
      </nav>
    </>
  );
}

export default Navbar;
