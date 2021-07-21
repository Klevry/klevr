const GET_LOGIN_STATUS = 'GET_LOGIN_STATUS';

const initialState = false;

const loginReducer = (state = initialState, action) => {
  switch (action.type) {
    case GET_LOGIN_STATUS:
      return action.payload;
    default:
      return state;
  }
};
export default loginReducer;
