import React, { useEffect, useState } from "react";
import "./AirlinesShortSummaryPage.scss";
import { getAirlineSummary } from "../../api/summaryApi";
import { useApi } from "../../context/apiContext";
import { subtractDays } from "../../global/dateUtils/dateUtils";

const AirlinesShortSummaryPage = () => {
  const { apiClient } = useApi();

  const [summary, setSummary] = useState(null);

  function fetchSummary() {
    let body = {
      start_date: "2017-11-29",
    };

    getAirlineSummary(apiClient, body)
      .then((response) => {
        if (response.status == 200) {
          console.log("Суммарный отчет получен");
          console.log(response.data);
          setSummary(response.data);
        }
      })
      .catch((error) => {
        console.error(error);
        if (error.response) {
          console.log(error.response.data);
        }
      });
  }

  useEffect(() => fetchSummary(), []);

  return (
    <section className="summary-page">
      {summary != null && (
        <div className="summary-report">
          <div className="content">
            <section className="section flights">
              <p className="common-text common-text_big">Flights</p>
              <p className="common-text">
                Number confirmed: [ {summary.number_confirmed_flights} ]
              </p>
              <p className="common-text">
                Number canceled: [ {summary.number_cancelled_flights} ]
              </p>
              <p className="common-text">
                Average daily flight time: [{" "}
                {summary.average_daily_flight_time_minutes} ] minutes
              </p>
            </section>

            <section className="section passengers">
              <p className="common-text common-text_big">
                Number of passengers flying
              </p>
              <p className="common-text">
                Busiest day: [ {summary.busiest_day} ] with [{" "}
                {summary.busiest_day_number_of_passengers} ] flying
              </p>
              <p className="common-text">
                Most quiet day: [ {summary.most_quiet_day} ] with [{" "}
                {summary.most_quiet_day_number_of_passengers} ] flying
              </p>
            </section>

            <section className="section top-customers">
              <p className="common-text common-text_big">
                Top Customers (Number of purchases)
              </p>
              <ol>
                {summary.top_customer_by_purchased_tickets.map((cuctomer) => (
                  <li className="common-text">{cuctomer}</li>
                ))}
              </ol>
            </section>

            <section className="section top-offices">
              <p className="common-text common-text_big">
                Top AMONIC Airlines Offices (Revenue)
              </p>
              <ol>
                {summary.top_offices.map((office) => (
                  <li className="common-text">{office}</li>
                ))}
              </ol>
            </section>

            <section className="section revenue">
              <p className="common-text common-text_big">
                Revenue from ticket sales
              </p>
              <p className="common-text">
                Yesterday: $[ {summary.revenue_from_ticket_sales[0]} ]
              </p>
              <p className="common-text">
                Two days ago: $[ {summary.revenue_from_ticket_sales[1]} ]
              </p>
              <p className="common-text">
                Three days ago: $[ {summary.revenue_from_ticket_sales[2]} ]
              </p>
            </section>

            <section className="section empty-seats">
              <p className="common-text common-text_big">
                Weekly report of percentage of empty seats
              </p>
              <p className="common-text">
                This week: [{" "}
                {Math.round(
                  summary.weekly_report_of_percentage_of_empty_seats[0] * 100,
                )}{" "}
                ]%
              </p>
              <p className="common-text">
                Last week: [{" "}
                {Math.round(
                  summary.weekly_report_of_percentage_of_empty_seats[1] * 100,
                )}{" "}
                ]%
              </p>
              <p className="common-text">
                Two weeks ago: [{" "}
                {Math.round(
                  summary.weekly_report_of_percentage_of_empty_seats[2] * 100,
                )}{" "}
                ]%
              </p>
            </section>
          </div>

          <p className="common-text report-time">
            Report generated in [ {summary.time_taken_to_generate_report / 1000}{" "}
            ] seconds
          </p>
        </div>
      )}
    </section>
  );
};

export default AirlinesShortSummaryPage;
