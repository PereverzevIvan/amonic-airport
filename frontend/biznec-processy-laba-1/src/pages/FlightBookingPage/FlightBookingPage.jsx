import React, { useState, useEffect } from "react";
import "./FlightBookingPage.scss";
import { useApi } from "../../context/apiContext.jsx";
import { useToast } from "../../context/ToastContext.jsx";
import { getAirports } from "../../api/airportApi.jsx";
import Button from "../../components/Button/Button.jsx";
import { getFlights } from "../../api/schedulesApi.jsx";
import { remainingSeatsCount } from "../../api/tickets.jsx";
import { SetPassengersModal } from "./SetPassengersModal.jsx";
import { SearchForm } from "./searchForm.jsx";
import { TableOfFlights } from "./TableOfFlights.jsx";
import { getCountries } from "../../api/countriesApi.jsx";
import { getScheduleIDsInFlight } from "../../global/scheduleUtils/scheduleUtils.jsx";

let searchParamsExample = {
  from: 0,
  to: 0,
  outboundDate: "",
  inboundDate: "",
  outboundWindow: false,
  inboundWindow: false,
  flightType: "one-way",
  cabinType: "economy",
};

function calculateBisnessAndFirstClassPrice(economyPrice) {
  if (!economyPrice) return;
  const businessPrice = Math.floor(economyPrice * 1.35); // Билет в бизнес-класс на 35% дороже, округляем вниз
  const firstClassPrice = Math.floor(businessPrice * 1.3); // Билет в первый класс на 30% дороже, округляем вниз

  return {
    businessClass: businessPrice,
    firstClass: firstClassPrice,
  };
}

function createAirportsMap(airports) {
  let airportsMap = new Map();

  for (let i = 0; i < airports.length; i++) {
    airportsMap.set(`${airports[i].id}`, airports[i]);
  }

  return airportsMap;
}

