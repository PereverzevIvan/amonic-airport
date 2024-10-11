import { useApi } from "../../context/apiContext";
import { AddNewUser } from "../../api/usersApi";
import Modal from "../../components/Modal/Modal";
import Button from "../../components/Button/Button";
import {
  allFieldsNotEmpty,
  validateEmail,
} from "../../global/formUtils/formUtils.jsx";
import { useState } from "react";

export function AddUserModal({
  show = false,
  handleClose,
  offices = [],
  addNewUserInState,
}) {
  const { apiClient } = useApi();

  const [message, setMessage] = useState("");
  const [messageType, setMessageType] = useState("");

  const [formData, setFormData] = useState({
    email: "",
    first_name: "",
    last_name: "",
    office_id: "",
    birthday: "",
    password: "",
  });

  function clearFormData() {
    setFormData({
      email: "",
      first_name: "",
      last_name: "",
      office_id: "",
      birthday: "",
      password: "",
    });
  }

  function clearMessage() {
    setMessageType("");
    setMessage("");
  }

  function setTimeoutForMessage() {
    setTimeout(clearMessage, 3 * 1000);
  }

  function handleChange(event) {
    const { name, value } = event.target;
    setFormData((prevData) => ({ ...prevData, [name]: value }));
  }

  async function handleSubmit(event) {
    event.preventDefault();

    const newUserData = {
      email: formData.email,
      first_name: formData.first_name,
      last_name: formData.last_name,
      office_id: parseInt(formData.office_id, 10),
      birthday: new Date(formData.birthday).toISOString(),
      password: formData.password,
    };

    AddNewUser(apiClient, newUserData)
      .then((response) => {
        console.log(response);
        if (response.status == 201) {
          const newUser = response.data;
          setMessage(`Пользователь c почтой ${newUser.email} успешно создан`);
          setMessageType("success");
          addNewUserInState(newUser);
          clearFormData();
        }
      })
      .catch((error) => {
        if (error.response) {
          error.response.data
            ? setMessage(error.response.data)
            : setMessage(`Произошла ошибка ${error.response.status}`);
        } else {
          setMessage("Ошибка сети");
        }
        setMessageType("error");
      });
    setTimeoutForMessage();
  }

  return (
    <>
      <Modal
        show={show}
        title="Добавление пользователя"
        handleClose={() => {
          clearMessage();
          handleClose();
        }}
      >
        {message && (
          <div className={`message message_${messageType}`}>{message}</div>
        )}
        <form className="form" onSubmit={handleSubmit}>
          <label className="form__label" htmlFor="email-field">
            Почта
          </label>
          <input
            type="email"
            id="email-field"
            name="email"
            value={formData.email}
            onChange={handleChange}
            className="form__input"
            required
          />
          <label className="form__label" htmlFor="first-name-field">
            Имя
          </label>
          <input
            type="text"
            id="first-name-field"
            name="first_name"
            value={formData.first_name}
            className="form__input"
            onChange={handleChange}
            required
          />
          <label className="form__label" htmlFor="last-name-field">
            Фамилия
          </label>
          <input
            type="text"
            id="last-name-field"
            name="last_name"
            value={formData.last_name}
            onChange={handleChange}
            className="form__input"
            required
          />
          <label className="form__label" htmlFor="select-office-field">
            Офис
          </label>
          <select
            id="select-office-field"
            name="office_id"
            value={formData.office_id}
            onChange={handleChange}
            required
            className="drop-down-list"
          >
            <option className="drop-down-list__item" value="">
              Выберите офис
            </option>
            {offices.map((curOffice, index) => {
              return (
                <option
                  className="drop-down-list__item"
                  key={index}
                  value={curOffice.id}
                >
                  {curOffice.title}
                </option>
              );
            })}
          </select>
          <label className="form__label" htmlFor="birthday-field">
            День рождения
          </label>
          <input
            type="date"
            id="birthday-field"
            name="birthday"
            value={formData.birthday}
            onChange={handleChange}
            className="form__input"
            required
          />
          <label className="form__label" htmlFor="password-field">
            Пароль
          </label>
          <input
            type="password"
            id="password-field"
            name="password"
            value={formData.password}
            onChange={handleChange}
            className="form__input"
            required
          />

          <div className="form__button-box">
            <Button
              type="submit"
              disabled={
                !validateEmail(formData.email) || !allFieldsNotEmpty(formData)
              }
              onClick={() => {}}
              color="green"
            >
              Сохранить
            </Button>
            <Button
              color="red"
              onClick={() => {
                clearMessage();
                clearFormData();
                handleClose();
              }}
            >
              Отменить
            </Button>
          </div>
        </form>
      </Modal>
    </>
  );
}
