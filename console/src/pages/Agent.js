import React from 'react';
import { Box, Container } from '@material-ui/core';
import 'antd/dist/antd.css';
import AgentOverview from 'src/components/overview/AgentOverview';

const AgentView = () => {
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
          <AgentOverview />
        </Container>
      </Box>
    </>
  );
};

export default AgentView;
