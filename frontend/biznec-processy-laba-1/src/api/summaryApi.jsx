export function getAirlineSummary(
  apiClient,
  params = { start_date: "", end_date: "" },
) {
  return apiClient.get("/summary", { params: params });
}
