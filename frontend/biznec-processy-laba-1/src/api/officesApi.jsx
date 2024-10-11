export async function getOffices(apiClient) {
  const response = await apiClient.get("/offices");
  return response;
}
