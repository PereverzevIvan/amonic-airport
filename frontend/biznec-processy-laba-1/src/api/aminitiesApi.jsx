export function getAmenities(apiClient) {
  return apiClient.get("/amenities");
}

export function getAmenitiesForCabinType(apiClient, cabinType) {
  return apiClient.get("/cabin-type-default-amenities", {
    params: { cabin_type_id: cabinType },
  });
}

export function getAmenitiesForTicket(apiClient, ticketID) {
  return apiClient.get("/ticket-amenities", {
    params: { ticket_id: ticketID },
  });
}

export function sendAmenitiesForTicket(
  apiClient,
  body = { ticket_id: 0, amenity_ids: [0] },
) {
  return apiClient.post("/ticket-amenities/edit", body);
}

export function getAmenitiesReport(
  apiClient,
  params = { flight_number: 0, from_date: "", to_date: "" },
) {
  return apiClient.get("/amenities/count", { params: params });
}
