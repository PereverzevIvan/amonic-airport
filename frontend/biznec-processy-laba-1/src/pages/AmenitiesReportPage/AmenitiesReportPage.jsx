import Button from "../../components/Button/Button";
import "./AmenitiesReportPage.scss";
import { useToast } from "../../context/ToastContext";
import { useApi } from "../../context/apiContext";
import { useState, useEffect } from "react";
import { getAmenities, getAmenitiesReport } from "../../api/aminitiesApi";
import { getAirports } from "../../api/airportApi";
import { createMapFromList } from "../../global/mapUtils/mapUtils";

function AmenitiesReportPage() {
  const { apiClient } = useApi();
  const { addToast } = useToast();

  const [amenities, setAmenities] = useState([]);
  const [flightNumber, setFlightNumber] = useState("");
  const [report, setReport] = useState(null);
  const [airports, setAirports] = useState([]);
  const [airportsMap, setAirportsMap] = useState(null);
  const [fromDate, setFromDate] = useState("");
  const [toDate, setToDate] = useState("");

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

  function fetchReport() {
    let body = {
      flight_number: flightNumber,
      from_date: fromDate,
      to_date: toDate,
    };

    console.log(body);

    getAmenitiesReport(apiClient, body)
      .then((response) => {
        if (response.status == 200) {
          console.log("Отчет по услугам получен");
          addToast("Отчет по услугам получен");
          console.log(response.data);
          setReport(response.data.amenities_count);
        }
      })
      .catch((error) => {
        console.error(error);
        if (error.response) {
          addToast(error.response.data, "error");
        }
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
    console.log("report", report);
  }, [report]);
  return (
    <>
      <section className="amenities-report-page">
        <form
          className="form"
          onSubmit={(e) => {
            e.preventDefault();
            fetchReport();
          }}
        >
          <fieldset className="form__fieldset">
            <legend className="form__legend">Введите информацию о рейсе</legend>
            <div className="form__container">
              <label className="form__label">Номер бронирования: </label>
              <input
                className="form__input"
                type="text"
                value={flightNumber}
                onChange={(e) => {
                  setFlightNumber(e.target.value);
                }}
              />
            </div>
            <div className="form__container">
              <label className="form__label" htmlFor="from">
                Дата отправления
              </label>
              <input
                className="form__input"
                type="date"
                id="from"
                name="outboundDate"
                value={fromDate}
                onChange={(e) => {
                  setFromDate(e.target.value);
                }}
                required
              />
            </div>
            <div className="form__container">
              <label className="form__label" htmlFor="to">
                Дата прибытия
              </label>
              <input
                className="form__input"
                type="date"
                id="to"
                name="outboundDate"
                value={toDate}
                onChange={(e) => {
                  setToDate(e.target.value);
                }}
                required
              />
            </div>
          </fieldset>
          <Button type="submit">Получить отчет</Button>
        </form>

        {report && (
          <div className="table-container">
            <table className="table">
              <thead className="table__head">
                <tr className="table__row">
                  <th className="table__header">Название услуги</th>
                  <th className="table__header">Эконом</th>
                  <th className="table__header">Бизнес</th>
                  <th className="table__header">Первый класс</th>
                </tr>
              </thead>
              <tbody className="table__body">
                {Object.entries(report).map(([key, value]) => {
                  return (
                    <tr className="table__row">
                      <td className="table__data">{key}</td>
                      <td className="table__data">{value[0]}</td>
                      <td className="table__data">{value[1]}</td>
                      <td className="table__data">{value[2]}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </>
  );
}

export default AmenitiesReportPage;
