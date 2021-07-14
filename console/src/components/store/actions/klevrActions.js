const FILTER_BY_ZONE = 'FILTER_BY_ZONE';
const GET_ZONE_LIST = 'GET_ZONE_LIST';
const GET_AGENT_LIST = 'GET_AGENT_LIST';
const GET_TASK_LIST = 'GET_TASK_LIST';

export const filterByZone = (payload) => ({
  type: FILTER_BY_ZONE,
  payload
});

export const getZoneList = (payload) => ({
  type: GET_ZONE_LIST,
  payload
});

export const getAgentList = (payload) => ({
  type: GET_AGENT_LIST,
  payload
});

export const getTaskList = (payload) => ({
  type: GET_TASK_LIST,
  payload
});
