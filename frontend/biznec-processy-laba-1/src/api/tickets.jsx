export function remainingSeatsCount(apiClient, body = { schedule_ids: [] }) {
  return apiClient.post("/tickets/remaining-seats-count", body);
}

export function bookingTickets(
  apiClient,
  body = {
    cabin_type: 0,
    inbound_schedules: [0],
    outbound_schedules: [0],
    passengers: [
      {
        birthday: "string",
        first_name: "string",
        last_name: "string",
        passport_country_id: 0,
        passport_number: "string",
        phone: "string",
      },
    ],
  },
) {
  return apiClient.post("/tickets/booking", body);
}

export function confirmTickets(
  apiClient,
  body = {
    tickets: [0],
  },
) {
  return apiClient.post("/tickets/confirm", body);
}

export function getTicketsByBooking(apiClient, booking) {
  return apiClient.get("/tickets", { params: { booking_reference: booking } });
}
