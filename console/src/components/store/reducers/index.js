import { combineReducers } from 'redux';

import zoneReducer from './zoneReducer';
import zoneListReducer from './zoneListReducer';
import agentListReducer from './agentListReducer';
import taskListReducer from './taskListReducer';
import zoneNameReducer from './zoneNameReducer';
import credentialReducer from './credentialReducer';
import taskLogReducer from './taskLogReducer';

export default combineReducers({
  zoneReducer,
  zoneListReducer,
  agentListReducer,
  taskListReducer,
  zoneNameReducer,
  credentialReducer,
  taskLogReducer
});
