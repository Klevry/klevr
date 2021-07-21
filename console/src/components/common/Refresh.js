import axios from 'axios';
import { Button as RefreshBtn } from 'antd';
import { RefreshCcw as RfreshIcon } from 'react-feather';
import { useDispatch, useSelector } from 'react-redux';
import { API_SERVER } from 'src/config';
import {
  getAgentList,
  getCredential,
  getTaskList,
  getTasklog,
  getZoneList
} from '../store/actions/klevrActions';

const Refresh = ({ from }) => {
  const dispatch = useDispatch();
  const currentZone = useSelector((store) => store.zoneReducer);

  const fetchTask = () => {
    let completed = false;

    async function get() {
      const result = await axios.get(
        `${API_SERVER}/inner/tasks?groupID=${currentZone}`
      );
      if (!completed) dispatch(getTaskList(result.data));
    }
    get();
    return () => {
      completed = true;
    };
  };

  const fetchAgent = () => {
    let completed = false;

    async function get() {
      const result = await axios.get(
        `${API_SERVER}/inner/groups/${currentZone}/agents`
      );
      if (!completed) dispatch(getAgentList(result.data));
    }
    get();
    return () => {
      completed = true;
    };
  };

  const fetchZone = () => {
    let completed = false;

    async function get() {
      const result = await axios.get(`${API_SERVER}/inner/groups`);
      if (!completed) dispatch(getZoneList(result.data));
    }
    get();
    return () => {
      completed = true;
    };
  };

  const fetchCredential = () => {
    let completed = false;

    async function get() {
      const result = await axios.get(
        `${API_SERVER}/inner/groups/${currentZone}/credentials`
      );
      if (!completed) dispatch(getCredential(result.data));
    }
    get();
    return () => {
      completed = true;
    };
  };

  const fetchTasklog = () => {
    let completed = false;

    async function get() {
      const result = await axios.get(
        `${API_SERVER}/inner/tasks/${currentZone}/logs`
      );
      if (!completed) dispatch(getTasklog(result.data));
    }
    get();
    return () => {
      completed = true;
    };
  };

  const handleRefresh = () => {
    from === 'task' && fetchTask();
    from === 'agent' && fetchAgent();
    from === 'zone' && fetchZone();
    from === 'credential' && fetchCredential();
    from === 'log' && fetchTasklog();
  };

  return (
    <RefreshBtn type="primary" onClick={handleRefresh}>
      <RfreshIcon size="14px" color="white" />
    </RefreshBtn>
  );
};

export default Refresh;