const FlightBookingPage = () => {
  const { apiClient } = useApi();
  const { addToast } = useToast();

  const [airports, setAirports] = useState([]);
  const [countries, setCountries] = useState([]);
  const [airportsMap, setAirportsMap] = useState(null);

  useEffect(() => {
    fetchAirports();
    fetchCountries();
  }, []);
  useEffect(() => {
    setAirportsMap(createAirportsMap(airports));
  }, [airports]);

  const [searchParams, setSearchParams] = useState(searchParamsExample);
  const handleFormChange = (e) => {
    const { name, value, type, checked } = e.target;
    setSearchParams({
      ...searchParams,
      [name]: type === "checkbox" ? checked : value,
    });
  };

  const [outboundFlights, setOutboundFlights] = useState(null);
  const [selectedOutboundFlight, setSelectedOutboundFlight] = useState(null);
  const [inboundFlights, setInboundFlights] = useState(null);
  const [selectedInboundFlight, setSelectedInboundFlight] = useState(null);

  const [passengersCount, setPassengersCount] = useState(0);
  const [passengersInfo, setPassengersInfo] = useState(null);

  const [showSetPassengersModal, setShowSetPassengersModal] = useState(false);

  const fetchAirports = async () => {
    try {
      const response = await getAirports(apiClient);
      setAirports(response.data);
      console.log("Аэропорты успешно получены");
    } catch (error) {
      console.error("Ошибка при получении списка аэропортов:", error);
    }
  };

  const fetchCountries = async () => {
    try {
      const response = await getCountries(apiClient);
      setCountries(response.data);
      console.log("Страны успешно получены");
    } catch (error) {
      console.error("Ошибка при получении списка стран:", error);
    }
  };

  function getFlightNumbers(flight) {
    let numbers = "";

    if (flight?.length > 0) {
      for (let index = 0; index < flight.length; index++) {
        const element = flight[index];
        numbers += element.flight_number;
        if (index != flight.length - 1) numbers += "-";
      }
    }

    return numbers;
  }

  function getFlightTransfersCount(flight) {
    let transfers = 0;

    if (flight?.length > 0) {
      transfers += flight.length - 1;
    }

    return transfers;
  }

  function getPriceOfFlight(flight) {
    let price = 0;
    let cabinType = searchParams?.cabinType;

    if (flight?.length > 0) {
      for (let i = 0; i < flight.length; i++) {
        let schedule = flight[i];
        if (cabinType === "business") {
          let { businessClass } = calculateBisnessAndFirstClassPrice(
            schedule.economy_price,
          );
          price += businessClass;
        } else if (cabinType === "first-class") {
          let { firstClass } = calculateBisnessAndFirstClassPrice(
            schedule.economy_price,
          );
          price += firstClass;
        } else {
          price += schedule.economy_price;
        }
      }
    }

    return price;
  }

  function prepareFlights(flights, outbound = true) {
    let prFlights = [];

    if (flights?.length > 0) {
      for (let i = 0; i < flights.length; i++) {
        let oldFlight = flights[i];
        let newFlight = {
          cabinType: searchParams.cabinType,
          schedules: oldFlight,
          departure: airportsMap.get(
            outbound ? searchParams.from : searchParams.to,
          ).iata_code,
          arrival: airportsMap.get(
            outbound ? searchParams.to : searchParams.from,
          ).iata_code,
          flightNumbers: getFlightNumbers(oldFlight),
          price: getPriceOfFlight(oldFlight),
          TransfersCount: getFlightTransfersCount(oldFlight),
        };

        prFlights.push(newFlight);
      }
    }

    return prFlights;
  }

  const handleSearchSubmit = (e) => {
    e.preventDefault();

    const outbound = {
      from: parseInt(searchParams.from, 10),
      increase_search_interval: searchParams.outboundWindow,
      outbound_date: searchParams.outboundDate,
      to: parseInt(searchParams.to, 10),
    };

    const inbound = {
      from: parseInt(searchParams.to, 10),
      increase_search_interval: searchParams.inboundWindow,
      outbound_date: searchParams.inboundDate,
      to: parseInt(searchParams.from, 10),
    };

    const requestBody = {
      outbound: outbound,
      inbound: searchParams?.flightType === "one-way" ? null : inbound,
    };

    setSelectedInboundFlight(null);
    setSelectedOutboundFlight(null);

    getFlights(apiClient, requestBody)
      .then((response) => {
        if (response?.data) {
          if (response?.data?.outbound_flights)
            setOutboundFlights(prepareFlights(response.data.outbound_flights));
          else setOutboundFlights(null);

          if (response?.data?.inbound_flights)
            setInboundFlights(prepareFlights(response.data.inbound_flights));
          else setInboundFlights(null);

          addToast("Рейсы получены");
        }
      })
      .catch((error) => {
        if (error.response) {
          addToast(error.response.data, "error");
        }
        console.error(error);
        setOutboundFlights(null);
        setInboundFlights(null);
      });
  };

  function OpenSetPassengersModal() {
    let scheduleIDs = getScheduleIDsInFlight(
      outboundFlights[selectedOutboundFlight],
    );

    if (searchParams.flightType != "one-way") {
      if (compareDatesOfSelectedFlights()) {
        scheduleIDs.concat(
          getScheduleIDsInFlight(inboundFlights[selectedInboundFlight]),
        );
      } else {
        addToast(
          "Дата обратного рейса не может быть раньше исходящего рейса",
          "error",
        );
        return;
      }
    }
    //console.log(outboundFlights[selectedOutboundFlight]);
    //console.log(inboundFlights[selectedInboundFlight]);

    remainingSeatsCount(apiClient, { schedule_ids: scheduleIDs })
      .then((response) => {
        let remainingSeats = 0;
        if (searchParams.cabinType === "economy")
          remainingSeats = response.data.economy_seats;
        if (searchParams.cabinType === "business")
          remainingSeats = response.data.business_seats;
        if (searchParams.cabinType === "first-class")
          remainingSeats = response.data.first_class_seats;

        if (remainingSeats > passengersCount) {
          addToast("На выбранные рейсы есть свободные места");
          setShowSetPassengersModal(true);
          console.log("Окно открылось");
        } else addToast(`На выбранные рейсы не хватает мест`, "error");
      })
      .catch((error) => {
        if (error.response) addToast(error.response.data, "error");
        else console.error(error);
      });
  }

  function compareDatesOfSelectedFlights() {
    if (searchParams.flightType !== "one-way") {
      const outbound =
        selectedOutboundFlight !== null &&
        outboundFlights[selectedOutboundFlight];
      const inbound =
        selectedInboundFlight !== null && inboundFlights[selectedInboundFlight];

      if (inbound && outbound) {
        return (
          new Date(outbound.schedules[0].outbound) <
          new Date(inbound.schedules[0].outbound)
        );
      }
    }
  }

  return (
    <section className="flight-booking-page">
      <SetPassengersModal
        show={showSetPassengersModal}
        handleClose={() => {
          setShowSetPassengersModal(false);
        }}
        outboundFlight={
          selectedOutboundFlight !== null &&
          outboundFlights[selectedOutboundFlight]
        }
        inboundFlight={
          selectedInboundFlight !== null &&
          inboundFlights[selectedInboundFlight]
        }
        airportsMap={airportsMap}
        countries={countries}
        maxPassangersCount={passengersCount}
        cabinType={searchParams.cabinType}
      />

      <SearchForm
        airports={airports}
        handleSubmit={handleSearchSubmit}
        formData={searchParams}
        handleChange={handleFormChange}
        setFormData={setSearchParams}
        setSelectedInbound={setSelectedInboundFlight}
      />
      <TableOfFlights
        formData={searchParams}
        handleChangeCheckbox={handleFormChange}
        flights={outboundFlights}
        selectedFlight={selectedOutboundFlight}
        setSelectedFlight={setSelectedOutboundFlight}
        inputName={"outboundWindow"}
        title={"Вылетающий рейс"}
        airportsMap={airportsMap}
      />

      {searchParams.flightType !== "one-way" && (
        <TableOfFlights
          formData={searchParams}
          handleChangeCheckbox={handleFormChange}
          flights={inboundFlights}
          selectedFlight={selectedInboundFlight}
          setSelectedFlight={setSelectedInboundFlight}
          inputName={"inboundWindow"}
          title={"Обратный рейс"}
          airportsMap={airportsMap}
        />
      )}

      <form className="form passengers-form">
        <fieldset className="form__fieldset">
          <legend className="form__legend">Подтвердить бронирование для</legend>
          <div className="form__container">
            <input
              className="form__input"
              type="number"
              name="passengersCount"
              value={passengersCount}
              step={1}
              min={0}
              onChange={(e) => {
                let value = e.target.value;
                value = value >= 0 ? value : 0;
                setPassengersCount(value);
              }}
            />
            <label className="form__label">Пассажиров</label>
          </div>
          <Button
            color="green"
            disabled={
              passengersCount < 1 ||
              selectedOutboundFlight == null ||
              (searchParams.flightType !== "one-way" &&
                selectedInboundFlight == null)
            }
            onClick={() => {
              OpenSetPassengersModal(true);
            }}
          >
            Подтвердить
          </Button>
        </fieldset>
      </form>
    </section>
  );
};

export default FlightBookingPage;
