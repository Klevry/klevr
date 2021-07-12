import { useState, useEffect } from 'react';
import axios from 'axios';
import { x } from '@xstyled/emotion';
import { API_SERVER, GROUP_ID } from '../../config';
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
  TableRow
} from '@material-ui/core';
import { useDispatch, useSelector } from 'react-redux';
import { Modal, Button } from 'antd';
import styled from '@emotion/styled/macro';
import Copy from 'react-copy-to-clipboard';
import { message } from 'antd';
import { CopyOutlined as CopyOutlinedIcon } from '@ant-design/icons';
import { getAgentList } from '../store/actions/klevrActions';

const AgentList = () => {
  const [data, setData] = useState(null);
  const dispatch = useDispatch();
  const currentZone = useSelector((store) => store.zoneReducer);
  const agentList = useSelector((store) => store.agentListReducer);

  const fetchAgent = () => {
    let completed = false;

    async function get() {
      const result = await axios.get(
        `${API_SERVER}/inner/groups/${currentZone}/agents`,
        {
          withCredentials: true
        }
      );
      if (!completed) dispatch(getAgentList(result.data));
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
      {agentList.map((item) => (
        <TableRow hover key={item.agentKey}>
          <TableCell>{item.agentKey}</TableCell>
          <TableCell>{item.ip}</TableCell>
          <TableCell>{item.disk}</TableCell>
          <TableCell>{item.memory}</TableCell>
        </TableRow>
      ))}
    </TableBody>
  );
};

const AddAgent = () => {
  const dispatch = useDispatch();
  const currentZone = useSelector((store) => store.zoneReducer);
  const [visible, setVisible] = useState(false);
  const [confirmLoading, setConfirmLoading] = useState(false);

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

  const showModal = () => {
    setVisible(true);
  };

  const handleOk = () => {
    setConfirmLoading(true);

    console.log('agent install 완료 후 모달 닫기');

    async function get() {
      const result = await axios.get(
        `${API_SERVER}/inner/groups/${currentZone}/agents`,
        {
          withCredentials: true
        }
      );
      dispatch(getAgentList(result.data));

      if (result.status === 200) {
        setVisible(false);
        setConfirmLoading(false);
      }
    }
    get();
  };

  const handleCancel = () => {
    console.log('click cancel');
    setVisible(false);
  };

  const content = `"curl -sL gg.gg/provbee | TAGPROV=0.5 TAGKLEVR=0.2.15-SNAPSHOT K3S_SET=Y K_API_KEY="45908457773441f38f32db09ed733eda" K_PLATFORM="baremetal" K_MANAGER_URL="https://dev.nexclipper.io:8080/klevr" K_ZONE_ID="900" K_CLUSTER_NAME="testasa" bash"`;

  return (
    <>
      <Button type="primary" onClick={showModal}>
        ADD AGENT
      </Button>
      <Modal
        title="Install Agent"
        centered
        okText="Done"
        visible={visible}
        onOk={handleOk}
        confirmLoading={confirmLoading}
        onCancel={handleCancel}
      >
        <Wrapper>
          <Content>{content}</Content>
          <x.div position="absolute" top="5px" right="5px">
            <Copy
              text={content}
              onCopy={() => {
                message.success('Copied');
              }}
            >
              <CopyOutlined />
            </Copy>
          </x.div>
        </Wrapper>
      </Modal>
    </>
  );
};

const AgentOverview = (props) => {
  return (
    <Card {...props}>
      <x.div
        display="flex"
        alignItems="center"
        justifyContent="space-between"
        paddingRight="10px"
      >
        <CardHeader title="Agent" />
        <AddAgent />
      </x.div>
      <Divider />
      <PerfectScrollbar>
        <Box sx={{ minWidth: 800 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Agent ID</TableCell>
                <TableCell>IP</TableCell>
                <TableCell>Disk</TableCell>
                <TableCell>Memory</TableCell>
              </TableRow>
            </TableHead>
            <AgentList />
          </Table>
        </Box>
      </PerfectScrollbar>
    </Card>
  );
};

export default AgentOverview;
