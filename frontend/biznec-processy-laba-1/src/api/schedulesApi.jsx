export function getSchedules(apiClient, params = {}) {
  return apiClient.get("/schedules", { params: params });
}

export function updateScheduleByID(apiClient, id, scheduleUpdateParams) {
  return apiClient.patch(`/schedule/${id}`, scheduleUpdateParams);
}

export function changeScheduleConfirmed(apiClient, id, isConfirmed) {
  return apiClient.put(`/schedule/${id}`, { confirmed: isConfirmed });
}

export function UploadFileWithSchedules(apiClient, file) {
  // Создаем FormData для отправки файла
  const formData = new FormData();
  formData.append("file", file);

  return apiClient.post("/schedules/upload", formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
}

export function getFlights(apiClient, searchParams) {
  return apiClient.post("/search-flights", searchParams);
}
