import React, { useState } from 'react';
import axios from 'axios';
import { Modal, Button, Form, Input, Divider } from 'antd';
import { API_SERVER } from 'src/config';
import { useDispatch, useSelector } from 'react-redux';
import 'antd/dist/antd.css';
import { x } from '@xstyled/emotion';
import { getCredential } from '../store/actions/klevrActions';
import { useEffect } from 'react';

const layout = {
  labelCol: {
    span: 6
  },
  wrapperCol: {
    span: 16
  }
};

const UpdateCredential = ({ CdKey }) => {
  const dispatch = useDispatch();
  const currentZone = useSelector((store) => store.zoneReducer);
  const [form] = Form.useForm();
  const [visible, setVisible] = useState(false);
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [keyValue, setKeyValue] = useState({
    key: '',
    value: '',
    zoneId: ''
  });

  useEffect(() => {
    setKeyValue({
      ...keyValue,
      zoneId: currentZone
    });
  }, []);

  useEffect(() => {
    setKeyValue({
      ...keyValue,
      zoneId: currentZone
    });
  }, [currentZone]);

  const onReset = () => {
    form.resetFields();
  };

  const showModal = () => {
    setVisible(true);
  };

  const handleOk = async () => {
    if (keyValue.value === '') {
      return;
    }

    setConfirmLoading(true);

    const headers = {
      'Content-Type': 'application/x-www-form-urlencoded'
    };

    const response = await axios.put(
      `${API_SERVER}/inner/credentials`,
      keyValue,
      {
        headers
      }
    );

    if (response.status === 200) {
      async function get() {
        const result = await axios.get(
          `${API_SERVER}/inner/groups/${currentZone}/credentials`
        );
        dispatch(getCredential(result.data));
      }
      get();
      setVisible(false);
      setConfirmLoading(false);
    }

    onReset();
  };

  const handleCancel = () => {
    setKeyValue({
      ...keyValue,
      value: ''
    });
    onReset();
    setVisible(false);
  };

  const handleChange = (e) => {
    setKeyValue(e.target.value);
    setKeyValue({
      ...keyValue,
      key: CdKey,
      [e.target.name]: e.target.value
    });
  };

  return (
    <>
      <Button type="dashed" onClick={showModal}>
        Update
      </Button>
      <Modal
        title="Update Credential"
        centered
        visible={visible}
        confirmLoading={confirmLoading}
        onCancel={handleCancel}
        footer={false}
      >
        <Form {...layout} name="control-ref" form={form} onFinish={handleOk}>
          <Form.Item name="key" label="Key">
            <Input name="key" placeholder={CdKey} disabled />
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
            <x.div display="flex" justifyContent="space-between" w="165px">
              <Button onClick={handleCancel}>Cancel</Button>
              <Button type="primary" htmlType="submit">
                Update
              </Button>
            </x.div>
          </x.div>
        </Form>
      </Modal>
    </>
  );
};

export default UpdateCredential;
