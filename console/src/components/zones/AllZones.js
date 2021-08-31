import { useState, useEffect } from 'react';
import axios from 'axios';
import { API_SERVER } from '../../config';
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
  TableHead,
  TableSortLabel
} from '@material-ui/core';

import { useDispatch, useSelector } from 'react-redux';
import { getZoneList } from '../store/actions/klevrActions';
import { Button, Modal } from 'antd';
import { ExclamationCircleOutlined } from '@ant-design/icons';
import AddZone from './AddZone';
import Refresh from '../common/Refresh';

const ZoneList = ({ sortedZoneList }) => {
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
      {sortedZoneList.map((item) => (
        <TableRow hover key={item.Id}>
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

const AllZones = ({ customers, ...rest }) => {
  const zoneList = useSelector((store) => store.zoneListReducer);
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
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'Id'}
                    direction={valueToOrderBy === 'Id' ? orderDirection : 'asc'}
                    onClick={createSortHandler('Id')}
                  >
                    ID
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'GroupName'}
                    direction={
                      valueToOrderBy === 'GroupName' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('GroupName')}
                  >
                    GroupName
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'CreatedAt'}
                    direction={
                      valueToOrderBy === 'CreatedAt' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('CreatedAt')}
                  >
                    Created At
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'Platform'}
                    direction={
                      valueToOrderBy === 'Platform' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('Platform')}
                  >
                    Platform
                  </TableSortLabel>
                </TableCell>
                <TableCell>Action</TableCell>
              </TableRow>
            </TableHead>
            <ZoneList
              sortedZoneList={
                zoneList &&
                stableSort(
                  zoneList,
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

export default AllZones;
