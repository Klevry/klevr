import React from 'react';
import { Box, Container } from '@material-ui/core';
import 'antd/dist/antd.css';

import TaskLog from 'src/components/task/TaskLog';

const LogView = () => {
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
          <TaskLog />
        </Container>
      </Box>
    </>
  );
};

export default LogView;
