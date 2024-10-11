export function getTimeFromString(string) {
  const dateObj = new Date(string);
  const time = dateObj.toLocaleTimeString("ru-RU", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });

  return time;
}

export function getDateFromString(string) {
  const dateObj = new Date(string);
  const date = dateObj.toLocaleDateString("ru-RU", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });

  return date;
}

export function ConvertMillisecondsToTimeString(timeDiff) {
  const diffInSeconds = Math.floor(timeDiff / 1000); // переводим миллисекунды в секунды
  const hours = Math.floor(diffInSeconds / 3600);
  const minutes = Math.floor((diffInSeconds % 3600) / 60);
  const seconds = diffInSeconds % 60;

  const padZero = (value) => (value < 10 ? `0${value}` : value);
  return `${padZero(hours)}:${padZero(minutes)}:${padZero(seconds)}`;
}

export function getTimeDiferenceFrom2Strings(string1, string2) {
  const dateObj1 = new Date(string1);
  const dateObj2 = new Date(string2);

  const timeDiff = dateObj2 - dateObj1;

  return ConvertMillisecondsToTimeString(timeDiff);
}

export function getYearDiferenceFrom2Strings(string1, string2) {
  const date1 = new Date(string1);
  const date2 = new Date(string2);

  // Получаем годы обеих дат
  const year1 = date1.getFullYear();
  const year2 = date2.getFullYear();

  // Вычисляем разницу в годах
  let ageDifference = year2 - year1;

  // Проверяем, произошел ли день рождения в текущем году
  if (
    date2.getMonth() < date1.getMonth() ||
    (date2.getMonth() === date1.getMonth() && date2.getDate() < date1.getDate())
  ) {
    ageDifference--;
  }

  return ageDifference;
}
