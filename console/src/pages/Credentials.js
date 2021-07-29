import React from 'react';
import { Box, Container } from '@material-ui/core';
import 'antd/dist/antd.css';

import AllCredentials from 'src/components/credentials/AllCredentials';

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
