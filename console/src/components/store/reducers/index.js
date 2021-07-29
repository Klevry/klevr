import { combineReducers } from 'redux';
import { persistReducer } from 'redux-persist';
import storage from 'redux-persist/lib/storage';
import zoneReducer from './zoneReducer';
import zoneListReducer from './zoneListReducer';
import agentListReducer from './agentListReducer';
import taskListReducer from './taskListReducer';
import zoneNameReducer from './zoneNameReducer';
import credentialReducer from './credentialReducer';
import taskLogReducer from './taskLogReducer';
import loginReducer from './loginReducer';

const persistConfig = {
  key: 'root',
  storage,
  whitelist: ['loginReducer']
};

const rootReducer = combineReducers({
  zoneReducer,
  zoneListReducer,
  agentListReducer,
  taskListReducer,
  zoneNameReducer,
  credentialReducer,
  taskLogReducer,
  loginReducer
});

export default persistReducer(persistConfig, rootReducer);
