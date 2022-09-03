import React, { Children, useRef, useState } from "react";
import { Col, Dropdown, FormControl } from "react-bootstrap";
import { AiOutlineInfoCircle, AiOutlinePlus } from "react-icons/ai";
import { CustomToggle } from "../../components/dropdown/Dropdown";
import { CenteredModal } from "../../components/modal/Modal";
import { Pop } from "../../components/pop/Pop";
import { InteractionBadge, NetworkBadge } from "../../components/utils/Common";
import { UtilsBar } from "../../components/utils/Utils";
import { getPage } from "../../constants/globals";

import { Service } from "../../recoil/atoms/services";
import { Profile, profiles } from "../../recoil/atoms/profiles";
import { useRecoilCallback, useRecoilValue } from "recoil";
import {
  instances,
  instanceIds,
  intanceFormFieldErrors,
  instanceFormFields,
  Instance,
} from "../../recoil/atoms/instances";
import { SimpleForm } from "../../components/forms/Form";
import { InstanceFormFields } from "./InstanceForm";

const CustomInstanceDropdownItem = () => {
  const pageName = "Instances";
  const page = getPage(pageName);

  const [modalShow, setModalShow] = React.useState(false);

  const ids = useRecoilValue(instanceIds);
  const onSubmit = useRecoilCallback(({ set }) => (instance: Instance) => {
    const id = ids.length;
    // Set the new id in the list
    set(instanceIds, [...ids, id]);

    const newInstnace = {
      ...instance,
      id: id,
    };

    set(instances(id), (prev) => ({ ...prev, ...newInstnace }));
  });

  const defaultValues = useRecoilValue(instanceFormFields);

  const content = (
    <SimpleForm
      create={true}
      defaultValues={defaultValues}
      errors={intanceFormFieldErrors}
      onSubmit={onSubmit}
      page={pageName}
      fields={InstanceFormFields}
    />
  );

  const props = {
    title: "New Custom Instance",
    onHide: () => setModalShow(false),
    show: modalShow,
    content: content,
    icon: page?.icon,
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
              !value || child.props.profile.name.toLowerCase().startsWith(value)
          )}
        </ul>
        <CustomInstanceDropdownItem />
      </div>
    );
  }
);

const ProfilePop = ({ service }: { service: Service }) => {
  return (
    <Col>
      {service.name}
      <NetworkBadge {...service.network} />
      <InteractionBadge {...service.interaction} />
    </Col>
  );
};

export const ProfileRowInfoPop = ({ services }: { services: Service[] }) => {
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
      {!!services.length && (
        <Pop target={target} show={show} placement="left">
          {services.map((service) => {
            return <ProfilePop service={service} />;
          })}
        </Pop>
      )}
    </>
  );
};

const ProfileDropdownRow = ({ profile }: { profile: Profile }) => {
  const ids = useRecoilValue(instanceIds);

  const insertInstance = useRecoilCallback(
    ({ set }) =>
      (id: number, prof: Profile) => {
        // Set the new id in the list
        set(instanceIds, [...ids, id]);

        // Set the instance in teh family
        const newInstance = { name: prof.name, id: id, profile: prof };
        set(instances(id), (prev) => ({ ...prev, ...newInstance }));
      }
  );

  return (
    <Dropdown.Item
      onClick={() => {
        insertInstance(ids.length, profile);
      }}
    >
      {profile.name}
      <ProfileRowInfoPop services={profile.services} />
    </Dropdown.Item>
  );
};

const AddButton = () => {
  // Get all the profiles
  const profs = useRecoilValue(profiles);

  return (
    <Dropdown drop="start">
      <Dropdown.Toggle drop="start" as={CustomToggle} id={`dropdown-row-add`}>
        <AiOutlinePlus />
      </Dropdown.Toggle>
      <Dropdown.Menu as={InstancesAddProfileDropdownMenu}>
        {profs.map((profile: Profile) => {
          return <ProfileDropdownRow profile={profile} key={profile.id} />;
        })}
      </Dropdown.Menu>
    </Dropdown>
  );
};

export const InstancesUtils = () => {
  const buttons = [<AddButton />];
  return <UtilsBar buttons={buttons} />;
};
