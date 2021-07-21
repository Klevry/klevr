/* eslint-disable object-curly-newline */
import { useState, useEffect } from 'react';
import axios from 'axios';
import { Link as RouterLink } from 'react-router-dom';
import { useNavigate } from 'react-router-dom';
import PropTypes from 'prop-types';
import { AppBar, Box, Hidden, IconButton, Toolbar } from '@material-ui/core';
import GitHubIcon from '@material-ui/icons/GitHub';
import InputIcon from '@material-ui/icons/Input';
import Logo from './Logo';
import MenuItem from '@material-ui/core/MenuItem';
import InputLabel from '@material-ui/core/InputLabel';
import NativeSelect from '@material-ui/core/Select';
import FormControl from '@material-ui/core/FormControl';
import { makeStyles } from '@material-ui/core/styles';
import { API_SERVER, GROUP_ID } from '../config';
import { useDispatch, useSelector } from 'react-redux';
import {
  filterByZone,
  getZoneList,
  getZoneName
} from './store/actions/klevrActions';

const useStyles = makeStyles((theme) => ({
  formControl: {
    margin: theme.spacing(1),
    width: 160
  },
  selectEmpty: {
    marginTop: theme.spacing(2)
  }
}));

const Zone = () => {
  const dispatch = useDispatch();
  const currentZone = useSelector((store) => store.zoneReducer);
  const zoneList = useSelector((store) => store.zoneListReducer);
  const [iz, setIz] = useState();

  const classes = useStyles();
  useEffect(() => {
    let completed = false;
    async function get() {
      const result = await axios.get(`${API_SERVER}/inner/groups`);
      if (!completed) dispatch(getZoneList(result.data));
      dispatch(filterByZone(result.data[0].Id));
      dispatch(getZoneName(result.data[0].GroupName));
      setIz(result.data[0].Id);
    }
    get();
    return () => {
      completed = true;
    };
  }, []);

  if (!zoneList) {
    return null;
  }

  if (!iz) {
    return null;
  }

  const selectZone = (id, groupName) => {
    dispatch(filterByZone(id));
    dispatch(getZoneName(groupName));
  };

  return (
    <FormControl className={classes.formControl}>
      {/* <InputLabel style={{ color: 'white', fontWeight: 'bold' }}>
        Zone
      </InputLabel> */}
      <NativeSelect defaultValue={iz}>
        {zoneList.map((item) => (
          <MenuItem
            value={item.Id}
            key={item.Id}
            onClick={() => selectZone(item.Id, item.GroupName)}
          >
            {`${item.GroupName} (${item.Id})`}
          </MenuItem>
        ))}
      </NativeSelect>
    </FormControl>
  );
};

const DashboardNavbar = ({ onMobileNavOpen, ...rest }) => {
  const navigate = useNavigate();
  const pageCheck = window.location.pathname !== '/login';

  const signOutHandler = () => {
    async function signOut() {
      const result = await axios.get(`${API_SERVER}/console/signout`);

      if (result.status === 200) {
        navigate('/login', { replace: true });
      }
    }
    signOut();
  };

  return (
    <AppBar elevation={0} {...rest}>
      <Toolbar>
        <RouterLink to={pageCheck ? '/app/overview' : '/'}>
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
          <IconButton color="inherit" onClick={signOutHandler}>
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
