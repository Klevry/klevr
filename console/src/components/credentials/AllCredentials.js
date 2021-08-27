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
  TableHead,
  TableSortLabel
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
import UpdateCredential from './UpdateCredential';

const CredentialList = ({ sortedList }) => {
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
      {sortedList.map((item) => (
        <TableRow hover key={item.agentKey}>
          <TableCell>{`${item.key}`}</TableCell>
          <TableCell>{`${item.hash}`}</TableCell>
          <TableCell>{`${item.updatedAt}`}</TableCell>
          <TableCell>{`${item.createdAt}`}</TableCell>
          <TableCell>
            <x.div w="160px" display="flex" justifyContent="space-between">
              <UpdateCredential CdKey={item.key} />
              <Button
                onClick={() => showDeleteConfirm(item.key, item.id)}
                type="dashed"
              >
                Delete
              </Button>
            </x.div>
          </TableCell>
        </TableRow>
      ))}
    </TableBody>
  );
};

const AllCredentials = ({ customers, ...rest }) => {
  const credentialList = useSelector((store) => store.credentialReducer);
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
    if (orderBy === 'update' || orderBy === 'create') {
      return 'time';
    }

    return order === 'desc'
      ? (a, b) => descendingComparator(a, b, orderBy)
      : (a, b) => -descendingComparator(a, b, orderBy);
  }

  function stableSort(array, comparator) {
    if (comparator === 'time') {
      switch (orderDirection) {
        case 'asc':
          return array;
          break;

        case 'desc':
          return [...array].reverse();
          break;
        default:
      }
    }

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
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'key'}
                    direction={
                      valueToOrderBy === 'key' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('key')}
                  >
                    Key
                  </TableSortLabel>
                </TableCell>
                <TableCell>Hash</TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'update'}
                    direction={
                      valueToOrderBy === 'update' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('update')}
                  >
                    Updated At
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'create'}
                    direction={
                      valueToOrderBy === 'create' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('create')}
                  >
                    Created At
                  </TableSortLabel>
                </TableCell>
                <TableCell>Action</TableCell>
              </TableRow>
            </TableHead>
            <CredentialList
              sortedList={
                credentialList &&
                stableSort(
                  credentialList,
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

export default AllCredentials;
