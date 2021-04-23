import { useState, useEffect } from 'react';
import axios from 'axios';
import { x } from '@xstyled/emotion';
import { API_SERVER, GROUP_ID } from '../../config';
import PerfectScrollbar from 'react-perfect-scrollbar';
import {
  Box,
  Card,
  CardHeader,
  Divider,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow
} from '@material-ui/core';

const AgentList = () => {
  const [data, setData] = useState(null);

  useEffect(() => {
    let completed = false;

    async function get() {
      const result = await axios.get(
        `${API_SERVER}/inner/groups/${GROUP_ID}/agents`,
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
          <TableCell>{item.agentKey}</TableCell>
          <TableCell>{item.ip}</TableCell>
          <TableCell>{item.disk}</TableCell>
          <TableCell>{item.memory}</TableCell>
        </TableRow>
      ))}
    </TableBody>
  );
};

const AgentOverview = (props) => {
  return (
    <Card {...props}>
      <x.div display="flex" alignItems="center">
        <CardHeader title="Agent" />
      </x.div>
      <Divider />
      <PerfectScrollbar>
        <Box sx={{ minWidth: 800 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Agent ID</TableCell>
                <TableCell>IP</TableCell>
                <TableCell>Disk</TableCell>
                <TableCell>Memory</TableCell>
              </TableRow>
            </TableHead>
            <AgentList />
          </Table>
        </Box>
      </PerfectScrollbar>
    </Card>
  );
};

export default AgentOverview;
