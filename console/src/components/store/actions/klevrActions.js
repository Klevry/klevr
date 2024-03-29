const FILTER_BY_ZONE = 'FILTER_BY_ZONE';
const GET_ZONE_LIST = 'GET_ZONE_LIST';
const GET_AGENT_LIST = 'GET_AGENT_LIST';
const GET_TASK_LIST = 'GET_TASK_LIST';
const GET_ZONE_NAME = 'GET_ZONE_NAME';
const GET_CREDENTIAL = 'GET_CREDENTIAL';
const GET_TASK_LOG = 'GET_TASK_LOG';
const GET_LOGIN_STATUS = 'GET_LOGIN_STATUS';

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

export const getZoneName = (payload) => ({
  type: GET_ZONE_NAME,
  payload
});

export const getCredential = (payload) => ({
  type: GET_CREDENTIAL,
  payload
});

export const getTasklog = (payload) => ({
  type: GET_TASK_LOG,
  payload
});

export const getLoginStatus = (payload) => ({
  type: GET_LOGIN_STATUS,
  payload
});
