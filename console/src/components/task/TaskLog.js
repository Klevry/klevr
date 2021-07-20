import { useEffect } from 'react';
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
  TableHead
} from '@material-ui/core';
import Refresh from '../common/Refresh';
import { useDispatch, useSelector } from 'react-redux';
import { getTasklog } from '../store/actions/klevrActions';

const TaskLogList = () => {
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
      {taskLog
        .filter((log) => log.log !== '')
        .map((item) => (
          <TableRow hover key={item.id}>
            <TableCell>{`${item.id}`}</TableCell>
            <TableCell>{`${item.log}`}</TableCell>
            <TableCell>{`${item.updatedAt}`}</TableCell>
          </TableRow>
        ))}
    </TableBody>
  );
};

const TaskLog = ({ customers, ...rest }) => {
  return (
    <Card
      sx={{
        marginTop: '25px'
      }}
    >
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
                <TableCell>ID</TableCell>
                <TableCell>Log</TableCell>
                <TableCell style={{ width: 260 }}>Updated At</TableCell>
              </TableRow>
            </TableHead>
            <TaskLogList />
          </Table>
        </Box>
      </PerfectScrollbar>
    </Card>
  );
};

export default TaskLog;
