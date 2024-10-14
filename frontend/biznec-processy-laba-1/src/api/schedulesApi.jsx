export function getSchedules(apiClient, params = {}) {
  return apiClient.get("/schedules", { params: params });
}

export function updateScheduleByID(apiClient, id, scheduleUpdateParams) {
  return apiClient.patch(`/schedule/${id}`, scheduleUpdateParams);
}

export function changeScheduleConfirmed(apiClient, id, isConfirmed) {
  return apiClient.put(`/schedule/${id}`, { confirmed: isConfirmed });
}
