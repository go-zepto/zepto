import React from 'react';
import { MenuItemLink, Responsive, Menu } from 'react-admin';
import { makeStyles } from '@material-ui/core/styles';
import { LayoutComponentGenerator, MenuLink, Schema } from '../../../../../types';


const useStyles = makeStyles({
  menu: {
      padding: '8px',
  },
});

const MenuWrapper = (props: any) => {
  const classes = useStyles();
  return (
    <div className={classes.menu}>
      {props.children}
    </div>
  );
}

const MenuGenerator: LayoutComponentGenerator = (s: Schema): React.FC => {
  if (s?.admin.menu?.links.length === 0) {
    return Menu;
  }
  return ({ onMenuClick, logout }: any) => {
    const renderLink = (link: MenuLink, idx: number) => {
      const to = link.link_to_path.length > 0 ? link.link_to_path : (
        '/' + s.admin.resources.find(r => r.name === link.link_to_resource_name)?.endpoint
      );
      return (
        <MenuItemLink
          sidebarIsOpen={false}
          key={idx}
          to={to}
          primaryText={link.label}
          leftIcon={<span className="material-icons">{link.icon || 'list_view'}</span>}
          onClick={onMenuClick}
        />
      )
    }
    return (
      <MenuWrapper>
        {s.admin.menu.links.map((link, idx)=> renderLink(link, idx))}
        <Responsive
          small={logout}
        />
      </MenuWrapper>
    );
  }
}

export default MenuGenerator;
