import {
  Box,
  Container,
  Grid,
  Card,
  CardHeader,
  Divider
} from '@material-ui/core';

import { x } from '@xstyled/emotion';

import TaskOverview from 'src/components/overview/TaskOverview';
import AgentOverview from 'src/components/overview/AgentOverview';

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
              <TaskOverview />
            </Grid>
            <Grid item lg={12} md={12} xl={9} xs={12}>
              <AgentOverview />
            </Grid>
          </Grid>
          <Card sx={{ mt: 3 }}>
            <x.div display="flex" alignItems="center">
              <CardHeader title="Credential" />
            </x.div>
            <Divider />
            <Box
              sx={{
                display: 'flex',
                justifyContent: 'flex-end',
                p: 2
              }}
            ></Box>
          </Card>
        </Container>
      </Box>
    </>
  );
};

export default Dashboard;
