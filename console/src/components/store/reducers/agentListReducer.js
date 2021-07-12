const GET_AGENT_LIST = 'GET_AGENT_LIST';

const initialState = null;

const agentListReducer = (state = initialState, action) => {
  switch (action.type) {
    case GET_AGENT_LIST:
      return action.payload;
    default:
      return state;
  }
};
export default agentListReducer;
