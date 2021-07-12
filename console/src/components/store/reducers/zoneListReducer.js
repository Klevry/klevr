const GET_ZONE_LIST = 'GET_ZONE_LIST';

const initialState = [];

const zoneListReducer = (state = initialState, action) => {
  switch (action.type) {
    case GET_ZONE_LIST:
      return action.payload;
    default:
      return state;
  }
};
export default zoneListReducer;
