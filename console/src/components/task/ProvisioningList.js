import { useEffect } from 'react';
import axios from 'axios';
import { API_SERVER } from '../../config';
import PerfectScrollbar from 'react-perfect-scrollbar';
import {
  Box,
  Card,
  Table,
  TableBody,
  TableCell,
  TableRow,
  TableHead
} from '@material-ui/core';
import { useDispatch, useSelector } from 'react-redux';
import { getTaskList } from '../store/actions/klevrActions';

const TaskList = () => {
  const dispatch = useDispatch();
  const currentZone = useSelector((store) => store.zoneReducer);
  const taskList = useSelector((store) => store.taskListReducer);

  const fetchTaskList = () => {
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

  useEffect(() => {
    fetchTaskList();
  }, []);

  useEffect(() => {
    fetchTaskList();
  }, [currentZone]);

  if (!taskList) {
    return null;
  }

  return (
    <TableBody>
      {taskList.map((item) => (
        <TableRow hover key={item.agentKey}>
          {item.taskType === 'longTerm' && (
            <>
              <TableCell>{`${item.id}`}</TableCell>
              <TableCell>{`${item.name}`}</TableCell>
              <TableCell>{`${item.exeAgentKey}`}</TableCell>
              <TableCell>{`${item.status}`}</TableCell>
              <TableCell>{`${item.taskType}`}</TableCell>
              <TableCell>{`${item.createdAt}`}</TableCell>
            </>
          )}
        </TableRow>
      ))}
    </TableBody>
  );
};

const ProvisioningList = ({ customers, ...rest }) => {
  return (
    <Card>
      <PerfectScrollbar>
        <Box sx={{ minWidth: 800 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>Name</TableCell>
                <TableCell>ExeAgent</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Task Type</TableCell>
                <TableCell>Created At</TableCell>
              </TableRow>
            </TableHead>
            <TaskList />
          </Table>
        </Box>
      </PerfectScrollbar>
    </Card>
  );
};

export default ProvisioningList;
