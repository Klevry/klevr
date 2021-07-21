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

import { useDispatch, useSelector } from 'react-redux';
import { getZoneList } from '../store/actions/klevrActions';
import { Button, Modal } from 'antd';
import { ExclamationCircleOutlined } from '@ant-design/icons';
import AddZone from './AddZone';
import Refresh from '../common/Refresh';

const ZoneList = () => {
  const dispatch = useDispatch();
  const zoneList = useSelector((store) => store.zoneListReducer);
  const currentZone = useSelector((store) => store.zoneReducer);
  const { confirm } = Modal;

  useEffect(() => {
    let completed = false;

    async function get() {
      const result = await axios.get(`${API_SERVER}/inner/groups`);
      if (!completed) dispatch(getZoneList(result.data));
    }
    get();
    return () => {
      completed = true;
    };
  }, []);

  if (!zoneList) {
    return null;
  }

  function showDeleteConfirm(id, groupName) {
    confirm({
      title: `Are you sure delete the ${groupName}(Id:${id}) zone?`,
      icon: <ExclamationCircleOutlined />,
      okText: 'Yes',
      okType: 'danger',
      cancelText: 'No',
      onOk() {
        async function deleteZone() {
          const headers = {
            accept: 'application/json',
            'Content-Type': 'application/json'
          };
          const response = await axios.delete(
            `${API_SERVER}/inner/groups/${id}`,
            { headers }
          );
          if (response.status === 200) {
            const result = await axios.get(`${API_SERVER}/inner/groups`);
            dispatch(getZoneList(result.data));
          }
        }
        deleteZone();
      },
      onCancel() {}
    });
  }
  return (
    <TableBody>
      {zoneList.map((item) => (
        <TableRow hover key={item.agentKey}>
          <TableCell>{`${item.Id}`}</TableCell>
          <TableCell>{`${item.GroupName}`}</TableCell>
          <TableCell>{`${item.CreatedAt}`}</TableCell>
          <TableCell>{`${item.Platform}`}</TableCell>
          <TableCell>
            <Button
              onClick={() => showDeleteConfirm(item.Id, item.GroupName)}
              type="dashed"
              disabled={item.Id === currentZone}
            >
              Delete
            </Button>
          </TableCell>
        </TableRow>
      ))}
    </TableBody>
  );
};

const Alltasks = ({ customers, ...rest }) => {
  return (
    <Card>
      <x.div
        display="flex"
        alignItems="center"
        justifyContent="space-between"
        paddingRight="10px"
      >
        <CardHeader title="Zone" />
        <x.div w="100px" display="flex" justifyContent="space-between">
          <AddZone />
          <Refresh from="zone" />
        </x.div>
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
                <TableCell>Action</TableCell>
              </TableRow>
            </TableHead>
            <ZoneList />
          </Table>
        </Box>
      </PerfectScrollbar>
    </Card>
  );
};

export default Alltasks;
