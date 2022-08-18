import "./Common.scss";

import { Badge, Button, Dropdown } from "react-bootstrap";
import { getPage, InteractionOption, Page } from "../../constants/globals";
import { FiEdit2, FiMoreHorizontal } from "react-icons/fi";
import { SimpleDropdown } from "../dropdown/Dropdown";
import { CgDetailsLess } from "react-icons/cg";
import React from "react";
import { RiDeleteBinLine } from "react-icons/ri";
import { CenteredModal } from "../modal/Modal";

type BaseBadgeProps = {
  text: string;
  className: string;
};

const BaseBadge = ({ text, className }: BaseBadgeProps) => {
  return <Badge className={className}>{text}</Badge>;
};

export const InteractionBadge = ({ value, label }: InteractionOption) => {
  return <BaseBadge className={`interaction ${value}`} text={label} />;
};

export const NetworkBadge = ({ value, label }: InteractionOption) => {
  return <BaseBadge className={`network ${value}`} text={label} />;
};

type AddressProps = {
  host: string;
  port: Number;
};

export const Address = ({ host, port }: AddressProps) => {
  const page = getPage("Services");

  return (
    <span className="address">
      {page && <page.icon />} {host} : 1111
    </span>
  );
};

export const OptionsDropdown = ({ children }: any) => {
  const links = [
    { text: "Edit", icon: FiEdit2 },
    { text: "Details", icon: CgDetailsLess },
  ];

  const props = {
    icon: <FiMoreHorizontal />,
    links: links,
  };

  return <SimpleDropdown {...props}>{children}</SimpleDropdown>;
};

export const DeleteButton = () => {
  const onclick = () => {};

  return (
    <Button onClick={onclick} variant="danger">
      Delete
    </Button>
  );
};

export type DeleteProps = {
  page: Page;
  name: string;
  note?: string;
};

export const DeleteDropdownItem = ({ page, note, name }: DeleteProps) => {
  const [modalShow, setModalShow] = React.useState(false);

  const content = () => {
    const msg = `Are you sure you want to delete the following ${page.verbose.toLowerCase()}?`;
    const pstyle = {
      color: "#FF8686",
    };
    const sub = ["This action is irreversible.", note].join(" ");

    return (
      <>
        <h5>{msg}</h5>
        <ul>
          <li style={pstyle}>{name}</li>
        </ul>
        <small>{sub}</small>
      </>
    );
  };

  const props = {
    title: `Delete ${page.page}`,
    onHide: () => setModalShow(false),
    icon: page.icon,
    show: modalShow,
    submit: <DeleteButton />,
    content: content(),
  };

  return (
    <>
      <Dropdown.Item onClick={() => setModalShow(true)}>
        <RiDeleteBinLine />
        Delete
      </Dropdown.Item>
      <CenteredModal props={props}></CenteredModal>
    </>
  );
};
