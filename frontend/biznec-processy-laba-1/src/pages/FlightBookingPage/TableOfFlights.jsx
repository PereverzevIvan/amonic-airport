import {
  getDateFromString,
  getTimeFromString,
} from "../../global/dateUtils/dateUtils.jsx";

export function TableOfFlights({
  formData,
  handleChangeCheckbox,
  flights,
  selectedFlight,
  setSelectedFlight,
  inputName,
  title,
}) {
  const handleSelectFlight = (flightIndex) => {
    if (selectedFlight === flightIndex) {
      setSelectedFlight(null);
    } else {
      setSelectedFlight(flightIndex);
    }
  };

  return (
    <>
      <div className="under-table-line">
        <p className="common-text common-text_big">{title}</p>
        <label className="form__label">
          <input
            type="checkbox"
            name={inputName}
            value={
              inputName === "outboundWindow"
                ? formData?.outboundWindow
                : formData?.inboundWindow
            }
            onChange={handleChangeCheckbox}
          />
          Окно в 3 дня до и после
        </label>
      </div>

      <div className="table-container">
        <table className="table">
          <thead className="table__head">
            <tr className="table__row">
              <th className="table__header">Дата</th>
              <th className="table__header">Время</th>
              <th className="table__header">Откуда</th>
              <th className="table__header">Куда</th>
              <th className="table__header">Номер(а) рейса(ов)</th>
              <th className="table__header">Цена билета</th>
              <th className="table__header">Пересадки</th>
            </tr>
          </thead>
          <tbody className="table__body">
            {flights?.map((flight, index) => (
              <tr
                key={index}
                className={`table__row ${selectedFlight === index ? "blue" : ""}`}
                onClick={() => handleSelectFlight(index)}
              >
                <td className="table__data">
                  {getDateFromString(flight?.schedules[0]?.outbound)}
                </td>
                <td className="table__data">
                  {getTimeFromString(flight?.schedules[0]?.outbound)}
                </td>
                <td className="table__data">{flight.departure}</td>
                <td className="table__data">{flight.arrival}</td>
                <td className="table__data">{flight?.flightNumbers}</td>
                <td className="table__data">${flight?.price}</td>
                <td className="table__data">{flight?.TransfersCount}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {!flights && <p>Рейсы не найдены. Пожалуйста, выполните поиск.</p>}
      </div>
    </>
  );
}
