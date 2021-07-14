import { useState, useEffect } from 'react';
import axios from 'axios';
import { API_SERVER, GROUP_ID } from '../../config';
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
          {item.taskType === 'atOnce' && (
            <>
              <TableCell>{`${item.id}`}</TableCell>
              <TableCell>{`${item.name}`}</TableCell>
              <TableCell>{`${item.createdAt}`}</TableCell>
              <TableCell>{`${item.status}`}</TableCell>
            </>
          )}
        </TableRow>
      ))}
    </TableBody>
  );
};

const OrderList = ({ customers, ...rest }) => {
  return (
    <Card>
      <PerfectScrollbar>
        <Box sx={{ minWidth: 800 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>Name</TableCell>
                <TableCell>Created At</TableCell>
                <TableCell>Status</TableCell>
              </TableRow>
            </TableHead>
            <TaskList />
          </Table>
        </Box>
      </PerfectScrollbar>
    </Card>
  );
};

export default OrderList;
