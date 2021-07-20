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
  TableHead
} from '@material-ui/core';

import { useSelector } from 'react-redux';

const TaskLogList = () => {
  const [data, setData] = useState(null);
  const zoneList = useSelector((store) => store.zoneListReducer);

  useEffect(() => {}, []);

  if (!data) {
    return null;
  }

  return (
    <TableBody>
      {data
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
        <x.div w="100px" display="flex" justifyContent="space-between">
          {/* <AddZone /> */}
          {/* <Refresh from="zone" /> */}
        </x.div>
      </x.div>
      <Divider />
      <PerfectScrollbar>
        <Box sx={{ minWidth: 800 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>Log</TableCell>
                <TableCell>Updated At</TableCell>
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
