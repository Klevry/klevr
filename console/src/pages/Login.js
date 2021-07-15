import { useNavigate } from 'react-router-dom';
import * as Yup from 'yup';
import { Formik } from 'formik';
import {
  Box,
  Button,
  Container,
  TextField,
  Typography
} from '@material-ui/core';
import axios from 'axios';
import { API_SERVER } from 'src/config';
import { useEffect } from 'react';
// import { useHistory } from 'react-router-dom';

const Login = () => {
  // const history = useHistory();
  const navigate = useNavigate();
  const SignupSchema = Yup.object().shape({
    userId: Yup.string().max(255).required('ID is required'),
    password: Yup.string().max(255).required('Password is required')
  });

  useEffect(() => {
    async function check() {
      const result = await axios.get(`${API_SERVER}/console/activated/admin`);

      if (result.data.status === 'initialized') {
        navigate('/activate', { replace: true });
        // history.push('/activate');
      } else if (result.data.status === 'activated') {
        return;
      }
    }
    check();
  }, []);

  return (
    <>
      <Box
        sx={{
          backgroundColor: 'background.default',
          display: 'flex',
          flexDirection: 'column',
          height: '100%',
          justifyContent: 'center'
        }}
      >
        <Container maxWidth="sm">
          <Formik
            initialValues={{
              userId: 'admin',
              password: 'admin'
            }}
            validationSchema={SignupSchema}
            onSubmit={async (touched) => {
              const headers = {
                'Content-Type': 'multipart/form-data'
              };

              let form = new FormData();
              form.append('id', touched.userId);
              form.append('pw', touched.password);

              const response = await axios.post(
                `${API_SERVER}/console/signin`,
                form,
                { headers }
              );

              response.data.token &&
                navigate('/app/overview', { replace: true });
            }}
          >
            {({
              errors,
              handleBlur,
              handleChange,
              handleSubmit,
              isSubmitting,
              touched,
              values
            }) => (
              <form onSubmit={handleSubmit}>
                <Box sx={{ mb: 3 }}>
                  <Typography color="textPrimary" variant="h2">
                    Sign in
                  </Typography>
                </Box>
                <TextField
                  fullWidth
                  label="Klevr Manager URL"
                  margin="normal"
                  onBlur={handleBlur}
                  value={API_SERVER}
                  disabled
                />
                <TextField
                  error={Boolean(touched.userId && errors.userId)}
                  fullWidth
                  helperText={touched.userId && errors.userId}
                  label="ID"
                  margin="normal"
                  name="userId"
                  onBlur={handleBlur}
                  onChange={handleChange}
                  type="userId"
                  value={values.userId}
                  variant="outlined"
                />
                <TextField
                  error={Boolean(touched.password && errors.password)}
                  fullWidth
                  helperText={touched.password && errors.password}
                  label="Password"
                  margin="normal"
                  name="password"
                  onBlur={handleBlur}
                  onChange={handleChange}
                  type="password"
                  value={values.password}
                  variant="outlined"
                />
                <Box sx={{ py: 2 }}>
                  <Button
                    color="primary"
                    disabled={isSubmitting}
                    fullWidth
                    size="large"
                    type="submit"
                    variant="contained"
                  >
                    Sign in
                  </Button>
                </Box>
              </form>
            )}
          </Formik>
        </Container>
      </Box>
    </>
  );
};

export default Login;
