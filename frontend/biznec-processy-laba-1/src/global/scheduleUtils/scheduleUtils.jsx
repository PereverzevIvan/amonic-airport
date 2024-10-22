export function getScheduleIDsInFlight(flight) {
  let scheduleIDs = [];

  for (let index = 0; index < flight?.schedules?.length; index++) {
    const element = flight.schedules[index];
    scheduleIDs.push(element.id);
  }
  return scheduleIDs;
}
