const GET_CREDENTIAL = 'GET_CREDENTIAL';

const initialState = null;

const credentialReducer = (state = initialState, action) => {
  switch (action.type) {
    case GET_CREDENTIAL:
      return action.payload;
    default:
      return state;
  }
};
export default credentialReducer;
