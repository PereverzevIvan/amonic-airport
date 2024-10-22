export function getCountries(apiClient) {
  return apiClient.get("/countries");
}
