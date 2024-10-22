import { useEffect, useState } from "react";
import { useToast } from "../../context/ToastContext";
import "./SchedulePage.scss";
import Button from "../../components/Button/Button";
import {
  changeScheduleConfirmed,
  getSchedules,
  updateScheduleByID,
  UploadFileWithSchedules,
} from "../../api/schedulesApi";
import { useApi } from "../../context/apiContext.jsx";
import {
  getDateFromString,
  getTimeFromString,
  getUTCDateFromString,
} from "../../global/dateUtils/dateUtils.jsx";
import { getAirports } from "../../api/airportApi.jsx";
import Modal from "../../components/Modal/Modal.jsx";
import { allFieldsNotEmpty } from "../../global/formUtils/formUtils.jsx";

function calculateBisnessAndFirstClassPrice(economyPrice) {
  if (!economyPrice) return;
  const businessPrice = Math.floor(economyPrice * 1.35); // Билет в бизнес-класс на 30% дороже, округляем вниз
  const firstClassPrice = Math.floor(businessPrice * 1.3); // Билет в первый класс на 35% дороже, округляем вниз

  return {
    businessClass: businessPrice,
    firstClass: firstClassPrice,
  };
}

function SchedulePage() {
  const [filtersAndSortData, setFiltersAndSortData] = useState({
    outbound: "",
    flight_number: "",
    from: "",
    to: "",
    sort_by: "",
  });

  const { apiClient } = useApi();

  const [airports, setAirports] = useState(null);

  const [schedules, setSchedules] = useState(null);
  const [editingShedule, setEditingShedule] = useState(null);

  const [selectedRow, setSelectedRow] = useState(null);

  const [showEditScheduleModal, setShowEditScheduleModal] = useState(false);
  const [showImportFileModal, setShowImportFileModal] = useState(false);

  const { addToast } = useToast();

  function handleClickOnRow(index) {
    if (index == selectedRow) {
      setSelectedRow(null);
      setEditingShedule(null);
    } else {
      if (schedules && schedules[index]) {
        setSelectedRow(index);
        setEditingShedule(schedules[index]);
      }
    }
  }

  function handleChangeFiltersAndSort(event) {
    const { name, value } = event.target;
    setFiltersAndSortData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  }

  function fetchAirports() {
    getAirports(apiClient)
      .then((response) => {
        if (response.status == 200) {
          setAirports(response.data);
          console.log("Аэропорты получены");
        }
      })
      .catch((error) => {
        if (error.response) {
          addToast(error.response.data, "error");
        } else {
          addToast("Ошибка сети", "error");
          console.error("Ошибка сети");
        }
      });
  }

  function fetchSchedules() {
    let params = {
      outbound: filtersAndSortData.outbound,
      flight_number: filtersAndSortData.flight_number,
      from: filtersAndSortData.from,
      to: filtersAndSortData.to,
      sort_by: filtersAndSortData.sort_by,
    };

    // Замена пустых полей в объекте на null, чтобы не входили в запрос
    params = Object.fromEntries(
      Object.entries(params).map(([key, value]) => {
        if (key === "from" || key === "to") {
          return [key, value === "" ? null : Number(value)]; // Преобразуем в число
        }
        if (key === "outbound") {
          return [
            key,
            value === ""
              ? null
              : getUTCDateFromString(
                  new Date(filtersAndSortData.outbound).toISOString(),
                ),
          ];
        }
        return [key, value === "" ? null : value]; // Для остальных полей
      }),
    );

    clearSelectedRow();

    getSchedules(apiClient, params)
      .then((response) => {
        if (response.status == 200) {
          setSchedules(response.data);
          console.log("Расписание получено");
        }
      })
      .catch((error) => {
        if (error.response) {
          addToast(error.response.data, "error");
        } else {
          addToast("Ошибка сети", "error");
          console.error("Ошибка сети");
        }
      });
  }

  function toggleScheduleStatus(scheduleID, isConfirmed) {
    changeScheduleConfirmed(apiClient, scheduleID, isConfirmed)
      .then((response) => {
        if (response.status == 200) {
          addToast("Статус рейса изменен");
          console.log("Статус рейса изменен");
          changeScheduleConfirmedInState(scheduleID, isConfirmed);
          clearSelectedRow();
        }
      })
      .catch((error) => {
        if (error.response) {
          console.log(error.response);
          addToast(error.response.data, "error");
        } else {
          addToast("Ошибка сети", "error");
        }
      });
  }

  function changeScheduleConfirmedInState(scheduleID, isConfirmed) {
    setSchedules((prevSchedules) =>
      prevSchedules.map((schedule) =>
        schedule.id === scheduleID
          ? { ...schedule, confirmed: isConfirmed }
          : schedule,
      ),
    );
  }

  function clearSelectedRow() {
    setSelectedRow(null);
    setEditingShedule(null);
  }

  useEffect(() => {
    fetchAirports();
    fetchSchedules();
  }, []);

  useEffect(() => {}, [schedules]);

  function changeScheduleInState(newScheduleData) {
    if (newScheduleData && newScheduleData.id > 0) {
      setSchedules((prevItems) =>
        prevItems.map(
          (schedule) =>
            schedule.id === newScheduleData.id
              ? {
                  ...schedule,
                  outbound: `${newScheduleData.date}T${newScheduleData.time}Z`,
                  economy_price: newScheduleData.economy_price,
                }
              : schedule, // Заменяем элемент на новое значение
        ),
      );
    }
  }

  return (
    <>
      <EditScheduleModal
        show={showEditScheduleModal}
        handleClose={() => {
          setShowEditScheduleModal(false);
        }}
        oldSchedule={editingShedule}
        changeScheduleInState={changeScheduleInState}
      />
      <ImportFileModal
        show={showImportFileModal}
        handleClose={() => {
          setShowImportFileModal(false);
        }}
      />
      <section className="schedule-page">
        <form className="form">
          <fieldset className="form__fieldset">
            <legend className="form__legend">Фильтрация и сортировка</legend>

            <div className="form__container">
              <label className="form__label" htmlFor="select-departure-airport">
                Откуда
              </label>
              <select
                className="drop-down-list"
                id="select-departure-airport"
                name="from"
                value={filtersAndSortData.from}
                onChange={handleChangeFiltersAndSort}
              >
                <option className="drop-down-list__item" value="">
                  Аэропорт
                </option>
                {airports?.map((curAirport, index) => (
                  <option
                    className="drop-down-list__item"
                    value={curAirport?.id}
                    key={index}
                  >
                    {curAirport?.iata_code}
                  </option>
                ))}
              </select>
            </div>

            <div className="form__container">
              <label className="form__label" htmlFor="select-arrival-airport">
                Куда
              </label>
              <select
                className="drop-down-list"
                id="select-arrival-airport"
                name="to"
                value={filtersAndSortData.to}
                onChange={handleChangeFiltersAndSort}
              >
                <option className="drop-down-list__item" value="">
                  Аэропорт
                </option>
                {airports?.map((curAirport, index) => (
                  <option
                    className="drop-down-list__item"
                    value={curAirport?.id}
                    key={index}
                  >
                    {curAirport?.iata_code}
                  </option>
                ))}
              </select>
            </div>

            <div className="form__container">
              <label className="form__label" htmlFor="select-sort-param">
                Сортировка
              </label>
              <select
                className="drop-down-list"
                id="select-sort-param"
                name="sort_by"
                value={filtersAndSortData.sort_by}
                onChange={handleChangeFiltersAndSort}
              >
                <option className="drop-down-list__item" value="date_time">
                  Дата и время
                </option>
                <option className="drop-down-list__item" value="ticket_price">
                  Цена эконома
                </option>
                <option className="drop-down-list__item" value="confirmed">
                  Подтверждение
                </option>
              </select>
            </div>

            <div className="form__container">
              <label className="form__label" htmlFor="input-outbound">
                Отправление
              </label>
              <input
                className="form__input"
                id="input-outbound"
                type="date"
                name="outbound"
                value={filtersAndSortData.outbound}
                onChange={handleChangeFiltersAndSort}
              />
            </div>

            <div className="form__container">
              <label className="form__label" htmlFor="input-flight-number">
                Номер рейса
              </label>
              <input
                className="form__input"
                id="input-flight-number"
                type="number"
                name="flight_number"
                value={filtersAndSortData.flight_number}
                step={1}
                min={0}
                onChange={handleChangeFiltersAndSort}
              />
            </div>
            <Button onClick={fetchSchedules}>Применить</Button>
          </fieldset>
        </form>

        <div className="table-container">
          <table className="table">
            <thead className="table__head">
              <tr className="table__row">
                <th className="table__header">Дата</th>
                <th className="table__header">Время</th>
                <th className="table__header">Откуда</th>
                <th className="table__header">Куда</th>
                <th className="table__header">Номер рейса</th>
                <th className="table__header">Тип самолета</th>
                <th className="table__header">Эконом класс</th>
                <th className="table__header">Бизнес класс</th>
                <th className="table__header">Первый класс</th>
              </tr>
            </thead>
            <tbody className="table__body">
              {schedules &&
                schedules.map((curSchedule, index) => {
                  return (
                    <tr
                      className={`table__row 
                                ${!curSchedule?.confirmed ? "red" : ""}
                                ${index == selectedRow ? "blue" : ""}`}
                      key={index}
                      onClick={() => {
                        handleClickOnRow(index);
                      }}
                    >
                      <td className="table__data">
                        {getDateFromString(curSchedule?.outbound)}
                      </td>
                      <td className="table__data">
                        {getTimeFromString(curSchedule?.outbound)}
                      </td>
                      <td className="table__data">
                        {curSchedule?.route?.departure_airport?.iata_code}
                      </td>
                      <td className="table__data">
                        {curSchedule?.route?.arrival_airport?.iata_code}
                      </td>
                      <td className="table__data">
                        {curSchedule?.flight_number}
                      </td>
                      <td className="table__data">
                        {curSchedule?.aircraft?.name}
                      </td>
                      <td className="table__data">
                        ${curSchedule?.economy_price}
                      </td>
                      <td className="table__data">
                        $
                        {
                          calculateBisnessAndFirstClassPrice(
                            curSchedule?.economy_price,
                          )?.businessClass
                        }
                      </td>
                      <td className="table__data">
                        $
                        {
                          calculateBisnessAndFirstClassPrice(
                            curSchedule?.economy_price,
                          )?.firstClass
                        }
                      </td>
                    </tr>
                  );
                })}
            </tbody>
          </table>
        </div>
        <div className="under-table-line">
          <Button
            disabled={selectedRow == null}
            color="blue"
            onClick={() => {
              setShowEditScheduleModal(true);
            }}
          >
            Изменить расписание
          </Button>
          <Button
            color="green"
            disabled={
              selectedRow == null ||
              !schedules ||
              schedules[selectedRow]?.confirmed === true
            }
            onClick={() => {
              toggleScheduleStatus(schedules[selectedRow]?.id, true);
            }}
          >
            Подтвердить
          </Button>
          <Button
            color="red"
            disabled={
              selectedRow == null ||
              !schedules ||
              schedules[selectedRow]?.confirmed === false
            }
            onClick={() => {
              toggleScheduleStatus(schedules[selectedRow]?.id, false);
            }}
          >
            Отменить
          </Button>
          <Button
            onClick={() => {
              setShowImportFileModal(true);
            }}
          >
            Импорт файла
          </Button>
        </div>
      </section>
    </>
  );
}

