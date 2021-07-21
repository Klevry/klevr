import { useState } from 'react';
import axios from 'axios';
import { API_SERVER } from 'src/config';
import {
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  Divider,
  TextField
} from '@material-ui/core';
import { Alert } from 'antd';
import { x } from '@xstyled/emotion';

const SettingsPassword = (props) => {
  const [pwValid, setPwValid] = useState(false);
  const [updated, setUpdated] = useState(false);
  const [cancelResult, setCancelResult] = useState('');
  const [values, setValues] = useState({
    current: '',
    new: ''
  });

  const handleChange = (event) => {
    setPwValid(false);

    setValues({
      ...values,
      [event.target.name]: event.target.value
    });
  };

  const updateHandler = async () => {
    const headers = {
      'Content-Type': 'multipart/form-data'
    };

    let form = new FormData();
    form.append('id', 'admin');
    form.append('pw', values.current);
    form.append('cpw', values.new);

    try {
      const response = await axios.post(
        `${API_SERVER}/console/changepassword`,
        form,
        { headers }
      );

      if (response.status === 200) {
        setCancelResult('success');
      }
    } catch (err) {
      setPwValid(true);
    }
  };

  return (
    <>
      <form {...props}>
        <Card>
          <CardHeader subheader="Update password" title="Password" />
          <Divider />
          <CardContent>
            <TextField
              disabled
              fullWidth
              label="ID"
              defaultValue="admin"
              variant="outlined"
            />
            <TextField
              fullWidth
              label="Current password"
              margin="normal"
              name="current"
              onChange={handleChange}
              type="password"
              value={values.password}
              variant="outlined"
              error={pwValid}
            />
            <TextField
              fullWidth
              label="New password"
              margin="normal"
              name="new"
              onChange={handleChange}
              type="password"
              value={values.confirm}
              variant="outlined"
            />
          </CardContent>
          <Divider />
          <Box
            sx={{
              display: 'flex',
              justifyContent: 'flex-end',
              p: 2
            }}
          >
            <Button color="primary" variant="contained" onClick={updateHandler}>
              Update
            </Button>
          </Box>
        </Card>
      </form>
      <x.div
        position="fixed"
        bottom="20px"
        right="20px"
        zIndex="9999"
        minWidth="400px"
      >
        {cancelResult === 'success' && (
          <Alert
            message="Success"
            description="Password update successful."
            type="success"
            showIcon
            closable
            onClose={() => setCancelResult('')}
          />
        )}
      </x.div>
    </>
  );
};

export default SettingsPassword;
