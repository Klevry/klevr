const GET_TASK_LOG = 'GET_TASK_LOG';

const initialState = null;

const taskLogReducer = (state = initialState, action) => {
  switch (action.type) {
    case GET_TASK_LOG:
      return action.payload;
    default:
      return state;
  }
};
export default taskLogReducer;