function EditScheduleModal({
  show,
  handleClose,
  oldSchedule,
  changeScheduleInState,
}) {
  const [formData, setFormData] = useState({
    id: 0,
    date: "",
    time: "",
    economy_price: "",
  });

  const { apiClient } = useApi();
  const { addToast } = useToast();

  useEffect(() => {
    if (!oldSchedule) handleClose();
    setFormData({
      id: oldSchedule?.id || "",
      date: oldSchedule ? getUTCDateFromString(oldSchedule?.outbound) : "",
      time: oldSchedule ? getTimeFromString(oldSchedule?.outbound) : "",
      economy_price: oldSchedule?.economy_price || "",
    });
  }, [oldSchedule]);

  function handleChange(event) {
    const { name, value } = event.target;
    setFormData((prevData) => ({ ...prevData, [name]: value || "" }));
  }

  function handleSubmit(event) {
    event.preventDefault();

    if (!allFieldsNotEmpty(formData)) {
      addToast("Форма не заполнена", "error");
      return;
    }

    const scheduleID = oldSchedule.id;
    const scheduleUpdateParams = {
      id: formData.id,
      date: formData.date,
      time: formData.time,
      economy_price: parseInt(formData.economy_price, 10),
    };

    updateScheduleByID(apiClient, scheduleID, scheduleUpdateParams)
      .then((response) => {
        if (response.status == 200) {
          addToast("Данные расписания успешно изменены");
          changeScheduleInState(scheduleUpdateParams);
        }
      })
      .catch((error) => {
        if (error.response) {
          addToast(`Ошибка: ${error.response.data}`, "error");
        } else {
          addToast("Ошибка сети", "error");
        }
      });
  }

  return (
    <>
      <Modal
        show={show}
        handleClose={handleClose}
        title="Редактирование расписания"
      >
        <form className="form edit-schedule-form">
          <fieldset className="form__fieldset">
            <legend className="form__legend">Маршрут рейса</legend>
            <label className="form__label">
              Откуда: {oldSchedule?.route?.departure_airport?.iata_code}
            </label>
            <label className="form__label">
              Куда: {oldSchedule?.route?.arrival_airport?.iata_code}
            </label>
            <label className="form__label">
              Тип самолета: {oldSchedule?.aircraft?.name}
            </label>
          </fieldset>
          <div className="form__layout">
            <div className="form__container">
              <label className="form__label" htmlFor="edit-date-input">
                Дата:
              </label>
              <input
                className="form__input"
                id="edit-date-input"
                type="date"
                name="date"
                value={formData.date || ""}
                onChange={handleChange}
              />
            </div>
            <div className="form__container">
              <label className="form__label" htmlFor="edit-time-input">
                Время:
              </label>
              <input
                className="form__input"
                type="time"
                id="edit-time-input"
                name="time"
                value={formData.time || ""}
                onChange={handleChange}
              />
            </div>
            <div className="form__container">
              <label className="form__label" htmlFor="edit-economy-input">
                Цена эконома: $
              </label>
              <input
                type="number"
                className="form__input"
                id="edit-economy-input"
                name="economy_price"
                value={formData.economy_price || ""}
                onChange={handleChange}
              />
            </div>

            <div className="form__button-box">
              <Button
                color="green"
                type="submit"
                onClick={handleSubmit}
                disabled={!allFieldsNotEmpty(formData)}
              >
                Отправить
              </Button>
              <Button color="red" onClick={handleClose}>
                Отменить
              </Button>
            </div>
          </div>
        </form>
      </Modal>
    </>
  );
}

