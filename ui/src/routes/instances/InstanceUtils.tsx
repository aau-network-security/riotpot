import React, { Children, useRef, useState } from "react";
import { Col, Dropdown, FormControl, Row } from "react-bootstrap";
import { AiOutlineInfoCircle, AiOutlinePlus } from "react-icons/ai";
import { RiProfileLine } from "react-icons/ri";
import { useRecoilCallback, useRecoilValue } from "recoil";
import { CustomToggle } from "../../components/dropdown/Dropdown";
import { Pop } from "../../components/pop/Pop";
import { useToast } from "../../components/toast/Toast";
import { InteractionBadge, NetworkBadge } from "../../components/utils/Common";
import { UtilsBar } from "../../components/utils/Utils";
import {
  DefaultInstanceProxy,
  GetInstanceAddress,
  Instance,
  instance,
  InstanceProxy,
  useInstanceProxy,
} from "../../recoil/atoms/instances";
import { Profile, profiles } from "../../recoil/atoms/profiles";
import { Service, services } from "../../recoil/atoms/services";
import { ErrorToastVariant } from "../../recoil/atoms/toast";
import {
  addFromProfile,
  addProxyService,
  deleteProxyService,
} from "./InstanceAPI";

type CustomMenuProps = {
  children?: React.ReactNode;
  style?: React.CSSProperties;
  className?: string;
  labeledBy?: string;
};

// Dropdown style for services with a string filter for names
const InstanceAddServiceDropdownMenu = React.forwardRef(
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
        <Dropdown.Divider />
        <ul className="list-unstyled">
          {Children.toArray(props.children).filter(
            (child: any) =>
              !value ||
              child.props.service?.name.toLowerCase().startsWith(value)
          )}
        </ul>
      </div>
    );
  }
);

// Dropdown item to add a service
const ServiceDropdownItem = ({ service }: { service: Service }) => {
  const { showToast } = useToast();
  const { registerProxy } = useInstanceProxy();
  const address = GetInstanceAddress();

  // Small handler to send a request to the backend to register this service
  const addService = () => {
    const response = addProxyService(address, service);

    response
      .then((data) => {
        if ("error" in data) {
          showToast(data["error"], ErrorToastVariant);
          return;
        }
        registerProxy(data);
      })
      .catch((err) => {
        console.log(err);
      });
  };

  return (
    <Dropdown.Item onClick={() => addService()}>
      <Col>{service.name}</Col>
      <Col className="badges">
        <NetworkBadge {...service.network} />
        <InteractionBadge {...service.interaction} />
      </Col>
    </Dropdown.Item>
  );
};

// Button that contains the dropdown to add a service
const AddServiceDropdownButton = () => {
  // Get the ids and the services
  const servs = useRecoilValue(services);

  return (
    <Dropdown drop="start">
      <Dropdown.Toggle drop="start" as={CustomToggle} id={`dropdown-row-add`}>
        <AiOutlinePlus />
      </Dropdown.Toggle>
      <Dropdown.Menu as={InstanceAddServiceDropdownMenu}>
        {servs.map((service: Service) => {
          return <ServiceDropdownItem service={service} key={service.id} />;
        })}
      </Dropdown.Menu>
    </Dropdown>
  );
};

const InstanceChangeProfileDropdownMenu = React.forwardRef(
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
              !value ||
              child.props.profile?.name.toLowerCase().startsWith(value)
          )}
        </ul>
      </div>
    );
  }
);

const ProfilePop = ({ service }: { service: Service }) => {
  return (
    <Row>
      <Col className="profile">{service.name}</Col>
      <Col className="badges">
        <NetworkBadge {...service.network} />
        <InteractionBadge {...service.interaction} />
      </Col>
    </Row>
  );
};

export const ProfileRowInfoPop = ({ services }: { services: Service[] }) => {
  const [show, setShow] = useState(false);
  const target = useRef(null);

  return (
    <Row>
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
            return <ProfilePop key={service.id} service={service} />;
          })}
        </Pop>
      )}
    </Row>
  );
};

const InstanceChangeProfileDropdownRow = ({
  profile,
  replaceProfileCallback,
}: {
  profile: Profile;
  replaceProfileCallback: Function;
}) => {
  return (
    <Dropdown.Item onClick={() => replaceProfileCallback(profile)}>
      {profile.name}
      <ProfileRowInfoPop services={profile.services} />
    </Dropdown.Item>
  );
};

// Button to change the profile of the instance
const InstanceChangeProfile = () => {
  const { removeProxyFromService, registerProxy } = useInstanceProxy();

  // Get all the profiles and the instance
  const profs = useRecoilValue(profiles);
  const inst = useRecoilValue(instance);
  const address = GetInstanceAddress();

  // Handler to remove the current profile proxies using its services
  // Note: It does not remove custom proxies
  const removeCurrentProfileServices = () => {
    const currProfile = inst.profile;

    // Delete all the services included in the current profile and registered
    if (currProfile) {
      for (const serv of currProfile?.services) {
        // Send a request to delete the proxy service
        const deleted = deleteProxyService(address, serv.id);

        deleted
          .then((data) => {
            if ("success" in data) {
              // If succeeded, remove the service from the list
              removeProxyFromService(serv.id);
            }
          })
          .catch((err) => {
            console.log(err);
          });
      }
    }
  };

  // Handler for registering proxies from profile services
  // It makes a call to the API to create them from the profile
  const registerProxies = async (profile: Profile) => {
    return await addFromProfile(address, profile.services).then(
      async (proxs) => {
        // Register all the new proxy
        proxs.forEach((prox) => {
          registerProxy(prox);
        });
        return {
          ...profile,
          services: proxs.map((prox) => prox.service),
        };
      }
    );
  };

  // Handler to replace the previous profile registered proxies with the new ones
  // This deletes ALL the proxies registered from the profile,
  // but it leaves the custom created ones
  const replaceProfileServices = useRecoilCallback(
    ({ set }) =>
      async (profile: Profile) => {
        removeCurrentProfileServices();
        const newProfile = await registerProxies(profile);

        set(instance, (prev: Instance) => ({ ...prev, profile: newProfile }));
      }
  );

  return (
    <Dropdown drop="start">
      <Dropdown.Toggle drop="start" as={CustomToggle} id={`dropdown-row-add`}>
        <RiProfileLine />
      </Dropdown.Toggle>
      <Dropdown.Menu as={InstanceChangeProfileDropdownMenu}>
        {profs.map((profile: Profile) => (
          <InstanceChangeProfileDropdownRow
            key={profile.id}
            profile={profile}
            replaceProfileCallback={replaceProfileServices}
          />
        ))}
      </Dropdown.Menu>
    </Dropdown>
  );
};

export const InstanceUtils = () => {
  const buttons = [
    <InstanceChangeProfile key={1} />,
    <AddServiceDropdownButton key={0} />,
  ];
  return <UtilsBar buttons={buttons} />;
};

export default InstanceUtils;
