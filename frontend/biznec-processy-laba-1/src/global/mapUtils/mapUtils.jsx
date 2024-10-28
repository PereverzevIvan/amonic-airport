export function createMapFromList(list) {
  let map = new Map();
  list.forEach((element) => {
    map.set(`${element.id}`, element);
  });

  return map;
}