function ImportFileModal({ show, handleClose }) {
  const [file, setFile] = useState(null);
  const [results, setResults] = useState({
    success: 0,
    duplicates: 0,
    missingFields: 0,
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const { apiClient } = useApi();
  const { addToast } = useToast();

  // Обработчик выбора файла
  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  // Обработчик отправки формы
  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!file) {
      addToast("Выберите файл", error);
      return;
    }

    setIsSubmitting(true);

    UploadFileWithSchedules(apiClient, file)
      .then((response) => {
        const {
          successful_rows_cnt,
          duplicated_rows_cnt,
          missing_fields_rows_cnt,
        } = response.data;

        setResults({
          success: successful_rows_cnt,
          duplicates: duplicated_rows_cnt,
          missingFields: missing_fields_rows_cnt,
        });

        addToast("Файл успешно импортирован");
      })
      .catch((error) => {
        console.error("Error uploading file:", error);
        addToast("Не удалось импортировать файл", "error");
      })
      .finally(() => {
        setIsSubmitting(false);
      });
  };

  return (
    <Modal show={show} handleClose={handleClose}>
      <form onSubmit={handleSubmit} className="form__container">
        <label htmlFor="file" className="form__label">
          Please select the text file with the changes:
        </label>
        <input
          type="file"
          id="file"
          name="file"
          className="form__input"
          onChange={handleFileChange}
          accept=".csv"
        />

        <div className="form__button-box">
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Importing..." : "Import"}
          </Button>
        </div>

        {/* Отображение результатов */}
        <fieldset className="form__fieldset">
          <legend className="form__legend">Results</legend>
          <p className="common-text">
            Successful Changes Applied: {results.success}
          </p>
          <p className="common-text">
            Duplicate Records Discarded: {results.duplicates}
          </p>
          <p className="common-text">
            Records with Missing Fields Discarded: {results.missingFields}
          </p>
        </fieldset>
      </form>
    </Modal>
  );
}

export default SchedulePage;
