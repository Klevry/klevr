import React from 'react';
import { Box, Container } from '@material-ui/core';
import 'antd/dist/antd.css';

import ApiKey from 'src/components/zones/ApiKey';
import AllZones from 'src/components/zones/AllZones';

const Dashboard = () => {
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
          <ApiKey />
          <AllZones />
        </Container>
      </Box>
    </>
  );
};

export default Dashboard;
