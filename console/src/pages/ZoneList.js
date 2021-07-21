import React from 'react';
import { Box, Container, Grid } from '@material-ui/core';
import AllZones from 'src/components/zones/AllZones';
import 'antd/dist/antd.css';
import ApiKey from 'src/components/zones/ApiKey';

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
