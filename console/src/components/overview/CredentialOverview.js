import { useEffect } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import axios from 'axios';
import PerfectScrollbar from 'react-perfect-scrollbar';
import { x } from '@xstyled/emotion';
import { API_SERVER } from '../../config';
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

import Refresh from '../common/Refresh';
import { getCredential } from '../store/actions/klevrActions';

const TaskList = () => {
  const dispatch = useDispatch();
  const currentZone = useSelector((store) => store.zoneReducer);
  const credentialList = useSelector((store) => store.credentialReducer);

  const fetchCredential = () => {
    let completed = false;

    async function get() {
      const result = await axios.get(
        `${API_SERVER}/inner/groups/${currentZone}/credentials`
      );
      if (!completed) dispatch(getCredential(result.data));
    }
    get();
    return () => {
      completed = true;
    };
  };

  useEffect(() => {
    fetchCredential();
  }, []);

  useEffect(() => {
    fetchCredential();
  }, [currentZone]);

  if (!credentialList) {
    return null;
  }
  return (
    <TableBody>
      {credentialList.slice(0, 5).map((item) => (
        <TableRow hover key={item.agentKey}>
          <TableCell>{`${item.key}`}</TableCell>
          <TableCell>{`${item.hash}`}</TableCell>
          <TableCell>{`${item.updatedAt}`}</TableCell>
          <TableCell>{`${item.createdAt}`}</TableCell>
        </TableRow>
      ))}
    </TableBody>
  );
};

const CredentialOverview = (props) => {
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
        <CardHeader title="Credential" />
        <Refresh from="credential" />
      </x.div>
      <Divider />
      <PerfectScrollbar>
        <Box sx={{ minWidth: 800 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Key</TableCell>
                <TableCell>Hash</TableCell>
                <TableCell>Updated At</TableCell>
                <TableCell>Created At</TableCell>
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
        <RouterLink to="/app/credentials">
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

export default CredentialOverview;
