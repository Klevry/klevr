/* eslint-disable object-curly-newline */
import { useState, useEffect } from 'react';
import axios from 'axios';
import { Link as RouterLink } from 'react-router-dom';
import PropTypes from 'prop-types';
import { AppBar, Box, Hidden, IconButton, Toolbar } from '@material-ui/core';
import GitHubIcon from '@material-ui/icons/GitHub';
import InputIcon from '@material-ui/icons/Input';
import Logo from './Logo';
import MenuItem from '@material-ui/core/MenuItem';
import InputLabel from '@material-ui/core/InputLabel';
import Select from '@material-ui/core/Select';
import FormControl from '@material-ui/core/FormControl';
import { makeStyles } from '@material-ui/core/styles';
import { API_SERVER, GROUP_ID } from '../config';

const useStyles = makeStyles((theme) => ({
  formControl: {
    margin: theme.spacing(1),
    width: 120
  },
  selectEmpty: {
    marginTop: theme.spacing(2)
  }
}));

const Zone = () => {
  const [data, setData] = useState(null);
  const classes = useStyles();
  useEffect(() => {
    let completed = false;
    async function get() {
      const result = await axios.get(`${API_SERVER}/inner/groups`, {
        withCredentials: true
      });
      if (!completed) setData(result.data);
    }
    get();
    return () => {
      completed = true;
    };
  }, []);

  if (!data) {
    return null;
  }
  return (
    <FormControl className={classes.formControl}>
      <InputLabel style={{ color: 'white', fontWeight: 'bold' }}>
        {GROUP_ID}
      </InputLabel>
      <Select disabled>
        {data.map((item) => (
          <MenuItem value={item.GroupName} key={item.Id}>
            {item.GroupName}
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );
};

const DashboardNavbar = ({ onMobileNavOpen, ...rest }) => {
  return (
    <AppBar elevation={0} {...rest}>
      <Toolbar>
        <RouterLink to="/">
          <Logo />
        </RouterLink>
        <Box sx={{ flexGrow: 1 }} />
        <Zone />
        <Hidden lgDown>
          <IconButton color="default">
            <a
              href="https://github.com/Klevry/klevr"
              target="_blank"
              rel="noreferrer"
            >
              <GitHubIcon />
            </a>
          </IconButton>
        </Hidden>
        <Hidden lgDown>
          <IconButton color="inherit">
            <InputIcon />
          </IconButton>
        </Hidden>
      </Toolbar>
    </AppBar>
  );
};

DashboardNavbar.propTypes = {
  onMobileNavOpen: PropTypes.func
};

export default DashboardNavbar;
