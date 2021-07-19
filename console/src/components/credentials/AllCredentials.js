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
import { getCredential } from '../store/actions/klevrActions';

const CredentialList = () => {
  const dispatch = useDispatch();
  const currentZone = useSelector((store) => store.zoneReducer);
  const credentialList = useSelector((store) => store.credentialReducer);

  const { confirm } = Modal;

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

  function showDeleteConfirm(key, id) {
    confirm({
      title: `Are you sure delete the credential? [key:${key}]`,
      icon: <ExclamationCircleOutlined />,
      okText: 'Yes',
      okType: 'danger',
      cancelText: 'No',
      onOk() {
        async function deleteKey() {
          const headers = {
            accept: 'application/json',
            'Content-Type': 'application/json'
          };
          const response = await axios.delete(
            `${API_SERVER}/inner/credentials/${id}`,
            { headers }
          );
          if (response.status === 200) {
            const result = await axios.get(
              `${API_SERVER}/inner/groups/${currentZone}/credentials`
            );
            dispatch(getCredential(result.data));
          }
        }
        deleteKey();
      },
      onCancel() {}
    });
  }

  return (
    <TableBody>
      {credentialList.map((item) => (
        <TableRow hover key={item.agentKey}>
          <TableCell>{`${item.key}`}</TableCell>
          <TableCell>{`${item.hash}`}</TableCell>
          <TableCell>{`${item.updatedAt}`}</TableCell>
          <TableCell>{`${item.createdAt}`}</TableCell>
          <TableCell>
            <Button onClick={() => showDeleteConfirm(item.key)} type="dashed">
              Update
            </Button>
            <Button
              onClick={() => showDeleteConfirm(item.key, item.id)}
              type="dashed"
            >
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
