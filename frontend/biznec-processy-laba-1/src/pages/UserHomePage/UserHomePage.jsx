import { useEffect, useMemo, useState } from "react";
import "./UserHomePage.scss";
import Button from "../../components/Button/Button.jsx";
import Modal from "../../components/Modal/Modal.jsx";
import {
  getDateFromString,
  getTimeDiferenceFrom2Strings,
  getTimeFromString,
  ConvertMillisecondsToTimeString,
} from "../../global/dateUtils/dateUtils.jsx";
import { useApi } from "../../context/apiContext.jsx";
import {
  getUserSessions,
  updateUserSession,
} from "../../api/userSessionsApi.jsx";
import { useAuth } from "../../context/authContext.jsx";
import { getUserById } from "../../api/usersApi.jsx";
import { allFieldsNotEmpty } from "../../global/formUtils/formUtils.jsx";
import { useToast } from "../../context/ToastContext.jsx";

function prepareDataToShow(data) {
  if (!data) return [[], 0];

  let userSessions = [];
  let crashCount = 0;

  for (let i = 1; i < data.length; i++) {
    let curSession = {
      id: 0,
      login_at: "",
      login_date: "",
      login_time: "",
      logout_at: "",
      logout_date: "",
      logout_time: "",
      spentTime: "",
      invalid_logout_reason: "",
      crashType: null,
    };

    curSession.id = data[i].id;
    curSession.login_at = data[i].login_at;
    curSession.login_date = getDateFromString(data[i].login_at);
    curSession.login_time = getTimeFromString(data[i].login_at);
    if (data[i].logout_at != null) {
      curSession.logout_at = data[i].logout_at;
      curSession.logout_date = getDateFromString(data[i].logout_at);
      curSession.logout_time = getTimeFromString(data[i].logout_at);
      curSession.spentTime = getTimeDiferenceFrom2Strings(
        data[i].login_at,
        data[i].logout_at,
      );
    }
    curSession.crashType = data[i].crash_reason_type;
    if (data[i].crash_reason_type != null) {
      crashCount++;
    }
    curSession.invalid_logout_reason = data[i].invalid_logout_reason;
    userSessions.push(curSession);
  }
  return [userSessions, crashCount];
}

function getSpentTimeFromSessions(userSessions, maxDays = 30) {
  const now = new Date();
  const dateMaxDaysAgo = new Date(now);
  dateMaxDaysAgo.setDate(now.getDate() - maxDays);

  let totalMilliseconds = 0;

  for (let i = 0; i < userSessions.length; i++) {
    let curSession = userSessions[i];
    //console.log(curSession.login_at, curSession.logout_at);
    if (!curSession.login_at || !curSession.logout_at) continue;

    let sessionStart = new Date(curSession.login_at);
    let sessionEnd = new Date(curSession.logout_at);

    if (sessionStart < dateMaxDaysAgo) break;
    else totalMilliseconds += sessionEnd - sessionStart;
  }

  return ConvertMillisecondsToTimeString(totalMilliseconds);
}

