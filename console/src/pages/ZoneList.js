import {
  Box,
  Container,
  Grid,
  Button,
  Card,
  CardHeader,
  Divider
} from '@material-ui/core';

import { x } from '@xstyled/emotion';

import TaskOverview from 'src/components/overview/TaskOverview';
import AgentOverview from 'src/components/overview/AgentOverview';

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
          <x.div
            display="flex"
            justifyContent="flex-end"
            alignItems="center"
            mb="20"
          >
            <Button color="primary" variant="contained" disabled>
              Add Zone
            </Button>
          </x.div>
          <Grid container spacing={3}>
            <Grid item lg={12} md={12} xl={9} xs={12}>
              <AllZones />
            </Grid>
          </Grid>
        </Container>
      </Box>
    </>
  );
};

export default Dashboard;
