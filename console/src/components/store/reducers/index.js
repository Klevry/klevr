import { combineReducers } from 'redux';

import zoneReducer from './zoneReducer';
import zoneListReducer from './zoneListReducer';
import agentListReducer from './agentListReducer';
import taskListReducer from './taskListReducer';
import zoneNameReducer from './zoneNameReducer';

export default combineReducers({
  zoneReducer,
  zoneListReducer,
  agentListReducer,
  taskListReducer,
  zoneNameReducer
});
