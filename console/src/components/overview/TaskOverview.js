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
  TableRow,
  TableSortLabel
} from '@material-ui/core';
import ArrowRightIcon from '@material-ui/icons/ArrowRight';
import { useDispatch, useSelector } from 'react-redux';
import { getTaskList } from '../store/actions/klevrActions';

import Refresh from '../common/Refresh';

const TaskList = ({ sortedTaskList }) => {
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
      {sortedTaskList.slice(0, 5).map((item) => (
        <TableRow hover key={item.id}>
          <TableCell>{`${item.id}`}</TableCell>
          <TableCell>{`${item.name}`}</TableCell>
          <TableCell>{`${item.status}`}</TableCell>
          <TableCell>{`${item.taskType}`}</TableCell>
          <TableCell>{`${item.createdAt}`}</TableCell>
        </TableRow>
      ))}
    </TableBody>
  );
};

const TaskOverview = (props) => {
  const taskList = useSelector((store) => store.taskListReducer);
  const [orderDirection, setOrderDirection] = useState('asc');
  const [valueToOrderBy, setValueToOrderBy] = useState('');

  const handleRequestSort = (e, property) => {
    const isAscending = valueToOrderBy === property && orderDirection === 'asc';
    setValueToOrderBy(property);
    setOrderDirection(isAscending ? 'desc' : 'asc');
  };

  const createSortHandler = (property) => (e) => {
    handleRequestSort(e, property);
  };

  function descendingComparator(a, b, orderBy) {
    if (b[orderBy] < a[orderBy]) {
      return -1;
    }
    if (b[orderBy] > a[orderBy]) {
      return 1;
    }
    return 0;
  }

  function getComparator(order, orderBy) {
    return order === 'desc'
      ? (a, b) => descendingComparator(a, b, orderBy)
      : (a, b) => -descendingComparator(a, b, orderBy);
  }

  function stableSort(array, comparator) {
    const stabilizedThis = array.map((el, index) => [el, index]);
    stabilizedThis.sort((a, b) => {
      const order = comparator(a[0], b[0]);
      if (order !== 0) return order;
      return a[1] - b[1];
    });
    return stabilizedThis.map((el) => el[0]);
  }

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
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'id'}
                    direction={valueToOrderBy === 'id' ? orderDirection : 'asc'}
                    onClick={createSortHandler('id')}
                  >
                    ID
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'name'}
                    direction={
                      valueToOrderBy === 'name' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('name')}
                  >
                    Name
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'status'}
                    direction={
                      valueToOrderBy === 'status' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('status')}
                  >
                    Status
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'taskType'}
                    direction={
                      valueToOrderBy === 'taskType' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('taskType')}
                  >
                    Task Type
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'createdAt'}
                    direction={
                      valueToOrderBy === 'createdAt' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('createdAt')}
                  >
                    Created At
                  </TableSortLabel>
                </TableCell>
              </TableRow>
            </TableHead>
            <TaskList
              sortedTaskList={
                taskList &&
                stableSort(
                  taskList,
                  getComparator(orderDirection, valueToOrderBy)
                )
              }
            />
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
