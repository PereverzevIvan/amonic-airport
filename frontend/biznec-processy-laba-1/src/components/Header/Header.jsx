import * as rp from "../../global/routerPaths/routerPaths";
import amonicLogo from "../../assets/images/logo.png";
import Navbar from "../Navbar/Navbar";
import { Link } from "react-router-dom";
import "./Header.scss";
import { useAuth } from "../../context/authContext";
import { LogoutButton } from "../LogoutButton/LogoutButton";

function Header() {
  const { isAuth } = useAuth();

  return (
    <>
      <header className="header">
        <div className="header__container">
          <a href="/">
            <img src={amonicLogo} alt="AmonicLogo" className="logo" />
          </a>
          <div style={{ display: "flex", gap: "10px" }}>
            <Navbar />
            {!isAuth ? (
              <Link className="link-button" to={rp.LOGIN_ROUTE}>
                Войти
              </Link>
            ) : (
              <LogoutButton />
            )}
          </div>
        </div>
      </header>
    </>
  );
}

export default Header;
