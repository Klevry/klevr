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
  TableHead,
  TableSortLabel
} from '@material-ui/core';
import { useDispatch, useSelector } from 'react-redux';
import { getTaskList } from '../store/actions/klevrActions';
import { Button, Modal, Alert } from 'antd';
import { ExclamationCircleOutlined } from '@ant-design/icons';
import { x } from '@xstyled/emotion';

const TaskList = ({ sortedTaskList }) => {
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
        {sortedTaskList.map((item) => (
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
  const taskList = useSelector((store) => store.taskListReducer);
  const [orderDirection, setOrderDirection] = useState('asc');
  const [valueToOrderBy, setValueToOrderBy] = useState('');

  const handleRequestSort = (e, property) => {
    const isAscending = valueToOrderBy === property && orderDirection === 'asc';
    setValueToOrderBy(property);
    setOrderDirection(isAscending ? 'desc' : 'asc');
  };

  const createSortHandler = (property) => (e) => {
    handleRequestSort(e, property);
  };

  function descendingComparator(a, b, orderBy) {
    if (b[orderBy] < a[orderBy]) {
      return -1;
    }
    if (b[orderBy] > a[orderBy]) {
      return 1;
    }
    return 0;
  }

  function getComparator(order, orderBy) {
    return order === 'desc'
      ? (a, b) => descendingComparator(a, b, orderBy)
      : (a, b) => -descendingComparator(a, b, orderBy);
  }

  function stableSort(array, comparator) {
    const stabilizedThis = array.map((el, index) => [el, index]);
    stabilizedThis.sort((a, b) => {
      const order = comparator(a[0], b[0]);
      if (order !== 0) return order;
      return a[1] - b[1];
    });
    return stabilizedThis.map((el) => el[0]);
  }

  return (
    <Card>
      <PerfectScrollbar>
        <Box sx={{ minWidth: 800 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'id'}
                    direction={valueToOrderBy === 'id' ? orderDirection : 'asc'}
                    onClick={createSortHandler('id')}
                  >
                    ID
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'name'}
                    direction={
                      valueToOrderBy === 'name' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('name')}
                  >
                    Name
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'exe'}
                    direction={
                      valueToOrderBy === 'exe' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('exe')}
                  >
                    ExeAgent
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'status'}
                    direction={
                      valueToOrderBy === 'status' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('status')}
                  >
                    Status
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'taskType'}
                    direction={
                      valueToOrderBy === 'taskType' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('taskType')}
                  >
                    Task Type
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'createdAt'}
                    direction={
                      valueToOrderBy === 'createdAt' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('createdAt')}
                  >
                    Created At
                  </TableSortLabel>
                </TableCell>
                <TableCell>Action</TableCell>
              </TableRow>
            </TableHead>
            <TaskList
              sortedTaskList={
                taskList &&
                stableSort(
                  taskList,
                  getComparator(orderDirection, valueToOrderBy)
                )
              }
            />
          </Table>
        </Box>
      </PerfectScrollbar>
    </Card>
  );
};

export default ProvisioningList;
