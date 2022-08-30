import "./Common.scss";

import { Badge, Button, Dropdown } from "react-bootstrap";
import { getPage, InteractionOption, Page } from "../../constants/globals";
import { FiEdit2, FiMoreHorizontal } from "react-icons/fi";
import { SimpleDropdown } from "../dropdown/Dropdown";
import { CgDetailsLess } from "react-icons/cg";
import { RiDeleteBinLine } from "react-icons/ri";
import { CenteredModal } from "../modal/Modal";
import { useState } from "react";

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
  port: string | Number;
};

export const Address = ({ host, port }: AddressProps) => {
  const page = getPage("Services");

  return (
    <span className="address">
      {page && <page.icon />} {`${host} : ${port}`}
    </span>
  );
};

export const OptionsDropdown = ({ children }: any) => {
  const links = [
    //{ text: "Edit", icon: FiEdit2 },
    { text: "Details", icon: CgDetailsLess },
  ];

  const props = {
    icon: <FiMoreHorizontal />,
    links: links,
  };

  return <SimpleDropdown {...props}>{children}</SimpleDropdown>;
};

type EditProps = {
  icon: any;
  title: string;
  form: any;
};

export const EditDropdownItem = ({ icon, title, form }: EditProps) => {
  const [modalShow, setModalShow] = useState(false);

  const modalprops = {
    content: form,
    icon: icon,
    title: "Edit " + title,
    onHide: () => setModalShow(false),
    show: modalShow,
  };

  return (
    <>
      <Dropdown.Item onClick={() => setModalShow(true)}>
        <FiEdit2 />
        Edit
      </Dropdown.Item>
      <CenteredModal props={modalprops}></CenteredModal>
    </>
  );
};

export type DeleteProps = {
  page: Page;
  name: string;
  note?: string;
  onClick: () => void;
};

export const DeleteButton = ({ onClick }: { onClick?: () => void }) => {
  return (
    <Button onClick={onClick} variant="danger">
      Delete
    </Button>
  );
};

export const DeleteDropdownItem = ({
  page,
  note,
  name,
  onClick,
}: DeleteProps) => {
  const [modalShow, setModalShow] = useState(false);

  const Content = () => {
    const msg = `Are you sure you want to delete the following ${page.verbose.toLowerCase()}?`;
    const pstyle = {
      color: "#FF8686",
    };
    const sub = ["This action is irreversible.", note].join(" ");

    return (
      <>
        <div>
          <h5>{msg}</h5>
          <ul>
            <li style={pstyle}>{name}</li>
          </ul>
          <small>{sub}</small>
        </div>
        <div>
          <DeleteButton
            onClick={() => {
              onClick();
              setModalShow(false);
            }}
          />
        </div>
      </>
    );
  };

  const props = {
    title: `Delete ${page.page}`,
    onHide: () => setModalShow(false),
    icon: page.icon,
    show: modalShow,
    content: <Content />,
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
