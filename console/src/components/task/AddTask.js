import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { x } from '@xstyled/emotion';
import { Modal, Form, Input, Select, Radio, Divider, Button } from 'antd';
import { useDispatch, useSelector } from 'react-redux';
import { API_SERVER } from 'src/config';
import { getAgentList } from '../store/actions/klevrActions';

const { TextArea } = Input;
const { Option } = Select;

const AddTask = () => {
  const dispatch = useDispatch();
  const currentZone = useSelector((store) => store.zoneReducer);
  const agentList = useSelector((store) => store.agentListReducer);

  const [form] = Form.useForm();
  const [visible, setVisible] = useState(false);
  const [confirmLoading, setConfirmLoading] = useState(false);

  const [taskValues, setTaskValues] = useState({
    exeAgentChangeable: false,
    hasRecover: false,
    name: '',
    taskType: '',
    steps: [],
    totalStepCount: 1,
    zoneId: 0
  });

  const [values, setValues] = useState({
    command: '',
    commandName: '',
    commandType: '',
    isRecover: false,
    seq: 1
  });

  const fetchAgent = () => {
    let completed = false;

    async function get() {
      const result = await axios.get(
        `${API_SERVER}/inner/groups/${currentZone}/agents`
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
  }, [currentZone]);

  useEffect(() => {
    setTaskValues({
      ...taskValues,
      steps: [values]
    });
  }, [values]);

  const onReset = () => {
    form.resetFields();
  };

  const showModal = () => {
    setVisible(true);

    setTaskValues({
      ...taskValues,
      zoneId: currentZone
    });
  };

  const handleOk = async () => {
    setConfirmLoading(true);

    //일단 이것만 살리기.. 지울거
    setVisible(false);
    setConfirmLoading(false);
    //

    onReset();
  };

  const handleCancel = () => {
    onReset();
    setVisible(false);
  };

  //task settings
  const ontaskChange = (e) => {
    console.log(e.target.name);
    console.log(e.target.value);

    setTaskValues({
      ...taskValues,
      [e.target.name]: e.target.value
    });
  };

  const onAgentChange = (value) => {
    if (value === 'none') {
      return;
    }
  };

  //task step setting
  const handleStepChange = (event) => {
    setValues({
      ...values,
      [event.target.name]: event.target.value
    });
  };

  const handleCmdType = (value) => {
    setValues({
      ...values,
      commandType: value
    });
  };

  return (
    <>
      <Button type="primary" onClick={showModal}>
        ADD TASK
      </Button>
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
            name="name"
            rules={[
              {
                required: true,
                message: 'Please put name'
              }
            ]}
          >
            <Input onChange={ontaskChange} name="name" />
          </Form.Item>
          <Form.Item label="Type" name="taskType" required>
            <Radio.Group allowClear onChange={ontaskChange} name="taskType">
              <Radio.Button value="atOnce">Order</Radio.Button>
              <Radio.Button value="iteration">Scheduler</Radio.Button>
              <Radio.Button value="provisioning" disabled>
                Provisioning
              </Radio.Button>
            </Radio.Group>
          </Form.Item>
          <Form.Item
            noStyle
            shouldUpdate={(prevValues, currentValues) =>
              prevValues.taskType !== currentValues.taskType
            }
          >
            {({ getFieldValue }) =>
              getFieldValue('taskType') === 'iteration' ? (
                <Form.Item
                  name="cron"
                  label="Iteration period"
                  rules={[
                    {
                      required: true
                    }
                  ]}
                >
                  <Input
                    placeholder="crontab"
                    onChange={ontaskChange}
                    name="cron"
                  />
                </Form.Item>
              ) : null
            }
          </Form.Item>
          <Form.Item label="Target Agent" required name="agentKey">
            <Select
              placeholder="Select a target agent"
              onChange={onAgentChange}
              allowClear
              name="agentKey"
            >
              <Option value="none">None</Option>
              {agentList &&
                agentList.map((item) => (
                  <Option value={item.agentKey}>{item.agentKey}</Option>
                ))}
            </Select>
          </Form.Item>
          <Form.Item label="Parameter" name="parameter">
            <Input onChange={ontaskChange} name="parameter" />
          </Form.Item>
          <Divider />
          <Form.Item
            required
            label="Command Name"
            name="commandName"
            rules={[
              {
                required: true,
                message: 'Please put step name'
              }
            ]}
          >
            <Input onChange={handleStepChange} name="commandName" />
          </Form.Item>
          <Form.Item label="Command Type" required name="commandType">
            <Select
              placeholder="Select a step type"
              onChange={handleCmdType}
              allowClear
            >
              <Select.Option value="inline">inline</Select.Option>
              <Select.Option value="reserved">reserved</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item
            label="Command"
            required
            name="command"
            onChange={handleStepChange}
          >
            <TextArea rows={4} name="command" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};

export default AddTask;
