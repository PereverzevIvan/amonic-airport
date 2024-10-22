import { useEffect, useState } from "react";
import Modal from "../../components/Modal/Modal.jsx";
import Button from "../../components/Button/Button.jsx";
import { getDateFromString } from "../../global/dateUtils/dateUtils.jsx";
import { allFieldsNotEmpty } from "../../global/formUtils/formUtils.jsx";
import { useToast } from "../../context/ToastContext.jsx";
import { getScheduleIDsInFlight } from "../../global/scheduleUtils/scheduleUtils.jsx";
import { bookingTickets } from "../../api/tickets.jsx";
import { useApi } from "../../context/apiContext.jsx";
import { ConfirmPaymentModal } from "./confirmPaymentModal.jsx";

function createMapFromList(list) {
  let map = new Map();
  list.forEach((element) => {
    map.set(`${element.id}`, element);
  });

  return map;
}

export function SetPassengersModal({
  show,
  handleClose,
  maxPassangersCount,
  outboundFlight,
  inboundFlight,
  countries,
  airportsMap,
  cabinType,
}) {
  const { addToast } = useToast();
  const { apiClient } = useApi();

  const countriesMap = createMapFromList(countries);

  const [showConfirmPaymentModal, setShowConfirmPaymentModal] = useState(false);

  const [passengers, setPassengers] = useState([]);
  const [tickets, setTickets] = useState([]);
  const [passenger, setPassenger] = useState({
    firstName: "",
    lastName: "",
    birthdate: "",
    passportNumber: "",
    passportCountry: "",
    phone: "",
  });

  function clearForm() {
    setPassenger({
      firstName: "",
      lastName: "",
      birthdate: "",
      passportNumber: "",
      passportCountry: "",
      phone: "",
    });
  }

  function clearPassengers() {
    setPassengers([]);
  }

  useEffect(() => {
    if (!show) {
      clearForm();
      clearPassengers();
      setTickets([]);
    }
  }, [show]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setPassenger({ ...passenger, [name]: value });
  };

  const addPassenger = () => {
    if (passenger.firstName && passenger.lastName && passenger.passportNumber) {
      setPassengers([...passengers, passenger]);
      clearForm();
    }
  };

  const removePassenger = (index) => {
    setPassengers(passengers.filter((_, i) => i !== index));
    setTickets([]);
  };

  function registerTickets() {
    if (tickets.length > 0) {
      setShowConfirmPaymentModal(true);
      return;
    }

    let outboundScheduleIDs = getScheduleIDsInFlight(outboundFlight);
    let inboundScheduleIDs = getScheduleIDsInFlight(inboundFlight);
    let cbType = 0;
    if (cabinType == "economy") cbType = 1;
    else if (cabinType == "business") cbType = 2;
    else if (cabinType == "first-class") cbType = 3;

    let requestBody = {
      cabin_type: cbType,
      inbound_schedules: inboundScheduleIDs,
      outbound_schedules: outboundScheduleIDs,
      passengers: passengers.map((passenger, index) => ({
        birthday: passenger.birthdate,
        first_name: passenger.firstName,
        last_name: passenger.lastName,
        passport_country_id: parseInt(passenger.passportCountry, 10),
        passport_number: passenger.passportNumber,
        phone: passenger.phone,
      })),
    };
    bookingTickets(apiClient, requestBody)
      .then((response) => {
        console.log(response);
        if (response.status == 201) {
          console.log(response.data);
          setTickets(response.data.tickets);
          setShowConfirmPaymentModal(true);
          addToast("Билеты созданы");
        }
      })
      .catch((error) => {
        console.log(error);
        if (error.response) {
          addToast(error.response.data, "error");
        }
      });
  }

  return (
    <>
      <ConfirmPaymentModal
        show={showConfirmPaymentModal}
        handleClose={() => {
          setShowConfirmPaymentModal(false);
        }}
        outboundFlight={outboundFlight}
        inboundFlight={inboundFlight}
        tickets={tickets}
        sub={true}
        closeParentModal={handleClose}
      />

      <Modal
        show={show}
        handleClose={handleClose}
        title="Подтверждение бронирования"
      >
        <div>
          {tickets.length > 0 && <p className="common-text">Билеты созданы</p>}
          <form className="form">
            <fieldset className="form__fieldset flights-list">
              <legend className="form__legend">Исходящие полеты</legend>
              {outboundFlight?.schedules?.length > 0 ? (
                outboundFlight?.schedules?.map((schedule, index) => (
                  <div className="flights-list__item" key={index}>
                    <label className="form__label">
                      From:{" "}
                      {
                        airportsMap.get(
                          `${schedule?.route?.departure_airport_id}`,
                        )?.iata_code
                      }
                    </label>
                    <label className="form__label">
                      To:{" "}
                      {
                        airportsMap.get(`${schedule?.route.arrival_airport_id}`)
                          ?.iata_code
                      }
                    </label>
                    <label className="form__label">
                      Date: {getDateFromString(schedule.outbound)}
                    </label>
                    <label className="form__label">
                      Flight number: {schedule.flight_number}
                    </label>
                  </div>
                ))
              ) : (
                <p>Рейсы не выбраны</p>
              )}
            </fieldset>
          </form>

          {/* Отображение обратных рейсов */}
          {inboundFlight && (
            <form className="form">
              <fieldset className="form__fieldset flights-list">
                <legend className="form__legend">Обратные полеты</legend>
                {inboundFlight?.schedules?.length > 0 ? (
                  inboundFlight.schedules.map((schedule, index) => (
                    <div className="flights-list__item" key={index}>
                      <label className="form__label">
                        From:{" "}
                        {
                          airportsMap.get(
                            `${schedule?.route?.departure_airport_id}`,
                          )?.iata_code
                        }
                      </label>
                      <label className="form__label">
                        To:{" "}
                        {
                          airportsMap.get(
                            `${schedule?.route.arrival_airport_id}`,
                          )?.iata_code
                        }
                      </label>
                      <label className="form__label">
                        Date: {getDateFromString(schedule.outbound)}
                      </label>
                      <label className="form__label">
                        Flight number: {schedule.flight_number}
                      </label>{" "}
                    </div>
                  ))
                ) : (
                  <p>Рейсы не выбраны</p>
                )}
              </fieldset>
            </form>
          )}
        </div>

        {/* Passenger details input */}
        <form className="form">
          <fieldset className="form__fieldset">
            <legend className="form__legend">Информация о пассажире</legend>
            <input
              className="form__input"
              type="text"
              name="firstName"
              placeholder="Firstname"
              value={passenger.firstName}
              onChange={handleInputChange}
            />
            <input
              className="form__input"
              type="text"
              name="lastName"
              placeholder="Lastname"
              value={passenger.lastName}
              onChange={handleInputChange}
            />
            <input
              className="form__input"
              type="date"
              name="birthdate"
              placeholder="Birthdate"
              value={passenger.birthdate}
              onChange={handleInputChange}
            />
            <input
              className="form__input"
              type="text"
              name="passportNumber"
              placeholder="Passport number"
              value={passenger.passportNumber}
              onChange={handleInputChange}
            />
            <select
              className="drop-down-list form__input"
              name="passportCountry"
              value={passenger.passportCountry}
              onChange={handleInputChange}
            >
              <option value="">Выберите страну</option>
              {countries.length > 0 &&
                countries.map((country, index) => (
                  <option value={country.id} key={index}>
                    {country.name}
                  </option>
                ))}
            </select>

            <input
              className="form__input"
              type="phone"
              name="phone"
              placeholder="Phone"
              value={passenger.phone}
              onChange={handleInputChange}
            />
          </fieldset>
          <Button
            disabled={!allFieldsNotEmpty(passenger)}
            onClick={() => {
              if (passengers.length < maxPassangersCount) {
                addPassenger();
              } else {
                addToast("Вы не можете добавить больше пассажиров", "error");
              }
            }}
          >
            Добавить
          </Button>
        </form>

        {/* Passengers list */}
        {passengers.length > 0 && (
          <div className="table-container">
            <table className="table">
              <thead className="table__head">
                <tr className="table__row">
                  <th className="table__header">Имя</th>
                  <th className="table__header">Фамилия</th>
                  <th className="table__header">Дата рождения</th>
                  <th className="table__header">Номер паспорта</th>
                  <th className="table__header">Страна</th>
                  <th className="table__header">Номер телефона</th>
                  <th className="table__header">Действие</th>
                </tr>
              </thead>
              <tbody className="table__body">
                {passengers.map((p, index) => (
                  <tr className="table__row" key={index}>
                    <td className="table__data">{p.firstName}</td>
                    <td className="table__data">{p.lastName}</td>
                    <td className="table__data">{p.birthdate}</td>
                    <td className="table__data">{p.passportNumber}</td>
                    <td className="table__data">
                      {countriesMap.get(p.passportCountry).name}
                    </td>
                    <td className="table__data">{p.phone}</td>
                    <td className="table__data">
                      <Button
                        color="red"
                        onClick={() => removePassenger(index)}
                      >
                        Удалить
                      </Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        <div className="button-box">
          <Button
            disabled={passengers?.length != maxPassangersCount}
            color="green"
            onClick={registerTickets}
          >
            Подтвердить
          </Button>
          <Button onClick={handleClose}>Вернуться</Button>
        </div>
      </Modal>
    </>
  );
}
