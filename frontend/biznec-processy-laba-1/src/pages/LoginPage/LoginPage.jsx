import { useEffect, useState } from "react";
import logo from "../../assets/images/logo.png";
import Button from "../../components/Button/Button.jsx";
import "./LoginPage.scss";
import { useNavigate } from "react-router-dom";
import * as rp from "../../global/routerPaths/routerPaths.jsx";
import { useAuth } from "../../context/authContext.jsx";
import { useApi } from "../../context/apiContext.jsx";
import { validateEmail } from "../../global/formUtils/formUtils.jsx";

function LoginPage() {
  return (
    <>
      <section className="login-page">
        <img src={logo} alt="Amonic logo" className="logo" />
        <LoginForm />
      </section>
    </>
  );
}

function LoginForm() {
  const [formData, setFormData] = useState({ email: "", password: "" }); // значения внутри формы
  const navigate = useNavigate(); // Хук для перенаправления пользователя на другую страницу

  const [faildTries, setFaildTries] = useState(0); // Количество неудачных попыток авторизоваться
  const [errorMessage, setErrorMessage] = useState("");
  const [timerIsVisible, setTimerIsVisible] = useState(false);

  const { login, isAuth } = useAuth();
  const { apiClient } = useApi();

  useEffect(() => {
    if (isAuth) navigate(rp.USER_HOME_ROUTE);
  }, [isAuth]);

  function handleChange(event) {
    const { name, value } = event.target;
    setFormData((prevData) => ({ ...prevData, [name]: value }));
  }

  function resetFailedTries() {
    setErrorMessage(
      "Достигнуто максимальное количество попыток входа. Повторите попытку через 10 секунд.",
    );
    setTimerIsVisible(true);

    setTimeout(() => {
      setFaildTries(0);
      setErrorMessage();
      setTimerIsVisible(false);
      console.log("Вы можете снова попытаться войти в систему");
    }, 10 * 1000);
  }

  useEffect(() => {
    if (faildTries == 3) resetFailedTries();
  }, [faildTries]);

  async function handleSubmit(event) {
    event.preventDefault();
    login(apiClient, formData).catch((error) => {
      if (!error.response) {
        setErrorMessage(
          "Произошла обшика сети. Попробуйте авторизоваться позже",
        );
      }
      if (error.response.status == 401) {
        setFaildTries(faildTries + 1);
        setErrorMessage(error.response.data);
      }
      if (error.response.status == 403) {
        setErrorMessage(error.response.data);
      }
    });

    setFormData({ email: "", password: "" });
  }

  function handleReturn() {
    navigate(rp.MAIN_ROUTE);
  }

  return (
    <>
      <form onSubmit={handleSubmit} className="form login-form">
        {errorMessage && (
          <div className="message message_error">{errorMessage}</div>
        )}
        <Timer seconds={10} isVisible={timerIsVisible} />
        <div className="form__container">
          <label className="form__label" htmlFor="email-field">
            Email:
          </label>
          <input
            disabled={faildTries == 3}
            className="form__input"
            id="email-field"
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            required
          />
          <label className="form__label" htmlFor="password-field">
            Password:
          </label>
          <input
            disabled={faildTries == 3}
            className="form__input"
            id="password-field"
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form__button-box">
          <Button
            disabled={!validateEmail(formData.email) || !formData.password}
            onClick={() => {}}
            type="submit"
            color="green"
          >
            Отправить
          </Button>
          <Button color="blue" onClick={handleReturn}>
            На главную
          </Button>
        </div>
      </form>
    </>
  );
}

function Timer({ seconds, isVisible = false }) {
  return (
    <>
      <div
        className="timer"
        style={isVisible ? { display: "block" } : { display: "none" }}
      >
        <div
          className="line"
          style={
            isVisible
              ? { animation: `countdown ${seconds}s linear forwards` }
              : {}
          }
        ></div>
      </div>
    </>
  );
}

export default LoginPage;
