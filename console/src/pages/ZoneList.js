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
          <Grid container spacing={3}>
            <Grid item lg={12} md={12} xl={9} xs={12}>
              <ApiKey />
              <AllZones />
            </Grid>
          </Grid>
        </Container>
      </Box>
    </>
  );
};

export default Dashboard;
