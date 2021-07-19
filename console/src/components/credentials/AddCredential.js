import React, { useState } from 'react';
import axios from 'axios';
import { Modal, Button, Form, Input, Divider } from 'antd';
import { API_SERVER } from 'src/config';
import { useDispatch } from 'react-redux';
import 'antd/dist/antd.css';
import { x } from '@xstyled/emotion';
import { Plus as AddIcon } from 'react-feather';

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

  const [keyValue, setKeyValue] = useState({
    key: '',
    value: ''
  });
  const [platform, setPlatform] = useState('');
  const dispatch = useDispatch();

  const onReset = () => {
    form.resetFields();
  };

  const showModal = () => {
    setVisible(true);
  };

  const handleOk = async () => {
    if (keyValue.key === '' || keyValue.value === '') {
      console.log('ERROR');
      return;
    }

    console.log('SUCCESS');
    setConfirmLoading(true);
    console.log(keyValue);

    // const headers = {
    //   'Content-Type': 'application/x-www-form-urlencoded'
    // };

    // const response = await axios.post(
    //   `${API_SERVER}/inner/groups`,
    //   {
    //     groupName: groupname,
    //     platform: platform
    //   },
    //   { headers }
    // );

    // if (response.status === 200) {
    //   async function get() {
    //     const result = await axios.get(`${API_SERVER}/inner/groups`);
    //     dispatch(getZoneList(result.data));
    //   }
    //   get();
    //   setVisible(false);
    //   setConfirmLoading(false);
    // }

    //지우기
    setVisible(false);
    setConfirmLoading(false);

    onReset();
  };

  const handleCancel = () => {
    console.log('cancel');
    setKeyValue({
      ...keyValue,
      key: '',
      value: ''
    });
    onReset();
    setVisible(false);
  };

  const handleChange = (e) => {
    setKeyValue(e.target.value);
    setKeyValue({
      ...keyValue,
      [e.target.name]: e.target.value
    });
  };

  return (
    <>
      <Button type="primary" onClick={showModal}>
        <AddIcon size="14px" />
      </Button>
      <Modal
        title="Add credential"
        centered
        visible={visible}
        confirmLoading={confirmLoading}
        onCancel={handleCancel}
        footer={false}
      >
        <Form {...layout} name="control-ref" form={form} onFinish={handleOk}>
          <Form.Item
            name="key"
            label="Key"
            rules={[
              {
                required: true,
                message: 'Please input Key'
              }
            ]}
          >
            <Input onChange={handleChange} name="key" />
          </Form.Item>
          <Form.Item
            required
            name="value"
            label="value"
            rules={[
              {
                required: true,
                message: 'Please input Value'
              }
            ]}
          >
            <Input onChange={handleChange} name="value" />
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
