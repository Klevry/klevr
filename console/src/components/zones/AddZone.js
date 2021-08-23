import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { Modal, Button, Form, Input, Select } from 'antd';
import { x } from '@xstyled/emotion';
import { API_SERVER } from 'src/config';
import {
  filterByZone,
  getZoneList,
  getZoneName
} from '../store/actions/klevrActions';
import { Plus as AddIcon } from 'react-feather';

const { Option } = Select;
const layout = {
  labelCol: {
    span: 6
  },
  wrapperCol: {
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
  const navigate = useNavigate();

  const onReset = () => {
    form.resetFields();
  };

  const showModal = () => {
    setVisible(true);
  };

  const handleOk = async () => {
    if (groupname === '' || platform === '') {
      return;
    }

    setConfirmLoading(true);

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
        selectZone(
          result.data[result.data.length - 1].Id,
          result.data[result.data.length - 1].GroupName
        );
      }
      get();
      setVisible(false);
      setConfirmLoading(false);
    }

    onReset();
  };

  const selectZone = (id, groupName) => {
    dispatch(filterByZone(id));
    dispatch(getZoneName(groupName));
  };

  const handleCancel = () => {
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
        <AddIcon size="14px" />
      </Button>
      <Modal
        title="Add zone"
        centered
        visible={visible}
        confirmLoading={confirmLoading}
        onCancel={handleCancel}
        footer={false}
      >
        <Form {...layout} name="control-ref" form={form} onFinish={handleOk}>
          <Form.Item
            required
            name="groupname"
            label="Name"
            rules={[
              {
                required: true,
                message: 'Please input Name'
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

          <x.div display="flex" justifyContent="flex-end" mt="40px">
            <x.div display="flex" justifyContent="space-between" w="145px">
              <Button onClick={handleCancel}>Cancel</Button>
              <Button type="primary" htmlType="submit">
                Add
              </Button>
            </x.div>
          </x.div>
        </Form>
      </Modal>
    </>
  );
};

export default AddZone;
