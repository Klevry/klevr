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

const Activate = () => {
  const navigate = useNavigate();
  const SignupSchema = Yup.object().shape({
    userId: Yup.string().max(255).required('ID is required'),
    password: Yup.string().max(255).required('Comfirm password is required')
  });

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
              form.append('cpw', touched.password);

              const response = await axios.post(
                `${API_SERVER}/console/changepassword`,
                form,
                { headers }
              );

              if (response.status === 200) {
                navigate('/login', { replace: true });
              }
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
                    Activate
                  </Typography>
                </Box>
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
                  label="Confirm Password"
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
                    Apply
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

export default Activate;
