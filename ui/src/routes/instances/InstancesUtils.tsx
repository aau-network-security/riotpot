import React, { Children, useRef, useState } from "react";
import { Col, Dropdown, FormControl } from "react-bootstrap";
import { AiOutlineInfoCircle, AiOutlinePlus } from "react-icons/ai";
import { CustomToggle } from "../../components/dropdown/Dropdown";
import { CenteredModal } from "../../components/modal/Modal";
import { Pop } from "../../components/pop/Pop";
import { InteractionBadge, NetworkBadge } from "../../components/utils/Common";
import { Submit, UtilsBar } from "../../components/utils/Utils";
import {
  getPage,
  InteractionOptions,
  NetworkOptions,
} from "../../constants/globals";
import { ProfileService, Profile } from "../profiles/ProfilesTable";
import InstanceForm from "./InstanceForm";

const CustomInstanceDropdownItem = () => {
  const [modalShow, setModalShow] = React.useState(false);
  const icon = getPage("Instances")?.icon;
  const content = <InstanceForm />;

  const props = {
    title: "New Custom Instance",
    icon: icon,
    onHide: () => setModalShow(false),
    show: modalShow,
    submit: <Submit />,
    content: content,
  };

  return (
    <>
      <Dropdown.Item onClick={() => setModalShow(true)} id="custom add">
        <AiOutlinePlus />
        Custom Instance
      </Dropdown.Item>
      <CenteredModal props={props} />
    </>
  );
};

type CustomMenuProps = {
  children?: React.ReactNode;
  style?: React.CSSProperties;
  className?: string;
  labeledBy?: string;
};

const InstancesAddProfileDropdownMenu = React.forwardRef(
  (props: CustomMenuProps, ref: React.Ref<HTMLDivElement>) => {
    const [value, setValue] = useState("");

    return (
      <div
        ref={ref}
        style={props.style}
        className={props.className}
        aria-labelledby={props.labeledBy}
      >
        <FormControl
          autoFocus
          className="mx-3 my-2 w-auto"
          placeholder="Type to filter..."
          onChange={(e) => setValue(e.target.value)}
          value={value}
        />
        <ul className="list-unstyled">
          {Children.toArray(props.children).filter(
            (child: any) =>
              !value || child.props.name.toLowerCase().startsWith(value)
          )}
        </ul>
        <CustomInstanceDropdownItem />
      </div>
    );
  }
);

const ProfilePop = ({ service }: any) => {
  return (
    <Col>
      {service.name}
      <NetworkBadge {...service.network} />
      <InteractionBadge {...service.interaction} />
    </Col>
  );
};

export const ProfileRowInfoPop = ({
  services,
}: {
  services: ProfileService[];
}) => {
  const [show, setShow] = useState(false);
  const target = useRef(null);

  return (
    <>
      <small
        className="info"
        ref={target}
        onMouseEnter={() => setShow(true)}
        onMouseLeave={() => setShow(false)}
      >
        <AiOutlineInfoCircle />
        {`${services.length} ${services.length === 1 ? "service" : "services"}`}
      </small>
      <Pop target={target} show={show} placement="left">
        {services.map((service) => {
          return <ProfilePop service={service} />;
        })}
      </Pop>
    </>
  );
};

const ProfileDropdownRow = ({ name, services }: any) => {
  return (
    <Dropdown.Item eventKey={name.toLowerCase()} key={name.toLowerCase()}>
      {name}
      <ProfileRowInfoPop services={services} />
    </Dropdown.Item>
  );
};

const AddButton = ({ profiles }: { profiles: Profile[] }) => {
  return (
    <Dropdown>
      <Dropdown.Toggle drop="start" as={CustomToggle} id={`dropdown-row-add`}>
        <AiOutlinePlus />
      </Dropdown.Toggle>
      <Dropdown.Menu as={InstancesAddProfileDropdownMenu}>
        {profiles.map((profile) => {
          return (
            <ProfileDropdownRow
              name={profile.name}
              services={profile.services}
            />
          );
        })}
      </Dropdown.Menu>
    </Dropdown>
  );
};

export const InstancesUtils = ({ view }: { view: string }) => {
  const profiles = [
    {
      name: "Wi-Fi Printer",
      services: [
        {
          name: "CoAP",
          interaction: InteractionOptions[0],
          network: NetworkOptions[0],
        },
      ],
    },
  ];

  const buttons = [<AddButton profiles={profiles} />];
  return <UtilsBar buttons={buttons} />;
};
