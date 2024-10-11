export function getUserSessions(
  apiClient,
  params = { user_id: 0, only_unmarked_invalid_sessions: false },
) {
  return apiClient.get("/user-sessions", { params: params });
}

export function updateUserSession(apiClient, newUserSession) {
  return apiClient.patch("/user-sessions", newUserSession);
}