function UserHomePage() {
  const { apiClient } = useApi();
  const { userID } = useAuth();

  const [user, setUser] = useState(null);
  const [userSessions, setUserSessions] = useState(null);
  const [invalidUserSessions, setInvalidUserSessions] = useState(null);
  const [userSessionsCount, setUserSessionsCount] = useState(0);
  const [crashCount, setCrashCount] = useState(0);
  const [curUserSession, setCurSession] = useState(null);
  const [editingInvalidUserSession, setEditingInvalidUserSession] =
    useState(null);

  const curSpentTime = useMemo(() => {
    if (userSessions && userSessions != []) {
      return getSpentTimeFromSessions(userSessions, 2);
    }
  }, [userSessions]);

  const [showModal, setShowModal] = useState(false);

  function changeSessionInState(newSessionData) {
    if (newSessionData) {
      setUserSessions((prevSessions) =>
        prevSessions.map((session) =>
          session.id == newSessionData.id
            ? {
                ...session,
                invalid_logout_reason: newSessionData.invalid_logout_reason,
                crashType: newSessionData.crash_reason_type,
              }
            : session,
        ),
      );
    }
  }

  useEffect(() => {
    // Получаем данные о пользователе
    function fetchUserInfo() {
      getUserById(apiClient, userID)
        .then((response) => {
          setUser(response.data);
          console.log("Данные о пользователе успешно получены");
        })
        .catch((error) => {
          console.error("Не удалось получить данные о пользователе");
        });
    }

    // Получаем список всех сессий пользователя
    function fetchUserSessions() {
      getUserSessions(apiClient, {
        user_id: userID,
        only_unmarked_invalid_sessions: false,
      })
        .then((response) => {
          if (!response.data) return;

          const [preparedSessions, preparedCrashCount] = prepareDataToShow(
            response.data,
          );

          setUserSessions(preparedSessions);
          setCrashCount(preparedCrashCount);
          setUserSessionsCount(preparedSessions.length);
          setCurSession(response.data[0] && response.data[0]);
          console.log("Данные о всех сессиях пользователя успешно получены");
        })
        .catch((error) => {
          console.error(
            "Не удалось получить данные о всех сессиях пользователя",
            error,
          );
        });
    }

    //Получаем весь список только невалидных сессий пользователя
    function fetchInvalidUserSessions() {
      getUserSessions(apiClient, {
        user_id: userID,
        only_unmarked_invalid_sessions: true,
      })
        .then((response) => {
          if (!response.data) {
            setInvalidUserSessions(null);
            return;
          }

          setInvalidUserSessions(response.data);
          console.log(
            "Данные о невалидных сессиях пользователя успешно получены",
          );
        })
        .catch((error) => {
          console.error(
            "Не удалось получить данные о невалидных сессиях пользователя",
            error,
          );
        });
    }

    fetchUserInfo();
    fetchUserSessions();
    fetchInvalidUserSessions();
  }, []);

  function leftPopInvalidUserSession() {
    setInvalidUserSessions((prevSessions) => prevSessions.slice(1));
  }

  useEffect(() => {
    if (invalidUserSessions && invalidUserSessions.length > 0) {
      setEditingInvalidUserSession(invalidUserSessions[0]);
      setShowModal(true);
    } else {
      setEditingInvalidUserSession(null);
      setShowModal(false);
    }
  }, [invalidUserSessions]);

  return (
    <>
      <EditUserSessionModal
        showModal={showModal}
        handleClose={() => {
          setShowModal(false);
        }}
        userSession={editingInvalidUserSession}
        popSession={leftPopInvalidUserSession}
        changeSessionInState={changeSessionInState}
      />
      <section className="user-home-page">
        <p className="common-text common-text_big">
          Привет, {user && user.first_name}, Добро пожаловать в авиакомпанию
          AMONIC.
        </p>
        <div className="above-table-container">
          <div className="statistic-container">
            <p className="common-text">
              Время в системе за 30 дней: {curSpentTime}
            </p>
            <p className="common-text">Всего сессий: {userSessionsCount}</p>
            <p className="common-text">Количество сбоев: {crashCount}</p>
          </div>
        </div>

        <div className="table-container">
          <table className="table">
            <thead className="table__head">
              <tr className="table__row">
                <th className="table__header">Дата</th>
                <th className="table__header">Время входа</th>
                <th className="table__header">Время выхода</th>
                <th className="table__header">Время в системе</th>
                <th className="table__header">Причина неудачного выхода</th>
              </tr>
            </thead>
            <tbody className="table__body">
              {userSessions &&
                userSessions.map((curSession, index) => (
                  <tr
                    className={`table__row 
                  ${curSession.crashType != null ? "red" : ""}`}
                    key={index}
                  >
                    <td className="table__data">{curSession.login_date}</td>
                    <td className="table__data">{curSession.login_time}</td>
                    <td className="table__data">{curSession.logout_time}</td>
                    <td className="table__data">{curSession.spentTime}</td>
                    <td className="table__data">
                      {curSession.invalid_logout_reason != "undefined"
                        ? curSession.invalid_logout_reason
                        : "Не указано"}
                    </td>
                  </tr>
                ))}
            </tbody>
          </table>
        </div>
      </section>
    </>
  );
}

