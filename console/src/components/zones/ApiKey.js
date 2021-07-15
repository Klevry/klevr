import React from 'react';
import { useEffect, useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import axios from 'axios';
import {
  Card,
  CardHeader,
  Divider,
  CardContent,
  TextField,
  Button
} from '@material-ui/core';
import { x } from '@xstyled/emotion';
import { useSelector } from 'react-redux';
import { API_SERVER } from 'src/config';

const useStyles = makeStyles({
  root: {
    minWidth: 275,
    marginBottom: 30
  },
  btn: {
    width: 150,
    height: 56
  }
});

const ApiKey = () => {
  const classes = useStyles();
  const [key, setKey] = useState(undefined);
  const currentZone = useSelector((store) => store.zoneReducer);
  const currentZoneName = useSelector((store) => store.zoneNameReducer);

  const fetchKey = () => {
    let completed = false;
    async function get() {
      try {
        const result = await axios.get(
          `${API_SERVER}/inner/groups/${currentZone}/apikey`
        );
        if (!completed) setKey(result.data);
      } catch (err) {
        setKey(undefined);
      }
    }
    get();
    return () => {
      completed = true;
    };
  };

  useEffect(() => {
    fetchKey();
  }, []);

  useEffect(() => {
    fetchKey();
  }, [currentZone]);

  return (
    <Card className={classes.root} variant="outlined" marginBottom="30px">
      <x.div
        display="flex"
        alignItems="center"
        justifyContent="space-between"
        paddingRight="10px"
      >
        <CardHeader
          title={`Current Zone > ${currentZoneName} (${currentZone})`}
        />
      </x.div>
      <Divider />
      <CardContent>
        <x.div display="flex" alignItems="center" padding="20px">
          <x.h3 w="100px" mr="80px" ml="25px">
            API Key
          </x.h3>
          <TextField
            id="outlined-full-width"
            label="API Key"
            style={{ margin: 8 }}
            placeholder="Please register by pressing the button next to the input."
            //   helperText="Full width!"
            fullWidth
            margin="normal"
            InputLabelProps={{
              shrink: true
            }}
            variant="outlined"
            disabled={key}
            value={key === undefined ? '' : key}
            size="medium"
          />
          <Button
            variant="contained"
            color="primary"
            disableElevation
            className={classes.btn}
            disabled={key}
          >
            ADD KEY
          </Button>
        </x.div>
      </CardContent>
    </Card>
  );
};

export default ApiKey;
