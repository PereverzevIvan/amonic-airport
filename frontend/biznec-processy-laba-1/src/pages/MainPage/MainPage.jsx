import * as rp from "../../global/routerPaths/routerPaths";
import cloud from "../../assets/images/mainPage/cloud.png";
import vector from "../../assets/images/mainPage/vector.svg";
import airplane from "../../assets/images/mainPage/airplane.png";
import "./MainPage.scss";
import { Link } from "react-router-dom";

function MainPage() {
  return (
    <>
      <section className="main-page">
        <h1 className="title">Путешествуй вместе с нами!</h1>
        <img className="main-vector-image" src={vector} />
        <img className="main-cloud-image" src={cloud} />
        <img className="main-airplane-image" src={airplane} />
        <Link className="link-button" to={rp.LOGIN_ROUTE}>
          Войти в личный кабинет
        </Link>
      </section>
    </>
  );
}

export default MainPage;
