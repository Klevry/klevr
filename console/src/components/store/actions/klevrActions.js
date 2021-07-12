const FILTER_BY_ZONE = 'FILTER_BY_ZONE';
const GET_ZONE_LIST = 'GET_ZONE_LIST';

export const filterByZone = (payload) => ({
  type: FILTER_BY_ZONE,
  payload
});

export const getZoneList = (payload) => ({
  type: GET_ZONE_LIST,
  payload
});
