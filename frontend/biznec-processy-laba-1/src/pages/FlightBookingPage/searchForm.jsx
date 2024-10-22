import { useEffect } from "react";
import Button from "../../components/Button/Button";

export function SearchForm({
  formData,
  setFormData,
  handleChange,
  airports,
  handleSubmit,
  setSelectedInbound,
}) {
  useEffect(() => {
    if (formData?.flightType === "one-way") {
      setFormData((prevValue) => ({ ...prevValue, inboundDate: "" }));
      setSelectedInbound(null);
    }
  }, [formData?.flightType]);

  return (
    <>
      <form onSubmit={handleSubmit} className="form">
        <fieldset className="form__fieldset">
          <legend className="form__legend">Поиск рейсов</legend>

          {/* Аэропорт отправления */}
          <div className="form__container">
            <label className="form__label" htmlFor="from">
              Откуда
            </label>
            <select
              className="form__input"
              id="from"
              name="from"
              value={formData?.from}
              onChange={handleChange}
              required
            >
              <option value="">Выберите аэропорт</option>
              {airports.map((airport) => (
                <option key={airport.id} value={airport.id}>
                  {airport.iata_code} - {airport.name}
                </option>
              ))}
            </select>
          </div>

          {/* Аэропорт прибытия */}
          <div className="form__container">
            <label className="form__label" htmlFor="to">
              Куда
            </label>
            <select
              className="form__input"
              id="to"
              name="to"
              value={formData?.to}
              onChange={handleChange}
              required
            >
              <option value="">Выберите аэропорт</option>
              {airports.map((airport) => (
                <option key={airport.id} value={airport.id}>
                  {airport.iata_code} - {airport.name}
                </option>
              ))}
            </select>
          </div>

          {/* Дата отправления */}
          <div className="form__container">
            <label className="form__label" htmlFor="departureDate">
              Отправление
            </label>
            <input
              className="form__input"
              type="date"
              id="departureDate"
              name="outboundDate"
              value={formData?.outboundDate}
              onChange={handleChange}
              required
            />
          </div>

          {/* Дата возвращения (если туда-обратно) */}
          <div className="form__container">
            <label className="form__label" htmlFor="returnDate">
              Возвращение
            </label>
            <input
              disabled={formData?.flightType !== "round-trip"}
              className="form__input"
              type="date"
              id="returnDate"
              name="inboundDate"
              value={
                formData?.flightType === "round-trip"
                  ? formData?.inboundDate
                  : ""
              }
              onChange={handleChange}
              required
            />
          </div>

          {/* Тип рейса */}
          <div className="form__container">
            <label className="form__label" htmlFor="flightType">
              Тип рейса
            </label>
            <select
              className="form__input"
              id="flightType"
              name="flightType"
              value={formData?.flightType}
              onChange={handleChange}
              required
            >
              <option value="one-way">В одну сторону</option>
              <option value="round-trip">Туда и обратно</option>
            </select>
          </div>
          {/* Класс билета */}
          <div className="form__container">
            <label className="form__label" htmlFor="cabinType">
              Класс билета
            </label>
            <select
              className="form__input"
              id="cabinType"
              name="cabinType"
              value={formData?.cabinType}
              onChange={handleChange}
              required
            >
              <option value="economy">Эконом</option>
              <option value="business">Бизнес</option>
              <option value="first-class">Первый класс</option>
            </select>
          </div>

          <div className="form__button-box"></div>
        </fieldset>
        <Button type="submit">Найти рейсы</Button>
      </form>
    </>
  );
}
