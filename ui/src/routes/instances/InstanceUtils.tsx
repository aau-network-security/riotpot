import React, { Children, useState } from "react";
import { Col, Dropdown, FormControl } from "react-bootstrap";
import { AiOutlinePlus } from "react-icons/ai";
import { useRecoilCallback, useRecoilValue } from "recoil";
import { CustomToggle } from "../../components/dropdown/Dropdown";
import { InteractionBadge, NetworkBadge } from "../../components/utils/Common";
import { UtilsBar } from "../../components/utils/Utils";
import {
  InstanceProxyService,
  instanceService,
  instanceServiceIDs,
} from "../../recoil/atoms/instances";
import { Service, services } from "../../recoil/atoms/services";
import { addProxyService } from "./InstanceAPI";

type CustomMenuProps = {
  children?: React.ReactNode;
  style?: React.CSSProperties;
  className?: string;
  labeledBy?: string;
};

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

const ServiceDropdownRow = ({
  handler,
  service,
}: {
  handler: (serv: Service) => void;
  service: Service;
}) => {
  return (
    <Dropdown.Item onClick={() => handler(service)}>
      <Col>{service.name}</Col>
      <Col className="badges">
        <NetworkBadge {...service.network} />
        <InteractionBadge {...service.interaction} />
      </Col>
    </Dropdown.Item>
  );
};

const AddButton = ({ host }: { host: string }) => {
  // Get the ids and the services
  const servs = useRecoilValue(services);

  // Get the list of proxy service ids
  const ids = useRecoilValue(instanceServiceIDs);

  // Callback to add a service to the list.
  // This is used to track and update the state of the proxies
  const addToList = useRecoilCallback(
    ({ set }) =>
      (proxyService: InstanceProxyService) => {
        // Set the new id in the list if it is not there yet
        if (!ids.includes(proxyService.id)) {
          set(instanceServiceIDs, [...ids, proxyService.id]);
        }

        set(instanceService(proxyService.id), proxyService);
      }
  );

  const handler = (service: Service) => {
    const response = addProxyService(host, service);

    response
      .then((data) => {
        if ("error" in data) {
          return;
        }
        addToList(data);
      })
      .catch((err) => {
        console.log(err);
      });
  };

  return (
    <Dropdown drop="start">
      <Dropdown.Toggle drop="start" as={CustomToggle} id={`dropdown-row-add`}>
        <AiOutlinePlus />
      </Dropdown.Toggle>
      <Dropdown.Menu as={InstanceAddServiceDropdownMenu}>
        {servs.map((service: Service) => {
          return (
            <ServiceDropdownRow
              service={service}
              key={service.id}
              handler={handler}
            />
          );
        })}
      </Dropdown.Menu>
    </Dropdown>
  );
};

export const InstanceUtils = ({ host }: { host: string }) => {
  const buttons = [<AddButton key={0} host={host} />];
  return <UtilsBar buttons={buttons} />;
};

export default InstanceUtils;
