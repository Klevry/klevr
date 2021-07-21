const FILTER_BY_ZONE = 'FILTER_BY_ZONE';

const initialState = 0;

const zoneReducer = (state = initialState, action) => {
  switch (action.type) {
    case FILTER_BY_ZONE:
      return action.payload;
    default:
      return state;
  }
};
export default zoneReducer;
