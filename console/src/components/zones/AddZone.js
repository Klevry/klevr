import React, { useState } from 'react';
import axios from 'axios';
import { Modal, Button, Form, Input, Select } from 'antd';
import { API_SERVER } from 'src/config';
import { useDispatch } from 'react-redux';
import { getZoneList } from '../store/actions/klevrActions';

const { Option } = Select;
const layout = {
  labelCol: {
    span: 6
  },
  wrapperCol: {
    span: 16
  }
};
const tailLayout = {
  wrapperCol: {
    offset: 8,
    span: 16
  }
};

const AddZone = () => {
  const [form] = Form.useForm();
  const [visible, setVisible] = useState(false);
  const [confirmLoading, setConfirmLoading] = useState(false);

  const [groupname, setGroupname] = useState('');
  const [platform, setPlatform] = useState('');
  const dispatch = useDispatch();

  const onReset = () => {
    form.resetFields();
  };

  const showModal = () => {
    setVisible(true);
  };

  const handleOk = async () => {
    setConfirmLoading(true);
    console.log(`groupname: ${groupname}, platform: ${platform}`);

    const headers = {
      'Content-Type': 'application/x-www-form-urlencoded'
    };

    const response = await axios.post(
      `${API_SERVER}/inner/groups`,
      {
        groupName: groupname,
        platform: platform
      },
      { headers }
    );

    if (response.status === 200) {
      async function get() {
        const result = await axios.get(`${API_SERVER}/inner/groups`);
        dispatch(getZoneList(result.data));
      }
      get();
      setVisible(false);
      setConfirmLoading(false);
    }

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

  return (
    <>
      <Button type="primary" onClick={showModal}>
        ADD ZONE
      </Button>
      <Modal
        title="Add zone"
        centered
        okText="Add"
        visible={visible}
        onOk={handleOk}
        confirmLoading={confirmLoading}
        onCancel={handleCancel}
      >
        <Form {...layout} name="control-ref" form={form}>
          <Form.Item
            required
            name="groupname"
            label="Groupname"
            rules={[
              {
                required: true,
                message: 'Please put Groupname'
              }
            ]}
          >
            <Input onChange={handleChange} />
          </Form.Item>
          <Form.Item
            required
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
        </Form>
      </Modal>
    </>
  );
};

export default AddZone;
