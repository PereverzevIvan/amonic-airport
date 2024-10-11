export const validateEmail = (email) => {
  const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

  if (!emailPattern.test(email)) {
    return false;
  }

  return true;
};

export function allFieldsNotEmpty(formData) {
  return Object.values(formData).every((value) => !!value);
}
