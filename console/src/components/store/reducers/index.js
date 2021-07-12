import { combineReducers } from 'redux';

import zoneReducer from './zoneReducer';
import zoneListReducer from './zoneListReducer';
import agentListReducer from './agentListReducer';

export default combineReducers({
  zoneReducer,
  zoneListReducer,
  agentListReducer
});
