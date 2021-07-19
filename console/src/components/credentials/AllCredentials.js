import { useState, useEffect } from 'react';
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
import axios from 'axios';
import PerfectScrollbar from 'react-perfect-scrollbar';
import { useDispatch, useSelector } from 'react-redux';
import { Button, Modal } from 'antd';
import { ExclamationCircleOutlined } from '@ant-design/icons';
import { API_SERVER } from 'src/config';
import Refresh from '../common/Refresh';
import AddCredential from './AddCredential';

const CredentialList = () => {
  const [data, setDate] = useState([
    {
      id: 2,
      zoneId: 932,
      key: 'test',
      value: '',
      hash: '21232f297a57a5a743894a0e4a801fc3',
      createdAt: '2021-07-19T03:59:17.000000Z',
      updatedAt: '2021-07-19T03:59:17.000000Z'
    },
    {
      id: 3,
      zoneId: 932,
      key: 'test2',
      value: '',
      hash: '21232f297a57a5a743894a0e4a801fc3',
      createdAt: '2021-07-19T03:59:17.000000Z',
      updatedAt: '2021-07-19T03:59:17.000000Z'
    }
  ]);
  const { confirm } = Modal;

  useEffect(() => {}, []);

  if (!data) {
    return null;
  }

  function showDeleteConfirm(key) {
    confirm({
      title: `Are you sure delete the credential? [key:${key}]`,
      icon: <ExclamationCircleOutlined />,
      okText: 'Yes',
      okType: 'danger',
      cancelText: 'No',
      onOk() {
        console.log('OK');
        console.log(`delete key : ${key}`);

        // async function deleteZone() {
        //   const headers = {
        //     accept: 'application/json',
        //     'Content-Type': 'application/json'
        //   };
        //   const response = await axios.delete(
        //     `${API_SERVER}/inner/groups/${id}`,
        //     { headers }
        //   );
        //   if (response.status === 200) {
        //     const result = await axios.get(`${API_SERVER}/inner/groups`);
        //     dispatch(getZoneList(result.data));
        //   }
        // }
        // deleteZone();
      },
      onCancel() {
        console.log('credential delete cancel');
      }
    });
  }
  return (
    <TableBody>
      {data.map((item) => (
        <TableRow hover key={item.agentKey}>
          <TableCell>{`${item.key}`}</TableCell>
          <TableCell>{`${item.hash}`}</TableCell>
          <TableCell>{`${item.updatedAt}`}</TableCell>
          <TableCell>{`${item.createdAt}`}</TableCell>
          <TableCell>
            <Button onClick={() => showDeleteConfirm(item.key)} type="dashed">
              Delete
            </Button>
          </TableCell>
        </TableRow>
      ))}
    </TableBody>
  );
};

const AllCredentials = ({ customers, ...rest }) => {
  return (
    <Card>
      <x.div
        display="flex"
        alignItems="center"
        justifyContent="space-between"
        paddingRight="10px"
      >
        <CardHeader title="Credential" />
        <x.div w="100px" display="flex" justifyContent="space-between">
          <AddCredential />
          <Refresh from="credential" />
        </x.div>
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
                <TableCell>Action</TableCell>
              </TableRow>
            </TableHead>
            <CredentialList />
          </Table>
        </Box>
      </PerfectScrollbar>
    </Card>
  );
};

export default AllCredentials;
