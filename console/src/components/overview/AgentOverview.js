import { useState, useEffect } from 'react';
import axios from 'axios';
import { x } from '@xstyled/emotion';
import { API_SERVER } from '../../config';
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
  TableRow,
  TableSortLabel
} from '@material-ui/core';
import { useDispatch, useSelector } from 'react-redux';
import { Modal, Button, Form, Input, Select, Tag, Alert } from 'antd';
import styled from '@emotion/styled/macro';
import Copy from 'react-copy-to-clipboard';
import { CopyOutlined as CopyOutlinedIcon } from '@ant-design/icons';
import { getAgentList } from '../store/actions/klevrActions';
import { Link as RouterLink } from 'react-router-dom';
import Refresh from '../common/Refresh';
import { CatchingPokemonSharp } from '@material-ui/icons';

const { Option } = Select;

const AgentList = ({ sortedAgentList }) => {
  const dispatch = useDispatch();
  const currentZone = useSelector((store) => store.zoneReducer);
  const agentList = useSelector((store) => store.agentListReducer);

  const fetchAgent = () => {
    let completed = false;

    async function get() {
      const agentResult = await axios.get(
        `${API_SERVER}/inner/groups/${currentZone}/agents`
      );
      const primaryResult = await axios.get(
        `${API_SERVER}/inner/groups/${currentZone}/primary`
      );
      if (!completed) {
        const roleAddedData = agentResult.data.filter((el) => {
          if (el.agentKey === primaryResult.data.agentKey) {
            return (el.role = 'Primary');
          } else {
            return (el.role = 'Secondary');
          }
        });

        dispatch(getAgentList(roleAddedData));
      }
    }
    get();
    return () => {
      completed = true;
    };
  };

  useEffect(() => {
    fetchAgent();
  }, []);

  useEffect(() => {
    fetchAgent();
  }, [currentZone]);

  if (!agentList) {
    return null;
  }

  return (
    <TableBody>
      {sortedAgentList.map((item) => (
        <TableRow hover key={item.agentKey}>
          <TableCell>{item.agentKey}</TableCell>
          <TableCell>{item.ip}</TableCell>
          <TableCell>{item.port}</TableCell>
          <TableCell>{item.core}</TableCell>
          <TableCell>{item.disk}</TableCell>
          <TableCell>{item.memory}</TableCell>
          <TableCell>{item.isActive ? 'Active' : 'Inactive'}</TableCell>
          <TableCell>
            {item.role === 'Primary' ? (
              <Tag color="blue">Primary</Tag>
            ) : (
              <Tag>Secondary</Tag>
            )}
          </TableCell>
        </TableRow>
      ))}
    </TableBody>
  );
};

