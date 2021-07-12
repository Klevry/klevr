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

const TaskList = () => {
  const [data, setData] = useState(null);
  const dispatch = useDispatch();
  const zoneList = useSelector((store) => store.zoneListReducer);
  const { confirm } = Modal;

  useEffect(() => {
    let completed = false;

    async function get() {
      const result = await axios.get(`${API_SERVER}/inner/groups`, {
        withCredentials: true
      });
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
        console.log('OK');
        console.log(`delete zone id ${id}`);

        async function deleteZone() {
          // const headers = {
          //   accept: 'application/json',
          //   'Content-Type': 'application/json'
          // };
          // const response = await axios.delete(
          //   `${API_SERVER}/inner/groups/${id}`,
          //   { headers },
          //   {
          //     withCredentials: true
          //   }
          // );
          // console.log(response);
        }
        deleteZone();
      },
      onCancel() {
        console.log('cancel');
      }
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
                <TableCell>Action</TableCell>
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
