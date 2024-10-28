import React, { useEffect, useState } from "react";
import "./AmenitiesPage.scss";
import {
  getAmenities,
  getAmenitiesForCabinType,
  getAmenitiesForTicket,
  sendAmenitiesForTicket,
} from "../../api/aminitiesApi";
import { useToast } from "../../context/ToastContext";
import { useApi } from "../../context/apiContext";
import Button from "../../components/Button/Button";
import { getTicketsByBooking } from "../../api/tickets";
import { createMapFromList } from "../../global/mapUtils/mapUtils.jsx";
import { getAirports } from "../../api/airportApi.jsx";
import {
  getDateFromString,
  getTimeFromString,
} from "../../global/dateUtils/dateUtils.jsx";

const cabinTypes = {
  1: "Эконом",
  2: "Бизнес",
  3: "Первый класс",
};

const AmenitiesPage = () => {
  const { apiClient } = useApi();
  const { addToast } = useToast();

  const [amenities, setAmenities] = useState([]);
  const [defaultAmenities, setDefaultAmenities] = useState([]);
  const [amenitiesForTicket, setAmenitiesForTicket] = useState([]);
  const [showAmenitiesForTicket, setShowAmenitiesForTicket] = useState(false);
  const [selectedAmenityIDs, setSelectedAmenityIDs] = useState([]);
  const [bookingNumber, setBookingNumber] = useState("");
  const [tickets, setTickets] = useState([]);
  const [selectedTicket, setSelectedTicket] = useState("");
  const [passengerInfo, setPassengerInfo] = useState(null);
  const [airports, setAirports] = useState([]);
  const [airportsMap, setAirportsMap] = useState(null);

  const [finalAmount, setFinalAmount] = useState({
    itemsCount: 0,
    duties: 0,
    total: 0,
  });

  const fetchAirports = async () => {
    try {
      const response = await getAirports(apiClient);
      setAirports(response.data);
      console.log("Аэропорты успешно получены");
    } catch (error) {
      console.error("Ошибка при получении списка аэропортов:", error);
    }
  };

  function fetchAminities() {
    getAmenities(apiClient)
      .then((response) => {
        console.log("Получены все дополнительные услуги");
        if (response.status == 200) {
          setAmenities(response.data);
        }
      })
      .catch((error) => {
        console.error(error);
        if (error.response) {
          addToast(error.response.data, "error");
        }
      });
  }

  function fetchAmenitiesForTicket() {
    if (selectedTicket !== "")
      getAmenitiesForTicket(apiClient, parseInt(selectedTicket, 10))
        .then((response) => {
          if (response.status == 200) {
            setAmenitiesForTicket(response.data);
            setShowAmenitiesForTicket(true);
            console.log(response.data);
            console.log("Ранее выбранные для данного билета услуги получены");
          }
        })
        .catch((error) => {
          if (error.response) addToast(error.response.data);
          console.error(error);
        });
    setSelectedAmenityIDs([]);
  }

  function fetchDefaultAmenities() {
    if (tickets?.length == 0) return;
    getAmenitiesForCabinType(apiClient, tickets[0].cabin_type_id)
      .then((response) => {
        if (response.status == 200) {
          setDefaultAmenities(response.data);
          console.log("Дефолтные для данного класса услуги получены");
        }
      })
      .catch((error) => {
        if (error.response) addToast(error.response.data);
        console.error(error);
      });
  }

  function setPassengerInfoFromData(data) {
    let info = {
      fullName: data.first_name + " " + data.last_name,
      passportNumber: data.passport_number,
      cabinType: cabinTypes[data.cabin_type_id],
    };

    setPassengerInfo(info);
  }

  function searchTickets() {
    setTickets([]);
    setDefaultAmenities([]);
    setShowAmenitiesForTicket(false);
    setSelectedTicket("");
    setSelectedAmenityIDs([]);
    getTicketsByBooking(apiClient, bookingNumber)
      .then((response) => {
        if (response.status == 200) {
          setTickets(response.data);
          setPassengerInfoFromData(response.data[0]);
          addToast("Билеты найдены");
        }
      })
      .catch((error) => {
        console.error(error);
        if (error.response) {
          addToast(error.response.data, "error");
        }
      });
  }

  function getScheduleDescription(schedule) {
    if (!schedule) return;
    let depAirportCode = airportsMap.get(
      `${schedule.route.departure_airport_id}`,
    ).iata_code;
    let arrAirportCode = airportsMap.get(
      `${schedule.route.arrival_airport_id}`,
    ).iata_code;

    return `${schedule.flight_number}, ${depAirportCode}-${arrAirportCode}, ${getDateFromString(schedule.outbound)}, ${getTimeFromString(schedule.outbound)}`;
  }

  function handleChangeSelectedAmenityIDs(event) {
    const { value, checked } = event.target;

    if (checked)
      setSelectedAmenityIDs((prevAllowedAmenities) => [
        ...prevAllowedAmenities,
        value,
      ]);
    else
      setSelectedAmenityIDs((prevAmenities) =>
        prevAmenities.filter((amenity_id) => amenity_id != value),
      );
  }

  // Заполнение списка выбранных услуг бесплатными услугами после получения списка разрешенных для данного рейса услуг
  function addToSelectedFreeAmenities() {
    if (defaultAmenities) {
      setSelectedAmenityIDs((prevAmenities) => [
        ...prevAmenities,
        ...amenities
          .filter(
            (amenity) =>
              amenity.price === 0 &&
              !selectedAmenityIDs.includes(`${amenity.id}`),
          )
          .map((amenity) => `${amenity.id}`),
      ]);
    }
  }

  // Добавление ранее выбранных услуг в список выбранных В ДАННЫЙ момент
  function addToSelectedPreviouslySelectedAmenities() {
    if (defaultAmenities && amenitiesForTicket) {
      setSelectedAmenityIDs((prevAmenities) => [
        ...prevAmenities,
        ...amenitiesForTicket
          .map((amenity) => `${amenity.amenity_id}`)
          .filter(
            (amenity) => !selectedAmenityIDs.includes(amenity.amenity_id),
          ),
      ]);
    }
  }

  function addToSelectedDefaultAmenities() {
    if (defaultAmenities) {
      setSelectedAmenityIDs((prevAmenities) => [
        ...prevAmenities,
        ...defaultAmenities
          .map((amenityID) => `${amenityID}`)
          .filter((amenityID) => !selectedAmenityIDs.includes(amenityID)),
      ]);
    }
  }

  function handleSendAmenitiesForTicket(event) {
    event.preventDefault();
    if (selectedTicket == "") {
      addToast("Билет не выбран", "error");
    }
    let body = {
      ticket_id: parseInt(selectedTicket, 10),
      amenity_ids: selectedAmenityIDs.map((amenityID) =>
        parseInt(amenityID, 10),
      ),
    };

    console.log(body);
    sendAmenitiesForTicket(apiClient, body)
      .then((response) => {
        if (response.status == 200) {
          console.log("Услуги для выбранного рейса успешно изменены");
          addToast("Услуги для выбранного рейса успешно изменены");
          fetchAmenitiesForTicket();
        }
      })
      .catch((error) => {
        if (error.response) addToast(error.response.data, "error");
      });
  }

  useEffect(() => {
    fetchAminities();
    fetchAirports();
  }, []);
  useEffect(() => {
    setAirportsMap(createMapFromList(airports));
  }, [airports]);
  useEffect(() => {
    fetchDefaultAmenities();
  }, [tickets]);
  useEffect(() => {
    //addToSelectedFreeAmenities();
    addToSelectedPreviouslySelectedAmenities();
    //addToSelectedDefaultAmenities();
  }, [amenitiesForTicket]);
  useEffect(() => {
    // Автоматический подсчет стоимости всех услуг
    let prevItemsPrice = 0;
    let curItemsPrice = 0;

    amenitiesForTicket.forEach((amenity) => (prevItemsPrice += amenity.price));

    amenities.forEach(
      (amenity) =>
        (curItemsPrice += selectedAmenityIDs.includes(`${amenity.id}`)
          ? amenity.price
          : 0),
    );

    let obj = {
      itemsCount: curItemsPrice - prevItemsPrice,
      duties: 0,
      total: 0,
    };

    obj.duties = Math.round(obj.itemsCount * 0.05 * 100) / 100;
    obj.total = obj.itemsCount + obj.duties;

    setFinalAmount(obj);
  }, [selectedAmenityIDs]);

  return (
    <>
      <section className="amenities-page">
        {/* Поиск по номеру бронирования */}
        <form className="form">
          <fieldset className="form__fieldset">
            <legend className="form__legend">
              Поиск по номеру бронирования
            </legend>
            <label className="form__label">Номер бронирования: </label>
            <input
              className="form__input"
              type="text"
              value={bookingNumber}
              onChange={(e) => {
                setBookingNumber(e.target.value);
              }}
            />
            <Button
              color="blue"
              disabled={bookingNumber == ""}
              onClick={searchTickets}
            >
              Найти
            </Button>
          </fieldset>
        </form>

        {/* Выбор рейса */}
        {tickets.length > 0 && defaultAmenities.length > 0 && (
          <form className="form">
            <fieldset className="form__fieldset">
              <legend className="form__legend">Выбор рейса</legend>
              <label className="form__label">Рейс: </label>
              <select
                className="form__input"
                type="text"
                value={selectedTicket}
                onChange={(e) => {
                  setSelectedTicket(e.target.value);
                }}
              >
                <option value="">Выберите рейс</option>
                {tickets.map((ticket, index) => {
                  let schedule = ticket.schedule;
                  return (
                    <option key={index} value={ticket.id}>
                      {getScheduleDescription(schedule)}
                    </option>
                  );
                })}
              </select>
              <Button
                color="blue"
                disabled={selectedTicket == ""}
                onClick={fetchAmenitiesForTicket}
              >
                Показать услуги
              </Button>
            </fieldset>
          </form>
        )}

        {/* Информация о пассажире*/}
        {tickets.length > 0 && (
          <form className="form">
            <fieldset className="form__fieldset form__fieldset_between">
              <label className="form__label">
                Полное имя:{" "}
                <span style={{ fontWeight: "bold" }}>
                  {passengerInfo.fullName}
                </span>
              </label>
              <label className="form__label">
                Номер паспорта:{" "}
                <span style={{ fontWeight: "bold" }}>
                  {passengerInfo.passportNumber}
                </span>
              </label>
              <label className="form__label">
                Тип кабины:{" "}
                <span style={{ fontWeight: "bold" }}>
                  {passengerInfo.cabinType}
                </span>
              </label>
            </fieldset>
          </form>
        )}

        {/* Выбор дополнительных услуг для выбранного рейса (билета) */}
        {amenities.length > 0 &&
          defaultAmenities.length > 0 &&
          showAmenitiesForTicket && (
            <form className="form">
              <fieldset className="form__fieldset form__fieldset_between">
                {amenities.map((amenity, index) => {
                  return (
                    <label
                      key={index}
                      className="form__label"
                      htmlFor={`amenity-${amenity.id}`}
                    >
                      <input
                        type="checkbox"
                        id={`amenity-${amenity.id}`}
                        disabled={
                          amenity.price == 0 ||
                          defaultAmenities.includes(amenity.id)
                        }
                        checked={
                          selectedAmenityIDs.includes(`${amenity.id}`) ||
                          amenity.price == 0 ||
                          defaultAmenities.includes(amenity.id)
                        }
                        value={amenity.id}
                        onChange={handleChangeSelectedAmenityIDs}
                      />{" "}
                      {amenity.service} ({" "}
                      {amenity.price == 0 ? "Бесплатно" : `$${amenity.price}`} )
                    </label>
                  );
                })}
              </fieldset>
              <div className="summary-amenities">
                <div className="summary-amenities__info">
                  <label>
                    Общая стоимость услуг: {finalAmount.itemsCount}$
                  </label>
                  <label>Пошлины и налоги: {finalAmount.duties}$</label>
                  <label>Конечная сумма: {finalAmount.total}$</label>
                  {finalAmount.itemsCount < 0 && (
                    <label>
                      Вам вернется: {Math.abs(finalAmount.itemsCount)}$
                    </label>
                  )}
                </div>
                <div className="summary-amenities__buttons">
                  <Button
                    color="blue"
                    type="submit"
                    onClick={handleSendAmenitiesForTicket}
                  >
                    Сохранить и подтвердить
                  </Button>
                </div>
              </div>
            </form>
          )}
      </section>
    </>
  );
};

export default AmenitiesPage;
