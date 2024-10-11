export function getUsers(apiClient, params = {}) {
  return apiClient.get("/users", { params: params });
}

export function AddNewUser(apiClient, userData = {}) {
  return apiClient.post("/user", userData);
}

export function EditUser(
  apiClient,
  userData = {
    email: "",
    first_name: "",
    last_name: "",
    office_id: 0,
    birthday: "",
    password: "",
  },
) {
  return apiClient.patch("/user", userData);
}

export function updateUserActive(
  apiClient,
  userData = { user_id: 0, is_active: false },
) {
  return apiClient.put("/user/active", userData);
}

export function getUserById(apiClient, id) {
  return apiClient.get(`/user/${id}`);
}
