const GET_ZONE_NAME = 'GET_ZONE_NAME';

const initialState = null;

const zoneNameReducer = (state = initialState, action) => {
  switch (action.type) {
    case GET_ZONE_NAME:
      return action.payload;
    default:
      return state;
  }
};
export default zoneNameReducer;
