import { useState, useEffect } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import axios from 'axios';
import PerfectScrollbar from 'react-perfect-scrollbar';
import { x } from '@xstyled/emotion';
import { API_SERVER, GROUP_ID } from '../../config';
import {
  Box,
  Button,
  Card,
  CardHeader,
  Divider,
  Table,
  TableHead,
  TableBody,
  TableCell,
  TableRow
} from '@material-ui/core';
import ArrowRightIcon from '@material-ui/icons/ArrowRight';
import { useDispatch, useSelector } from 'react-redux';
import { getTaskList } from '../store/actions/klevrActions';

import Refresh from '../common/Refresh';

const TaskList = () => {
  const dispatch = useDispatch();
  const currentZone = useSelector((store) => store.zoneReducer);
  const taskList = useSelector((store) => store.taskListReducer);

  useEffect(() => {
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
  }, []);

  useEffect(() => {
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
  }, [currentZone]);

  if (!taskList) {
    return null;
  }
  return (
    <TableBody>
      {taskList.slice(0, 5).map((item) => (
        <TableRow hover key={item.agentKey}>
          <TableCell>{`${item.id}`}</TableCell>
          <TableCell>{`${item.name}`}</TableCell>
          <TableCell>{`${item.createdAt}`}</TableCell>
          <TableCell>{`${item.status}`}</TableCell>
          <TableCell>{`${item.taskType}`}</TableCell>
        </TableRow>
      ))}
    </TableBody>
  );
};

const TaskOverview = (props) => {
  return (
    <Card
      {...props}
      sx={{
        marginBottom: '25px'
      }}
    >
      <x.div
        display="flex"
        alignItems="center"
        justifyContent="space-between"
        paddingRight="10px"
      >
        <CardHeader title="Task" />
        <Refresh from="task" />
      </x.div>
      <Divider />
      <PerfectScrollbar>
        <Box sx={{ minWidth: 800 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>Name</TableCell>
                <TableCell>Created At</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Task Type</TableCell>
              </TableRow>
            </TableHead>
            <TaskList />
          </Table>
        </Box>
      </PerfectScrollbar>
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'flex-end',
          p: 2
        }}
      >
        <RouterLink to="/app/tasks">
          <Button
            color="primary"
            endIcon={<ArrowRightIcon />}
            size="small"
            variant="text"
          >
            View all
          </Button>
        </RouterLink>
      </Box>
    </Card>
  );
};

export default TaskOverview;
