import { Box, Container } from '@material-ui/core';

import TaskOverview from 'src/components/overview/TaskOverview';
import AgentOverview from 'src/components/overview/AgentOverview';
import CredentialOverview from 'src/components/overview/CredentialOverview';

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
          <TaskOverview />
          <CredentialOverview />
          <AgentOverview />
        </Container>
      </Box>
    </>
  );
};

export default Dashboard;
