import { Box, Card, CardContent, Grid, Typography } from '@material-ui/core';
import HighlightOffIcon from '@material-ui/icons/HighlightOff';

const DockerCredential = (props) => (
  <Card sx={{ height: '100%' }} {...props}>
    <CardContent>
      <Grid container spacing={3} sx={{ justifyContent: 'space-between' }}>
        <Grid item>
          <Typography color="textSecondary" gutterBottom variant="h6">
            Credential
          </Typography>
          <Typography color="textPrimary" variant="h3">
            Docker
          </Typography>
        </Grid>
        <Grid item>
          <HighlightOffIcon
            onClick={() => {
              console.log('editCredential');
            }}
          />
        </Grid>
      </Grid>
      <Box
        sx={{
          pt: 2,
          display: 'flex',
          alignItems: 'center'
        }}
      />
    </CardContent>
  </Card>
);

export default DockerCredential;
