export function getAirports(apiClient) {
  return apiClient.get("/airports");
}
