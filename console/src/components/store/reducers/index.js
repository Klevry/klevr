import { combineReducers } from 'redux';

import zoneReducer from './zoneReducer';
import zoneListReducer from './zoneListReducer';

export default combineReducers({
  zoneReducer,
  zoneListReducer
});