const AddAgent = () => {
  const currentZone = useSelector((store) => store.zoneReducer);
  const [visible, setVisible] = useState(false);
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [form] = Form.useForm();
  const [agentValues, setAgentValues] = useState({
    apiKey: '',
    managerUrl: API_SERVER,
    platform: '',
    zoneId: ''
  });

  const [key, setKey] = useState(undefined);
  const [platform, setPlatform] = useState(undefined);
  const [agentScript, setAgentScript] = useState(undefined);

  const [copyResult, setCopyResult] = useState('');

  const Wrapper = styled.div`
    width: 100%;
    padding: 15px 20px;
    border-radius: 3px;
    background-color: #e3e4e4;
    color: #353535;
    position: relative;
    border: 1px solid #cccccc;
  `;

  const Content = styled.pre`
    line-height: 1.2em;
    white-space: break-spaces;
  `;

  const CopyOutlined = styled(CopyOutlinedIcon)`
    cursor: pointer;
    svg {
      font-size: 1.2em;
      color: #1890ff;
    }
  `;

  const fetchKey = () => {
    let completed = false;
    async function get() {
      try {
        const result = await axios.get(
          `${API_SERVER}/inner/groups/${currentZone}/apikey`
        );
        if (!completed) setKey(result.data);
        setAgentValues({
          ...agentValues,
          apiKey: result.data,
          zoneId: currentZone
        });
      } catch (err) {
        setKey(undefined);
      }
    }
    get();
    return () => {
      completed = true;
    };
  };

  useEffect(() => {
    fetchKey();
  }, []);

  useEffect(() => {
    fetchKey();
  }, [currentZone]);

  const onReset = () => {
    form.resetFields();
  };

  const showModal = () => {
    setVisible(true);
  };

  const handleCancel = () => {
    setPlatform(undefined);
    setAgentScript(undefined);
    setAgentValues({
      ...agentValues,
      platform: undefined
    });
    onReset();
    setVisible(false);
  };

  const onPlatformChange = (value) => {
    setPlatform(value);

    setAgentValues({
      ...agentValues,
      platform: value
    });
  };

  const makeScript = async () => {
    if (platform === undefined) {
      return;
    }

    const headers = {
      'Content-Type': 'application/x-www-form-urlencoded'
    };

    const response = await axios.post(
      `${API_SERVER}/install/agents/bootstrap`,
      agentValues,
      { headers }
    );

    setAgentScript(response.data);
  };

  const handleMsg = () => {
    setCopyResult('success');

    setTimeout(function () {
      setCopyResult('');
    }, 1500);
  };

  return (
    <>
      <Button size="small" onClick={showModal}>
        +
      </Button>
      <Modal
        title="Install Agent"
        centered
        visible={visible}
        confirmLoading={confirmLoading}
        onCancel={handleCancel}
        footer={null}
      >
        <Form
          form={form}
          name="control-ref"
          labelCol={{
            span: 5
          }}
          wrapperCol={{
            span: 17
          }}
        >
          <Form.Item required label="API Key" name="apikey" disabled>
            {key && <Input name="apikey" placeholder={key} disabled />}
            {!key && (
              <RouterLink to="/app/zones">
                <Button type="link">Register API Key</Button>
              </RouterLink>
            )}
          </Form.Item>
          <Form.Item
            name="platform"
            label="Platform"
            rules={[
              {
                required: true,
                message: 'Please select a Platform'
              }
            ]}
          >
            <Select
              placeholder="Select a platform"
              onChange={onPlatformChange}
              allowClear
            >
              <Option value="linux">linux</Option>
              <Option value="baremetal">baremetal</Option>
              <Option value="kubernetes">kubernetes</Option>
            </Select>
          </Form.Item>
          <Form.Item required label="Manager" name="manager">
            <Input name="manager" disabled placeholder={API_SERVER} />
          </Form.Item>
          <Form.Item required label="Zone Id" name="zondId">
            <Input name="zoneId" disabled placeholder={currentZone} />
          </Form.Item>
          <x.div textAlign="center" pt="10px" pb="40px">
            <Button onClick={makeScript} htmlType="submit">
              Create agent setup script
            </Button>
          </x.div>
        </Form>
        {agentScript && (
          <>
            <Wrapper>
              <Content>{agentScript}</Content>
              <x.div position="absolute" top="5px" right="5px">
                <Copy
                  text={agentScript}
                  onCopy={() => {
                    handleMsg();
                  }}
                >
                  <CopyOutlined />
                </Copy>
              </x.div>
              <x.div position="absolute" top="10px" right="185px">
                {copyResult === 'success' && (
                  <Alert message="Copied" type="success" showIcon />
                )}
              </x.div>
            </Wrapper>
          </>
        )}
      </Modal>
    </>
  );
};

const AgentOverview = (props) => {
  const agentList = useSelector((store) => store.agentListReducer);
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
    <Card
      {...props}
      sx={{
        marginBottom: '25px'
      }}
    >
      <x.div
        display="flex"
        alignItems="center"
        paddingRight="10px"
        justifyContent="space-between"
      >
        <x.div display="flex" alignItems="center">
          <CardHeader title="Agent" />
          <AddAgent />
        </x.div>
        <Refresh from="agent" />
      </x.div>
      <Divider />
      <PerfectScrollbar>
        <Box sx={{ minWidth: 800 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'agentKey'}
                    direction={
                      valueToOrderBy === 'agentKey' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('agentKey')}
                  >
                    Agent ID
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'ip'}
                    direction={valueToOrderBy === 'ip' ? orderDirection : 'asc'}
                    onClick={createSortHandler('ip')}
                  >
                    IP
                  </TableSortLabel>
                </TableCell>
                <TableCell>Port</TableCell>
                <TableCell>CPU</TableCell>
                <TableCell>Disk</TableCell>
                <TableCell>Memory</TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'isActive'}
                    direction={
                      valueToOrderBy === 'isActive' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('isActive')}
                  >
                    Status
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={valueToOrderBy === 'role'}
                    direction={
                      valueToOrderBy === 'role' ? orderDirection : 'asc'
                    }
                    onClick={createSortHandler('role')}
                  >
                    Role
                  </TableSortLabel>
                </TableCell>
              </TableRow>
            </TableHead>
            <AgentList
              sortedAgentList={
                agentList &&
                stableSort(
                  agentList,
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

export default AgentOverview;
