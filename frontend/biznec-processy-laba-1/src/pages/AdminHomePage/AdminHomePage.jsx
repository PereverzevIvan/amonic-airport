import "./AdminHomePage.scss";
import { useEffect, useState } from "react";
import { getYearDiferenceFrom2Strings } from "../../global/dateUtils/dateUtils.jsx";
import Button from "../../components/Button/Button.jsx";
import { getUsers, updateUserActive } from "../../api/usersApi.jsx";
import { useApi } from "../../context/apiContext.jsx";
import { getOffices } from "../../api/officesApi.jsx";
import { AddUserModal } from "./AddUserModal.jsx";
import { EditUserModal } from "./EditUserModal.jsx";

function createOfficesMap(offices) {
  let officesMap = new Map();

  for (let i = 0; i < offices.length; i++) {
    officesMap.set(offices[i].id, offices[i]);
  }

  return officesMap;
}

function createRolesMap() {
  let rolesMap = new Map();
  rolesMap.set(1, "Администратор");
  rolesMap.set(2, "Пользователь");

  return rolesMap;
}

function AdminHomePage() {
  const [users, setUsers] = useState([]); // Список пользователей, который получается от сервера
  const [offices, setOffices] = useState([]); // Список офисов, который получается от сервера
  const { apiClient } = useApi();

  const [selectedRow, setSelectedRow] = useState(null);
  const [userDataForEdit, setUserDataForEdit] = useState(null);

  const [selectedOffice, setSelectedOffice] = useState("");

  const [showAddUserModal, setShowAddUserModal] = useState(false);
  const [showEditUserModal, setShowEditUserModal] = useState(false);

  function changeOfficeFilter(event) {
    setSelectedOffice(event.target.value);
  }

  const officesMap = createOfficesMap(offices);
  const rolesMap = createRolesMap();

  function addNewUserInState(userData) {
    setUsers((prevItems) => [...prevItems, userData]);
  }

  function editUserInState(editedUserData) {
    if (editedUserData && editedUserData.id > 0) {
      setUsers((prevItems) =>
        prevItems.map(
          (user) => (user.id === editedUserData.id ? editedUserData : user), // Заменяем элемент на новое значение
        ),
      );
    }
  }

  const [message, setMessage] = useState("");
  const [messageType, setMessageType] = useState("");

  function clearMessage() {
    setMessageType("");
    setMessage("");
  }

  function setTimeoutForMessage() {
    setTimeout(clearMessage, 3 * 1000);
  }

  function handleClickOnRow(index) {
    if (index == selectedRow) {
      setSelectedRow(null);
      setUserDataForEdit(null);
    } else {
      if (users[index]) {
        setSelectedRow(index);
        setUserDataForEdit(users[index]);
      }
    }
  }

  function toggleUserActive(index, isActive) {
    let userData = {
      id: users[index].id,
      is_active: isActive,
    };

    updateUserActive(apiClient, userData)
      .then((response) => {
        if (response.status == 200) {
          updateUserActiveInState(index, isActive);
          setMessage("Статус пользователя успешно изменен");
          setMessageType("success");
        }
      })
      .catch((error) => {
        if (error.response) {
          if (error.response.data) setMessage(`Ошибка: ${error.response.data}`);
          else setMessage("Ошибка сети");
        }
        setMessageType("error");
      });

    setSelectedRow(null);
    setTimeoutForMessage();
  }

  const updateUserActiveInState = (index, isActive) => {
    setUsers((prevUsers) =>
      prevUsers.map((user, userIndex) =>
        userIndex === index ? { ...user, active: isActive } : user,
      ),
    );
  };

  useEffect(() => {
    function fetchOffices() {
      getOffices(apiClient)
        .then((response) => {
          setOffices(response.data);
          console.log("Офисы успешно получены");
        })
        .catch((error) => {
          console.error("Не удалось получить офисы");
        });
    }

    fetchOffices();
  }, []);

  useEffect(() => {
    const fetchUsers = () => {
      getUsers(apiClient, {
        office_id: selectedOffice,
      })
        .then((response) => {
          setUsers(response.data);
          console.log("Пользователи успешно получены");
        })
        .catch((error) => {
          if (error.response) {
            console.error(
              "Не удалось получить данные о пользователях",
              error.response.status,
            );
          }
        });
    };

    setSelectedRow(null);
    fetchUsers();
  }, [selectedOffice]);

  return (
    <>
      <section className="admin-home-page">
        <AddUserModal
          show={showAddUserModal}
          handleClose={() => {
            setShowAddUserModal(false);
          }}
          offices={offices}
          addNewUserInState={addNewUserInState}
        />
        <EditUserModal
          show={showEditUserModal}
          handleClose={() => {
            setShowEditUserModal(false);
          }}
          userData={userDataForEdit}
          offices={offices}
          editUserInState={editUserInState}
        />
        <div className="above-table-line">
          <label htmlFor="select-office" className="select-label">
            Офис:
            <select
              id="select-office"
              value={selectedOffice}
              onChange={changeOfficeFilter}
              className="drop-down-list"
            >
              <option className="drop-down-list__item" value="">
                Все офисы
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
          </label>
          {message && (
            <div className={`message message_${messageType}`}>{message}</div>
          )}
          <Button
            color="blue"
            onClick={() => {
              setShowAddUserModal(true);
            }}
          >
            Добавить пользователя
          </Button>
        </div>
        <div className="table-container">
          <table className="table">
            <thead className="table__head">
              <tr className="table__row">
                <th className="table__header">Имя</th>
                <th className="table__header">Фамилия</th>
                <th className="table__header">Возраст</th>
                <th className="table__header">Роль пользователя</th>
                <th className="table__header">Почта</th>
                <th className="table__header">Офис</th>
              </tr>
            </thead>
            <tbody className="table__body">
              {users.map((curUser, index) => {
                if (
                  selectedOffice != "" &&
                  curUser.office_id != selectedOffice
                ) {
                  return "";
                }
                return (
                  <tr
                    className={`table__row
                          ${curUser.role_id == 1 ? "green" : ""} 
                          ${curUser.active == false ? "red" : ""} 
                          ${index == selectedRow ? "blue" : ""}`}
                    key={index}
                    onClick={() => {
                      handleClickOnRow(index);
                    }}
                  >
                    <td className="table__data">{curUser.first_name}</td>
                    <td className="table__data">{curUser.last_name}</td>
                    <td className="table__data">
                      {getYearDiferenceFrom2Strings(
                        curUser.birthday,
                        new Date().toISOString(),
                      )}
                    </td>
                    <td className="table__data">
                      {rolesMap.get(curUser.role_id)}
                    </td>
                    <td className="table__data">{curUser.email}</td>
                    <td className="table__data">
                      {officesMap.has(curUser.office_id)
                        ? officesMap.get(curUser.office_id).title
                        : "Неизвестно"}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
        <div className="button-box">
          <Button
            disabled={selectedRow == null}
            color="blue"
            onClick={() => {
              setShowEditUserModal(true);
            }}
          >
            Изменить роль
          </Button>
          <Button
            disabled={selectedRow == null || users[selectedRow].active}
            color="green"
            onClick={() => {
              toggleUserActive(selectedRow, true);
            }}
          >
            Разблокировать
          </Button>
          <Button
            disabled={selectedRow == null || !users[selectedRow].active}
            color="red"
            onClick={() => {
              toggleUserActive(selectedRow, false);
            }}
          >
            Заблокировать
          </Button>
        </div>
      </section>
    </>
  );
}

export default AdminHomePage;
