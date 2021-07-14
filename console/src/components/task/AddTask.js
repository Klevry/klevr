import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { x } from '@xstyled/emotion';
import { Modal, Form, Input, Select, Radio, Divider, Button } from 'antd';
import { useDispatch, useSelector } from 'react-redux';
import { API_SERVER } from 'src/config';
import { getAgentList, getTaskList } from '../store/actions/klevrActions';

const { TextArea } = Input;
const { Option } = Select;

const AddTask = () => {
  const dispatch = useDispatch();
  const currentZone = useSelector((store) => store.zoneReducer);
  const agentList = useSelector((store) => store.agentListReducer);

  const [form] = Form.useForm();
  const [visible, setVisible] = useState(false);
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [cronValid, setCronValid] = useState();

  const [taskValues, setTaskValues] = useState({
    exeAgentChangeable: true,
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

    const headers = {
      'Content-Type': 'application/x-www-form-urlencoded'
    };

    const response = await axios.post(`${API_SERVER}/inner/tasks`, taskValues, {
      headers
    });

    if (response.status === 200) {
      async function get() {
        const result = await axios.get(
          `${API_SERVER}/inner/tasks?groupID=${currentZone}`
        );
        dispatch(getTaskList(result.data));
      }
      get();
      setVisible(false);
      setConfirmLoading(false);
    }

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
      setTaskValues({
        ...taskValues,
        agentKey: '',
        exeAgentChangeable: true
      });
      return;
    }

    setTaskValues({
      ...taskValues,
      agentKey: value,
      exeAgentChangeable: false
    });

    console.log(value);
  };

  const onCronValidator = (e) => {
    const cron = require('cron-validator');

    setCronValid(cron.isValidCron(e.target.value, { seconds: true }));

    if (cron.isValidCron(e.target.value, { seconds: true })) {
      setTaskValues({
        ...taskValues,
        [e.target.name]: e.target.value
      });
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
                  label="Fail"
                  validateStatus={cronValid ? 'success' : 'error'}
                  help="Please adjust the crontab format."
                >
                  <Input
                    placeholder="crontab"
                    // onChange={ontaskChange}
                    onChange={onCronValidator}
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
          <x.div border="1px solid #e4e4e4" pt="15px" pb="15px" mb="20px">
            <x.h5
              fontWeight="bold"
              pb="10px"
              pl="15px"
              pr="15px"
              mb="20px"
              borderBottom="1px solid #e4e4e4"
              fontSize="1.1rem"
            >
              Step
            </x.h5>
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
          </x.div>
        </Form>
      </Modal>
    </>
  );
};

export default AddTask;