function EditUserSessionModal({
  showModal,
  handleClose,
  userSession,
  popSession,
  changeSessionInState,
}) {
  const { apiClient } = useApi();
  const { addToast } = useToast();

  const [formData, setFormData] = useState({
    crash_reason_type: null,
    id: 0,
    invalid_logout_reason: "",
    login_at: "",
    logout_at: "",
    user_id: 0,
  });

  useEffect(() => {
    if (userSession) {
      setFormData({
        crash_reason_type: null,
        id: userSession.id,
        invalid_logout_reason: "",
        login_at: userSession.login_at,
        logout_at: userSession.logout_at,
        user_id: userSession.user_id,
      });
    }
  }, [userSession]);

  function handleChange(event) {
    const { name, value } = event.target;
    setFormData((prevData) => ({ ...prevData, [name]: value }));
  }

  function handleSubmit(event) {
    event.preventDefault();

    const changeUserSessionData = {
      crash_reason_type: `${formData.crash_reason_type}`,
      id: `${formData.id}`,
      reason: `${formData.invalid_logout_reason}`,
      login_at: `${formData.login_at}`,
      logout_at: `${formData.logout_at}`,
      user_id: `${formData.user_id}`,
      crash_reason_type: `${formData.crash_reason_type}`,
    };

    updateUserSession(apiClient, changeUserSessionData)
      .then((response) => {
        console.log("Данные невалидной сессии успешно изменены");
        changeSessionInState(formData);
        handleClose();
        popSession();
      })
      .catch((error) => {
        if (error.response) {
          addToast(error.response.data, "error");
        } else {
          addToast("Ошибка сети", "error");
        }
      });
  }

  return (
    <>
      <Modal
        show={showModal}
        handleClose={handleClose}
        title="Выход из системы не обнаружен"
      >
        <p className="common-text common-text_big">
          При вашем последнем входе в систему не был выполнен выход из системы.
          Что произошло{" "}
          {formData.login_at &&
            `${getDateFromString(formData.login_at)} в ${getTimeFromString(formData.login_at)}`}
          ?
        </p>
        <form className="form" onSubmit={handleSubmit}>
          <div className="form__container">
            <label
              className="form__label"
              htmlFor="invalid_logout_reason_textarea"
            >
              Укажите причину:
            </label>
            <textarea
              className="form__input"
              id="invalid_logout_reason_textarea"
              required
              name="invalid_logout_reason"
              value={formData.invalid_logout_reason}
              onChange={handleChange}
            ></textarea>
          </div>
          <div className="form__button-box">
            <div className="button-box">
              <label className="form__radio-label">
                <input
                  type="radio"
                  name="crash_reason_type"
                  value={1}
                  checked={formData.crash_reason_type == 1}
                  onChange={handleChange}
                  required
                />
                Сбой программного обеспечения
              </label>
              <label className="form__radio-label">
                <input
                  type="radio"
                  name="crash_reason_type"
                  value={2}
                  checked={formData.crash_reason_type == 2}
                  onChange={handleChange}
                  required
                />
                Сбой системы
              </label>
            </div>
            <Button
              color="green"
              type="submit"
              disabled={
                !allFieldsNotEmpty({
                  invalid_logout_reason: formData.invalid_logout_reason,
                  crash_reason_type: formData.crash_reason_type,
                })
              }
            >
              Отправить
            </Button>
          </div>
        </form>
      </Modal>
    </>
  );
}

export default UserHomePage;
