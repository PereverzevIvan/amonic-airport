import { useApi } from "../../context/apiContext";
import { EditUser } from "../../api/usersApi";
import Modal from "../../components/Modal/Modal";
import Button from "../../components/Button/Button";
import {
  allFieldsNotEmpty,
  validateEmail,
} from "../../global/formUtils/formUtils.jsx";
import { useEffect, useState } from "react";

export function EditUserModal({
  show = false,
  handleClose,
  offices = [],
  userData,
  editUserInState,
}) {
  const { apiClient } = useApi();

  const [message, setMessage] = useState("");
  const [messageType, setMessageType] = useState("");

  const [formData, setFormData] = useState({
    id: "",
    email: "",
    first_name: "",
    last_name: "",
    office_id: "",
    role_id: "",
  });

  function clearFormData() {
    setFormData({
      id: "",
      email: "",
      first_name: "",
      last_name: "",
      office_id: "",
      role_id: "",
    });
  }

  function clearMessage() {
    setMessageType("");
    setMessage("");
  }

  useEffect(() => {
    if (userData) {
      setFormData({
        id: userData.id,
        email: userData.email,
        first_name: userData.first_name,
        last_name: userData.last_name,
        office_id: userData.office_id,
        role_id: userData.role_id,
      });
    } else {
      clearFormData();
    }
  }, [userData]);

  function handleChange(event) {
    const { name, value } = event.target;
    setFormData((prevData) => ({ ...prevData, [name]: value }));
  }

  async function handleSubmit(event) {
    event.preventDefault();

    const editUserData = {
      id: parseInt(formData.id, 10),
      email: formData.email,
      first_name: formData.first_name,
      last_name: formData.last_name,
      office_id: parseInt(formData.office_id, 10),
      role_id: parseInt(formData.role_id, 10),
    };

    EditUser(apiClient, editUserData)
      .then((response) => {
        if (response.status == 200) {
          const editedUser = response.data;
          setMessage(`Пользователь успешно изменён`);
          setMessageType("success");
          editUserInState(editedUser);
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

  function setTimeoutForMessage() {
    setTimeout(clearMessage, 3 * 1000);
  }

  return (
    <>
      <Modal
        show={show}
        title="Редактирование пользователя"
        handleClose={() => {
          clearMessage();
          handleClose();
        }}
      >
        {message && (
          <div className={`message message_${messageType}`}>{message}</div>
        )}
        <form className="form" onSubmit={handleSubmit}>
          <label className="form__label" htmlFor="edit-email-field">
            Почта
          </label>
          <input
            type="email"
            id="edit-email-field"
            name="email"
            value={formData.email}
            onChange={handleChange}
            className="form__input"
            required
          />
          <label className="form__label" htmlFor="edit-first-name-field">
            Имя
          </label>
          <input
            type="text"
            id="edit-first-name-field"
            name="first_name"
            value={formData.first_name}
            className="form__input"
            onChange={handleChange}
            required
          />
          <label className="form__label" htmlFor="edit-last-name-field">
            Фамилия
          </label>
          <input
            type="text"
            id="edit-last-name-field"
            name="last_name"
            value={formData.last_name}
            onChange={handleChange}
            className="form__input"
            required
          />
          <label className="form__label" htmlFor="edit-select-office-field">
            Офис
          </label>
          <select
            id="edit-select-office-field"
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
          <label className="form__label" htmlFor="edit-role-field">
            Роль
          </label>
          <div
            id="edit-role-field"
            className="form__radio-button-box form__radio-button-box_v"
          >
            <label className="form__radio-label">
              <input
                type="radio"
                name="role_id"
                value="1"
                onChange={handleChange}
                checked={formData.role_id == 1}
              />
              Администратор
            </label>
            <label className="form__radio-label">
              <input
                type="radio"
                name="role_id"
                value="2"
                onChange={handleChange}
                checked={formData.role_id == 2}
              />
              Пользователь
            </label>
          </div>

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
