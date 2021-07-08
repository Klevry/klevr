import { useState, useEffect } from 'react';
import axios from 'axios';
import { API_SERVER, GROUP_ID } from '../../config';
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

const TaskList = () => {
  const createTime = Date.now();
  console.log(createTime);

  const [data, setData] = useState(null);

  useEffect(() => {
    let completed = false;

    async function get() {
      const result = await axios.get(`${API_SERVER}/inner/groups`, {
        withCredentials: true
      });
      if (!completed) setData(result.data);
    }
    get();
    return () => {
      completed = true;
    };
  }, []);

  if (!data) {
    return null;
  }
  return (
    <TableBody>
      {data.map((item) => (
        <TableRow hover key={item.agentKey}>
          <TableCell>{`${item.Id}`}</TableCell>
          <TableCell>{`${item.GroupName}`}</TableCell>
          <TableCell>{`${item.CreatedAt}`}</TableCell>
          <TableCell>{`${item.Platform}`}</TableCell>
        </TableRow>
      ))}
    </TableBody>
  );
};

const Alltasks = ({ customers, ...rest }) => {
  return (
    <Card>
      <x.div display="flex" alignItems="center">
        <CardHeader title="Zone" />
      </x.div>
      <Divider />
      <PerfectScrollbar>
        <Box sx={{ minWidth: 800 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>GroupName</TableCell>
                <TableCell>Created At</TableCell>
                <TableCell>Platform</TableCell>
              </TableRow>
            </TableHead>
            <TaskList />
          </Table>
        </Box>
      </PerfectScrollbar>
    </Card>
  );
};

export default Alltasks;
