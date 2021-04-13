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

const TaskList = () => {
  const [data, setData] = useState(null);

  useEffect(() => {
    let completed = false;

    async function get() {
      const result = await axios.get(
        `${API_SERVER}/inner/tasks?groupID=${GROUP_ID}`,
        {
          withCredentials: true
        }
      );
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
          <TableCell>{`${item.id}`}</TableCell>
          <TableCell>{`${item.name}`}</TableCell>
          <TableCell>{`${item.createdAt}`}</TableCell>
          <TableCell>{`${item.status}`}</TableCell>
        </TableRow>
      ))}
    </TableBody>
  );
};

const Alltasks = ({ customers, ...rest }) => {
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

export default Alltasks;
