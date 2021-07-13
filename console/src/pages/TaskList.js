import React, { useState } from 'react';
import { x } from '@xstyled/emotion';
import { Box, Container, Button } from '@material-ui/core';
import AllTasks from 'src/components/task/AllTasks';
import OrderList from 'src/components/task/OrderList';
import SchedulerList from 'src/components/task/SchedulerList';
import { Modal, Form, Input, Select, Radio } from 'antd';
import { Button as AddBtn } from 'antd';
import { Plus } from 'react-feather';

const { Option } = Select;

const content = [
  {
    tab: 'All',
    content: <AllTasks />
  },
  {
    tab: 'Order',
    content: <OrderList />
  },
  {
    tab: 'Scheduler',
    content: <SchedulerList />
  }
];

const useTabs = (initialTabs, allTabs) => {
  const [contentIndex, setContentIndex] = useState(initialTabs);
  return {
    contentItem: allTabs[contentIndex],
    contentChange: setContentIndex
  };
};

const AddTask = () => {
  const [form] = Form.useForm();
  const [visible, setVisible] = useState(false);
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [addStep, setAddStep] = useState(false);

  // const [componentSize, setComponentSize] = useState('default');

  // const onFormLayoutChange = ({ size }) => {
  //   setComponentSize(size);
  // };

  const onReset = () => {
    form.resetFields();
  };

  const showModal = () => {
    setVisible(true);
  };

  const handleOk = async () => {
    setConfirmLoading(true);
    // console.log(`groupname: ${groupname}, platform: ${platform}`);

    // const headers = {
    //   'Content-Type': 'application/x-www-form-urlencoded'
    // };

    // const response = await axios.post(
    //   `${API_SERVER}/inner/groups`,
    //   {
    //     groupName: groupname,
    //     platform: platform
    //   },
    //   { headers },
    //   {
    //     withCredentials: true
    //   }
    // );

    // console.log(response.status === 200);
    // if (response.status === 200) {
    //   async function get() {
    //     const result = await axios.get(`${API_SERVER}/inner/groups`, {
    //       withCredentials: true
    //     });
    //     dispatch(getZoneList(result.data));
    //   }
    //   get();
    //   setVisible(false);
    //   setConfirmLoading(false);
    // }

    //일단 이것만 살리기.. 지울거
    setVisible(false);
    setConfirmLoading(false);
    //

    onReset();
  };

  const handleCancel = () => {
    console.log('cancel');
    onReset();
    setVisible(false);
  };

  const onPlatformChange = (value) => {
    setPlatform(value);
  };

  const handleChange = (e) => {
    setGroupname(e.target.value);
  };

  const handleAddStep = () => {
    setAddStep(!addStep);
    console.log(addStep);
  };

  return (
    <>
      <AddBtn type="primary" onClick={showModal}>
        ADD TASK
      </AddBtn>
      <Modal
        title="Add task"
        centered
        okText="Add"
        visible={visible}
        onOk={handleOk}
        confirmLoading={confirmLoading}
        onCancel={handleCancel}
        width={700}
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
          <Form.Item
            required
            label="Task name"
            name="taskname"
            rules={[
              {
                required: true,
                message: 'Please put taskname'
              }
            ]}
          >
            <Input />
          </Form.Item>
          <Form.Item label="Type" name="type" required>
            <Radio.Group allowClear>
              <Radio.Button value="order">Order</Radio.Button>
              <Radio.Button value="scheduler">Scheduler</Radio.Button>
              <Radio.Button value="provisioning">Provisioning</Radio.Button>
            </Radio.Group>
          </Form.Item>
          <Form.Item
            noStyle
            shouldUpdate={(prevValues, currentValues) =>
              prevValues.type !== currentValues.type
            }
          >
            {({ getFieldValue }) =>
              getFieldValue('type') === 'scheduler' ? (
                <Form.Item
                  name="Iteration period"
                  label="Iteration period"
                  rules={[
                    {
                      required: true
                    }
                  ]}
                >
                  <Input placeholder="crontab" />
                </Form.Item>
              ) : null
            }
          </Form.Item>
          <Form.Item label="Target Agent" required name="targetAgent">
            <Select
              placeholder="Select a taget agent"
              //  onChange={onPlatformChange}
              allowClear
            >
              <Select.Option value="none">None</Select.Option>
              <Select.Option value="agent1">Agent1</Select.Option>
              <Select.Option value="agent2">Agent2</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item label="Step" name="step">
            <AddBtn onClick={handleAddStep}>Add Step +</AddBtn>
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};

const TaskList = () => {
  const { contentItem, contentChange } = useTabs(0, content);
  return (
    <>
      <Box
        sx={{
          backgroundColor: 'background.default',
          minHeight: '100%',
          py: 3
        }}
      >
        <Container maxWidth={false}>
          <x.div
            display="flex"
            justifyContent="space-between"
            alignItems="center"
            mb="20"
          >
            <div>
              {content.map((section, index) => (
                <Button onClick={() => contentChange(index)}>
                  {section.tab}
                </Button>
              ))}
            </div>
            {/* <Button color="primary" variant="contained" disabled>
              Add Task
            </Button> */}
            <AddTask />
          </x.div>
          {contentItem.content}
        </Container>
      </Box>
    </>
  );
};

export default TaskList;
