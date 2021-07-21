import { useState, useEffect } from 'react';
import axios from 'axios';
import { API_SERVER } from '../../config';
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
import { useDispatch, useSelector } from 'react-redux';
import { getTaskList } from '../store/actions/klevrActions';
import { Button, Modal, Alert } from 'antd';
import { ExclamationCircleOutlined } from '@ant-design/icons';
import { x } from '@xstyled/emotion';

const TaskList = () => {
  const dispatch = useDispatch();
  const currentZone = useSelector((store) => store.zoneReducer);
  const taskList = useSelector((store) => store.taskListReducer);
  const { confirm } = Modal;
  const [cancelResult, setCancelResult] = useState('');

  const fetchTaskList = () => {
    let completed = false;

    async function get() {
      const result = await axios.get(
        `${API_SERVER}/inner/tasks?groupID=${currentZone}`
      );
      if (!completed) dispatch(getTaskList(result.data));
    }
    get();
    return () => {
      completed = true;
    };
  };

  useEffect(() => {
    fetchTaskList();
  }, []);

  useEffect(() => {
    fetchTaskList();
  }, [currentZone]);

  if (!taskList) {
    return null;
  }

  function showDeleteConfirm(id, taskName) {
    setCancelResult('');

    confirm({
      title: `Are you sure cancel the ${taskName}(Id:${id}) task?`,
      icon: <ExclamationCircleOutlined />,
      okText: 'Yes',
      okType: 'danger',
      cancelText: 'No',
      onOk() {
        async function cancelTask() {
          const headers = {
            accept: 'application/json',
            'Content-Type': 'application/json'
          };
          const response = await axios.delete(
            `${API_SERVER}/inner/tasks/${id}`,
            { headers }
          );

          if (response.data.canceled) {
            setCancelResult('success');
            fetchTaskList();
          } else {
            setCancelResult('error');
          }
        }
        cancelTask();
      },
      onCancel() {}
    });
  }

  return (
    <>
      <TableBody>
        {taskList.map((item) => (
          <TableRow hover key={item.agentKey}>
            {item.taskType === 'longTerm' && (
              <>
                <TableCell>{`${item.id}`}</TableCell>
                <TableCell>{`${item.name}`}</TableCell>
                <TableCell>{`${item.exeAgentKey}`}</TableCell>
                <TableCell>{`${item.status}`}</TableCell>
                <TableCell>{`${item.taskType}`}</TableCell>
                <TableCell>{`${item.createdAt}`}</TableCell>
                <TableCell>
                  <Button
                    onClick={() => showDeleteConfirm(item.id, item.name)}
                    type="dashed"
                    disabled={
                      item.status === 'scheduled' ||
                      item.status === 'wait-polling'
                        ? false
                        : true
                    }
                  >
                    Cancel
                  </Button>
                </TableCell>
              </>
            )}
          </TableRow>
        ))}
      </TableBody>
      <x.div
        position="fixed"
        bottom="20px"
        right="20px"
        zIndex="9999"
        minWidth="400px"
      >
        {cancelResult === 'error' && (
          <Alert
            message="Error"
            description="Cancel task failed. Please refresh the list and try again."
            type="error"
            showIcon
            closable
            onClose={() => setCancelResult('')}
          />
        )}
        {cancelResult === 'success' && (
          <Alert
            message="Success"
            description="Task canceled successfully."
            type="success"
            showIcon
            closable
            onClose={() => setCancelResult('')}
          />
        )}
      </x.div>
    </>
  );
};

const ProvisioningList = ({ customers, ...rest }) => {
  return (
    <Card>
      <PerfectScrollbar>
        <Box sx={{ minWidth: 800 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>Name</TableCell>
                <TableCell>ExeAgent</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Task Type</TableCell>
                <TableCell>Created At</TableCell>
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

export default ProvisioningList;
