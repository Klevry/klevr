import React from 'react';
import { Box, Container } from '@material-ui/core';
import AllCredentials from 'src/components/credentials/AllCredentials';
import 'antd/dist/antd.css';

const CredentialView = () => {
  return (
    <>
      <Box
        sx={{
          backgroundColor: 'background.default',
          minHeight: '100%',
          py: 3
        }}
      >
        <Container maxWidth={false}>
          <AllCredentials />
        </Container>
      </Box>
    </>
  );
};

export default CredentialView;
