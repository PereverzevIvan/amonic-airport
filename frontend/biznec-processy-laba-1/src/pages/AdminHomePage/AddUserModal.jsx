import { useApi } from "../../context/apiContext";
import { AddNewUser } from "../../api/usersApi";
import Modal from "../../components/Modal/Modal";
import Button from "../../components/Button/Button";
import {
  allFieldsNotEmpty,
  validateEmail,
} from "../../global/formUtils/formUtils.jsx";
import { useState } from "react";
import { useToast } from "../../context/ToastContext.jsx";

export function AddUserModal({
  show = false,
  handleClose,
  offices = [],
  addNewUserInState,
}) {
  const { apiClient } = useApi();
  const { addToast } = useToast();

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
          addToast(`Пользователь c почтой ${newUser.email} успешно создан`);
          addNewUserInState(newUser);
          clearFormData();
        }
      })
      .catch((error) => {
        if (error.response) {
          error.response.data
            ? addToast(error.response.data, "error")
            : addToast(`Произошла ошибка ${error.response.status}`, "error");
        } else {
          addToast("Ошибка сети", "error");
        }
      });
  }

  return (
    <>
      <Modal
        show={show}
        title="Добавление пользователя"
        handleClose={() => {
          handleClose();
        }}
      >
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
