import { useState } from 'react';
import axios from 'axios';
import {
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  Divider,
  TextField
} from '@material-ui/core';
import { API_SERVER } from 'src/config';

const SettingsPassword = (props) => {
  const [values, setValues] = useState({
    current: '',
    new: ''
  });

  const handleChange = (event) => {
    setValues({
      ...values,
      [event.target.name]: event.target.value
    });
  };

  const updateHandler = async () => {
    console.log(values.new);
    console.log(values.confirm);

    const headers = {
      'Content-Type': 'multipart/form-data'
    };

    let form = new FormData();
    form.append('id', 'admin');
    form.append('pw', values.current);
    form.append('cpw', values.new);

    const response = await axios.post(
      `${API_SERVER}/console/changepassword`,
      form,
      { headers },
      {
        withCredentials: true
      }
    );

    console.log(response);
  };

  return (
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
  );
};

export default SettingsPassword;
