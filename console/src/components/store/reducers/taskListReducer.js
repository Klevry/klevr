const GET_TASK_LIST = 'GET_TASK_LIST';

const initialState = null;

const taskListReducer = (state = initialState, action) => {
  switch (action.type) {
    case GET_TASK_LIST:
      return action.payload;
    default:
      return state;
  }
};
export default taskListReducer;
