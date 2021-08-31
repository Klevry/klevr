import { useState, useEffect } from 'react';
import axios from 'axios';
import { API_SERVER } from '../../config';
import PerfectScrollbar from 'react-perfect-scrollbar';
import { x } from '@xstyled/emotion';
import {
  Box,
  Card,
  CardHeader,
  Divider,
  Table,
  TableBody,
  TableCell,
  TableRow,
  TableHead,
  TableSortLabel
} from '@material-ui/core';
import Refresh from '../common/Refresh';
import { useDispatch, useSelector } from 'react-redux';
import { getTasklog } from '../store/actions/klevrActions';

const TaskLogList = ({ sortedList }) => {
  const dispatch = useDispatch();
  const taskLog = useSelector((store) => store.taskLogReducer);
  const currentZone = useSelector((store) => store.zoneReducer);

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

  useEffect(() => {
    fetchTasklog();
  }, []);

  useEffect(() => {
    fetchTasklog();
  }, [currentZone]);

  if (!taskLog) {
    return null;
  }

  return (
    <TableBody>
      {sortedList
        .filter((log) => log.log !== '')
        .map((item) => (
          <TableRow hover key={item.id}>
            <TableCell>{`${item.id}`}</TableCell>
            <TableCell>{`${item.name}`}</TableCell>
            <TableCell>{`${item.log}`}</TableCell>
            <TableCell>{`${item.updatedAt}`}</TableCell>
          </TableRow>
        ))}
    </TableBody>
  );
};

const TaskLog = ({ customers, ...rest }) => {
  const taskLog = useSelector((store) => store.taskLogReducer);
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
    <Card>
      <x.div
        display="flex"
        alignItems="center"
        justifyContent="space-between"
        paddingRight="10px"
      >
        <CardHeader title="Logs" />
        <Refresh from="log" />
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
                <TableCell>Log</TableCell>
                <TableCell style={{ width: 260 }}>
                  <TableSortLabel
                    active={valueToOrderBy === 'updatedAt'}
                    direction={
                      valueToOrderBy === 'updatedAt' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('updatedAt')}
                  >
                    Updated At
                  </TableSortLabel>
                </TableCell>
              </TableRow>
            </TableHead>
            <TaskLogList
              sortedList={
                taskLog &&
                stableSort(
                  taskLog,
                  getComparator(orderDirection, valueToOrderBy)
                )
              }
            />
          </Table>
        </Box>
      </PerfectScrollbar>
    </Card>
  );
};

export default TaskLog;
