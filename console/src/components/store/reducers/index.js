import { combineReducers } from 'redux';

import zoneReducer from './zoneReducer';
import zoneListReducer from './zoneListReducer';
import agentListReducer from './agentListReducer';
import taskListReducer from './taskListReducer';

export default combineReducers({
  zoneReducer,
  zoneListReducer,
  agentListReducer,
  taskListReducer
});
