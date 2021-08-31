/* eslint-disable object-curly-newline */
import { useState, useEffect } from 'react';
import axios from 'axios';
import { Link as RouterLink } from 'react-router-dom';
import { useNavigate } from 'react-router-dom';
import styled from '@emotion/styled/macro';
import PropTypes from 'prop-types';
import {
  AppBar,
  Box,
  Hidden,
  IconButton,
  Toolbar,
  Divider,
  List
} from '@material-ui/core';
import { Drawer, Button } from 'antd';
import GitHubIcon from '@material-ui/icons/GitHub';
import InputIcon from '@material-ui/icons/Input';
import Logo from './Logo';
import MenuItem from '@material-ui/core/MenuItem';
import InputLabel from '@material-ui/core/InputLabel';
import NativeSelect from '@material-ui/core/Select';
import FormControl from '@material-ui/core/FormControl';
import { makeStyles } from '@material-ui/styles';
import { API_SERVER, GROUP_ID } from '../config';
import { useDispatch, useSelector } from 'react-redux';
import {
  filterByZone,
  getLoginStatus,
  getZoneList,
  getZoneName
} from './store/actions/klevrActions';

import NavItem from './NavItem';

import {
  BarChart as BarChartIcon,
  Settings as SettingsIcon,
  FileText as TaskIcon,
  Grid as ZoneIcon,
  Key as CredentialIcon,
  AlignLeft as LogIcon,
  UserCheck as AgentIcon,
  Menu as MenuIcon
} from 'react-feather';
import { justifyContent, x } from '@xstyled/emotion';

const items = [
  {
    href: '/app/overview',
    icon: BarChartIcon,
    title: 'Overview'
  },
  {
    href: '/app/tasks',
    icon: TaskIcon,
    title: 'Tasks'
  },
  {
    href: '/app/agent',
    icon: AgentIcon,
    title: 'Agent'
  },
  {
    href: '/app/credentials',
    icon: CredentialIcon,
    title: 'Credentials'
  },
  {
    href: '/app/logs',
    icon: LogIcon,
    title: 'Logs'
  },
  {
    href: '/app/settings',
    icon: SettingsIcon,
    title: 'Settings'
  }
];

const useStyles = makeStyles((theme) => ({
  formControl: {
    margin: 8,
    width: 160
  },
  selectEmpty: {
    marginTop: 16
  }
}));

const FirstZone = () => {
  const zoneList = useSelector((store) => store.zoneListReducer);

  const BlinkWrapper = styled.div`
    width: 200px;
    display: flex;
    justify-content: center;
    align-items: center;
    margin-right: 10px;
    animation: blink-effect 1.5s step-end infinite;

    @keyframes blink-effect {
      50% {
        opacity: 0;
      }
    }
  `;

  if (zoneList) {
    return null;
  }

  return (
    <BlinkWrapper>
      <NavItem
        href="/app/zones"
        key="Zones"
        title="Please add a zone first."
        icon={ZoneIcon}
      />
    </BlinkWrapper>
  );
};

const Zone = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const currentZone = useSelector((store) => store.zoneReducer);
  const zoneList = useSelector((store) => store.zoneListReducer);

  const classes = useStyles();
  useEffect(() => {
    let completed = false;
    async function get() {
      const result = await axios.get(`${API_SERVER}/inner/groups`);

      if (!completed)
        if (result.data === null) {
          dispatch(getZoneList(null));
          return;
        }
      dispatch(getZoneList(result.data));
      dispatch(filterByZone(result.data[0].Id));
      dispatch(getZoneName(result.data[0].GroupName));
    }
    get();
    return () => {
      completed = true;
    };
  }, []);

  if (!zoneList) {
    return null;
  }

  const selectZone = (id, groupName) => {
    dispatch(filterByZone(id));
    dispatch(getZoneName(groupName));
    navigate('/app/overview', { replace: true });
  };

  return (
    <>
      <x.div
        w="100px"
        display="flex"
        justifyContent="center"
        alignItems="center"
        mr="10px"
      >
        <NavItem href="/app/zones" key="Zones" title="Zones" icon={ZoneIcon} />
      </x.div>
      <FormControl className={classes.formControl}>
        <NativeSelect value={currentZone} variant="standard">
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
    </>
  );
};

const MobileMenu = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const [visible, setVisible] = useState(false);

  const showDrawer = () => {
    setVisible(true);
  };

  const onClose = () => {
    setVisible(false);
  };

  const signOutHandler = () => {
    async function signOut() {
      const result = await axios.get(`${API_SERVER}/console/signout`);

      if (result.status === 200) {
        dispatch(getLoginStatus(false));
        navigate('/login', { replace: true });
      }
    }
    signOut();
  };

  return (
    <>
      <Button type="text" onClick={visible ? onClose : showDrawer}>
        <x.div position="relative" top="-5px" left="-10px">
          <MenuIcon color="#50A1D6" size="30px" />
        </x.div>
      </Button>
      <Drawer
        title="Basic Drawer"
        placement="left"
        closable={false}
        onClose={onClose}
        visible={visible}
      >
        <Box
          sx={{
            display: 'flex',
            flexDirection: 'column',
            height: '100%'
          }}
        >
          <List>
            {items.map((item) => {
              if (item.title === 'Overview') {
                return (
                  <x.div>
                    <NavItem
                      href={item.href}
                      key={item.title}
                      title={item.title}
                      icon={item.icon}
                      onClick={onClose}
                    />
                  </x.div>
                );
              } else if (item.title === 'Settings') {
                return (
                  <x.div>
                    <NavItem
                      href={item.href}
                      key={item.title}
                      title={item.title}
                      icon={item.icon}
                      onClick={onClose}
                    />
                  </x.div>
                );
              } else {
                return (
                  <x.div ml="25px">
                    <NavItem
                      href={item.href}
                      key={item.title}
                      title={item.title}
                      icon={item.icon}
                      onClick={onClose}
                    />
                  </x.div>
                );
              }
            })}
          </List>
          <x.div position="absolute" bottom="30px">
            <IconButton color="inherit" onClick={signOutHandler}>
              <InputIcon />
            </IconButton>
          </x.div>
          <Box sx={{ flexGrow: 1 }} />
        </Box>
      </Drawer>
    </>
  );
};

const DashboardNavbar = ({ onMobileNavOpen, ...rest }) => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const pageCheck = window.location.pathname !== '/login';

  const signOutHandler = () => {
    async function signOut() {
      const result = await axios.get(`${API_SERVER}/console/signout`);

      if (result.status === 200) {
        dispatch(getLoginStatus(false));
        navigate('/login', { replace: true });
      }
    }
    signOut();
  };

  return (
    <AppBar elevation={0} {...rest}>
      <Toolbar>
        <Hidden lgUp>
          <MobileMenu />
        </Hidden>
        <RouterLink to={pageCheck ? '/app/overview' : '/'}>
          <Logo />
        </RouterLink>
        <Box sx={{ flexGrow: 1 }} />
        <Zone />
        <FirstZone />
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
