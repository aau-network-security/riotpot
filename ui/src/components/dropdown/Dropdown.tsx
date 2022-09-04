import "./Dropdown.scss";
import { Dropdown } from "react-bootstrap";
import React from "react";

type DropdownLink = {
  text: string;
  icon?: any;

  href?: string;
  onClick?: () => void;
};

type SimpleDropdownProps = {
  id?: string;
  icon: any;
  links?: DropdownLink[];
  children?: any;
};

export const CustomToggle = React.forwardRef(
  ({ children, onClick }: any, ref: any) => (
    // eslint-disable-next-line jsx-a11y/anchor-is-valid
    <a
      href=""
      ref={ref}
      onClick={(e) => {
        e.preventDefault();
        onClick(e);
      }}
    >
      {children}
    </a>
  )
);

export const SimpleDropdown = ({
  id,
  icon,
  links,
  children,
}: SimpleDropdownProps) => {
  return (
    <Dropdown>
      <Dropdown.Toggle as={CustomToggle} id={`dropdown-row-${id}`}>
        {icon}
      </Dropdown.Toggle>

      <Dropdown.Menu>
        {links?.map((link) => {
          return (
            <Dropdown.Item {...link}>
              <link.icon />
              {link.text}
            </Dropdown.Item>
          );
        })}
        {children}
      </Dropdown.Menu>
    </Dropdown>
  );
};
